package mkpstorage

import (
	"context"
	"service-nft-marketplace-200lab/common"
	mkpmodel "service-nft-marketplace-200lab/modules/marketplace/model"
	"strings"
)

func (s *sqlStore) ListSellingWithCondition(
	ctx context.Context,
	filter *mkpmodel.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]mkpmodel.NFTPet, error) {
	db := s.db

	var result []mkpmodel.NFTPet

	db = db.Where("status in ('selling')")

	if filter.Element != "" {
		db = db.Where("pet_id in (select id from pets where element in (?))", strings.Split(filter.Element, ","))
	}

	if err := db.Table(mkpmodel.NFTPet{}.TableName()).Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	for i := range moreKeys {
		db = db.Preload(moreKeys[i]) // for auto preload
	}

	offset := (paging.Page - 1) * paging.Limit
	db = db.Offset(offset)

	switch filter.SortedBy {
	case "latest":
		db = db.Order("listed_at desc")
	case "oldest":
		db = db.Order("listed_at asc")
	case "cheapest":
		db = db.Order("listing_price asc")
	case "expensive":
		db = db.Order("listing_price desc")
	default:
		db = db.Order("listed_at desc")
	}

	if err := db.
		Limit(paging.Limit).
		Find(&result).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return result, nil
}

func (s *sqlStore) ListNFTWithCondition(
	ctx context.Context,
	filter *mkpmodel.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]mkpmodel.NFTPet, error) {
	db := s.db

	var result []mkpmodel.NFTPet

	db = db.Where("status not in ('deleted', 'banned')")

	if filter.Element != "" {
		db = db.Where("pet_id in (select id from pets where element in (?))", strings.Split(filter.Element, ","))
	}

	if filter.OwnerId > 0 {
		db = db.Where("owner_id = ?", filter.OwnerId)
	}

	if filter.IsSelling {
		db = db.Where("status = ?", "selling")
	}

	if err := db.Table(mkpmodel.NFTPet{}.TableName()).Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	for i := range moreKeys {
		db = db.Preload(moreKeys[i]) // for auto preload
	}

	offset := (paging.Page - 1) * paging.Limit
	db = db.Offset(offset).Order("updated_at desc")

	if err := db.
		Limit(paging.Limit).
		Find(&result).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return result, nil
}
