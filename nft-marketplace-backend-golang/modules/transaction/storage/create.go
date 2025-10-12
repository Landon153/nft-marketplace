package txstorage

import (
	"context"
	"service-nft-marketplace-200lab/common"
	txmodel "service-nft-marketplace-200lab/modules/transaction/model"
)

func (s *sqlStore) CreateData(
	ctx context.Context,
	tx *txmodel.Transaction,
) error {
	db := s.db

	if err := db.Create(tx).Error; err != nil {

		return common.ErrDB(err)
	}

	return nil
}
