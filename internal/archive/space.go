package archive

import (
	"fmt"
	"path/filepath"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Space struct {
	Parent   *Workspace
	Data     api.Space
	Children []*Folder
}

func NewSpace(parent *Workspace, dir string) (*Space, error) {
	log.Fatal("not implemented")
	return nil, fmt.Errorf("not implemented")
}

func (s *Space) GetDir() string {
	return filepath.Join(s.Parent.GetDir(), s.Data.ID)
}
