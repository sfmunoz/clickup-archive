package api

type FolderSpace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Folder struct {
	ID               string      `json:"id"`
	Name             string      `json:"name"`
	Orderindex       int         `json:"orderindex"`
	OverrideStatuses bool        `json:"override_statuses"`
	Hidden           bool        `json:"hidden"`
	Space            FolderSpace `json:"space"`
	TaskCount        string      `json:"task_count"`
}

type FoldersResponse struct {
	Folders []Folder `json:"folders"`
}
