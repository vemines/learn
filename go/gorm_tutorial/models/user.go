package models

import (
	"time"
)

type User struct {
	UserId        uint       `gorm:"primary_key;column:user_id"`
	Username      string     `gorm:"column:username"`
	Password      string     `gorm:"column:passwd"`
	Email         string     `gorm:"column:email"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	LastUpdatedAt time.Time  `gorm:"column:last_updated_at"`
	Status        UserStatus `gorm:"column:user_status;type:string"`
}

type UserStatus string // enum represents the status of a user.

const (
	UserStatusInactive UserStatus = "Inactive"
	UserStatusActive   UserStatus = "Active"
	UserStatusDeleted  UserStatus = "Deleted"
)

type UserCountGroupedByStatus struct {
	Status UserStatus `gorm:"column:user_status;type:string"`
	Count  int        `json:"count"`
}
