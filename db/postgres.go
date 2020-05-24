package db

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/infrastructure/model"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	_ "github.com/lib/pq"
)

// PostgresConn godoc
type PostgresConn struct {
	db *gorm.DB
}

// NewPostgresConn godoc
func NewPostgresConn(cfg config.Configuration) (GormConnection, error) {
	gormDB, err := gorm.Open(cfg.Postgres.Driver,
		"host="+cfg.Postgres.Host+
			" port="+cfg.Postgres.Port+
			" user="+cfg.Postgres.User+
			" dbname="+cfg.Postgres.Name+
			" password="+cfg.Postgres.Password+
			" sslmode=disable")

	if err != nil {
		return nil, errors.Wrap(err, "db connect failed...")
	}

	gormDB.AutoMigrate(
		&model.Todo{},
		&model.User{},
	)

	return &PostgresConn{
		db: gormDB,
	}, nil
}

// Healthy return database connection status if connection connected, method return true
func (conn *PostgresConn) Healthy() bool {
	return conn.db.DB().Ping() == nil
}

// DB return database connection.
func (conn *PostgresConn) DB() *gorm.DB {
	return conn.db
}

// Close close db session.
func (conn *PostgresConn) Close() {
	conn.DB().Close()
}
