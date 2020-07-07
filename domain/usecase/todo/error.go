package todo

import (
	"errors"
)

var (
	// ErrEmptyTodoID godoc
	ErrEmptyTodoID = errors.New("Empty Todo ID")

	// ErrNotFoundTodo godoc
	ErrNotFoundTodo = errors.New("Not Found Todo")
)
