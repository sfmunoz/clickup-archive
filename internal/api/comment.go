package api

type CommentUser struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	Initials       string `json:"initials"`
	Email          string `json:"email"`
	Color          string `json:"color"`
	ProfilePicture string `json:"profilePicture"`
}

type CommentAttributes struct {
	Bold      bool   `json:"bold,omitempty"`
	Italic    bool   `json:"italic,omitempty"`
	Code      bool   `json:"code,omitempty"`
	Link      string `json:"link,omitempty"`
	Indent    int    `json:"indent,omitempty"`
}

type CommentItem struct {
	Text       string            `json:"text,omitempty"`
	Attributes CommentAttributes `json:"attributes,omitempty"`
	Type       string            `json:"type,omitempty"`
	Emoticon   *struct {
		Code string `json:"code"`
	} `json:"emoticon,omitempty"`
	User *struct {
		ID int `json:"id"`
	} `json:"user,omitempty"`
}

type Comment struct {
	ID          string        `json:"id"`
	Comment     []CommentItem `json:"comment"`
	Text        string        `json:"comment_text"`
	User        CommentUser   `json:"user"`
	Resolved    bool          `json:"resolved"`
	Assignee    *CommentUser  `json:"assignee"`
	AssignedBy  *CommentUser  `json:"assigned_by"`
	Reactions   []interface{} `json:"reactions"`
	Date        string        `json:"date"`
	ReplyCount  int           `json:"reply_count"`
}

type CommentsResponse struct {
	Comments []Comment `json:"comments"`
}
