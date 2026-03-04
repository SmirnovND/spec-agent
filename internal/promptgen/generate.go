package promptgen

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

func GenerateAll(sourceDir string, check bool) error {
	specFiles, err := findSpecFiles(sourceDir)
	if err != nil {
		return err
	}

	for _, specPath := range specFiles {
		if err := generateFromFile(specPath, check); err != nil {
			return err
		}
	}

	return nil
}

func findSpecFiles(sourceDir string) ([]string, error) {
	files := make([]string, 0, 16)
	err := filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if strings.EqualFold(filepath.Base(path), "prompt.yaml") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Strings(files)
	return files, nil
}

func generateFromFile(specPath string, check bool) error {
	spec, err := loadSpec(specPath)
	if err != nil {
		return fmt.Errorf("invalid prompt spec %s: %w", specPath, err)
	}

	outputPath := filepath.Join(filepath.Dir(specPath), spec.OutputFile)
	content := []byte(RenderMarkdown(spec))

	if check {
		existing, err := os.ReadFile(outputPath)
		if err != nil {
			return fmt.Errorf("read generated file %s: %w", outputPath, err)
		}
		if !bytes.Equal(existing, content) {
			return fmt.Errorf("generated file is outdated: %s", outputPath)
		}
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(outputPath, content, 0644)
}

func loadSpec(specPath string) (PromptSpec, error) {
	data, err := os.ReadFile(specPath)
	if err != nil {
		return PromptSpec{}, err
	}

	var spec PromptSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return PromptSpec{}, err
	}

	if err := validateSpec(spec); err != nil {
		return PromptSpec{}, err
	}
	return spec, nil
}

func validateSpec(spec PromptSpec) error {
	if strings.TrimSpace(spec.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(spec.OutputFile) == "" {
		return errors.New("output_file is required")
	}
	if len(spec.Persona) == 0 {
		return errors.New("persona is required")
	}
	if len(spec.WorkflowSteps) == 0 {
		return errors.New("workflow_steps is required")
	}
	for _, step := range spec.WorkflowSteps {
		if strings.TrimSpace(step.Name) == "" {
			return errors.New("workflow_steps.name is required")
		}
		if len(step.Steps) == 0 {
			return errors.New("workflow_steps.steps is required")
		}
	}
	if len(spec.MustRules) == 0 {
		return errors.New("must_rules is required")
	}
	if len(spec.ToolPolicy.Allowed) == 0 && len(spec.ToolPolicy.Forbidden) == 0 && len(spec.ToolPolicy.Notes) == 0 {
		return errors.New("tool_policy is required")
	}
	if len(spec.OutputContract.RequiredSections) == 0 && len(spec.OutputContract.Rules) == 0 {
		return errors.New("output_contract is required")
	}

	return nil
}
