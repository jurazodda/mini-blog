package entities

import "time"

type Repost struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	PostID    int        `json:"post_id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	ImageURL  *string    `json:"image_url"`
	Comments  int        `json:"comments"`
	Likes     int        `json:"likes"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
