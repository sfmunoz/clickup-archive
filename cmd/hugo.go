package cmd

import (
	"github.com/spf13/cobra"
)

var hugoCmd = &cobra.Command{
	Use:   "hugo",
	Short: "Manage the Hugo documentation site",
	Long: `hugo provides subcommands for building and serving the Hugo
documentation site bundled with this project.

Use 'hugo build' to produce a static site or 'hugo run' to start a
local development server with live reload.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(hugoCmd)
}
