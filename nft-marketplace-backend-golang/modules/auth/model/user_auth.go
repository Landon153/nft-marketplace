package authmodel

import (
	"errors"
	"service-nft-marketplace-200lab/common"
	"strings"
)

var (
	ErrWalletAddressInvalid       = errors.New("wallet address is invalid")
	ErrNonceMustBeGreaterThanZero = errors.New("nonce must be greater than zero")
	ErrSignatureInvalid           = errors.New("signature is invalid")
)

const EntityName = "User"

type AuthData struct {
	Id            int    `json:"-" gorm:"column:id;"`
	WalletAddress string `json:"wallet_address" gorm:"column:wallet_address;"`
	Nonce         int    `json:"nonce" gorm:"column:nonce;"`
}

type AuthVerifyData struct {
	AuthData
	Signature string `json:"signature"`
}

func (AuthData) TableName() string { return "users" }

type AuthDataCreation struct {
	common.SQLModel
	WalletAddress string `json:"wallet_address" gorm:"column:wallet_address;"`
	Nonce         int    `json:"nonce" gorm:"column:nonce;"`
}

func (AuthDataCreation) TableName() string { return "users" }

func (data *AuthDataCreation) PrepareForCreating() {
	data.SQLModel = common.NewSQLModel()
	data.WalletAddress = strings.ToLower(strings.TrimSpace(data.WalletAddress))
}

func (data *AuthVerifyData) Validate() error {
	data.WalletAddress = strings.ToLower(strings.TrimSpace(data.WalletAddress))
	data.Signature = strings.ToLower(strings.TrimSpace(data.Signature))

	if data.WalletAddress == "" { // need validate more
		return ErrWalletAddressInvalid
	}

	if data.Signature == "" {
		return ErrSignatureInvalid
	}

	if data.Nonce <= 0 {
		return ErrNonceMustBeGreaterThanZero
	}

	return nil
}

type AuthDataUpdating struct {
	Nonce       *int    `json:"-" gorm:"column:nonce;"`
	Status      *string `json:"-" gorm:"column:status;"`
	DisplayName *string `json:"-" gorm:"column:display_name;"`
}

func (AuthDataUpdating) TableName() string { return "users" }
