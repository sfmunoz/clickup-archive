package cmd

import (
	"github.com/sfmunoz/clickup-archive/internal/generate"
	"github.com/spf13/cobra"
)

var generateContentCmd = &cobra.Command{
	Use:   "content",
	Short: "Generate content from fetched ClickUp data",
	Long:  `Processes previously fetched ClickUp data and generates derived content.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		g, err := generate.NewGenerateContent(clickupDir(), "content")
		if err != nil {
			return err
		}
		return g.Run()
	},
}

func init() {
	generateCmd.AddCommand(generateContentCmd)
}
