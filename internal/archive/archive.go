package archive

import (
	"fmt"
	"os"

	"github.com/sfmunoz/clickup-archive/internal/api"
	"github.com/sfmunoz/logit"
)

var log = logit.Logit().WithLevel(logit.LevelInfo)

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

func (a *Archive) SaveWorkspace(w *api.Workspace, update bool) error {
	var wOld *Workspace = nil
	for _, c := range a.Children {
		if c.Data.ID != w.ID {
			continue
		}
		if !update {
			return fmt.Errorf("workspace '%s=%s' already exists and 'update' is false", c.Data.ID, c.Data.Name)
		}
		wOld = c
		break
	}
	if wOld == nil {
		log.Info("creating workspace", "id", w.ID, "name", w.Name)
	} else {
		log.Warn("updating workspace", "id_old", wOld.Data.ID, "name_old", wOld.Data.Name, "id_new", w.ID, "name_new", w.Name)
	}
	dir := workspaceDir(a.GetDir(), w.ID)
	if err := jsonSave(w, dir); err != nil {
		return err
	}
	if wOld == nil {
		a.Children = append(a.Children, &Workspace{Parent: a, Data: w, Children: make([]*Space, 0)})
	} else {
		wOld.Data = w
	}
	return nil
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
