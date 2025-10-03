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

func ListNFTItemsByUser(appCtx appctx.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {
		ownerUID, err := common.FromBase58(c.Param("user-id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var filter mkpmodel.Filter

		if err := c.ShouldBind(&filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		filter.OwnerId = int(ownerUID.GetLocalID())
		filter.IsSelling = c.DefaultQuery("is_selling", "false") == "true"

		if err := paging.Process(); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := mkpstorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := mkpbiz.NewItemUserBiz(store)

		result, err := biz.ListItemsByUser(c.Request.Context(), &filter, &paging)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask(appCtx.GetAssetDomain())
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, filter))
	}
}
