package archive

import (
	"fmt"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Space struct {
	Parent   *Workspace
	Data     api.Space
	Children []*Folder
}

func loadSpace(dir string, parent *Workspace) (*Space, error) {
	log.Fatal("not implemented")
	return nil, fmt.Errorf("not implemented")
}
