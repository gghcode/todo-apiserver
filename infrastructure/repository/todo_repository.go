package repository

import (
	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/gghcode/apas-todo-apiserver/domain/todo"
	"github.com/jinzhu/gorm"
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
	var result []todo.Todo

	err := r.dbConn.DB().
		Where("assignor_id = ?", userID).
		Find(&result).
		Error

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *gormTodoRepository) TodoByTodoID(todoID string) (todo.Todo, error) {
	var result todo.Todo

	err := r.dbConn.DB().
		Where("id=?", todoID).
		First(&result).Error

	if err == gorm.ErrRecordNotFound {
		return result, todo.ErrNotFoundTodo
	} else if err != nil {
		return result, err
	}

	return result, nil
}

func (r *gormTodoRepository) UpdateTodo(todoID string, data map[string]interface{}) (todo.Todo, error) {
	t, err := r.TodoByTodoID(todoID)
	if err != nil {
		return t, err
	}

	r.dbConn.DB().Model(&t).UpdateColumns(data)

	return t, nil
}

func (r *gormTodoRepository) RemoveTodo(todoID string) error {
	t, err := r.TodoByTodoID(todoID)
	if err != nil {
		return err
	}

	if err := r.dbConn.DB().Delete(t).Error; err != nil {
		return err
	}

	return nil
}
