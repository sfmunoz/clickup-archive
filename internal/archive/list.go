package archive

import (
	"fmt"
	"path/filepath"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type List struct {
	Parent   *Folder
	Data     api.List
	Children []*Task
}

func NewList(parent *Folder, dir string) (*List, error) {
	log.Fatal("not implemented")
	return nil, fmt.Errorf("not implemented")
}

func (l *List) GetDir() string {
	return filepath.Join(l.Parent.GetDir(), l.Data.ID)
}
