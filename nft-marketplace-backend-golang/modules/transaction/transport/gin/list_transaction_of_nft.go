package txgin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service-nft-marketplace-200lab/common"
	"service-nft-marketplace-200lab/component/appctx"
	txbiz "service-nft-marketplace-200lab/modules/transaction/biz"
	txmodel "service-nft-marketplace-200lab/modules/transaction/model"
	txstorage "service-nft-marketplace-200lab/modules/transaction/storage"
)

func ListTxOfNFT(appCtx appctx.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {
		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var filter txmodel.Filter

		if err := c.ShouldBind(&filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		filter.NFTId = c.DefaultQuery("nft_id", "")

		if err := paging.Process(); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := txstorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := txbiz.NewListTxBiz(store)

		result, err := biz.ListTx(c.Request.Context(), &filter, &paging)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask(common.DbTypeTx)

			if i == len(result)-1 {
				paging.NextCursor = result[i].FakeId.String()
			}
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, filter))
	}
}
