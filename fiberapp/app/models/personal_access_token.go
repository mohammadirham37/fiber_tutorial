package models

import "time"

// PersonalAccessToken mewakili tabel "personal_access_tokens".
// Mirip tabel token milik Laravel Sanctum.
type PersonalAccessToken struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"not null;index" json:"user_id"`
	User       User       `json:"-"`
	Name       string     `gorm:"size:100;not null" json:"name"`
	Token      string     `gorm:"size:64;uniqueIndex;not null" json:"-"` // hash SHA-256, bukan token asli
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`
}
