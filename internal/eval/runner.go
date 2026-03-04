package eval

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type RunOptions struct {
	RepoRoot    string
	TasksDir    string
	FixturesDir string
	Thresholds  string
	TaskIDs     map[string]bool
	OutputPath  string
}

func Run(opts RunOptions) (RunResult, error) {
	thresholds, err := loadThresholds(opts.Thresholds)
	if err != nil {
		return RunResult{}, err
	}

	taskDirs, err := discoverTaskDirs(opts.TasksDir, opts.TaskIDs)
	if err != nil {
		return RunResult{}, err
	}
	if len(taskDirs) == 0 {
		return RunResult{}, fmt.Errorf("no tasks found")
	}

	results := make([]TaskResult, 0, len(taskDirs))
	for _, td := range taskDirs {
		r, err := runTask(opts, td, thresholds)
		if err != nil {
			return RunResult{}, err
		}
		results = append(results, r)
	}

	run := aggregate(results, thresholds)
	if opts.OutputPath != "" {
		if err := writeResult(opts.OutputPath, run); err != nil {
			return RunResult{}, err
		}
	}

	return run, nil
}

func runTask(opts RunOptions, taskDir string, thresholds Thresholds) (TaskResult, error) {
	cfg, checks, err := loadTaskFiles(taskDir)
	if err != nil {
		return TaskResult{}, err
	}

	taskID := filepath.Base(taskDir)
	fixturePath := filepath.Join(opts.FixturesDir, cfg.Fixture)
	if _, err := os.Stat(fixturePath); err != nil {
		return TaskResult{}, fmt.Errorf("fixture %q not found for task %q: %w", cfg.Fixture, taskID, err)
	}

	workspace, err := os.MkdirTemp("", "spec-agent-eval-")
	if err != nil {
		return TaskResult{}, err
	}
	defer os.RemoveAll(workspace)

	if err := copyDir(fixturePath, workspace); err != nil {
		return TaskResult{}, err
	}

	before, err := snapshot(workspace)
	if err != nil {
		return TaskResult{}, err
	}

	started := time.Now()
	runErr := ""
	if strings.TrimSpace(cfg.RunCommand) != "" {
		cmd := exec.Command("sh", "-lc", cfg.RunCommand)
		cmd.Env = append(os.Environ(),
			"WORKSPACE="+workspace,
			"TASK_DIR="+taskDir,
			"REPO_ROOT="+opts.RepoRoot,
		)
		cmd.Dir = opts.RepoRoot
		if out, err := cmd.CombinedOutput(); err != nil {
			runErr = fmt.Sprintf("%v\n%s", err, strings.TrimSpace(string(out)))
		}
	}
	finished := time.Now()

	after, err := snapshot(workspace)
	if err != nil {
		return TaskResult{}, err
	}

	changed := changedPaths(before, after)
	scopeViolations, scopePaths := countScopeViolations(changed, cfg.AllowedPaths)
	completedChecks, totalChecks := evaluateChecks(workspace, checks)
	completeness := 1.0
	if totalChecks > 0 {
		completeness = float64(completedChecks) / float64(totalChecks)
	}
	passedTests, totalTests, testFailures := runTestCommands(workspace, cfg.TestCommands)
	testPassRate := 1.0
	if totalTests > 0 {
		testPassRate = float64(passedTests) / float64(totalTests)
	}
	specViolations := countSpecRuleViolations(workspace, changed)

	res := TaskResult{
		TaskID:              taskID,
		TaskName:            cfg.Name,
		Fixture:             cfg.Fixture,
		StartedAt:           started,
		FinishedAt:          finished,
		RunCommand:          cfg.RunCommand,
		RunError:            runErr,
		ChangedPaths:        changed,
		Completeness:        completeness,
		CompletedChecks:     completedChecks,
		TotalChecks:         totalChecks,
		ScopeViolations:     scopeViolations,
		ScopeViolationPaths: scopePaths,
		TestPassRate:        testPassRate,
		TestsPassed:         passedTests,
		TestsTotal:          totalTests,
		TestFailures:        testFailures,
		SpecRuleViolations:  specViolations,
	}
	res.GatingPass = evaluateTaskGating(res, thresholds)
	return res, nil
}

