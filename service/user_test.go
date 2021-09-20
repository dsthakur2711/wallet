package service

import (
	"context"
	"github.com/dsthakur2711/wallet/dto"
	"github.com/dsthakur2711/wallet/model"
	"github.com/dsthakur2711/wallet/pkg/local_errors"
	"github.com/dsthakur2711/wallet/store/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var mockUserRepo= mocks.UserRepo{}
var usererviceImpl= NewUserService(&mockUserRepo)



func TestUserService_CreateUserFailedOnGetUserByUsername(t *testing.T) {

	//line 50 we need a mock
	mockUserRepo.On("GetUserByUsername", mock.Anything, mock.Anything).
		Return(model.User{}, local_errors.ErrSomethingWrong).Once()

	userDto, err := usererviceImpl.CreateUser(context.Background(), dto.CreateUserDto{})

	//assertion
	assert.Equal(t, local_errors.ErrSomethingWrong, err)
	assert.Equal(t, dto.UserDto{}, userDto)
}

func TestUserService_CreateUserFailedOnErrUsernameAlreadyTaken(t *testing.T) {

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

func TestUserService_LoginUserFailedOnGetUserByUsername(t *testing.T) {

	//line 74 we need a mock
	mockUserRepo.On("GetUserByUsername", mock.Anything, mock.Anything).
		Return(model.User{}, local_errors.Err).Once()

	loggedInDto, err := usererviceImpl.LoginUser(context.Background(), dto.LoginCredentialsDto{})

	//assertion
	assert.NotNil(t, err)
	assert.Equal(t, dto.LoggedInUserDto{}, loggedInDto)
}