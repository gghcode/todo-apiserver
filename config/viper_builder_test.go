package config_test

import (
	"os"
	"testing"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/stretchr/testify/assert"
)

func TestBuildExpectedDefaultConfig(t *testing.T) {
	expected := config.DefaultConfig()

	actual, err := config.NewViperBuilder().Build()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, actual)
}

func TestBindEnvs(t *testing.T) {
	expected := config.DefaultConfig()
	expected.Addr = ":3000"

	os.Setenv("TEST2_ADDR", expected.Addr)

	actual, err := config.NewViperBuilder().
		BindEnvs("TEST2").
		Build()

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expected, actual)
}
