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
	f, err := fetch.NewFetchTree()
	if err != nil {
		log.Fatal("Failed", "err", err)
	}
	if err := f.Run(outputDir); err != nil {
		log.Fatal("Failed", "err", err)
	}
}
