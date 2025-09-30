package mkpbiz

import (
	"context"
	"service-nft-marketplace-200lab/common"
	mkpmodel "service-nft-marketplace-200lab/modules/marketplace/model"
	petmodel "service-nft-marketplace-200lab/modules/pet/model"
)

type ListItemStore interface {
	ListSellingWithCondition(
		ctx context.Context,
		filter *mkpmodel.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]mkpmodel.NFTPet, error)
}

type listItemBiz struct {
	store ListItemStore
}

func NewSellingItemPetBiz(store ListItemStore) *listItemBiz {
	return &listItemBiz{store: store}
}

func (biz *listItemBiz) ListSellingItems(
	ctx context.Context,
	filter *mkpmodel.Filter,
	paging *common.Paging,
) ([]mkpmodel.NFTPet, error) {
	result, err := biz.store.ListSellingWithCondition(ctx, filter, paging, "User", "Pet")

	if err != nil {
		return nil, common.ErrCannotListEntity(petmodel.EntityName, err)
	}

	return result, nil
}
