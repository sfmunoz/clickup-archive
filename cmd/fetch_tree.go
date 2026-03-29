package cmd

import (
	"os"
	"path/filepath"

	"github.com/sfmunoz/clickup-archive/internal/fetch"
	"github.com/spf13/cobra"
)

var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Traverse and dump the full ClickUp hierarchy",
	Long: `Walks the complete ClickUp hierarchy and writes each entity as
index.json under $HOME/src/clickup/<id>/:

  workspaces → spaces → folders → lists → tasks (+ subtasks)

Tasks are fetched with subtasks=true and paginated until exhausted.
Subtasks are stored alongside top-level tasks, not nested under them.

Requires the CLICKUP_TOKEN environment variable (personal API token).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		outputDir := filepath.Join(os.Getenv("HOME"), "src", "clickup")
		f, err := fetch.NewFetchTree()
		if err != nil {
			return err
		}
		return f.Run(outputDir)
	},
}

func init() {
	fetchCmd.AddCommand(treeCmd)
}
