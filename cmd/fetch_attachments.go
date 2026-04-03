package cmd

import (
	"github.com/sfmunoz/clickup-archive/internal/archive"
	"github.com/sfmunoz/clickup-archive/internal/fetch"
	"github.com/spf13/cobra"
)

var attachmentsCmd = &cobra.Command{
	Use:   "attachments",
	Short: "Download all task attachments from the ClickUp archive",
	Long: `Walks the fetch-tree output directory and downloads all attachments for every task.

For each task, attachment metadata is saved as index.json under <task-id>/attachments/<attachment-id>/
and the binary file is saved alongside it. A <task-id>/attachments.done marker file is created on
success; if it already exists the task is skipped. If it is absent, <task-id>/attachments/ is
deleted and fully re-downloaded.

Attachment metadata is read from the task's existing index.json — no extra API calls are needed.

Requires the CLICKUP_TOKEN environment variable (personal API token).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := archive.LoadArchive(clickupDir())
		if err != nil {
			return err
		}
		f, err := fetch.NewFetchAttachments(a)
		if err != nil {
			return err
		}
		return f.Run()
	},
}

func init() {
	fetchCmd.AddCommand(attachmentsCmd)
}
