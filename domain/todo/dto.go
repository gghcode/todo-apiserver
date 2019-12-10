package todo

import "time"

type (
	// AddTodoRequest godoc
	// struct tag use swagger
	AddTodoRequest struct {
		Title      string
		Contents   string
		AssignorID int64
	}

	// UpdateTodoRequest godoc
	UpdateTodoRequest struct {
		Title    *string `json:"title"`
		Contents *string `json:"contents"`
		DueDate  *string `json:"due_date"`
	}

	// TodoResponse godoc
	TodoResponse struct {
		ID         string    `json:"id"`
		Title      string    `json:"title"`
		Contents   string    `json:"contents"`
		AssignorID int64     `json:"assignor_id"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		DueDate    string    `json:"due_date"`
	}
)
