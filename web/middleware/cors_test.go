package middleware_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/gghcode/apas-todo-apiserver/web/middleware"

	"github.com/stretchr/testify/assert"
)

func TestNewCors(t *testing.T) {
	handlerFunc := middleware.NewCors(config.Configuration{
		CorsAllowOrigins: []string{"*"},
	})

	assert.NotNil(t, handlerFunc)
}
