package promptgen

import (
	"fmt"
	"strings"
)

func RenderMarkdown(spec PromptSpec) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# %s\n\n", spec.Title))
	b.WriteString("<!-- Generated from prompt.yaml. Do not edit manually. -->\n\n")

	if spec.Purpose != "" {
		b.WriteString("## Purpose\n\n")
		b.WriteString(spec.Purpose)
		b.WriteString("\n\n")
	}

	b.WriteString("## Persona\n\n")
	for _, item := range spec.Persona {
		b.WriteString("- ")
		b.WriteString(item)
		b.WriteString("\n")
	}
	b.WriteString("\n")

	b.WriteString("## Workflow Steps\n\n")
	for i, step := range spec.WorkflowSteps {
		b.WriteString(fmt.Sprintf("### %d. %s\n\n", i+1, step.Name))
		for j, item := range step.Steps {
			b.WriteString(fmt.Sprintf("%d. %s\n", j+1, item))
		}
		b.WriteString("\n")
	}

	b.WriteString("## Must Rules\n\n")
	for i, rule := range spec.MustRules {
		b.WriteString(fmt.Sprintf("%d. %s\n", i+1, rule))
	}
	b.WriteString("\n")

	b.WriteString("## Tool Policy\n\n")
	if len(spec.ToolPolicy.Allowed) > 0 {
		b.WriteString("### Allowed\n\n")
		for _, item := range spec.ToolPolicy.Allowed {
			b.WriteString("- ")
			b.WriteString(item)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
	if len(spec.ToolPolicy.Forbidden) > 0 {
		b.WriteString("### Forbidden\n\n")
		for _, item := range spec.ToolPolicy.Forbidden {
			b.WriteString("- ")
			b.WriteString(item)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
	if len(spec.ToolPolicy.Notes) > 0 {
		b.WriteString("### Notes\n\n")
		for _, item := range spec.ToolPolicy.Notes {
			b.WriteString("- ")
			b.WriteString(item)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	b.WriteString("## Output Contract\n\n")
	if len(spec.OutputContract.RequiredSections) > 0 {
		b.WriteString("### Required Sections\n\n")
		for _, section := range spec.OutputContract.RequiredSections {
			b.WriteString("- ")
			b.WriteString(section)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
	if len(spec.OutputContract.Rules) > 0 {
		b.WriteString("### Rules\n\n")
		for _, rule := range spec.OutputContract.Rules {
			b.WriteString("- ")
			b.WriteString(rule)
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}
