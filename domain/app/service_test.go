package app_test

import (
	"testing"

	"github.com/gghcode/apas-todo-apiserver/domain/app"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAppService_Version(t *testing.T) {
	testCases := []struct {
		description     string
		stubVersion     string
		stubErr         error
		expectedVersion string
	}{
		{
			description:     "ShouldReturnTestVersion",
			stubVersion:     "test version",
			stubErr:         nil,
			expectedVersion: "test version",
		},
		{
			description:     "ShouldReturnDefaultVersion",
			stubVersion:     "",
			stubErr:         fake.ErrStub,
			expectedVersion: "dev version",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			fakeFileReader := fake.NewFileReader()
			fakeFileReader.
				On("ReadString", mock.Anything).
				Return(tc.stubVersion, tc.stubErr)

			srv := app.NewService(fakeFileReader)

			actualVersion := srv.Version()

			assert.Equal(t, tc.expectedVersion, actualVersion)
		})
	}
}
