package cli

import (
	"github.com/SmirnovND/spec-agent/internal/fs"
	"github.com/spf13/cobra"
)

func init() {
	initPHPCmd.Flags().Bool("zenflow", false, "создать custom workflow для Zenflow в .zenflow/workflows/")
	rootCmd.AddCommand(initPHPCmd)
}

var initPHPCmd = &cobra.Command{
	Use:     "init-php",
	Aliases: []string{"init_php"},
	Short:   "Инициализировать spec-agent для PHP-проекта",
	Long: `Инициализирует структуру spec-agent для PHP-проекта.

Создает:
- .spec_agent/config.yaml — конфиг с default roots для PHP
- .spec_agent/php/examples/ — примеры PHP-спецификаций
- .spec_agent/php/prompts/base/ — базовые PHP-промты
- .spec_agent/php/prompts/zenflow/ — добавляются при --zenflow
- spec_changes/ — директория для планов изменений

Все ресурсы встраиваются в бинарь и автоматически распаковываются.
Примечание: повторный запуск не перезаписывает существующие файлы; добавляются только отсутствующие.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		zenflow, _ := cmd.Flags().GetBool("zenflow")
		if err := fs.InitPHPSpecAgent(zenflow); err != nil {
			return err
		}
		return nil
	},
}
