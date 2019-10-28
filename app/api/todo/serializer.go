package todo

import "time"

type (
	// TodoResponse godoc
	TodoResponse struct {
		ID         string    `json:"id"`
		Title      string    `json:"title"`
		Contents   string    `json:"contents"`
		AssignorID int64     `json:"assignor_id"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
		DueDate    string    `json:"due_date"`
	}

	// TodoSerializer godoc
	TodoSerializer struct {
		Model Todo
	}

	// TodosSerializer godoc
	TodosSerializer struct {
		Model []Todo
	}
)

// Response godoc
func (serializer TodoSerializer) Response() TodoResponse {
	return serializeTodoResponse(serializer.Model)
}

// Response godoc
func (serializer TodosSerializer) Response() []TodoResponse {
	result := []TodoResponse{}

	for _, model := range serializer.Model {
		result = append(result, serializeTodoResponse(model))
	}

	return result
}

func serializeTodoResponse(todo Todo) TodoResponse {
	date := ""
	if todo.DueDate.Valid {
		date = todo.DueDate.Time.String()
	}

	return TodoResponse{
		ID:         todo.ID.String(),
		Title:      todo.Title,
		Contents:   todo.Contents,
		AssignorID: todo.AssignorID,
		CreatedAt:  todo.CreatedAt,
		UpdatedAt:  todo.UpdatedAt,
		DueDate:    date,
	}
}
