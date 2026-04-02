package archive

import (
	"fmt"
	"path/filepath"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Comment struct {
	Parent   *Task
	Data     api.Comment
	Children []*struct{}
}

func NewComment(parent *Task, dir string) (*Comment, error) {
	log.Fatal("not implemented")
	return nil, fmt.Errorf("not implemented")
}

func (c *Comment) GetDir() string {
	return filepath.Join(c.Parent.GetDir(), "comments", c.Data.ID)
}
