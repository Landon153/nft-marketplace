package mkpstorage

import (
	"context"
	"gorm.io/gorm"
	"service-nft-marketplace-200lab/common"
	mkpmodel "service-nft-marketplace-200lab/modules/marketplace/model"
)

func (s *sqlStore) GetDataWithCondition(
	ctx context.Context,
	condition map[string]interface{},
	moreKeys ...string,
) (*mkpmodel.NFTPet, error) {
	db := s.db

	var result mkpmodel.NFTPet

	for i := range moreKeys {
		db = db.Preload(moreKeys[i]) // for auto preload
	}

	if err := db.
		Where(condition).
		First(&result).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, common.ErrRecordNotFound
		}
		return nil, common.ErrDB(err)
	}

	return &result, nil
}
