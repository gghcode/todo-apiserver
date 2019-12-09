package todo

// UsecaseInteractor is todo usecase interface
type UsecaseInteractor interface {
	AddTodo(AddTodoRequest) (TodoResponse, error)
	GetTodosByUserID(userID int64) ([]TodoResponse, error)
	RemoveTodo(todoID string) error
}
