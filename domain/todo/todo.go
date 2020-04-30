package todo

import "time"

// Todo is todo entity model
type Todo struct {
	ID string
	Title string
	Contents string
	AssignorID int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
