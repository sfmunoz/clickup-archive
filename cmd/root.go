package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func clickupDir() string {
	return filepath.Join(os.Getenv("HOME"), "src", "clickup")
}

var rootCmd = &cobra.Command{
	Use:   "clickup-archive",
	Short: "Archive ClickUp workspaces to local JSON files",
	Long: `clickup-archive fetches the full ClickUp hierarchy via API v2
and writes each entity as index.json under $HOME/src/clickup/<id>/.

The traversal order is: workspaces → spaces → folders → lists → tasks.
Subtasks are fetched recursively and stored alongside their siblings.

Requires the CLICKUP_TOKEN environment variable (personal API token).`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
