package gorm

import (
	"github.com/jinzhu/gorm"
)

// Connection is gorm connection
type Connection interface {
	Healthy() bool
	DB() *gorm.DB

	Close()
}
