package data

import (
	"time"
)

type Todo struct {
	ID          string    `json:"todoId"`
	Task        string    `json:"task"`
	IsCompleted bool      `json:"isCompleted"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
