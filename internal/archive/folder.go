package archive

import (
	"fmt"
	"path/filepath"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Folder struct {
	Parent   *Space
	Data     api.Folder
	Children []*List
}

func NewFolder(parent *Space, dir string) (*Folder, error) {
	log.Fatal("not implemented")
	return nil, fmt.Errorf("not implemented")
}

func (f *Folder) GetDir() string {
	return filepath.Join(f.Parent.GetDir(), f.Data.ID)
}
