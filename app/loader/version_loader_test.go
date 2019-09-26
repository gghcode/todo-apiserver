package loader_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/app/loader"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGetVersionWhenNotExistsVersionFile(t *testing.T) {
	expectedVersion := "dev version"

	memFs := afero.NewMemMapFs()

	versionLoader := loader.NewVersionLoader(memFs)
	assert.NotNil(t, versionLoader)

	actualVersion := versionLoader.GetVersion()
	assert.Equal(t, expectedVersion, actualVersion)
}

func TestGetVersionWhenExistsVersionFile(t *testing.T) {
	expectedVersion := "test version"

	memFs := afero.NewMemMapFs()
	afero.WriteFile(memFs, "VERSION", []byte(expectedVersion), 0644)

	versionLoader := loader.NewVersionLoader(memFs)
	assert.NotNil(t, versionLoader)

	actualVersion := versionLoader.GetVersion()
	assert.Equal(t, expectedVersion, actualVersion)
}
