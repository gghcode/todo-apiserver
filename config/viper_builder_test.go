package config_test

import (
	"os"
	"testing"

	"github.com/gghcode/apas-todo-apiserver/config"
	"github.com/stretchr/testify/suite"
)

type ViperBuilderUnit struct {
	suite.Suite
}

func TestConfigViperBuilderUnit(t *testing.T) {
	suite.Run(t, new(ViperBuilderUnit))
}

func (suite *ViperBuilderUnit) TestBuildReturnDefaultConfig() {
	expected := config.DefaultConfig()

	actual, err := config.NewViperBuilder().Build()

	suite.NoError(err)
	suite.Equal(expected, actual)
}

func (suite *ViperBuilderUnit) TestAddConfigFileReturnPathErrorWhenOptionalFalse() {
	_, actualErr := config.NewViperBuilder().
		AddConfigFile("NOT_EXISTS_CONFIG_FILE.json", false).
		Build()

	if _, ok := actualErr.(*os.PathError); !ok {
		suite.T().Fatal(actualErr)
	}
}

func (suite *ViperBuilderUnit) TestAddConfigFileReturnNilWhenOptionalTrue() {
	_, actualErr := config.NewViperBuilder().
		AddConfigFile("NOT_EXISTS_CONFIG_FILE.json", true).
		Build()

	suite.NoError(actualErr)
}

func (suite *ViperBuilderUnit) TestBindEnvs() {
	expected := config.DefaultConfig()
	expected.Addr = ":3000"

	os.Setenv("TEST2_ADDR", expected.Addr)

	actual, err := config.NewViperBuilder().
		BindEnvs("TEST2").
		Build()

	suite.NoError(err)
	suite.Equal(expected, actual)
}
