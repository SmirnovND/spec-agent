package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "spec-agent",
	Short: "CLI для работы со spec-driven архитектурой",
	Long: `spec-agent — инструмент для работы со спецификациями,
которые управляют архитектурой и изменениями в коде.`,
}

func Execute() error {
	return rootCmd.Execute()
}
