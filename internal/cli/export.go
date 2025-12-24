package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/SmirnovND/spec-agent/internal/config"
	"github.com/SmirnovND/spec-agent/internal/spec"
)

func init() {
	rootCmd.AddCommand(exportCmd)
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Экспортировать спецификации в статичный HTML",
	Long: `
Команда export:
- читает .spec_agent/config.yaml
- находит спеки рядом с указанными roots
- определяет root-спеки (на которые никто не ссылается)
- строит граф зависимостей от этих корней
- генерирует статичный HTML с навигацией и оглавлением
- сохраняет результат в .spec_agent/build/
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("не удалось загрузить config.yaml: %w", err)
		}

		if len(cfg.Roots) == 0 {
			return fmt.Errorf("в config.yaml не указаны roots")
		}

		specFiles, err := findSpecsNearRoots(cfg.Roots)
		if err != nil {
			return err
		}

		if len(specFiles) == 0 {
			return fmt.Errorf("не найдено ни одной спецификации рядом с roots")
		}

		referenced := spec.CollectAllReferences(specFiles)

		rootSpecs := spec.FindRootSpecs(specFiles, referenced)
		if len(rootSpecs) == 0 {
			return fmt.Errorf("не удалось определить корневые спецификации")
		}

		fmt.Printf("🌳 Найдено %d корневых спецификаций:\n", len(rootSpecs))
		for _, root := range rootSpecs {
			fmt.Printf("  - %s\n", root)
		}
		fmt.Println()

		graph, err := spec.BuildGraphFromRoots(rootSpecs)
		if err != nil {
			return err
		}

		fmt.Printf("📊 Граф содержит %d узлов и %d ребер\n", len(graph.Nodes), len(graph.Edges))
		fmt.Println()

		outputDir := filepath.Join(".spec_agent", "build")
		fmt.Printf("📝 Генерирую HTML в %s...\n", outputDir)

		if err := spec.ExportToHTML(graph, outputDir); err != nil {
			return fmt.Errorf("ошибка при экспорте: %w", err)
		}

		indexPath := filepath.Join(outputDir, "index.html")
		absPath, _ := filepath.Abs(indexPath)

		fmt.Println()
		fmt.Printf("✅ HTML экспортирован успешно!\n")
		fmt.Printf("📂 Файлы находятся в: %s\n", outputDir)
		fmt.Printf("🌐 Откройте в браузере: file://%s\n", absPath)

		return nil
	},
}
