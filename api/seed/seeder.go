package seed

import (
	"github.com/dsthakur2711/wallet/model"
	"github.com/jinzhu/gorm"
)

var users = []model.User{
	model.User{
		Username: "deepak",
		Email:    "deepak@gmail.com",
		HashedPassword: "password",
	},
	model.User{
		Username: "dk",
		Email:    "dsthakur@gmail.com",
		HashedPassword:  "password",
	},
}



func Load(db *gorm.DB) {

}