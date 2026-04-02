package cmd

import (
	"github.com/sfmunoz/clickup-archive/internal/archive"
	"github.com/sfmunoz/clickup-archive/internal/fetch"
	"github.com/spf13/cobra"
)

var commentsCmd = &cobra.Command{
	Use:   "comments",
	Short: "Fetch comments for all tasks in the ClickUp tree",
	Long: `Walks the fetch-tree output directory and fetches all comments for every task.

For each task, comments are saved as index.json under <task-id>/comments/<comment-id>/.
A <task-id>/comments.done marker file is created on success; if it already exists the
task is skipped. If it is absent, <task-id>/comments/ is deleted and fully re-fetched.

Requires the CLICKUP_TOKEN environment variable (personal API token).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := archive.LoadArchive(clickupDir())
		if err != nil {
			return err
		}
		f, err := fetch.NewFetchComments(a)
		if err != nil {
			return err
		}
		return f.Run()
	},
}

func init() {
	fetchCmd.AddCommand(commentsCmd)
}
