package cmd

import (
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate content from fetched ClickUp data",
	Long: `generate provides subcommands that process previously fetched ClickUp
data and produce derived output.

Requires data previously downloaded with the fetch command.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
