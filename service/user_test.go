package service

import (
	"context"
	"github.com/dsthakur2711/wallet/dto"
	"github.com/dsthakur2711/wallet/model"
	"github.com/dsthakur2711/wallet/pkg/local_errors"
	"github.com/dsthakur2711/wallet/store/mocks"
	"github.com/dsthakur2711/wallet/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

var mockUserRepo= mocks.UserRepo{}
var usererviceImpl= NewUserService(&mockUserRepo)


var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

//func TestUserService_CreateUserFailedOnHashPassword(t *testing.T) {
//	// line 33 we need a mock
//	mockWalletRepo.On("HashPassword", mock.Anything).
//		Return("", local_errors.Err).Once()
//
//	userDto, err := usererviceImpl.CreateUser(context.Background(), dto.CreateUserDto{})
//
//	//assertion
//	assert.Error(t, err)
//	assert.Equal(t, dto.UserDto{}, userDto)
//}


func TestUserService_CreateUserFailedOnGetUserByUsername(t *testing.T) {
	// line 33 we need a mock
	password := RandomString(8)
	hashedPassword1, err := util.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	//line 50 we need a mock
	mockUserRepo.On("GetUserByUsername", mock.Anything, mock.Anything).
		Return(model.User{}, local_errors.ErrSomethingWrong).Once()

	userDto, err := usererviceImpl.CreateUser(context.Background(), dto.CreateUserDto{})

	//assertion
	assert.Equal(t, local_errors.ErrSomethingWrong, err)
	assert.Equal(t, dto.UserDto{}, userDto)
}

func TestUserService_CreateUserFailedOnErrUsernameAlreadyTaken(t *testing.T) {
	// line 33 we need a mock
	password := RandomString(8)
	hashedPassword1, err := util.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	//line 50 we need a mock
	u := model.User{
		Username: "deepak",
		HashedPassword: "password",
		Status: model.UserStatusACTIVE,
		Email: "deepak@email.com",
	}
	mockUserRepo.On("GetUserByUsername", mock.Anything, mock.Anything).
		Return(u, nil).Once()

	userDto, err := usererviceImpl.CreateUser(context.Background(), dto.CreateUserDto{})

	//assertion
	assert.Equal(t, local_errors.ErrUsernameAlreadyTaken, err)
	assert.Equal(t, dto.UserDto{}, userDto)
}

func TestUserService_CreateUserFailedOnErrSomethingWrong(t *testing.T) {
	// line 33 we need a mock
	password := RandomString(8)
	hashedPassword1, err := util.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	//line 50 we need a mock
	mockUserRepo.On("GetUserByUsername", mock.Anything, mock.Anything).
		Return(model.User{}, local_errors.ErrUserNotFound).Once()

	//line 59 we need a mock
	mockUserRepo.On("CreateUser", mock.Anything, mock.Anything).
		Return(model.User{}, local_errors.ErrSomethingWrong).Once()

	userDto, err := usererviceImpl.CreateUser(context.Background(), dto.CreateUserDto{})

	//assertion
	assert.Equal(t, local_errors.ErrSomethingWrong, err)
	assert.Equal(t, dto.UserDto{}, userDto)
}


func TestUserService_CreateUserSuccess(t *testing.T) {
	// line 33 we need a mock
	password := RandomString(8)
	hashedPassword1, err := util.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword1)

	//line 50 we need a mock
	mockUserRepo.On("GetUserByUsername", mock.Anything, mock.Anything).
		Return(model.User{}, local_errors.ErrUserNotFound).Once()

	//line 59 we need a mock
	u := model.User{
		Username: "deepak",
		HashedPassword: "password",
		Status: model.UserStatusACTIVE,
		Email: "deepak@email.com",
	}
	mockUserRepo.On("CreateUser", mock.Anything, mock.Anything).
		Return(u, nil).Once()

	userDto, err := usererviceImpl.CreateUser(context.Background(), dto.CreateUserDto{})

	//assertion
	assert.Nil(t, err)
	assert.Equal(t, dto.NewUserDto(u), userDto)
}