package models

import "time"

// Post mewakili tabel "posts".
type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:200;not null" json:"title"`
	Body      string    `gorm:"type:text;not null" json:"body"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	User      User      `json:"user"` // relasi belongsTo
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
