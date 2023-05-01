package data

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID          uuid.UUID `json:"todoId"`
	Task        string    `json:"task"`
	IsCompleted bool      `json:"isCompleted"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
