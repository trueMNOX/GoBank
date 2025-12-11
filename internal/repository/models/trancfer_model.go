package models

import "time"

type TransferModel struct {
	ID            int64     `gorm:"primaryKey;autoIncrement"`
	FromAccountID int64     `gorm:"not null;index"`
	ToAccountID   int64     `gorm:"not null;index"`
	Amount        int64     `gorm:"not null"`
	Currency      string    `gorm:"not null;size:3"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`

	FromAccount AccountModel `gorm:"foreignKey:FromAccountID;constraint:OnDelete:CASCADE"`
	ToAccount   AccountModel `gorm:"foreignKey:ToAccountID;constraint:OnDelete:CASCADE"`
}

func (TransferModel) TableName() string {
	return "transfers"
}

type EntryModel struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	AccountID int64     `gorm:"not null;index"`
	Amount    int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	Account AccountModel `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE"`
}

func (EntryModel) TableName() string {
	return "entries"
}
