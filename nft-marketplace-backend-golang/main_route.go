package main

import (
	"github.com/gin-gonic/gin"
	"service-nft-marketplace-200lab/component/appctx"
	"service-nft-marketplace-200lab/middleware"
	authgin "service-nft-marketplace-200lab/modules/auth/transport/gin"
	mkpgin "service-nft-marketplace-200lab/modules/marketplace/transport/gin"
	petgin "service-nft-marketplace-200lab/modules/pet/transport/gin"
	txgin "service-nft-marketplace-200lab/modules/transaction/transport/gin"
	userstorage "service-nft-marketplace-200lab/modules/user/storage"
	usergin "service-nft-marketplace-200lab/modules/user/transport/gin"
)

func mainRoute(g *gin.RouterGroup, appCtx appctx.AppContext) {
	pets := g.Group("/pets")
	{
		pets.GET("", petgin.ListAllPets(appCtx))
	}

	marketplaces := g.Group("/marketplaces")
	{
		marketplaces.GET("", mkpgin.ListSellingNFTItem(appCtx))
		marketplaces.GET("/:id", mkpgin.GetNFTItem(appCtx))
	}

	auth := g.Group("/auth")
	{
		auth.GET("/nonce", authgin.RequestNonce(appCtx))
		auth.POST("/verify_signature", authgin.VerifySignature(appCtx))
	}

	users := g.Group("users")
	{
		users.GET("/:user-id/nfts", mkpgin.ListNFTItemsByUser(appCtx))
	}

	transactions := g.Group("transactions")
	{
		transactions.GET("", txgin.ListTxOfNFT(appCtx))
	}

	authStore := userstorage.NewSQLStore(appCtx.GetMainDBConnection())

	g.GET("/profile", middleware.RequiredAuth(appCtx, authStore), usergin.GetProfile(appCtx))
	g.PUT("/profile", middleware.RequiredAuth(appCtx, authStore), usergin.UpdateProfile(appCtx))
}
