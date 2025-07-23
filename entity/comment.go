package entity

import "time"

type Comment struct {
	ID          int        `json:"id" gorm:"primaryKey"`
	UserID      int        `json:"user_id"`
	PostID      int        `json:"post_id"`
	Comment     string     `json:"comment"`
	CommentLike int        `json:"comment_like"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
