package mkpbiz

import (
	"context"
	"service-nft-marketplace-200lab/common"
	mkpmodel "service-nft-marketplace-200lab/modules/marketplace/model"
	petmodel "service-nft-marketplace-200lab/modules/pet/model"
)

type ListItemUserStore interface {
	ListNFTWithCondition(
		ctx context.Context,
		filter *mkpmodel.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]mkpmodel.NFTPet, error)
}

type listItemUserBiz struct {
	store ListItemUserStore
}

func NewItemUserBiz(store ListItemUserStore) *listItemUserBiz {
	return &listItemUserBiz{store: store}
}

func (biz *listItemUserBiz) ListItemsByUser(
	ctx context.Context,
	filter *mkpmodel.Filter,
	paging *common.Paging,
) ([]mkpmodel.NFTPet, error) {
	result, err := biz.store.ListNFTWithCondition(ctx, filter, paging, "User", "Pet")

	if err != nil {
		return nil, common.ErrCannotListEntity(petmodel.EntityName, err)
	}

	return result, nil
}
