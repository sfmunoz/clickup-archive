package api

type WorkspaceMemberUser struct {
	ID                int    `json:"id"`
	Username          string `json:"username"`
	Email             string `json:"email"`
	Color             string `json:"color"`
	ProfilePicture    string `json:"profilePicture"`
	Initials          string `json:"initials"`
	WeekStartDay      int    `json:"week_start_day"`
	GlobalFontSupport bool   `json:"global_font_support"`
	Timezone          string `json:"timezone"`
}

type WorkspaceMember struct {
	User WorkspaceMemberUser `json:"user"`
}

type Workspace struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	Color   string            `json:"color"`
	Avatar  string            `json:"avatar"`
	Members []WorkspaceMember `json:"members"`
}

type WorkspacesResponse struct {
	Workspaces []Workspace `json:"teams"`
}
