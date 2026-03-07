package models

import (
	"time"
)

type Session struct {
	ID        uint      `gorm:"primarykey"`
	Token     string    `gorm:"type:varchar(128);uniqueIndex"`
	Username  string    `gorm:"type:varchar(255);index"`
	IPAddress string    `gorm:"type:varchar(64)"`
	UserAgent string    `gorm:"type:varchar(512)"`
	ExpiresAt time.Time `gorm:"index"`
	CreatedAt time.Time
}
