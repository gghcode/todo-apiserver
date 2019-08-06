package user

import (
	"time"

	"github.com/jinzhu/gorm"
	pg "github.com/lib/pq"
	"gitlab.com/gyuhwan/apas-todo-apiserver/db"
)

// Repository godoc
type Repository interface {
	CreateUser(User) (User, error)
	AllUsers() ([]User, error)
	UserByID(userID int64) (User, error)
	UserByUserName(username string) (User, error)
	UpdateUserByID(user User) (User, error)
	RemoveUserByID(userID int64) (User, error)
}

type repository struct {
	dbConn *db.PostgresConn
}

// NewRepository godoc
func NewRepository(postgres *db.PostgresConn) Repository {
	postgres.DB().AutoMigrate(User{})

	return &repository{
		dbConn: postgres,
	}
}

func (repo *repository) CreateUser(user User) (User, error) {
	user.CreatedAt = time.Now().Unix()

	err := repo.dbConn.DB().
		Create(&user).
		Error

	if pgErr, ok := err.(*pg.Error); ok {
		if pgErr.Code == "23505" {
			return EmptyUser, ErrAlreadyExistUser
		}

		return EmptyUser, err
	} else if err != nil {
		return EmptyUser, err
	}

	return user, nil
}

func (repo *repository) AllUsers() ([]User, error) {
	var result []User

	err := repo.dbConn.DB().
		Find(&result).
		Error

	if err != nil {
		return nil, err
	}

	return result, err
}

func (repo *repository) UserByID(userID int64) (User, error) {
	var result User

	err := repo.dbConn.DB().
		Where("id=?", userID).
		First(&result).
		Error

	if err == gorm.ErrRecordNotFound {
		return EmptyUser, ErrUserNotFound
	} else if err != nil {
		return EmptyUser, err
	}

	return result, nil
}

func (repo *repository) UserByUserName(username string) (User, error) {
	var result User

	err := repo.dbConn.DB().
		Where("user_name=?", username).
		First(&result).
		Error

	if err == gorm.ErrRecordNotFound {
		return EmptyUser, ErrUserNotFound
	} else if err != nil {
		return EmptyUser, err
	}

	return result, nil
}

func (repo *repository) UpdateUserByID(user User) (User, error) {
	entity, err := repo.UserByID(user.ID)
	if err != nil {
		return EmptyUser, err
	}

	err = repo.dbConn.DB().
		Model(&entity).
		Updates(&user).
		Error

	if err != nil {
		return EmptyUser, err
	}

	return entity, nil
}

func (repo *repository) RemoveUserByID(userID int64) (User, error) {
	entity, err := repo.UserByID(userID)
	if err != nil {
		return EmptyUser, err
	}

	err = repo.dbConn.DB().
		Delete(&entity).
		Error

	if err != nil {
		return EmptyUser, err
	}

	return entity, nil
}
