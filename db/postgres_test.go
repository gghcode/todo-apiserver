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

	cfg, err := config.NewViperBuilder().
		BindEnvs("TEST").
		Build()

	postgresConn, err := db.NewPostgresConn(cfg)
	if err != nil {
		t.Error(err)
	}

	if !postgresConn.Healthy() {
		t.Error("Postgres connection must be healthy")
	}

	postgresConn.Close()

	if postgresConn.Healthy() {
		t.Error("Postgres connection must be unhealthy after close connection")
	}
}
