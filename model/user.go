package model

import (
	"time"
)

type UserStatus string

const (
	UserStatusACTIVE  UserStatus = "ACTIVE"
	UserStatusBLOCKED UserStatus = "BLOCKED"
)
type User struct {

	ID                int64      `gorm:"primary_key;AUTO_INCREMENT;not_null" json:"id"`
	Username          string     `json:"username;unique"`
	HashedPassword    string     `json:"hashed_password"`
	Status            UserStatus `json:"status"`
	Email             string     `json:"email"`
	Address 		  string     `json:"address"`
	Nationality		  string	 `json:"nationality"`
	AadharNo		  string 	 `json:"aadhar_no"`
	PasswordChangedAt time.Time  `json:"password_changed_at"`
	CreatedAt         time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}



