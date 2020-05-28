package db_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/stretchr/testify/assert"
)

func TestPostgresConnIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg, err := config.FromEnvs()
	if err != nil {
		t.Fatal(err)
	}

	expectedHealthy := true
	expectedHealthyAfterClose := false

	postgresConn, cleanup, err := db.NewPostgresConn(cfg)
	if err != nil {
		t.Error(err)
	}

	actualHealthy := postgresConn.Healthy()
	assert.Equal(t, expectedHealthy, actualHealthy)

	cleanup()

	actualHealthyAfterClose := postgresConn.Healthy()
	assert.Equal(t, expectedHealthyAfterClose, actualHealthyAfterClose)
}
