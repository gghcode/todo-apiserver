package repo

import (
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/todo"
	myGorm "github.com/gghcode/apas-todo-apiserver/infra/gorm"
	"github.com/gghcode/apas-todo-apiserver/infra/gorm/model"

	"github.com/jinzhu/gorm"
)

type todoRepo struct {
	dbConn myGorm.Connection
}

// NewTodoRepository return new todo repository
func NewTodoRepository(dbConn myGorm.Connection) todo.Repository {
	return &todoRepo{
		dbConn: dbConn,
	}
}

func (r *todoRepo) AddTodo(t todo.Todo) (todo.Todo, error) {
	newTodo := model.FromTodoEntity(t)

	err := r.dbConn.DB().
		Create(&newTodo).
		Error

	if err != nil {
		return todo.Todo{}, err
	}

	return model.ToTodoEntity(newTodo), nil
}
func (r *todoRepo) AllTodosByUserID(userID int64) ([]todo.Todo, error) {
	var t []model.Todo

	err := r.dbConn.DB().
		Where("assignor_id = ?", userID).
		Find(&t).
		Error

	if err != nil {
		return nil, err
	}

	return model.ToTodoEntityArray(t), nil
}

func (r *todoRepo) TodoByTodoID(todoID string) (todo.Todo, error) {
	var t model.Todo

	err := r.dbConn.DB().
		Where("id=?", todoID).
		First(&t).Error

	if err == gorm.ErrRecordNotFound {
		return todo.Todo{}, todo.ErrNotFoundTodo
	} else if err != nil {
		return todo.Todo{}, err
	}

	return model.ToTodoEntity(t), nil
}

func (r *todoRepo) UpdateTodo(todoID string, data map[string]interface{}) (todo.Todo, error) {
	t, err := r.TodoByTodoID(todoID)
	if err != nil {
		return t, err
	}

	r.dbConn.DB().Model(&t).UpdateColumns(data)

	return t, nil
}

func (r *todoRepo) RemoveTodo(todoID string) error {
	t, err := r.TodoByTodoID(todoID)
	if err != nil {
		return err
	}

	if err := r.dbConn.DB().Delete(t).Error; err != nil {
		return err
	}

	return nil
}
