package todo

import (
	"time"
)

type TodoAddReq struct {
	Title       string    `json:"title" validate:"required,min=4,max=50"`
	Description string    `json:"description" validate:"required,min=0,max=320"`
	DueDate     time.Time `json:"dueDate" validate:"required,not-stale"`
}

type Todo struct {
	TodoId      int64     `json:"todo_id"`
	OwnerId     int64     `json:"-"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      int16     `json:"status"`
	DueDate     time.Time `json:"dueDate"`
	CreatedAt   time.Time `json:"createdAt"`
}

type TodoUpdateReq struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *int16     `json:"status"`
	DueDate     *time.Time `json:"dueDate"`
}

func NewFromAdd(t *TodoAddReq, ownerID int64) *Todo {
	return &Todo{
		OwnerId:     ownerID,
		Title:       t.Title,
		Description: t.Description,
		Status:      0,
		DueDate:     t.DueDate,
	}
}
