package repository

import (
	"time"

	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/gghcode/apas-todo-apiserver/domain/user"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/model"
	"github.com/jinzhu/gorm"
	pg "github.com/lib/pq"
)

type repository struct {
	dbConn db.GormConnection
}

// NewUserRepository godoc
func NewUserRepository(dbConn db.GormConnection) user.Repository {
	return &repository{
		dbConn: dbConn,
	}
}

func (repo *repository) CreateUser(usr user.User) (user.User, error) {
	newUser := model.FromUserEntity(usr)
	newUser.CreatedAt = time.Now().Unix()

	err := repo.dbConn.DB().
		Create(&newUser).
		Error

	if pgErr, ok := err.(*pg.Error); ok && pgErr.Code == "23505" {
		return user.User{}, user.ErrAlreadyExistUser
	} else if err != nil {
		return user.User{}, err
	}

	return model.ToUserEntity(newUser), nil
}

func (repo *repository) AllUsers() ([]user.User, error) {
	var u []model.User

	err := repo.dbConn.DB().
		Find(&u).
		Error

	if err != nil {
		return nil, err
	}

	return nil, err
}

func (repo *repository) UserByID(userID int64) (user.User, error) {
	var u model.User

	err := repo.dbConn.DB().
		Where("id=?", userID).
		First(&u).
		Error

	if err == gorm.ErrRecordNotFound {
		return user.User{}, user.ErrUserNotFound
	} else if err != nil {
		return user.User{}, err
	}

	return model.ToUserEntity(u), nil
}

func (repo *repository) UserByUserName(username string) (user.User, error) {
	var u model.User

	err := repo.dbConn.DB().
		Where("user_name=?", username).
		First(&u).
		Error

	if err == gorm.ErrRecordNotFound {
		return user.User{}, user.ErrUserNotFound
	} else if err != nil {
		return user.User{}, err
	}

	return model.ToUserEntity(u), nil
}

func (repo *repository) UpdateUserByID(usr user.User) (user.User, error) {
	u, err := repo.userByID(usr.ID)
	if err != nil {
		return user.User{}, err
	}

	err = repo.dbConn.DB().
		Model(&u).
		Updates(model.FromUserEntity(usr)).
		Error

	if err != nil {
		return user.User{}, err
	}

	return model.ToUserEntity(u), nil
}

func (repo *repository) RemoveUserByID(userID int64) (user.User, error) {
	u, err := repo.userByID(userID)
	if err != nil {
		return user.User{}, err
	}

	err = repo.dbConn.DB().
		Delete(&u).
		Error

	if err != nil {
		return user.User{}, err
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
