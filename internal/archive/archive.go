package archive

import (
	"fmt"
	"os"
)

type ArchiveData struct {
	Dir string
}

type Archive struct {
	Parent   *struct{}
	Data     *ArchiveData
	Children []*Workspace
}

func LoadArchive(dir string) (*Archive, error) {
	if err := isFolder(dir); err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	a := &Archive{
		Parent:   nil,
		Data:     &ArchiveData{Dir: dir},
		Children: make([]*Workspace, 0),
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		w, err := LoadWorkspace(a, e.Name())
		if err != nil {
			return nil, err
		}
		a.Children = append(a.Children, w)
	}
	return a, nil
}

func (a *Archive) GetDir() string {
	return a.Data.Dir
}

func isFolder(dir string) error {
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return fmt.Errorf("'%s' folder does not exist", dir)
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("'%s' path exists but it's not a folder", dir)
	}
	return nil
}
