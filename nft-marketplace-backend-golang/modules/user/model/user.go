package usermodel

import (
	"errors"
	"service-nft-marketplace-200lab/common"
)

const EntityName = "User"

type User struct {
	common.SQLModel `json:",inline"`
	Username        string `json:"username" gorm:"column:username;"`
	Password        string `json:"-" gorm:"column:password;"`
	DisplayName     string `json:"display_name" gorm:"column:display_name;"`
	Salt            string `json:"-" gorm:"column:salt;"`
	WalletAddress   string `json:"wallet_address" gorm:"column:wallet_address;"`
	Role            string `json:"role" gorm:"column:role;"`
}

type UserUpdate struct {
	DisplayName *string `json:"display_name" gorm:"column:display_name;"`
}

func (UserUpdate) TableName() string { return User{}.TableName() }

func (u *User) GetUserId() int {
	return u.Id
}

func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) GetRole() string {
	return u.Role
}

func (User) TableName() string {
	return "users"
}

var (
	ErrUsernameOrPasswordInvalid = common.NewCustomError(
		errors.New("username or password invalid"),
		"username or password invalid",
		"ErrUsernameOrPasswordInvalid",
	)

	ErrUsernameExisted = common.NewCustomError(
		errors.New("username has already existed"),
		"username has already existed",
		"ErrUsernameExisted",
	)
)
