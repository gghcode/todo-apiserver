package user

import (
	"time"

	"github.com/gghcode/apas-todo-apiserver/domain/user"
)

type (
	userResponseDTO struct {
		ID        int64     `json:"id"`
		UserName  string    `json:"username"`
		CreatedAt time.Time `json:"created_at"`
	}

	userResponseSerializer struct {
		model user.UserResponse
	}
)

func newUserResponseSerializer(model user.UserResponse) *userResponseSerializer {
	return &userResponseSerializer{
		model: model,
	}
}

func (s *userResponseSerializer) Response() userResponseDTO {
	return userResponseDTO{
		ID:        s.model.ID,
		UserName:  s.model.UserName,
		CreatedAt: s.model.CreatedAt,
	}
}
