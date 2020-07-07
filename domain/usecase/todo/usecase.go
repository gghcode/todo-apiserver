package todo

// UseCase is todo usecase interface
type UseCase interface {
	AddTodo(AddTodoRequest) (TodoResponse, error)
	GetTodosByUserID(userID int64) ([]TodoResponse, error)
	RemoveTodo(todoID string) error
}
