package cmd

import (
	"github.com/sfmunoz/clickup-archive/internal/hugo"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the Hugo development server",
	Long: `Removes the public/ output directory, then starts a local server:

  hugo server -D --disableFastRender --noHTTPCache

Draft pages are included (-D). The server watches for changes and
reloads automatically. The theme is read from HUGO_THEME (default: picocss).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		h, err := hugo.NewHugoRun(".")
		if err != nil {
			return err
		}
		return h.Run()
	},
}

func init() {
	hugoCmd.AddCommand(runCmd)
}
