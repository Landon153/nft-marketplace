package usergin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service-nft-marketplace-200lab/common"
	"service-nft-marketplace-200lab/component/appctx"
)

func GetProfile(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := c.MustGet(common.CurrentUser)
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(u))
	}
}
