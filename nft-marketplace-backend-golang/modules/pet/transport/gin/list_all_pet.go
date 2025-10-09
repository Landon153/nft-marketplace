package petgin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service-nft-marketplace-200lab/common"
	"service-nft-marketplace-200lab/component/appctx"
	petbiz "service-nft-marketplace-200lab/modules/pet/biz"
	petmodel "service-nft-marketplace-200lab/modules/pet/model"
	petstorage "service-nft-marketplace-200lab/modules/pet/storage"
)

func ListAllPets(appCtx appctx.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {
		var paging common.Paging

		if err := c.ShouldBind(&paging); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		var filter petmodel.Filter

		if err := c.ShouldBind(&filter); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		//if err := paging.Process(); err != nil {
		//	panic(common.ErrInvalidRequest(err))
		//}

		paging.NextCursor = ""
		paging.Page = 1     // always page 1
		paging.Limit = 1000 // enough all pet types

		store := petstorage.NewSQLStore(appCtx.GetMainDBConnection())
		biz := petbiz.NewListAllPetBiz(store)

		result, err := biz.ListAll(c.Request.Context(), &filter, &paging)

		if err != nil {
			panic(err)
		}

		for i := range result {
			result[i].Mask(common.DbTypePet)
			result[i].Image.Fulfill(appCtx.GetAssetDomain())

			if i == len(result)-1 {
				paging.NextCursor = result[i].FakeId.String()
			}
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, paging, filter))
	}
}
