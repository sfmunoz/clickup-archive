package cmd

import (
	"github.com/sfmunoz/clickup-archive/internal/archive"
	"github.com/sfmunoz/clickup-archive/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the terminal user interface",
	Long:  `Launches an interactive TUI for browsing the local ClickUp archive.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := archive.LoadArchive(clickupDir())
		if err != nil {
			return err
		}
		t, err := tui.NewTui(a)
		if err != nil {
			return err
		}
		return t.Run()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
