package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gghcode/apas-todo-apiserver/domain/entity"
	"github.com/gghcode/apas-todo-apiserver/domain/usecase/user"
	"github.com/gghcode/apas-todo-apiserver/internal/testutil/fake"
)

func TestUserService_GetUserByUserName(t *testing.T) {
	testCases := []struct {
		description string
		argUserName string
		stubUser    entity.User
		stubErr     error
		expectedRes user.UserResponse
		expectedErr error
	}{
		{
			description: "ShouldGetUser",
			argUserName: "test name",
			stubUser: entity.User{
				ID:       10,
				UserName: "test name",
			},
			stubErr: nil,
			expectedRes: user.UserResponse{
				ID:       10,
				UserName: "test name",
			},
			expectedErr: nil,
		},
		{
			description: "ShouldReturnErrStub",
			argUserName: "NOT_EXIST_USER_NAME",
			stubUser:    entity.User{},
			stubErr:     fake.ErrStub,
			expectedRes: user.UserResponse{},
			expectedErr: fake.ErrStub,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			fakePassport := fake.NewPassport()
			fakeUserRepo := fake.NewUserRepository()
			fakeUserRepo.
				On("UserByUserName", tc.argUserName).
				Return(tc.stubUser, tc.stubErr)

			srv := user.NewService(fakeUserRepo, fakePassport)

			actualRes, actualErr := srv.GetUserByUserName(tc.argUserName)

			assert.Equal(t, tc.expectedRes, actualRes)
			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}

func TestUserService_GetUserByUserID(t *testing.T) {
	testCases := []struct {
		description string
		argUserID   int64
		stubUser    user.User
		stubErr     error
		expectedRes user.UserResponse
		expectedErr error
	}{
		{
			description: "ShouldGetUser",
			argUserID:   10,
			stubUser: user.User{
				ID:       10,
				UserName: "test name",
			},
			stubErr: nil,
			expectedRes: user.UserResponse{
				ID:       10,
				UserName: "test name",
			},
			expectedErr: nil,
		},
		{
			description: "ShouldReturnErrStub",
			argUserID:   -1,
			stubUser:    user.User{},
			stubErr:     fake.ErrStub,
			expectedRes: user.UserResponse{},
			expectedErr: fake.ErrStub,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			fakePassport := fake.NewPassport()
			fakeUserRepo := fake.NewUserRepository()
			fakeUserRepo.
				On("UserByID", tc.argUserID).
				Return(tc.stubUser, tc.stubErr)

			srv := user.NewService(fakeUserRepo, fakePassport)

			actualRes, actualErr := srv.GetUserByUserID(tc.argUserID)

			assert.Equal(t, tc.expectedRes, actualRes)
			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
	testCases := []struct {
		description      string
		argCreateUserReq user.CreateUserRequest
		stubUser         user.User
		stubErr          error
		expectedRes      user.UserResponse
		expectedErr      error
	}{
		{
			description: "ShouldCreateUser",
			argCreateUserReq: user.CreateUserRequest{
				UserName: "test",
				Password: "testtest",
			},
			stubUser: user.User{
				UserName: "test",
			},
			expectedRes: user.UserResponse{
				UserName: "test",
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			fakePassport := fake.NewPassport()
			fakePassport.
				On("HashPassword", tc.argCreateUserReq.Password).
				Return([]byte{}, nil)

			fakeUserRepo := fake.NewUserRepository()
			fakeUserRepo.
				On("CreateUser", mock.Anything).
				Return(tc.stubUser, tc.stubErr)

			srv := user.NewService(fakeUserRepo, fakePassport)

			actualRes, actualErr := srv.CreateUser(tc.argCreateUserReq)

			assert.Equal(t, tc.expectedRes, actualRes)
			assert.Equal(t, tc.expectedErr, actualErr)
		})
	}
}
