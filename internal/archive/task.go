package archive

import (
	"fmt"
	"path/filepath"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Task struct {
	Parent   *List
	Data     api.Task
	Children []*Comment
}

func NewTask(parent *List, dir string) (*Task, error) {
	log.Fatal("not implemented")
	return nil, fmt.Errorf("not implemented")
}

func (t *Task) GetDir() string {
	return filepath.Join(t.Parent.GetDir(), t.Data.ID)
}
