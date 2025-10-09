package petstorage

import (
	"context"
	"service-nft-marketplace-200lab/common"
	petmodel "service-nft-marketplace-200lab/modules/pet/model"
)

func (s *sqlStore) ListDataWithCondition(
	ctx context.Context,
	filter *petmodel.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]petmodel.Pet, error) {
	db := s.db

	var result []petmodel.Pet

	db = db.Where("status in ('activated')")

	if err := db.Table(petmodel.Pet{}.TableName()).Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	for i := range moreKeys {
		db = db.Preload(moreKeys[i]) // for auto preload
	}

	if v := paging.FakeCursor; v != "" {
		uid, err := common.FromBase58(v)

		if err != nil {
			return nil, common.ErrDB(err)
		}

		db = db.Where("id < ?", uid.GetLocalID())
	} else {
		offset := (paging.Page - 1) * paging.Limit
		db = db.Offset(offset)
	}

	if err := db.
		Limit(paging.Limit).
		Order("id desc").
		Find(&result).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return result, nil
}
