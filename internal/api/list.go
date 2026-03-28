package api

type ListFolder struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Hidden bool   `json:"hidden"`
	Access bool   `json:"access"`
}

type ListSpace struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Access bool   `json:"access"`
}

type List struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Orderindex int        `json:"orderindex"`
	TaskCount  int        `json:"task_count"`
	Folder     ListFolder `json:"folder"`
	Space      ListSpace  `json:"space"`
	Archived   bool       `json:"archived"`
}

type ListsResponse struct {
	Lists []List `json:"lists"`
}
