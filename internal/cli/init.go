package cli

import (
	"github.com/SmirnovND/spec-agent/v3/internal/fs"
	"github.com/spf13/cobra"
)

func init() {
	initCmd.Flags().Bool("zenflow", false, "создать custom workflow для Zenflow в .zenflow/workflows/")
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
  - prompts/base/ — базовые промты для LLM
  - prompts/zenflow/ — добавляются при --zenflow
- spec_changes/ — директория для планов изменений

Все ресурсы встраиваются в бинарь и автоматически распаковываются.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		zenflow, _ := cmd.Flags().GetBool("zenflow")
		if err := fs.InitSpecAgent(zenflow); err != nil {
			return err
		}
		return nil
	},
}
