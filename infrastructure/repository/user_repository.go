package repository

import (
	"time"

	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/gghcode/apas-todo-apiserver/domain/user"
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
	usr.CreatedAt = time.Now().Unix()

	err := repo.dbConn.DB().
		Create(&usr).
		Error

	if pgErr, ok := err.(*pg.Error); ok && pgErr.Code == "23505" {
		return user.User{}, user.ErrAlreadyExistUser
	} else if err != nil {
		return user.User{}, err
	}

	return usr, nil
}

func (repo *repository) AllUsers() ([]user.User, error) {
	var result []user.User

	err := repo.dbConn.DB().
		Find(&result).
		Error

	if err != nil {
		return nil, err
	}

	return result, err
}

func (repo *repository) UserByID(userID int64) (user.User, error) {
	var result user.User

	err := repo.dbConn.DB().
		Where("id=?", userID).
		First(&result).
		Error

	if err == gorm.ErrRecordNotFound {
		return result, user.ErrUserNotFound
	} else if err != nil {
		return result, err
	}

	return result, nil
}

func (repo *repository) UserByUserName(username string) (user.User, error) {
	var result user.User

	err := repo.dbConn.DB().
		Where("user_name=?", username).
		First(&result).
		Error

	if err == gorm.ErrRecordNotFound {
		return result, user.ErrUserNotFound
	} else if err != nil {
		return result, err
	}

	return result, nil
}

func (repo *repository) UpdateUserByID(user user.User) (user.User, error) {
	entity, err := repo.UserByID(user.ID)
	if err != nil {
		return entity, err
	}

	err = repo.dbConn.DB().
		Model(&entity).
		Updates(&user).
		Error

	if err != nil {
		return entity, err
	}

	return entity, nil
}

func (repo *repository) RemoveUserByID(userID int64) (user.User, error) {
	entity, err := repo.UserByID(userID)
	if err != nil {
		return entity, err
	}

	err = repo.dbConn.DB().
		Delete(&entity).
		Error

	if err != nil {
		return entity, err
	}

	return entity, nil
}
