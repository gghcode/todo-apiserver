package todo

import (
	"errors"
	"net/http"

	"github.com/gghcode/apas-todo-apiserver/app/api"
)

var (
	// ErrEmptyTodoID godoc
	ErrEmptyTodoID = api.NewHandledError(
		http.StatusBadRequest,
		errors.New("Empty Todo ID"),
	)

	// ErrNotFoundTodo godoc
	ErrNotFoundTodo = api.NewHandledError(
		http.StatusNotFound,
		errors.New("Not Found Todo"),
	)
)
