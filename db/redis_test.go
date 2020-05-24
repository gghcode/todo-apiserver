package db_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/db"
	"github.com/stretchr/testify/suite"
)

type RedisIntegration struct {
	suite.Suite

	cfg config.Configuration
}

func (suite *RedisIntegration) SetupSuite() {
	cfg, err := config.NewViperBuilder().
		BindEnvs("TEST").
		Build()

	suite.NoError(err)

	suite.cfg = cfg
}

func TestRedisConnIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	suite.Run(t, new(RedisIntegration))
}

func (suite *RedisIntegration) TestNewRedisConn() {
	expectedHealthy := true
	expectedHealthyAfterClose := false

	conn := db.NewRedisConn(suite.cfg)
	actualHealthy := conn.Healthy()

	suite.Equal(expectedHealthy, actualHealthy)
	suite.NoError(conn.Close())

	actualHealthyAfterClose := conn.Healthy()

	suite.Equal(expectedHealthyAfterClose, actualHealthyAfterClose)
}
