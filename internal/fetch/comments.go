package fetch

import (
	"fmt"

	"github.com/sfmunoz/clickup-archive/internal/api"
	"github.com/sfmunoz/clickup-archive/internal/archive"
)

type FetchComments struct {
	archive *archive.Archive
	client  *Client
}

func NewFetchComments(a *archive.Archive) (*FetchComments, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	return &FetchComments{
		archive: a,
		client:  client,
	}, nil
}

func (f *FetchComments) fetchAllComments(task *archive.Task) (int, error) {
	taskID := task.Data.ID
	log.Info("Fetching comments", "task_id", taskID)
	startID := ""
	startDate := ""
	total := 0
	for {
		var resp api.CommentsResponse
		path := "/task/" + taskID + "/comment"
		if startID != "" {
			path += "?start=" + startDate + "&start_id=" + startID
		}
		log.Info("fetching", "path", path)
		if err := f.client.HttpGet(path, &resp); err != nil {
			return total, fmt.Errorf("fetch comments for task %s: %w", taskID, err)
		}
		if len(resp.Comments) == 0 {
			break
		}
		for _, comment := range resp.Comments {
			if err := task.SaveComment(&comment, false); err != nil {
				return total, fmt.Errorf("dump comment %s for task %s: %w", comment.ID, taskID, err)
			}
		}
		total += len(resp.Comments)
		last := resp.Comments[len(resp.Comments)-1]
		startID = last.ID
		startDate = last.Date
	}
	return total, nil
}

func (f *FetchComments) processTask(task *archive.Task) error {
	if task.IsCommentsDone() {
		log.Info("Comments already done", "task_id", task.Data.ID)
		return nil
	}

	if err := task.ClearComments(); err != nil {
		return fmt.Errorf("remove comments dir for task %s: %w", task.Data.ID, err)
	}

	total, err := f.fetchAllComments(task)
	if err != nil {
		return err
	}

	log.Info("Comments", "task_id", task.Data.ID, "count", total)

	if err := task.MarkCommentsDone(); err != nil {
		return fmt.Errorf("write comments.done for task %s: %w", task.Data.ID, err)
	}
	return nil
}

func (f *FetchComments) Run() error {
	for _, ws := range f.archive.Children {
		for _, sp := range ws.Children {
			for _, fo := range sp.Children {
				for _, li := range fo.Children {
					for _, ta := range li.Children {
						if err := f.processTask(ta); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}
