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
		Title    *string
		Contents *string
		DueDate  *string
	}

	// TodoResponse godoc
	TodoResponse struct {
		ID         string
		Title      string
		Contents   string
		AssignorID int64
		CreatedAt  time.Time
		UpdatedAt  time.Time
		DueDate    string
	}
)
