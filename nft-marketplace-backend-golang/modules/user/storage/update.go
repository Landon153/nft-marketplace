package userstorage

import (
	"context"
	"service-nft-marketplace-200lab/common"
	usermodel "service-nft-marketplace-200lab/modules/user/model"
)

func (s *sqlStore) UpdateUser(ctx context.Context, condition map[string]interface{}, data *usermodel.UserUpdate) error {
	db := s.db.Table(data.TableName())

	if err := db.Where(condition).Updates(data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
