package models

import (
	"gorm.io/gorm"
	"time"
)

type AccountModel struct {
	ID        int64          `gorm:"primaryKey;autoIncrement"`
	OwnerID   int64          `gorm:"not null;index"`
	Balance   int64          `gorm:"not null;default:0"`
	Currency  string         `gorm:"not null;size:3;default:'USD'"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Owner UserModel `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
}

func (AccountModel) TableName() string {
	return "accounts"
}
