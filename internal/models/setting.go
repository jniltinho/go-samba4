package models

import (
	"time"
)

type Setting struct {
	ID        uint   `gorm:"primarykey"`
	Key       string `gorm:"type:varchar(100);uniqueIndex"`
	Value     string `gorm:"type:text"`
	UpdatedAt time.Time
}
