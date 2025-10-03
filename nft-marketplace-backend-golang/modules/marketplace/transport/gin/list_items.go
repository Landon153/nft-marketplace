package mkpgin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service-nft-marketplace-200lab/common"
	"service-nft-marketplace-200lab/component/appctx"
	mkpbiz "service-nft-marketplace-200lab/modules/marketplace/biz"
	mkpmodel "service-nft-marketplace-200lab/modules/marketplace/model"
	mkpstorage "service-nft-marketplace-200lab/modules/marketplace/storage"
)

func ListSellingNFTItem(appCtx appctx.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {
		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var filter mkpmodel.Filter

		if err := c.ShouldBind(&filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		filter.Status = "selling" // listing selling status only

		if err := paging.Process(); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := mkpstorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := mkpbiz.NewSellingItemPetBiz(store)

		result, err := biz.ListSellingItems(c.Request.Context(), &filter, &paging)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask(appCtx.GetAssetDomain())

			//if i == len(result)-1 {
			//	paging.NextCursor = base58.Encode([]byte(fmt.Sprintf("%v",
			//		result[i].ListedAt.Format(common.DateTimeFmt))))
			//}
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, filter))
	}
}
