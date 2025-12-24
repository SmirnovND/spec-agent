package cli

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/SmirnovND/spec-agent/internal/config"
	"github.com/SmirnovND/spec-agent/internal/spec"
)

func init() {
	rootCmd.AddCommand(graphCmd)
}

var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Построить граф зависимостей спецификаций",
	Long: `
Команда graph:
- читает .spec_agent/config.yaml
- находит спеки рядом с указанными roots
- определяет root-спеки (на которые никто не ссылается)
- строит граф зависимостей от этих корней
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

		return nil
	},
}

func findSpecsNearRoots(roots []string) ([]string, error) {
	var specs []string

	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if filepath.Ext(path) == ".md" {
				specs = append(specs, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return specs, nil
}
