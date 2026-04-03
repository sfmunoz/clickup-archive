package archive

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sfmunoz/clickup-archive/internal/api"
)

type Attachment struct {
	Parent *Task
	Data   *api.TaskAttachment
}

func (a *Attachment) GetDir() string {
	return attachmentDir(a.Parent.GetDir(), a.Data.ID)
}

func LoadAttachment(parent *Task, id string) (*Attachment, error) {
	dir := attachmentDir(parent.GetDir(), id)
	if err := isFolder(dir); err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(indexFile(dir))
	if err != nil {
		return nil, err
	}
	var data api.TaskAttachment
	if err := json.Unmarshal(buf, &data); err != nil {
		return nil, err
	}
	return &Attachment{
		Parent: parent,
		Data:   &data,
	}, nil
}

func SaveAttachment(parent *Task, a *api.TaskAttachment, update bool) (*Attachment, error) {
	var aOld *Attachment = nil
	for _, ch := range parent.Attachments {
		if ch.Data.ID != a.ID {
			continue
		}
		if !update {
			return nil, fmt.Errorf("attachment '%s' already exists and 'update' is false", ch.Data.ID)
		}
		aOld = ch
		break
	}
	dir := attachmentDir(parent.GetDir(), a.ID)
	if err := jsonSave(a, dir); err != nil {
		return nil, err
	}
	if aOld == nil {
		log.Info("attachment created", "id", a.ID, "title", a.Title)
		aNew := &Attachment{Parent: parent, Data: a}
		parent.Attachments = append(parent.Attachments, aNew)
		return aNew, nil
	}
	log.Warn("attachment updated", "id_old", aOld.Data.ID, "id_new", a.ID)
	aOld.Data = a
	return aOld, nil
}
