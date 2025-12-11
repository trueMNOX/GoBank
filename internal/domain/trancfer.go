package domain

import "time"

type Transfer struct {
	ID            int64
	FromAccountID int64
	ToAccountID   int64
	Amount        int64
	Currency      string
	CreatedAt     time.Time
}

type Entry struct {
	ID        int64
	AccountID int64
	Amount    int64
	CreatedAt time.Time
}
