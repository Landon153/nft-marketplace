package authgin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service-nft-marketplace-200lab/common"
	"service-nft-marketplace-200lab/component/appctx"
	authbiz "service-nft-marketplace-200lab/modules/auth/biz"
	authstorage "service-nft-marketplace-200lab/modules/auth/storage"
)

func RequestNonce(appCtx appctx.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {
		walletAddress := c.DefaultQuery("address", "")

		store := authstorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := authbiz.NewRequestNonceBiz(store)

		result, err := biz.RequestNonce(c.Request.Context(), walletAddress)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(result))
	}
}
