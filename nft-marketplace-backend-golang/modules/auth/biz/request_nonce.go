package authbiz

import (
	"context"
	"math/rand"
	"service-nft-marketplace-200lab/common"
	authmodel "service-nft-marketplace-200lab/modules/auth/model"
	"strings"
	"time"
)

type RequestNonceStore interface {
	FindUserWithCondition(
		ctx context.Context,
		condition map[string]interface{},
	) (*authmodel.AuthData, error)
	CreateUser(ctx context.Context, data *authmodel.AuthDataCreation) error
}

type requestNonceItemBiz struct {
	store RequestNonceStore
}

func NewRequestNonceBiz(store RequestNonceStore) *requestNonceItemBiz {
	return &requestNonceItemBiz{store: store}
}

func (biz *requestNonceItemBiz) RequestNonce(
	ctx context.Context,
	walletAddress string,
) (*authmodel.AuthData, error) {
	walletAddress = strings.ToLower(strings.TrimSpace(walletAddress))

	if walletAddress == "" {
		return nil, authmodel.ErrWalletAddressInvalid
	}

	result, err := biz.store.FindUserWithCondition(ctx, map[string]interface{}{"wallet_address": walletAddress})

	if err == common.ErrRecordNotFound {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		nonce := r1.Intn(9999) + 10000

		newData := authmodel.AuthDataCreation{
			WalletAddress: walletAddress,
			Nonce:         nonce,
		}

		newData.PrepareForCreating()

		newData.Status = "not_verified"

		if err := biz.store.CreateUser(ctx, &newData); err != nil {
			return nil, common.ErrCannotCreateEntity(authmodel.EntityName, err)
		}

		result = &authmodel.AuthData{
			WalletAddress: walletAddress,
			Nonce:         nonce,
		}

		return result, nil
	}

	if err != nil {
		return nil, common.ErrCannotGetEntity(authmodel.EntityName, err)
	}

	return result, nil
}
