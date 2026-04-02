package archive

import (
	"fmt"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type List struct {
	Parent   *Folder
	Data     api.List
	Children []*Task
}

func loadList(dir string, parent *Folder, space *Space) (*List, error) {
	log.Fatal("not implemented")
	return nil, fmt.Errorf("not implemented")
}
