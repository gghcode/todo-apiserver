package todo

// UsecaseInteractor is todo usecase interface
type UsecaseInteractor interface {
	AddTodo(AddTodoRequest) (TodoResponse, error)
}
