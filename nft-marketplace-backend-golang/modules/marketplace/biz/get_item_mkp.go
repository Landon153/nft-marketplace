package mkpbiz

import (
	"context"
	"service-nft-marketplace-200lab/common"
	mkpmodel "service-nft-marketplace-200lab/modules/marketplace/model"
	petmodel "service-nft-marketplace-200lab/modules/pet/model"
)

type GetItemStore interface {
	GetDataWithCondition(
		ctx context.Context,
		condition map[string]interface{},
		moreKeys ...string,
	) (*mkpmodel.NFTPet, error)
}

type getItemBiz struct {
	store GetItemStore
}

func NewGetItemPetBiz(store GetItemStore) *getItemBiz {
	return &getItemBiz{store: store}
}

func (biz *getItemBiz) GetItem(
	ctx context.Context,
	id int,
) (*mkpmodel.NFTPet, error) {
	result, err := biz.store.GetDataWithCondition(ctx, map[string]interface{}{"id": id}, "User", "Pet")

	if err != nil {
		return nil, common.ErrCannotGetEntity(petmodel.EntityName, err)
	}

	return result, nil
}
