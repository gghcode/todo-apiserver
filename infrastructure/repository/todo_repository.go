package repository

import (
	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/gghcode/apas-todo-apiserver/domain/todo"
)

type gormTodoRepository struct {
	dbConn db.GormConnection
}

// NewGormTodoRepository return new todo repository
func NewGormTodoRepository(dbConn db.GormConnection) todo.Repository {
	dbConn.DB().AutoMigrate(todo.Todo{})

	return &gormTodoRepository{
		dbConn: dbConn,
	}
}

func (r *gormTodoRepository) AddTodo(t todo.Todo) (todo.Todo, error) {
	err := r.dbConn.DB().
		Create(&t).
		Error

	if err != nil {
		return todo.Todo{}, err
	}

	return t, nil
}
func (r *gormTodoRepository) AllTodosByUserID(userID int64) ([]todo.Todo, error) {
	return nil, nil
}

func (r *gormTodoRepository) TodoByTodoID(todoID string, todo *todo.Todo) error {
	return nil
}

func (r *gormTodoRepository) UpdateTodo(todoID string, data map[string]interface{}) (todo.Todo, error) {
	return todo.Todo{}, nil
}

func (r *gormTodoRepository) RemoveTodo(todoID string) error {
	return nil
}
