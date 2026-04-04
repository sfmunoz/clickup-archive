package fetch

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sfmunoz/clickup-archive/internal/archive"
)

type FetchAttachments struct {
	archive *archive.Archive
	client  *Client
}

func NewFetchAttachments(a *archive.Archive) (*FetchAttachments, error) {
	client, err := NewClient()
	if err != nil {
		return nil, err
	}
	return &FetchAttachments{
		archive: a,
		client:  client,
	}, nil
}

func (f *FetchAttachments) processTask(task *archive.Task) error {
	if task.IsAttachmentsDone() {
		log.Info("Attachments already done", "task_id", task.Data.ID)
		return nil
	}

	if err := task.ClearAttachments(); err != nil {
		return fmt.Errorf("remove attachments dir for task %s: %w", task.Data.ID, err)
	}

	// Fetch fresh task data — the list endpoint omits attachment metadata
	freshTask, err := f.client.GetTask(task.Data.ID)
	if err != nil {
		return fmt.Errorf("fetch task %s: %w", task.Data.ID, err)
	}

	total := 0
	for _, att := range freshTask.Attachments {
		if att.Deleted {
			continue
		}
		a, err := archive.SaveAttachment(task, &att, false)
		if err != nil {
			return fmt.Errorf("save attachment %s for task %s: %w", att.ID, task.Data.ID, err)
		}
		if att.URLWithQuery == "" {
			log.Warn("attachment has no URL, skipping binary download", "id", att.ID, "title", att.Title)
			total++
			continue
		}
		data, err := f.client.HttpGet(att.URLWithQuery)
		if err != nil {
			return fmt.Errorf("download attachment %s for task %s: %w", att.ID, task.Data.ID, err)
		}
		binPath := filepath.Join(a.GetDir(), att.Title)
		if err := os.WriteFile(binPath, data, 0o644); err != nil {
			return fmt.Errorf("write attachment file %s: %w", binPath, err)
		}
		total++
	}

	log.Info("Attachments", "task_id", task.Data.ID, "count", total)

	if err := task.MarkAttachmentsDone(); err != nil {
		return fmt.Errorf("write attachments.done for task %s: %w", task.Data.ID, err)
	}
	return nil
}

func (f *FetchAttachments) Run() error {
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
