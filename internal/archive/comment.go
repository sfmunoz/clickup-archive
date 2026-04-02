package archive

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Comment struct {
	Parent   *Task
	Data     *api.Comment
	Children []*struct{}
}

func (c *Comment) GetDir() string {
	return commentDir(c.Parent.GetDir(), c.Data.ID)
}

func LoadComment(parent *Task, id string) (*Comment, error) {
	dir := commentDir(parent.GetDir(), id)
	if err := isFolder(dir); err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(indexFile(dir))
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

func SaveComment(parent *Task, c *api.Comment, update bool) (*Comment, error) {
	var cOld *Comment = nil
	for _, ch := range parent.Children {
		if ch.Data.ID != c.ID {
			continue
		}
		if !update {
			return nil, fmt.Errorf("comment '%s' already exists and 'update' is false", ch.Data.ID)
		}
		cOld = ch
		break
	}
	dir := commentDir(parent.GetDir(), c.ID)
	if err := jsonSave(c, dir); err != nil {
		return nil, err
	}
	if cOld == nil {
		log.Info("comment created", "id", c.ID)
		cNew := &Comment{Parent: parent, Data: c, Children: make([]*struct{}, 0)}
		parent.Children = append(parent.Children, cNew)
		return cNew, nil
	}
	log.Warn("comment updated", "id_old", cOld.Data.ID, "id_new", c.ID)
	cOld.Data = c
	return cOld, nil
}
