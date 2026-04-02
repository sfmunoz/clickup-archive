package archive

import (
	"fmt"
	"path/filepath"

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

func (w *Workspace) GetDir() string {
	return filepath.Join(w.Parent.GetDir(), w.Data.ID)
}
