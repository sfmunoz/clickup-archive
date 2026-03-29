package cmd

import (
	"github.com/sfmunoz/clickup-archive/internal/fetch"
	"github.com/spf13/cobra"
)

var commentsCmd = &cobra.Command{
	Use:   "comments <task-id>",
	Short: "Fetch comments for a ClickUp task",
	Long: `Fetches all comments for the given task ID and writes them as
index.json under $HOME/src/clickup/<task-id>/comments/:

Requires the CLICKUP_TOKEN environment variable (personal API token).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := fetch.NewFetchComments(clickupDir())
		if err != nil {
			return err
		}
		return f.Run(args[0])
	},
}

func init() {
	fetchCmd.AddCommand(commentsCmd)
}
