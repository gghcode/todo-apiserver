package file_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/infrastructure/file"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestAferoReader_ReadString(t *testing.T) {
	testCases := []struct {
		description string
		argFilePath string
		stubPath    string
		stubFile    string
		expectedStr string
		expectedErr error
	}{
		{
			description: "ShouldReturnReadString",
			argFilePath: "/app/VERSION",
			stubPath:    "/app/VERSION",
			stubFile:    "test version",
			expectedStr: "test version",
			expectedErr: nil,
		},
		{
			description: "ShouldReturnErrNotFound",
			argFilePath: "/NOT_EXIST_PATH",
			stubPath:    "/app/VERSION",
			stubFile:    "test version",
			expectedStr: "",
			expectedErr: GetPathError("/NOT_EXIST_PATH"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			memFs := afero.NewMemMapFs()
			aferoReader := file.NewAferoFileReader(memFs)
			afero.WriteFile(memFs,
				tc.stubPath,
				[]byte(tc.stubFile),
				0644,
			)

			actualStr, actualErr := aferoReader.ReadString(tc.argFilePath)

			assert.Equal(t, tc.expectedStr, actualStr)
			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}

func GetPathError(path string) error {
	_, err := afero.NewMemMapFs().Open(path)
	return err
}

func TestGetVersionWhenNotExistsVersionFile(t *testing.T) {

	// expectedVersion := "dev version"

	// memFs := afero.NewMemMapFs()

	// versionLoader := loader.NewVersionLoader(memFs)
	// assert.NotNil(t, versionLoader)

	// actualVersion := versionLoader.GetVersion()
	// assert.Equal(t, expectedVersion, actualVersion)
}

// func TestGetVersionWhenExistsVersionFile(t *testing.T) {
// 	expectedVersion := "test version"
// 	currentPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

// 	memFs := afero.NewMemMapFs()
// 	afero.WriteFile(memFs, filepath.Join(currentPath, "VERSION"), []byte(expectedVersion), 0644)

// 	reader := file.NewAferoFileReader(memFs)

// 	actual := reader.ReadString(currentPath)
// 	versionLoader := loader.NewVersionLoader(memFs)
// 	assert.NotNil(t, versionLoader)

// 	actualVersion := versionLoader.GetVersion()
// 	assert.Equal(t, expectedVersion, actualVersion)
// }
