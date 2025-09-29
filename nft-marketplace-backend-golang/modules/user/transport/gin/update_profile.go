package usergin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service-nft-marketplace-200lab/common"
	"service-nft-marketplace-200lab/component/appctx"
	userbiz "service-nft-marketplace-200lab/modules/user/biz"
	usermodel "service-nft-marketplace-200lab/modules/user/model"
	userstorage "service-nft-marketplace-200lab/modules/user/storage"
)

func UpdateProfile(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		requester := c.MustGet(common.CurrentUser).(common.Requester)

		var data usermodel.UserUpdate

		if err := c.ShouldBind(&data); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		store := userstorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := userbiz.NewUpdateUserProfileBiz(store, requester)

		if err := biz.UpdateProfile(c.Request.Context(), &data); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}
