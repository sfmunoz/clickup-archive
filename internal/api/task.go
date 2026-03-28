package api

type Task struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Subtasks []Task `json:"subtasks"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}
