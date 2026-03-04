package eval

import "time"

type TaskConfig struct {
	Name         string   `yaml:"name"`
	Fixture      string   `yaml:"fixture"`
	RunCommand   string   `yaml:"run_command"`
	AllowedPaths []string `yaml:"allowed_paths"`
	TestCommands []string `yaml:"test_commands"`
}

type ChecksConfig struct {
	Checkpoints []Checkpoint `yaml:"checkpoints"`
}

type Checkpoint struct {
	ID      string `yaml:"id"`
	Type    string `yaml:"type"`
	Path    string `yaml:"path"`
	Pattern string `yaml:"pattern"`
}

type Thresholds struct {
	CompletenessMin       float64 `yaml:"completeness_min"`
	TestPassRateMin       float64 `yaml:"test_pass_rate_min"`
	ScopeViolationsMax    int     `yaml:"scope_violations_max"`
	SpecRuleViolationsMax int     `yaml:"spec_rule_violations_max"`
}

type TaskResult struct {
	TaskID              string    `json:"task_id"`
	TaskName            string    `json:"task_name"`
	Fixture             string    `json:"fixture"`
	StartedAt           time.Time `json:"started_at"`
	FinishedAt          time.Time `json:"finished_at"`
	RunCommand          string    `json:"run_command"`
	RunError            string    `json:"run_error,omitempty"`
	ChangedPaths        []string  `json:"changed_paths"`
	Completeness        float64   `json:"completeness"`
	CompletedChecks     int       `json:"completed_checks"`
	TotalChecks         int       `json:"total_checks"`
	ScopeViolations     int       `json:"scope_violations"`
	ScopeViolationPaths []string  `json:"scope_violation_paths"`
	TestPassRate        float64   `json:"test_pass_rate"`
	TestsPassed         int       `json:"tests_passed"`
	TestsTotal          int       `json:"tests_total"`
	TestFailures        []string  `json:"test_failures,omitempty"`
	SpecRuleViolations  int       `json:"spec_rule_violations"`
	GatingPass          bool      `json:"gating_pass"`
}

type RunResult struct {
	GeneratedAt             time.Time    `json:"generated_at"`
	Thresholds              Thresholds   `json:"thresholds"`
	Tasks                   []TaskResult `json:"tasks"`
	CompletenessAvg         float64      `json:"completeness_avg"`
	TestPassRateAvg         float64      `json:"test_pass_rate_avg"`
	ScopeViolationsTotal    int          `json:"scope_violations_total"`
	SpecRuleViolationsTotal int          `json:"spec_rule_violations_total"`
	GatingPass              bool         `json:"gating_pass"`
}
