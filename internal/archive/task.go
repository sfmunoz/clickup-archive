package archive

import (
	"fmt"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Task struct {
	Parent   *List
	Data     api.Task
	Children []*Comment
}

func loadTask(dir string, parent *List) (*Task, error) {
	log.Fatal("not implemented")
	return nil, fmt.Errorf("not implemented")
}
