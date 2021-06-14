package db

import "time"

type Todo struct {
	ID        uint64
	UserID    uint64
	User      User
	Text      string
	Done      bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
