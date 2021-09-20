package store

import (
	"context"
	"errors"
	"github.com/dsthakur2711/wallet/model"
	"github.com/dsthakur2711/wallet/pkg/local_errors"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"time"
)


type UserRepo interface {
	CreateUser(ctx context.Context, arg CreateUserParams) (model.User, error)
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
	//UpdateUserStatus(ctx context.Context, arg UpdateUserStatusParams) (model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepo(client *gorm.DB) UserRepo {
	return &userRepository{
		db: client,
	}
}

type CreateUserParams struct {
	Id             int64 			 `gorm:"primary_key;AUTO_INCREMENT;not_null" json:"id"`
	Username       string            `json:"username"`
	HashedPassword string            `json:"hashed_password"`
	Status         model.UserStatus  `json:"status"`
	Fullname	   string			 `json:"fullname" `
	Email          string            `json:"email"`
	Address        string    	 	 `json:"address"`
	Nationality    string   		 `json:"nationality"`
	AadharNo       string    	 	 `json:"aadhar_no"`
}

func (q *userRepository) CreateUser(ctx context.Context, arg CreateUserParams) (model.User, error) {

	logrus.Println("log create user in store/user/CreateUser ")

	var u model.User

	u = model.User{
			ID: arg.Id,
			Username: arg.Username,
			HashedPassword: arg.HashedPassword,
			Status: arg.Status,
			Fullname: arg.Fullname,
			Email: arg.Email,
			Address: arg.Address,
			Nationality: arg.Nationality,
			AadharNo: arg.AadharNo,
			PasswordChangedAt: time.Now(),
	}
	res := q.db.Create(&u) // pass pointer of data to Create
	if res.Error != nil{
		return u, local_errors.ErrSomethingWrong
	}

	return u, nil
}


func (q *userRepository) GetUserByUsername(ctx context.Context, username string) (model.User, error) {

	logrus.Println("log  Login user in store/user/")

	var u model.User
	res := q.db.Where("username = ?", username).Take(&u)
	// SELECT * FROM users WHERE username = "jinzhu";

	// check error ErrRecordNotFound
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		logrus.Println("username not found !! ")
		return u, local_errors.ErrUserNotFound
	}
	if res.Error != nil {
		return u, local_errors.ErrSomethingWrong
	}
		return u, nil
}