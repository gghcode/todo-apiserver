package db_test

import (
	"testing"

	"gitlab.com/gyuhwan/apas-todo-apiserver/config"
	"gitlab.com/gyuhwan/apas-todo-apiserver/db"
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
