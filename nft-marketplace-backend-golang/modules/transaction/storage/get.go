package txstorage

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"service-nft-marketplace-200lab/common"
	txmodel "service-nft-marketplace-200lab/modules/transaction/model"
	"strconv"
)

func (s *sqlStore) GetDataWithCondition(
	ctx context.Context,
	condition map[string]interface{},
	moreKeys ...string,
) (*txmodel.Transaction, error) {
	db := s.db

	var result txmodel.Transaction

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

func (s *sqlStore) GetLatestBlockNumber(ctx context.Context) (int64, error) {
	db := s.db

	rs := struct {
		KeyName string `gorm:"column:key_name;"`
		Value   string `gorm:"column:value;"`
	}{}

	if err := db.Table("configs").
		Where("key_name = ?", "BLOCK_SCANNED_NUMBER").
		First(&rs).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return 1, common.ErrRecordNotFound // it means no record in DB, return 1
		}
		return 0, common.ErrDB(err) // other errors need to handle
	}

	n, _ := strconv.ParseInt(rs.Value, 10, 64)

	return n, nil
}

func (s *sqlStore) UpdateLatestBlockNumber(ctx context.Context, num int64) error {
	db := s.db

	if err := db.Table("configs").
		Where("key_name = ?", "BLOCK_SCANNED_NUMBER").
		Updates(map[string]interface{}{"value": fmt.Sprintf("%d", num)}).Error; err != nil {

		return common.ErrDB(err)
	}

	return nil
}
