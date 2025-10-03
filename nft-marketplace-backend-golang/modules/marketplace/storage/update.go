package mkpstorage

import (
	"context"
	"service-nft-marketplace-200lab/common"
	mkpmodel "service-nft-marketplace-200lab/modules/marketplace/model"
)

func (s *sqlStore) UpdateDataWithCondition(
	ctx context.Context,
	condition map[string]interface{},
	data map[string]interface{},
) error {
	db := s.db

	if err := db.Table(mkpmodel.NFTPet{}.TableName()).
		Where(condition).
		Updates(&data).Error; err != nil {
		return common.ErrDB(err)
	}

	return nil
}
