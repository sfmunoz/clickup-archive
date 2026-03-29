package api

type Comment struct {
	ID   string `json:"id"`
	Text string `json:"comment_text"`
}

type CommentsResponse struct {
	Comments []Comment `json:"comments"`
}
