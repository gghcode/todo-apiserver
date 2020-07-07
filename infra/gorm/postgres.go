package gorm

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/infra/gorm/model"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	_ "github.com/lib/pq"
)

// PostgresConn godoc
type PostgresConn struct {
	db *gorm.DB
}

// NewPostgresConnection godoc
func NewPostgresConnection(cfg config.Configuration) (Connection, func(), error) {
	gormDB, err := gorm.Open("postgres",
		"host="+cfg.PostgresHost+
			" port="+cfg.PostgresPort+
			" user="+cfg.PostgresUser+
			" dbname="+cfg.PostgresName+
			" password="+cfg.PostgresPassword+
			" sslmode=disable")

	if err != nil {
		return nil, nil, errors.Wrap(err, "db connect failed...")
	}

	gormDB.AutoMigrate(
		&model.Todo{},
		&model.User{},
	)

	conn := &PostgresConn{
		db: gormDB,
	}

	cleanupFunc := func() {
		conn.Close()
	}

	return conn, cleanupFunc, nil
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
