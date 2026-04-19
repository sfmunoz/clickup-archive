package cmd

import (
	"github.com/spf13/cobra"
)

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch data from the ClickUp API",
	Long: `fetch provides subcommands that retrieve data from the ClickUp API v2
and persist it as JSON files under $HOME/.archive/clickup/.

Requires the CLICKUP_TOKEN environment variable (personal API token).`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
