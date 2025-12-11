package models

import (
	"gorm.io/gorm"
	"time"
)

type UserModel struct {
	ID             int64          `gorm:"primaryKey;autoIncrement"`
	Username       string         `gorm:"uniqueIndex;not null;size:50"`
	FullName       string         `gorm:"not null;size:100"`
	Email          string         `gorm:"uniqueIndex;not null;size:100"`
	HashedPassword string         `gorm:"not null"`
	CreatedAt      time.Time      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (UserModel) TableName() string {
	return "users"
}
