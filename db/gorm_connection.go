package db

import (
	"github.com/jinzhu/gorm"
)

// GormConnection connection
type GormConnection interface {
	DB() *gorm.DB
	Close()
}
