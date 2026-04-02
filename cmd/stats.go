package cmd

import (
	"github.com/sfmunoz/clickup-archive/internal/archive"
	"github.com/sfmunoz/clickup-archive/internal/stats"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Print statistics about fetched ClickUp data",
	Long:  `Walks the local ClickUp archive and prints entity counts by level.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := archive.LoadArchive(clickupDir())
		if err != nil {
			return err
		}
		s, err := stats.NewStats(a)
		if err != nil {
			return err
		}
		return s.Run()
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
