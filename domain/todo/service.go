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

func (srv *todoService) GetTodosByUserID(userID int64) ([]TodoResponse, error) {
	todos, err := srv.todoRepo.AllTodosByUserID(userID)
	if err != nil {
		return nil, err
	}

	var res []TodoResponse
	for _, todo := range todos {
		res = append(res, TodoResponse{
			ID:         todo.ID.String(),
			Title:      todo.Title,
			Contents:   todo.Contents,
			AssignorID: todo.AssignorID,
		})
	}

	return res, nil
}

func (srv *todoService) RemoveTodo(todoID string) error {
	return srv.todoRepo.RemoveTodo(todoID)
}
