package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SmirnovND/spec-agent/internal/eval"
)

func main() {
	repoRoot, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to detect working directory: %v\n", err)
		os.Exit(1)
	}

	tasksDir := flag.String("tasks-dir", "eval/tasks", "path to tasks directory")
	fixturesDir := flag.String("fixtures-dir", "eval/fixtures", "path to fixtures directory")
	thresholdsPath := flag.String("thresholds", "eval/baselines/thresholds.yaml", "path to thresholds file")
	tasks := flag.String("tasks", "", "comma-separated task ids; empty = all")
	output := flag.String("output", "", "output JSON path; default eval/results/<timestamp>.json")
	flag.Parse()

	outputPath := strings.TrimSpace(*output)
	if outputPath == "" {
		outputPath = filepath.Join("eval", "results", fmt.Sprintf("%s.json", time.Now().Format("20060102_150405")))
	}

	taskSet := map[string]bool{}
	if s := strings.TrimSpace(*tasks); s != "" {
		for _, t := range strings.Split(s, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				taskSet[t] = true
			}
		}
	}

	res, err := eval.Run(eval.RunOptions{
		RepoRoot:    repoRoot,
		TasksDir:    *tasksDir,
		FixturesDir: *fixturesDir,
		Thresholds:  *thresholdsPath,
		TaskIDs:     taskSet,
		OutputPath:  outputPath,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "eval-runner error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("eval completed: tasks=%d, gating_pass=%t\n", len(res.Tasks), res.GatingPass)
	fmt.Printf("completeness_avg=%.3f, test_pass_rate_avg=%.3f, scope_violations_total=%d, spec_rule_violations_total=%d\n",
		res.CompletenessAvg,
		res.TestPassRateAvg,
		res.ScopeViolationsTotal,
		res.SpecRuleViolationsTotal,
	)
	fmt.Printf("result_json=%s\n", outputPath)
	for _, task := range res.Tasks {
		if task.GatingPass {
			continue
		}
		fmt.Printf("task=%s gating_pass=false\n", task.TaskID)
		if strings.TrimSpace(task.RunError) != "" {
			fmt.Printf("run_error:\n%s\n", task.RunError)
		}
		if len(task.TestFailures) > 0 {
			fmt.Printf("test_failures (%d):\n", len(task.TestFailures))
			for i, failure := range task.TestFailures {
				fmt.Printf("[%d]\n%s\n", i+1, failure)
			}
		}
	}

	if !res.GatingPass {
		os.Exit(2)
	}
}
