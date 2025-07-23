package entity

import "time"

type Post struct {
	ID        int        `json:"id" gorm:"primaryKey"`
	UserID    int        `json:"user_id"`
	Title     string     `json:"title"`
	ImageURL  string     `json:"image_url"`
	Content   string     `json:"content"`
	IsPublic  bool       `json:"is_public"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Comments  []*Comment `json:"comments"`
	Likes     []*Like    `json:"likes"`
}