func evaluateTaskGating(r TaskResult, t Thresholds) bool {
	if r.Completeness < t.CompletenessMin {
		return false
	}
	if r.TestPassRate < t.TestPassRateMin {
		return false
	}
	if r.ScopeViolations > t.ScopeViolationsMax {
		return false
	}
	if r.SpecRuleViolations > t.SpecRuleViolationsMax {
		return false
	}
	if strings.TrimSpace(r.RunError) != "" {
		return false
	}
	return true
}

func aggregate(tasks []TaskResult, thresholds Thresholds) RunResult {
	res := RunResult{GeneratedAt: time.Now(), Thresholds: thresholds, Tasks: tasks, GatingPass: true}
	if len(tasks) == 0 {
		return res
	}
	for _, t := range tasks {
		res.CompletenessAvg += t.Completeness
		res.TestPassRateAvg += t.TestPassRate
		res.ScopeViolationsTotal += t.ScopeViolations
		res.SpecRuleViolationsTotal += t.SpecRuleViolations
		if !t.GatingPass {
			res.GatingPass = false
		}
	}
	res.CompletenessAvg /= float64(len(tasks))
	res.TestPassRateAvg /= float64(len(tasks))
	return res
}

func discoverTaskDirs(root string, filter map[string]bool) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		id := e.Name()
		if len(filter) > 0 && !filter[id] {
			continue
		}
		taskDir := filepath.Join(root, id)
		if _, err := os.Stat(filepath.Join(taskDir, "config.yaml")); err == nil {
			out = append(out, taskDir)
		}
	}
	sort.Strings(out)
	return out, nil
}

func loadThresholds(path string) (Thresholds, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Thresholds{}, err
	}
	var t Thresholds
	if err := yaml.Unmarshal(data, &t); err != nil {
		return Thresholds{}, err
	}
	if t.CompletenessMin == 0 {
		t.CompletenessMin = 0.9
	}
	if t.TestPassRateMin == 0 {
		t.TestPassRateMin = 0.95
	}
	return t, nil
}

func loadTaskFiles(taskDir string) (TaskConfig, ChecksConfig, error) {
	cfgData, err := os.ReadFile(filepath.Join(taskDir, "config.yaml"))
	if err != nil {
		return TaskConfig{}, ChecksConfig{}, err
	}
	checksData, err := os.ReadFile(filepath.Join(taskDir, "checks.yaml"))
	if err != nil {
		return TaskConfig{}, ChecksConfig{}, err
	}
	var cfg TaskConfig
	if err := yaml.Unmarshal(cfgData, &cfg); err != nil {
		return TaskConfig{}, ChecksConfig{}, err
	}
	var checks ChecksConfig
	if err := yaml.Unmarshal(checksData, &checks); err != nil {
		return TaskConfig{}, ChecksConfig{}, err
	}
	return cfg, checks, nil
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if strings.HasPrefix(rel, ".git") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Close()
}

func snapshot(root string) (map[string][32]byte, error) {
	m := map[string][32]byte{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		m[filepath.ToSlash(rel)] = sha256.Sum256(data)
		return nil
	})
	return m, err
}

