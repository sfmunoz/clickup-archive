package archive

import (
	"fmt"

	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

type Archive struct {
	Parent   *struct{}
	Data     string
	Children []*Workspace
}

func NewArchive(dir string) (*Archive, error) {
	log.Fatal("not implemented")
	return nil, fmt.Errorf("not implemented")
}

func (a *Archive) GetDir() string {
	return a.Data
}
