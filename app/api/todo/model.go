package todo

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// EmptyTodo godoc
var EmptyTodo = Todo{}

// Todo godoc
type Todo struct {
	ID         uuid.UUID `gorm:"type:uuid"`
	Title      string    `gorm:"not null;"`
	Contents   string    `gorm:"not null;"`
	AssignorID int64     `gorm:"not null;"`
	CreatedAt  time.Time `sql:"DEFAULT:current_timestamp"`
	UpdatedAt  time.Time `sql:"DEFAULT:current_timestamp"`
}

// BeforeSave godoc
func (todo *Todo) BeforeSave() (err error) {
	todo.ID = uuid.NewV4()
	// time.RFC1123
	// todo.CreatedAt, _ = time.Parse(time.RFC1123, time.Now().Format(time.RFC1123))
	// todo.UpdatedAt, _ = time.Parse(time.RFC1123, time.Now().Format(time.RFC1123))

	return
}

// BeforeUpdate godoc
func (todo *Todo) BeforeUpdate() (err error) {
	todo.UpdatedAt = time.Now()

	return
}
