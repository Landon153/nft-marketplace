package authgin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service-nft-marketplace-200lab/common"
	"service-nft-marketplace-200lab/component/appctx"
	"service-nft-marketplace-200lab/component/tokenprovider/jwt"
	authbiz "service-nft-marketplace-200lab/modules/auth/biz"
	authmodel "service-nft-marketplace-200lab/modules/auth/model"
	authstorage "service-nft-marketplace-200lab/modules/auth/storage"
)

func VerifySignature(appCtx appctx.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {
		var data authmodel.AuthVerifyData

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := authstorage.NewSQLStore(appCtx.GetMainDBConnection())
		tokenProvider := jwt.NewTokenJWTProvider(appCtx.SecretKey())

		biz := authbiz.NewVerifySignatureBiz(store, tokenProvider)

		accessToken, err := biz.VerifySignature(c.Request.Context(), &data)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(accessToken))
	}
}
