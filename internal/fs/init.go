package fs

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed assets
var embeddedAssets embed.FS

type initProfile struct {
	config            string
	examplesAssetDir  string
	examplesTargetDir string
	basePromptsDir    string
	basePromptsTarget string
	zenflowPromptsDir string
	zenflowTargetDir  string
	workflowName      string
	workflowPromptDir string
}

func InitSpecAgent(withZenflow bool) error {
	return initSpecAgent(withZenflow, initProfile{
		config: `roots:
  - internal/controllers
  - cmd
`,
		examplesAssetDir:  "assets/examples",
		examplesTargetDir: ".spec_agent/examples",
		basePromptsDir:    "assets/go/prompts/base",
		basePromptsTarget: ".spec_agent/prompts/base",
		zenflowPromptsDir: "assets/go/prompts/zenflow",
		zenflowTargetDir:  ".spec_agent/prompts/zenflow",
		workflowName:      "spec-agent-spec-driven.md",
		workflowPromptDir: ".spec_agent/prompts/zenflow",
	})
}

func InitPHPSpecAgent(withZenflow bool) error {
	return initSpecAgent(withZenflow, initProfile{
		config: `roots:
  - app/Http/Controllers
  - app/Console/Commands
`,
		examplesAssetDir:  "assets/php/examples",
		examplesTargetDir: ".spec_agent/php/examples",
		basePromptsDir:    "assets/php/prompts/base",
		basePromptsTarget: ".spec_agent/php/prompts/base",
		zenflowPromptsDir: "assets/php/prompts/zenflow",
		zenflowTargetDir:  ".spec_agent/php/prompts/zenflow",
		workflowName:      "spec-agent-php-spec-driven.md",
		workflowPromptDir: ".spec_agent/php/prompts/zenflow",
	})
}

func initSpecAgent(withZenflow bool, profile initProfile) error {
	if err := os.MkdirAll(".spec_agent", 0755); err != nil {
		return err
	}
	if err := os.MkdirAll("spec_changes", 0755); err != nil {
		return err
	}

	if err := os.WriteFile(".spec_agent/config.yaml", []byte(profile.config), 0644); err != nil {
		return err
	}

	if err := copyAssetDir(profile.examplesAssetDir, profile.examplesTargetDir); err != nil {
		return err
	}
	if err := copyAssetDir(profile.basePromptsDir, profile.basePromptsTarget); err != nil {
		return err
	}

	if withZenflow {
		if err := copyAssetDir(profile.zenflowPromptsDir, profile.zenflowTargetDir); err != nil {
			return err
		}
	}

	if withZenflow {
		if err := createZenflowWorkflow(profile.workflowName, profile.workflowPromptDir); err != nil {
			return err
		}
	}

	return nil
}

func copyAssetDir(assetDir, targetDir string) error {
	return fs.WalkDir(embeddedAssets, assetDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == assetDir {
			return nil
		}

		rel := strings.TrimPrefix(path, assetDir+"/")
		destPath := filepath.Join(targetDir, rel)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		data, err := embeddedAssets.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}
		return os.WriteFile(destPath, data, 0644)
	})
}

func createZenflowWorkflow(workflowName, zenflowPromptPath string) error {
	workflowDir := filepath.Join(".zenflow", "workflows")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return err
	}

	workflowPath := filepath.Join(workflowDir, workflowName)
	if _, err := os.Stat(workflowPath); err == nil {
		return nil
	}

	workflow := fmt.Sprintf(`# Spec-Agent Spec-Driven Workflow

## Configuration
- **Artifacts Path**: {@artifacts_path} -> .zenflow/tasks/{task_id}

---

## Workflow Steps

### [ ] Step: Planning
1. Обязательно прочитай: %s/planning.md.
2. Зафиксируй контекст, цель и ограничения задачи.
3. Составь план в {@artifacts_path}/plan.md.
Acceptance criteria:
- В plan.md описаны scope, затронутые компоненты и этапы.

### [ ] Step: Technical Specification
1. Обязательно прочитай: %s/technical_specification.md.
2. Подготовь/обнови спецификацию в {@artifacts_path}/spec.md.
3. Опиши Business Logic, Flow, Links и Dependencies.
Acceptance criteria:
- spec.md содержит полный сценарий изменений и ссылки на связанные спеки.

### [ ] Step: Specification Review
1. Обязательно прочитай: %s/specification_review.md.
2. Проверь spec.md на полноту, непротиворечивость и трассируемость к задаче.
3. Зафиксируй решение в {@artifacts_path}/report.md:
   - approve: можно переходить к реализации;
   - change requested: что нужно исправить в spec.md.
Acceptance criteria:
- Есть явное решение review (approve или change requested) и обоснование.

### [ ] Step: Implementation
1. Обязательно прочитай: %s/implementation.md.
2. Внеси изменения в код строго по спецификации.
3. Обнови связанные спеки рядом с кодом (если поведение меняется).
4. Зафиксируй результат в {@artifacts_path}/report.md:
   - что изменено;
   - какие проверки выполнены;
   - ограничения и риски.
Acceptance criteria:
- Код и спецификации синхронизированы, report.md заполнен.

### [ ] Step: Review & Wrap-Up
1. Обязательно прочитай: %s/review_wrap_up.md.
2. Проверь соответствие реализации спецификации и plan.md.
3. Проверь ссылки между спецификациями и итоговую документацию.
4. Если есть замечания, обнови report.md и вернись к нужному шагу.
Acceptance criteria:
- Изменения готовы к финальному ревью и передаче.
`, zenflowPromptPath, zenflowPromptPath, zenflowPromptPath, zenflowPromptPath, zenflowPromptPath)

	return os.WriteFile(workflowPath, []byte(workflow), 0644)
}
