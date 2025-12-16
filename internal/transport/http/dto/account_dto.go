package dto

import "time"

type CreateAccountRequest struct {
	Currency string `json:"currency" binding:"required,len=3"`
}

type AccountResponse struct {
	ID        int64     `json:"id"`
	OwnerID   int64     `json:"owner_id"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type ListAccountsRequest struct {
	Page     int `form:"page" binding:"required,min=1"`
	PageSize int `form:"page_size" binding:"required,min=5,max=100"`
}
