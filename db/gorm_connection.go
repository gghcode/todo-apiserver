package db

import (
	"github.com/jinzhu/gorm"
)

type GormConnection interface {
	DB() *gorm.DB
	Close()
}