package repository

import (
	"time"

	"github.com/gghcode/apas-todo-apiserver/domain/entity"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/user"
	myGorm "github.com/gghcode/apas-todo-apiserver/infra/gorm"
	"github.com/gghcode/apas-todo-apiserver/infra/gorm/model"
	"github.com/jinzhu/gorm"
	pg "github.com/lib/pq"
)

type repository struct {
	dbConn myGorm.Connection
}

// NewUserRepository godoc
func NewUserRepository(dbConn myGorm.Connection) user.Repository {
	return &repository{
		dbConn: dbConn,
	}
}

func (repo *repository) CreateUser(usr entity.User) (entity.User, error) {
	newUser := model.FromUserEntity(usr)
	newUser.CreatedAt = time.Now().Unix()

	err := repo.dbConn.DB().
		Create(&newUser).
		Error

	if pgErr, ok := err.(*pg.Error); ok && pgErr.Code == "23505" {
		return entity.User{}, user.ErrAlreadyExistUser
	} else if err != nil {
		return entity.User{}, err
	}

	return model.ToUserEntity(newUser), nil
}

func (repo *repository) UserByID(userID int64) (entity.User, error) {
	var u model.User

	err := repo.dbConn.DB().
		Where("id=?", userID).
		First(&u).
		Error

	if err == gorm.ErrRecordNotFound {
		return entity.User{}, user.ErrUserNotFound
	} else if err != nil {
		return entity.User{}, err
	}

	return model.ToUserEntity(u), nil
}

func (repo *repository) UserByUserName(username string) (entity.User, error) {
	var u model.User

	err := repo.dbConn.DB().
		Where("user_name=?", username).
		First(&u).
		Error

	if err == gorm.ErrRecordNotFound {
		return entity.User{}, user.ErrUserNotFound
	} else if err != nil {
		return entity.User{}, err
	}

	return entity.User{
		ID:           u.ID,
		UserName:     u.UserName,
		PasswordHash: u.PasswordHash,
	}, nil
	// return model.ToUserEntity(u), nil
}

func (repo *repository) UpdateUserByID(usr entity.User) (entity.User, error) {
	u, err := repo.userByID(usr.ID)
	if err != nil {
		return entity.User{}, err
	}

	err = repo.dbConn.DB().
		Model(&u).
		Updates(model.FromUserEntity(usr)).
		Error

	if err != nil {
		return entity.User{}, err
	}

	return model.ToUserEntity(u), nil
}

func (repo *repository) RemoveUserByID(userID int64) (entity.User, error) {
	u, err := repo.userByID(userID)
	if err != nil {
		return entity.User{}, err
	}

	err = repo.dbConn.DB().
		Delete(&u).
		Error

	if err != nil {
		return entity.User{}, err
	}

	return model.ToUserEntity(u), nil
}

func (repo *repository) userByID(userID int64) (model.User, error) {
	var u model.User

	err := repo.dbConn.DB().
		Where("id=?", userID).
		First(&u).
		Error

	if err == gorm.ErrRecordNotFound {
		return u, user.ErrUserNotFound
	} else if err != nil {
		return u, err
	}

	return u, nil
}
