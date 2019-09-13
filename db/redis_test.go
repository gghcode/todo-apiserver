package db_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gitlab.com/gyuhwan/apas-todo-apiserver/config"
	"gitlab.com/gyuhwan/apas-todo-apiserver/db"
)

type RedisIntegration struct {
	suite.Suite

	cfg config.Configuration
}

func (suite *RedisIntegration) SetupSuite() {
	cfg, err := config.NewBuilder().
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
	conn := db.NewRedisConn(suite.cfg)
	pong, err := conn.Client().Ping().Result()

	suite.Equal(pong, "PONG")
	suite.NoError(err)
}