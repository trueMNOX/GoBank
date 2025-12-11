package domain

import "time"

type User struct {
	ID        int64
	Username  string
	FullName  string
	Email     string
	CreatedAt time.Time
}

type UserResponse struct {
    ID        int64     `json:"id"`
    Username  string    `json:"username"`
    FullName  string    `json:"full_name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}
