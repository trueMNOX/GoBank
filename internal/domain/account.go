package domain

import "time"

type Account struct {
	ID        int64
	OwnerID   int64
	Balance   int64
	Currency  string
	CreatedAt time.Time
}
