package todo

import "time"

// TodoResponse godoc
type TodoResponse struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Contents   string    `json:"contents"`
	AssignorID int64     `json:"assignor_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TodoSerializer godoc
type TodoSerializer struct {
	Model Todo
}

// Response godoc
func (serializer TodoSerializer) Response() TodoResponse {
	return serializeTodoResponse(serializer.Model)
}

// TodosSerializer godoc
type TodosSerializer struct {
	Model []Todo
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
	return TodoResponse{
		ID:         todo.ID.String(),
		Title:      todo.Title,
		Contents:   todo.Contents,
		AssignorID: todo.AssignorID,
		CreatedAt:  todo.CreatedAt,
		UpdatedAt:  todo.UpdatedAt,
	}
}
