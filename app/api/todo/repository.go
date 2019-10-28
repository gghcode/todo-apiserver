package todo

import (
	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/jinzhu/gorm"
)

// Repository godoc
type Repository interface {
	AddTodo(todo Todo) (Todo, error)
	AllTodosByUserID(userID int64) ([]Todo, error)
	TodoByTodoID(todoID string, todo *Todo) error
	UpdateTodo(todoID string, todo map[string]interface{}) (Todo, error)
	RemoveTodo(todoID string) error
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

func (repo *repository) TodoByTodoID(todoID string, todo *Todo) error {
	err := repo.pgConn.DB().
		Where("id=?", todoID).
		First(&todo).Error

	if err == gorm.ErrRecordNotFound {
		return ErrNotFoundTodo
	} else if err != nil {
		return err
	}

	return nil
}

func (repo *repository) UpdateTodo(todoID string, todoData map[string]interface{}) (Todo, error) {
	var todo Todo
	if err := repo.TodoByTodoID(todoID, &todo); err != nil {
		return todo, err
	}

	repo.pgConn.DB().Model(&todo).UpdateColumns(todoData)

	return todo, nil
}

func (repo *repository) RemoveTodo(todoID string) error {
	var todo Todo
	if err := repo.TodoByTodoID(todoID, &todo); err != nil {
		return err
	}

	if err := repo.pgConn.DB().Delete(todo).Error; err != nil {
		return err
	}

	return nil
}
