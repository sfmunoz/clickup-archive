package fetch

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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

func (f *FetchComments) fetchAllComments(taskID, commentsDir string) (int, error) {
	log.Info("Fetching comments", "task_id", taskID)
	startID := ""
	total := 0
	for {
		var resp api.CommentsResponse
		path := "/task/" + taskID + "/comment"
		if startID != "" {
			path += "?start_id=" + startID
		}
		log.Info("fetching", "path", path)
		if err := f.client.HttpGet(path, &resp); err != nil {
			return total, fmt.Errorf("fetch comments for task %s: %w", taskID, err)
		}
		if len(resp.Comments) == 0 {
			break
		}
		for _, comment := range resp.Comments {
			commentDir := filepath.Join(commentsDir, comment.ID)
			if err := jsonDump(comment, commentDir); err != nil {
				return total, fmt.Errorf("dump comment %s for task %s: %w", comment.ID, taskID, err)
			}
		}
		total += len(resp.Comments)
		startID = resp.Comments[len(resp.Comments)-1].ID
	}
	return total, nil
}

func (f *FetchComments) processTask(taskID, taskDir string) error {
	doneFile := filepath.Join(taskDir, "comments.done")
	commentsDir := filepath.Join(taskDir, "comments")

	if _, err := os.Stat(doneFile); err == nil {
		log.Info("Comments already done", "task_id", taskID)
		return nil
	}

	if err := os.RemoveAll(commentsDir); err != nil {
		return fmt.Errorf("remove comments dir for task %s: %w", taskID, err)
	}

	total, err := f.fetchAllComments(taskID, commentsDir)
	if err != nil {
		return err
	}

	log.Info("Comments", "task_id", taskID, "count", total)

	if err := os.WriteFile(doneFile, nil, 0o644); err != nil {
		return fmt.Errorf("write comments.done for task %s: %w", taskID, err)
	}
	return nil
}

func (f *FetchComments) Run() error {
	return filepath.WalkDir(f.clickupDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(f.clickupDir, path)
		if err != nil || rel == "." {
			return nil
		}
		depth := len(strings.Split(rel, string(filepath.Separator)))
		if depth < 5 {
			return nil
		}
		if depth > 5 {
			return fs.SkipDir
		}
		// depth == 5: task directory (workspace/space/folder/list/task)
		return f.processTask(filepath.Base(path), path)
	})
}
