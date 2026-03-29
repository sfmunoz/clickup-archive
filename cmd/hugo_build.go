package cmd

import (
	"github.com/sfmunoz/clickup-archive/internal/hugo"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the Hugo documentation site",
	Long: `Removes the public/ output directory, then runs:

  hugo build --gc --panicOnWarning [--minify]

Minification is enabled by default; set MINIFY=0 to disable it.
The theme is read from HUGO_THEME (default: picocss).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := hugo.NewHugoBuild(".")
		if err != nil {
			return err
		}
		return h.Run()
	},
}

func init() {
	hugoCmd.AddCommand(buildCmd)
}
