package config_test

import (
	"os"
	"testing"

	"github.com/gghcode/apas-todo-apiserver/config"
)

func TestFromEnvs(t *testing.T) {
	testKey := "JWT_SECRET_KEY"
	expectedTestValue := "expected-test-value"

	os.Setenv(testKey, expectedTestValue)
	defer func() {
		os.Unsetenv(testKey)
	}()

	cfg, err := config.FromEnvs()
	if err != nil {
		t.Error(err)
	}

	if cfg.JwtSecretKey != expectedTestValue {
		t.Error("Port is not equal")
	}
}
