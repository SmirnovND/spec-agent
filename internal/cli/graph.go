package cli

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

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

		specFiles, err := findSpecsNearRoots(cfg.Roots, cfg.Exclude)
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

func findSpecsNearRoots(roots, exclude []string) ([]string, error) {
	var specs []string
	normalizedExclude := normalizePaths(exclude)

	for _, root := range roots {
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if isPathExcluded(path, normalizedExclude) {
				if d.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if d.IsDir() || filepath.Ext(path) != ".md" {
				return nil
			}
			isSpec, err := spec.IsSpecFile(path)
			if err != nil {
				return nil
			}
			if isSpec {
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

func normalizePaths(paths []string) []string {
	normalized := make([]string, 0, len(paths))
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		absPath, err := filepath.Abs(p)
		if err != nil {
			continue
		}
		normalized = append(normalized, filepath.ToSlash(filepath.Clean(absPath)))
	}
	return normalized
}

func isPathExcluded(path string, excluded []string) bool {
	if len(excluded) == 0 {
		return false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	normalizedPath := filepath.ToSlash(filepath.Clean(absPath))
	for _, ex := range excluded {
		if normalizedPath == ex || strings.HasPrefix(normalizedPath, ex+"/") {
			return true
		}
	}
	return false
}
