package db

import (
	"github.com/jinzhu/gorm"
)

// GormConnection connection
type GormConnection interface {
	Healthy() bool
	DB() *gorm.DB

	Close()
}
