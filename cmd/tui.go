package cmd

import (
	"github.com/sfmunoz/clickup-archive/internal/tui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the terminal user interface",
	Long:  `Launches an interactive TUI for browsing the local ClickUp archive.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := tui.NewTui(clickupDir())
		if err != nil {
			return err
		}
		return t.Run()
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
