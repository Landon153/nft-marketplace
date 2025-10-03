package mkpgin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service-nft-marketplace-200lab/common"
	"service-nft-marketplace-200lab/component/appctx"
	mkpbiz "service-nft-marketplace-200lab/modules/marketplace/biz"
	mkpstorage "service-nft-marketplace-200lab/modules/marketplace/storage"
)

func GetNFTItem(appCtx appctx.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {
		uid, err := common.FromBase58(c.Param("id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := mkpstorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := mkpbiz.NewGetItemPetBiz(store)

		result, err := biz.GetItem(c.Request.Context(), int(uid.GetLocalID()))

		if err != nil {
			panic(err)
		}

		result.Mask(appCtx.GetAssetDomain())

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
