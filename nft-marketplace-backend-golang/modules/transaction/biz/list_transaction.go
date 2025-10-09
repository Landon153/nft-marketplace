package txbiz

import (
	"context"
	"service-nft-marketplace-200lab/common"
	txmodel "service-nft-marketplace-200lab/modules/transaction/model"
)

type ListTxStore interface {
	ListDataWithCondition(
		ctx context.Context,
		filter *txmodel.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]txmodel.Transaction, error)
}

type listTxBiz struct {
	store ListTxStore
}

func NewListTxBiz(store ListTxStore) *listTxBiz {
	return &listTxBiz{store: store}
}

func (biz *listTxBiz) ListTx(
	ctx context.Context,
	filter *txmodel.Filter,
	paging *common.Paging,
) ([]txmodel.Transaction, error) {
	result, err := biz.store.ListDataWithCondition(ctx, filter, paging)

	if err != nil {
		return nil, common.ErrCannotListEntity(txmodel.EntityName, err)
	}

	return result, nil
}
