package tui

import (
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

type Tui struct {
	clickupDir string
}

func NewTui(clickupDir string) (*Tui, error) {
	return &Tui{clickupDir: clickupDir}, nil
}

func (t *Tui) Run() error {
	log.Warn("Tui.Run() not yet implemented")
	return nil
}
