package db_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/db"

	"github.com/stretchr/testify/assert"
)

func TestRedisConnIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	cfg, err := config.FromEnvs()
	if err != nil {
		t.Fatal(err)
	}

	expectedHealthy := true
	expectedHealthyAfterClose := false
	expectedPong := "ping: PONG"

	redisConn, cleanup := db.NewRedisConn(cfg)

	actualHealthy := redisConn.Healthy()
	assert.Equal(t, expectedHealthy, actualHealthy)

	actualPong := redisConn.Client().Ping().String()
	assert.Equal(t, expectedPong, actualPong)

	cleanup()

	actualHealthyAfterClose := redisConn.Healthy()
	assert.Equal(t, expectedHealthyAfterClose, actualHealthyAfterClose)
}
