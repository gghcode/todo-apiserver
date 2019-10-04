package middleware_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/app/middleware"
	"github.com/gghcode/apas-todo-apiserver/config"
)

func TestNewCors(t *testing.T) {
	handlerFunc := middleware.NewCors(config.CorsConfig{
		AllowOrigins: []string{"*"},
	})

	if handlerFunc == nil {
		t.Errorf("handlerFunc is nil")
	}
}
