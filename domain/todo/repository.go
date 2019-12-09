package todo

// Repository godoc
type Repository interface {
	AddTodo(todo Todo) (Todo, error)
	AllTodosByUserID(userID int64) ([]Todo, error)
	TodoByTodoID(todoID string, todo *Todo) error
	UpdateTodo(todoID string, todo map[string]interface{}) (Todo, error)
	RemoveTodo(todoID string) error
}
