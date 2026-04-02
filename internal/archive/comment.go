package archive

import (
	"path/filepath"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Comment struct {
	Parent   *Task
	Data     api.Comment
	Children []*struct{}
}

func NewComment(parent *Task, dir string) (*Comment, error) {
	dir = filepath.Join(parent.GetDir(), "comments", dir)
	if err := isFolder(dir); err != nil {
		return nil, err
	}
	return &Comment{
		Parent:   parent,
		Data:     api.Comment{},
		Children: make([]*struct{}, 0),
	}, nil
}

func (c *Comment) GetDir() string {
	return filepath.Join(c.Parent.GetDir(), "comments", c.Data.ID)
}
