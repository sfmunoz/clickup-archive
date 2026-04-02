package archive

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Comment struct {
	Parent   *Task
	Data     *api.Comment
	Children []*struct{}
}

func LoadComment(parent *Task, id string) (*Comment, error) {
	dir := commentDir(parent.GetDir(), id)
	if err := isFolder(dir); err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(filepath.Join(dir, "index.json"))
	if err != nil {
		return nil, err
	}
	var data api.Comment
	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, err
	}
	return &Comment{
		Parent:   parent,
		Data:     &data,
		Children: make([]*struct{}, 0),
	}, nil
}

func (c *Comment) GetDir() string {
	return commentDir(c.Parent.GetDir(), c.Data.ID)
}
