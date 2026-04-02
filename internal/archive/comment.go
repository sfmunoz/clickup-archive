package archive

import (
	"fmt"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Comment struct {
	Parent   *Task
	Data     api.Comment
	Children []*struct{}
}

func loadComment(dir string, parent *Task) (*Comment, error) {
	log.Fatal("not implemented")
	return nil, fmt.Errorf("not implemented")
}
