package main

import (
	"os"
	"path/filepath"

	"github.com/sfmunoz/clickup-archive/internal/fetch"
	"github.com/sfmunoz/logit"
)

var (
	log       = logit.Logit().WithLevel(logit.LevelInfo)
	outputDir = filepath.Join(os.Getenv("HOME"), "src", "clickup")
)

func main() {
	token := os.Getenv("CLICKUP_TOKEN")
	if token == "" {
		log.Fatal("CLICKUP_TOKEN env var is required")
	}
	c := fetch.NewClient(token)
	if err := c.GetWorkspaces(outputDir); err != nil {
		log.Fatal("Failed", "err", err)
	}
}
