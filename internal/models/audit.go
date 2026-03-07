package models

import (
	"time"
)

type AuditLog struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"index"`
	AdminUser string    `gorm:"type:varchar(255);index"`
	IPAddress string    `gorm:"type:varchar(64)"`
	Action    string    `gorm:"type:varchar(100)"`
	ObjectDN  string    `gorm:"type:text"`
	Details   string    `gorm:"type:text"`
}
