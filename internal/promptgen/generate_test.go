package promptgen

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateAll_WritesMarkdown(t *testing.T) {
	tmp := t.TempDir()
	specDir := filepath.Join(tmp, "spec")
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	yaml := `title: Test Prompt
purpose: Test purpose
output_file: ../out.md
persona:
  - persona item
workflow_steps:
  - name: step
    steps:
      - do work
must_rules:
  - must rule
tool_policy:
  allowed:
    - allowed
output_contract:
  required_sections:
    - section
`
	if err := os.WriteFile(filepath.Join(specDir, "prompt.yaml"), []byte(yaml), 0644); err != nil {
		t.Fatalf("write yaml: %v", err)
	}

	if err := GenerateAll(tmp, false); err != nil {
		t.Fatalf("GenerateAll: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmp, "out.md"))
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("output is empty")
	}
}

func TestGenerateAll_CheckDetectsOutdatedFile(t *testing.T) {
	tmp := t.TempDir()
	specDir := filepath.Join(tmp, "spec")
	if err := os.MkdirAll(specDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	yaml := `title: Test Prompt
purpose: Test purpose
output_file: ../out.md
persona:
  - persona item
workflow_steps:
  - name: step
    steps:
      - do work
must_rules:
  - must rule
tool_policy:
  allowed:
    - allowed
output_contract:
  required_sections:
    - section
`
	if err := os.WriteFile(filepath.Join(specDir, "prompt.yaml"), []byte(yaml), 0644); err != nil {
		t.Fatalf("write yaml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "out.md"), []byte("stale"), 0644); err != nil {
		t.Fatalf("write stale output: %v", err)
	}

	if err := GenerateAll(tmp, true); err == nil {
		t.Fatalf("expected check error for outdated generated file")
	}
}
