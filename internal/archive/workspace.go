package archive

import (
	"fmt"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Workspace struct {
	Parent   *Archive
	Data     api.Workspace
	Children []*Space
}

func NewWorkspace(parent *Archive, dir string) (*Workspace, error) {
	log.Fatal("not implemented")
	return nil, fmt.Errorf("not implemented")
}