func changedPaths(before, after map[string][32]byte) []string {
	changes := map[string]bool{}
	for p, h := range after {
		if old, ok := before[p]; !ok || old != h {
			changes[p] = true
		}
	}
	for p := range before {
		if _, ok := after[p]; !ok {
			changes[p] = true
		}
	}
	out := make([]string, 0, len(changes))
	for p := range changes {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

func countScopeViolations(paths, allowed []string) (int, []string) {
	if len(allowed) == 0 {
		return 0, nil
	}
	violations := make([]string, 0)
	for _, p := range paths {
		if !isAllowedPath(p, allowed) {
			violations = append(violations, p)
		}
	}
	return len(violations), violations
}

func isAllowedPath(path string, allowed []string) bool {
	for _, a := range allowed {
		a = filepath.ToSlash(strings.TrimSpace(a))
		if a == "" {
			continue
		}
		if strings.HasSuffix(a, "/**") {
			prefix := strings.TrimSuffix(a, "/**")
			if path == prefix || strings.HasPrefix(path, prefix+"/") {
				return true
			}
			continue
		}
		if ok, _ := filepath.Match(a, path); ok {
			return true
		}
		if path == a {
			return true
		}
	}
	return false
}

func evaluateChecks(workspace string, checks ChecksConfig) (int, int) {
	passed := 0
	total := len(checks.Checkpoints)
	for _, c := range checks.Checkpoints {
		if checkpointPasses(workspace, c) {
			passed++
		}
	}
	return passed, total
}

func checkpointPasses(workspace string, c Checkpoint) bool {
	switch c.Type {
	case "file_exists":
		_, err := os.Stat(filepath.Join(workspace, c.Path))
		return err == nil
	case "glob_exists":
		matches, _ := filepath.Glob(filepath.Join(workspace, c.Path))
		return len(matches) > 0
	case "file_contains":
		data, err := os.ReadFile(filepath.Join(workspace, c.Path))
		if err != nil {
			return false
		}
		return strings.Contains(string(data), c.Pattern)
	default:
		return false
	}
}

func runTestCommands(workspace string, cmds []string) (int, int, []string) {
	if len(cmds) == 0 {
		return 0, 0, nil
	}
	passed := 0
	failures := make([]string, 0)
	total := 0
	for _, c := range cmds {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		total++
		cmd := exec.Command("sh", "-lc", c)
		cmd.Dir = workspace
		out, err := cmd.CombinedOutput()
		if err == nil {
			passed++
			continue
		}
		details := strings.TrimSpace(string(out))
		if details == "" {
			details = "(no output)"
		}
		failures = append(failures, fmt.Sprintf("command: %s\nerror: %v\noutput:\n%s", c, err, details))
	}
	return passed, total, failures
}

func countSpecRuleViolations(workspace string, changed []string) int {
	violations := 0
	for _, rel := range changed {
		if !strings.HasSuffix(rel, ".md") {
			continue
		}
		if strings.HasPrefix(rel, ".spec_agent/") {
			continue
		}
		if strings.HasPrefix(rel, "spec_changes/") {
			continue
		}
		full := filepath.Join(workspace, rel)
		data, err := os.ReadFile(full)
		if err != nil {
			continue
		}
		content := string(data)
		if !strings.Contains(content, "<!-- SPEC:FILE=true -->") {
			continue
		}
		violations += validateSpecContent(content)
	}
	return violations
}

func validateSpecContent(content string) int {
	violations := 0
	requiredMeta := []string{
		"<!-- SPEC:START -->",
		"<!-- SPEC:FILE=true -->",
		"<!-- SPEC:ID=",
		"<!-- SPEC:KIND=",
		"<!-- SPEC:MENU=",
		"<!-- SPEC:END -->",
	}
	for _, m := range requiredMeta {
		if !strings.Contains(content, m) {
			violations++
		}
	}

	requiredSections := []string{
		"# ",
		"## Responsibility",
		"## Inputs",
		"## Outputs",
		"## Business Logic",
		"## Flow",
		"## Links",
		"## Dependencies",
		"## Errors",
	}
	prev := -1
	for _, s := range requiredSections {
		idx := strings.Index(content, s)
		if idx < 0 {
			violations++
			continue
		}
		if idx < prev {
			violations++
		}
		prev = idx
	}

	linksSection, err := extractSection(content, "## Links")
	if err == nil {
		re := regexp.MustCompile(`(?m)^-\s+(uses|reads|writes|calls|validates):\s+\[[^\]]+\]\([^)#]+\.md#[^)]+\)$`)
		lines := strings.Split(strings.TrimSpace(linksSection), "\n")
		for _, ln := range lines {
			ln = strings.TrimSpace(ln)
			if ln == "" {
				continue
			}
			if !re.MatchString(ln) {
				violations++
			}
		}
	}
	return violations
}

func extractSection(content, title string) (string, error) {
	start := strings.Index(content, title)
	if start < 0 {
		return "", errors.New("section not found")
	}
	rest := content[start+len(title):]
	next := strings.Index(rest, "\n## ")
	if next < 0 {
		return rest, nil
	}
	return rest[:next], nil
}

func writeResult(path string, result RunResult) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
