package todo

type todoService struct {
	todoRepo Repository
}

// NewTodoService return new todo service
func NewTodoService(todoRepo Repository) UsecaseInteractor {
	return &todoService{
		todoRepo: todoRepo,
	}
}

func (srv *todoService) AddTodo(req AddTodoRequest) (TodoResponse, error) {
	todo := Todo{
		Title:    req.Title,
		Contents: req.Contents,
	}

	insertedTodo, err := srv.todoRepo.AddTodo(todo)
	if err != nil {
		return TodoResponse{}, err
	}

	return TodoResponse{
		ID:         insertedTodo.ID.String(),
		Title:      insertedTodo.Title,
		Contents:   insertedTodo.Contents,
		AssignorID: insertedTodo.AssignorID,
		// DueDate:    insertedTodo.DueDate,
	}, nil
}
