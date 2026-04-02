package archive

import (
	"fmt"

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
