package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Показать версию spec-agent",
	Run: func(cmd *cobra.Command, args []string) {
		v := strings.TrimSpace(Version)
		if v == "" {
			v = "dev"
		}
		_, _ = fmt.Fprintln(cmd.OutOrStdout(), v)
	},
}
