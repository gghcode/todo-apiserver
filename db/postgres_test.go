package db_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/db"
)

func TestPostgresConnIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg, err := config.NewBuilder().
		BindEnvs("TEST").
		Build()

	postgresConn, err := db.NewPostgresConn(cfg)
	if err != nil {
		t.Error(err)
	}

	if err := postgresConn.DB().DB().Ping(); err != nil {
		t.Error(err)
	}
}
