package authstorage

import (
	"context"
	"service-nft-marketplace-200lab/common"
	authmodel "service-nft-marketplace-200lab/modules/auth/model"
)

func (s *sqlStore) CreateUser(
	ctx context.Context,
	data *authmodel.AuthDataCreation,
) error {

	db := s.db

	if err := db.Table(data.TableName()).Create(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
