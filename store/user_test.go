package store_test

import (
	"context"
	"github.com/dsthakur2711/wallet/model"
	"github.com/dsthakur2711/wallet/pkg/local_errors"
	"github.com/dsthakur2711/wallet/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)


func InitUserRepo(t *testing.T) store.UserRepo{
	testDB 		 = DbConn()
	userRepo 	:= store.NewUserRepo(testDB)

	require.NotEmpty(t, userRepo)
	return userRepo
}

func TestUserRepository_GetUserByUsernameFailedOnErrUserNotFound(t *testing.T) {
	userRepoImpl := InitUserRepo(t)

	i, err := userRepoImpl.GetUserByUsername(context.Background(), "thisUsernameDoesntExist")
	var i1 model.User

	//assertion
	assert.Error(t, local_errors.ErrUserNotFound, err)
	assert.Equal(t, i1, i)
}

func TestUserRepository_CreateUser(t *testing.T) {

}
