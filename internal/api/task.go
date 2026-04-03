package api

import "encoding/json"

type TaskStatus struct {
	Status     string `json:"status"`
	Color      string `json:"color"`
	Type       string `json:"type"`
	Orderindex int    `json:"orderindex"`
}

type TaskPriority struct {
	ID         string `json:"id"`
	Priority   string `json:"priority"`
	Color      string `json:"color"`
	Orderindex string `json:"orderindex"`
}

type TaskTag struct {
	Name    string `json:"name"`
	TagFg   string `json:"tag_fg"`
	TagBg   string `json:"tag_bg"`
	Creator int    `json:"creator"`
}

type TaskList struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Access bool   `json:"access"`
}

type TaskFolder struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Hidden bool   `json:"hidden"`
	Access bool   `json:"access"`
}

type TaskSpace struct {
	ID string `json:"id"`
}

type TaskCustomField struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Type           string          `json:"type"`
	TypeConfig     json.RawMessage `json:"type_config"`
	DateCreated    string          `json:"date_created"`
	HideFromGuests bool            `json:"hide_from_guests"`
	Value          json.RawMessage `json:"value"`
	Required       bool            `json:"required"`
}

type TaskChecklistItem struct {
	ID          string              `json:"id"`
	TaskID      string              `json:"task_id"`
	Name        string              `json:"name"`
	DateCreated string              `json:"date_created"`
	Orderindex  float64             `json:"orderindex"`
	Creator     int                 `json:"creator"`
	Resolved    bool                `json:"resolved"`
	Parent      *string             `json:"parent"`
	Assignee    *CommentUser        `json:"assignee"`
	Items       []TaskChecklistItem `json:"items"`
}

type TaskChecklist struct {
	ID          string              `json:"id"`
	TaskID      string              `json:"task_id"`
	Name        string              `json:"name"`
	DateCreated string              `json:"date_created"`
	Orderindex  float64             `json:"orderindex"`
	Creator     int                 `json:"creator"`
	Resolved    int                 `json:"resolved"`
	Unresolved  int                 `json:"unresolved"`
	Items       []TaskChecklistItem `json:"items"`
}

type TaskAttachment struct {
	ID               string      `json:"id"`
	Version          json.RawMessage `json:"version"`
	Date             string      `json:"date"`
	Title            string      `json:"title"`
	Extension        string      `json:"extension"`
	ThumbnailSmall   string      `json:"thumbnail_small"`
	ThumbnailMedium  string      `json:"thumbnail_medium"`
	ThumbnailLarge   string      `json:"thumbnail_large"`
	IsFolder         bool        `json:"is_folder"`
	Mimetype         string      `json:"mimetype"`
	Hidden           bool        `json:"hidden"`
	ParentID         string      `json:"parent_id"`
	Size             int64       `json:"size"`
	TotalComments    int         `json:"total_comments"`
	ResolvedComments int         `json:"resolved_comments"`
	User             CommentUser `json:"user"`
	Deleted          bool        `json:"deleted"`
	URL              string      `json:"url"`
	URLWithQuery     string      `json:"url_w_query"`
	URLWithHost      string      `json:"url_w_host"`
}

type TaskDependency struct {
	TaskID      string `json:"task_id"`
	DependsOn   string `json:"depends_on"`
	Type        int    `json:"type"`
	DateCreated string `json:"date_created"`
	UserID      string `json:"userid"`
	WorkspaceID string `json:"workspace_id"`
}

type TaskLinkedTask struct {
	TaskID      string `json:"task_id"`
	LinkID      string `json:"link_id"`
	DateCreated string `json:"date_created"`
	UserID      string `json:"userid"`
}

type Task struct {
	ID             string            `json:"id"`
	CustomID       *string           `json:"custom_id"`
	CustomItemID   *int              `json:"custom_item_id"`
	Name           string            `json:"name"`
	TextContent    string            `json:"text_content"`
	Description    string            `json:"description"`
	Status         TaskStatus        `json:"status"`
	Orderindex     string            `json:"orderindex"`
	DateCreated    string            `json:"date_created"`
	DateUpdated    string            `json:"date_updated"`
	DateClosed     *string           `json:"date_closed"`
	DateDone       *string           `json:"date_done"`
	Archived       bool              `json:"archived"`
	Creator        CommentUser       `json:"creator"`
	Assignees      []CommentUser     `json:"assignees"`
	Watchers       []CommentUser     `json:"watchers"`
	Checklists     []TaskChecklist   `json:"checklists"`
	Tags           []TaskTag         `json:"tags"`
	Parent         *string           `json:"parent"`
	Priority       *TaskPriority     `json:"priority"`
	DueDate        *string           `json:"due_date"`
	StartDate      *string           `json:"start_date"`
	Points         *float64          `json:"points"`
	TimeEstimate   *int64            `json:"time_estimate"`
	TimeSpent      int64             `json:"time_spent"`
	CustomFields   []TaskCustomField `json:"custom_fields"`
	List           TaskList          `json:"list"`
	Folder         TaskFolder        `json:"folder"`
	Space          TaskSpace         `json:"space"`
	URL            string            `json:"url"`
	PermissionLevel string           `json:"permission_level"`
	TeamID         string            `json:"team_id"`
	Subtasks       []Task            `json:"subtasks"`
	Attachments    []TaskAttachment  `json:"attachments"`
	Dependencies   []TaskDependency  `json:"dependencies"`
	LinkedTasks    []TaskLinkedTask  `json:"linked_tasks"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}
