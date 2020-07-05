package user

import (
	"github.com/gghcode/apas-todo-apiserver/domain/entity"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/security"
)

type userService struct {
	userRepo Repository
	passport security.Passport
}

// NewService return user service
func NewService(userRepo Repository, passport security.Passport) UseCase {
	return &userService{
		userRepo: userRepo,
		passport: passport,
	}
}

func (srv *userService) CreateUser(req CreateUserRequest) (UserResponse, error) {
	var res UserResponse

	hashPassword, err := srv.passport.HashPassword(req.Password)
	if err != nil {
		return res, err
	}

	usr := entity.User{
		UserName:     req.UserName,
		PasswordHash: hashPassword,
	}

	createdUser, err := srv.userRepo.CreateUser(usr)
	if err != nil {
		return res, err
	}

	return UserResponse{
		ID:       createdUser.ID,
		UserName: createdUser.UserName,
	}, nil
}

func (srv *userService) GetUserByUserID(userID int64) (UserResponse, error) {
	var res UserResponse

	usr, err := srv.userRepo.UserByID(userID)
	if err != nil {
		return res, err
	}

	res = UserResponse{
		ID:       usr.ID,
		UserName: usr.UserName,
	}

	return res, nil
}

func (srv *userService) GetUserByUserName(userName string) (UserResponse, error) {
	var res UserResponse

	usr, err := srv.userRepo.UserByUserName(userName)
	if err != nil {
		return res, err
	}

	res = UserResponse{
		ID:       usr.ID,
		UserName: usr.UserName,
	}

	return res, nil
}
