package hugo

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type HugoBuild struct {
	siteDir string
}

func NewHugoBuild(siteDir string) (*HugoBuild, error) {
	return &HugoBuild{siteDir: siteDir}, nil
}

func (h *HugoBuild) Run() error {
	theme := os.Getenv("HUGO_THEME")
	if theme == "" {
		theme = "picocss"
	}

	// cleanup: hugo doesn't do it
	if err := os.RemoveAll(filepath.Join(h.siteDir, "public")); err != nil {
		return err
	}

	args := []string{"hugo", "build", "--gc", "--panicOnWarning"}
	if os.Getenv("MINIFY") != "0" {
		args = append(args, "--minify")
	}
	log.Info("$ " + strings.Join(args, " "))

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = h.siteDir
	cmd.Env = append(os.Environ(), "HUGO_THEME="+theme)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
