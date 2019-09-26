package db

import (
	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	_ "github.com/lib/pq"
)

// PostgresProviderFunc godoc
type PostgresProviderFunc func(*gorm.DB) *gorm.DB

// PostgresConn godoc
type PostgresConn struct {
	db        *gorm.DB
	provideDB PostgresProviderFunc
}

// NewPostgresConn godoc
func NewPostgresConn(cfg config.Configuration) (*PostgresConn, error) {
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

	return &PostgresConn{
		db: gormDB,
		provideDB: func(db *gorm.DB) *gorm.DB {
			return db
		},
	}, nil
}

// SetDbProvider godoc
func (conn *PostgresConn) SetDbProvider(provider PostgresProviderFunc) {
	conn.provideDB = provider
}

// DB return database connection.
func (conn *PostgresConn) DB() *gorm.DB {
	return conn.provideDB(conn.db)
}

// Close close db session.
func (conn *PostgresConn) Close() {
	conn.DB().Close()
}
