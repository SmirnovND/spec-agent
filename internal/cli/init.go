package cli

import (
	"github.com/SmirnovND/spec-agent/internal/fs"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Инициализировать spec-agent",
	Long: `Инициализирует структуру spec-agent в проекте.

Создает:
- .spec_agent/ — директория конфигурации
  - config.yaml — файл с корневыми путями (roots)
  - examples/ — примеры спецификаций
  - prompts/ — промты для LLM
  - README.md — документация по ресурсам
- spec_changes/ — директория для планов изменений

Все ресурсы встраиваются в бинарь и автоматически распаковываются.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := fs.InitSpecAgent(); err != nil {
			return err
		}
		return nil
	},
}
