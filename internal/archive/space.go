package archive

import (
	"encoding/json"
	"os"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Space struct {
	Parent   *Workspace
	Data     *api.Space
	Children []*Folder
}

func LoadSpace(parent *Workspace, id string) (*Space, error) {
	dir := spaceDir(parent.GetDir(), id)
	if err := isFolder(dir); err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(indexFile(dir))
	if err != nil {
		return nil, err
	}
	var data api.Space
	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	s := &Space{
		Parent:   parent,
		Data:     &data,
		Children: make([]*Folder, 0),
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		f, err := LoadFolder(s, e.Name())
		if err != nil {
			return nil, err
		}
		s.Children = append(s.Children, f)
	}
	return s, nil
}

func (s *Space) GetDir() string {
	return spaceDir(s.Parent.GetDir(), s.Data.ID)
}
