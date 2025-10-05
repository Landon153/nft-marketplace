package petbiz

import (
	"context"
	"service-nft-marketplace-200lab/common"
	petmodel "service-nft-marketplace-200lab/modules/pet/model"
)

type ListPetStore interface {
	ListDataWithCondition(
		ctx context.Context,
		filter *petmodel.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]petmodel.Pet, error)
}

type listAllPetBiz struct {
	store ListPetStore
}

func NewListAllPetBiz(store ListPetStore) *listAllPetBiz {
	return &listAllPetBiz{store: store}
}

func (biz *listAllPetBiz) ListAll(
	ctx context.Context,
	filter *petmodel.Filter,
	paging *common.Paging,
) ([]petmodel.Pet, error) {
	result, err := biz.store.ListDataWithCondition(ctx, filter, paging)

	if err != nil {
		return nil, common.ErrCannotListEntity(petmodel.EntityName, err)
	}

	return result, nil
}
