package todos

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	TodoID      uuid.UUID `json:"todoId"`
	Task        string    `json:"task"`
	IsCompleted bool      `json:"isCompleted"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
