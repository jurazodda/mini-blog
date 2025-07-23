package entity

type PostList struct {
	PostID   int    `json:"post_id"`
	UserID   int    `json:"user_id"`
	Title    string `json:"title"`
	Comments int    `json:"comments"`
	Likes    int    `json:"likes"`
}
