package fetch

import (
	"fmt"
	"path/filepath"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type FetchComments struct {
	clickupDir string
	client     *Client
}

func NewFetchComments(clickupDir string) (*FetchComments, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	return &FetchComments{
		clickupDir: clickupDir,
		client:     client,
	}, nil
}

func (f *FetchComments) getComments(taskID, taskDir string) error {
	var resp api.CommentsResponse
	if err := f.client.HttpGet("/task/"+taskID+"/comment", &resp); err != nil {
		return fmt.Errorf("fetch comments for task %s: %w", taskID, err)
	}
	log.Info("Comments", "task_id", taskID, "count", len(resp.Comments))
	dir := filepath.Join(taskDir, "comments")
	if err := jsonDump(resp, dir); err != nil {
		return fmt.Errorf("dump comments for task %s: %w", taskID, err)
	}
	return nil
}

func (f *FetchComments) Run(taskID string) error {
	taskDir := filepath.Join(f.clickupDir, taskID)
	return f.getComments(taskID, taskDir)
}
