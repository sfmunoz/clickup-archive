package hugo

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type HugoRun struct {
	siteDir string
}

func NewHugoRun(siteDir string) (*HugoRun, error) {
	return &HugoRun{siteDir: siteDir}, nil
}

func (h *HugoRun) Run() error {
	theme := os.Getenv("HUGO_THEME")
	if theme == "" {
		theme = "picocss"
	}

	// cleanup: hugo doesn't do it
	if err := os.RemoveAll(filepath.Join(h.siteDir, "public")); err != nil {
		return err
	}

	args := []string{"hugo", "server", "-D", "--disableFastRender", "--noHTTPCache"}
	log.Info("$ " + strings.Join(args, " "))

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = h.siteDir
	cmd.Env = append(os.Environ(), "HUGO_THEME="+theme)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
