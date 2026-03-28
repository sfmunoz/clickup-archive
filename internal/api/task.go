package api

type Task struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}
