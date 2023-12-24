package models

import "time"

type Token struct {
	UserId    uint      `gorm:"primary_key;column:user_id"`
	Token     string    `gorm:"column:token"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
}
