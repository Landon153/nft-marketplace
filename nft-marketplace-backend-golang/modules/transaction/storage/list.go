package txstorage

import (
	"context"
	"service-nft-marketplace-200lab/common"
	txmodel "service-nft-marketplace-200lab/modules/transaction/model"
)

func (s *sqlStore) ListDataWithCondition(
	ctx context.Context,
	filter *txmodel.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]txmodel.Transaction, error) {
	db := s.db

	var result []txmodel.Transaction

	db = db.Where("event_name = ?", "order_matched")

	if filter.NFTId != "" {
		db = db.Where("nft_id = ?", filter.NFTId)
	}

	if err := db.Table(txmodel.Transaction{}.TableName()).Count(&paging.Total).Error; err != nil {
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
