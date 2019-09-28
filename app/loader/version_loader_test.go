package loader_test

import (
	"os"
	"path/filepath"
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
	currentPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	memFs := afero.NewMemMapFs()
	afero.WriteFile(memFs, filepath.Join(currentPath, "VERSION"), []byte(expectedVersion), 0644)

	versionLoader := loader.NewVersionLoader(memFs)
	assert.NotNil(t, versionLoader)

	actualVersion := versionLoader.GetVersion()
	assert.Equal(t, expectedVersion, actualVersion)
}
