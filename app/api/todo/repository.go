package todo

import "gitlab.com/gyuhwan/apas-todo-apiserver/db"

// Repository godoc
type Repository interface {
	AddTodo(todo Todo) (Todo, error)
	AllTodosByUserID(userID int64) ([]Todo, error)
}

type repository struct {
	pgConn *db.PostgresConn
}

// NewRepository godoc
func NewRepository(pgConn *db.PostgresConn) Repository {
	pgConn.DB().AutoMigrate(Todo{})

	return &repository{
		pgConn: pgConn,
	}
}

func (repo *repository) AddTodo(todo Todo) (Todo, error) {
	err := repo.pgConn.DB().
		Create(&todo).
		Error

	if err != nil {
		return EmptyTodo, err
	}

	return todo, nil
}

func (repo *repository) AllTodosByUserID(userID int64) ([]Todo, error) {
	var result []Todo

	err := repo.pgConn.DB().
		Where("assignor_id = ?", userID).
		Find(&result).
		Error

	if err != nil {
		return nil, err
	}

	return result, nil
}
