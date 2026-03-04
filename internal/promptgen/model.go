package promptgen

type PromptSpec struct {
	Title          string         `yaml:"title"`
	Purpose        string         `yaml:"purpose"`
	OutputFile     string         `yaml:"output_file"`
	Persona        []string       `yaml:"persona"`
	WorkflowSteps  []WorkflowStep `yaml:"workflow_steps"`
	MustRules      []string       `yaml:"must_rules"`
	ToolPolicy     ToolPolicy     `yaml:"tool_policy"`
	OutputContract OutputContract `yaml:"output_contract"`
}

type WorkflowStep struct {
	Name  string   `yaml:"name"`
	Steps []string `yaml:"steps"`
}

type ToolPolicy struct {
	Allowed   []string `yaml:"allowed"`
	Forbidden []string `yaml:"forbidden"`
	Notes     []string `yaml:"notes"`
}

type OutputContract struct {
	RequiredSections []string `yaml:"required_sections"`
	Rules            []string `yaml:"rules"`
}
