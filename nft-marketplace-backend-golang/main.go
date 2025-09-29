package main

import (
	"context"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"service-nft-marketplace-200lab/blockchain"
	"service-nft-marketplace-200lab/component/appctx"
	"service-nft-marketplace-200lab/middleware"
	authstorage "service-nft-marketplace-200lab/modules/auth/storage"
	mkpstorage "service-nft-marketplace-200lab/modules/marketplace/storage"
	txstorage "service-nft-marketplace-200lab/modules/transaction/storage"
	"strconv"
)

func main() {
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.JSONFormatter{})

	dsn := os.Getenv("DB_CONN_STR")
	assetDomain := os.Getenv("ASSET_DOMAIN")
	secretKey := os.Getenv("SECRET_KEY")

	rpcURL := os.Getenv("RPC_URL")
	contractAddress := os.Getenv("MKP_CONTRACT_ADDRESS")
	blockStart, _ := strconv.ParseInt(os.Getenv("BLOCK_START"), 10, 64)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	appCtx := appctx.NewAppContext(db, assetDomain, secretKey)

	//
	txStore := txstorage.NewSQLStore(appCtx.GetMainDBConnection())
	nftStore := mkpstorage.NewSQLStore(appCtx.GetMainDBConnection())
	userStore := authstorage.NewSQLStore(appCtx.GetMainDBConnection())

	logCrawler := blockchain.NewMarketplaceLogCrawler(rpcURL, blockStart, contractAddress, txStore)
	logHandler := blockchain.NewMkpHdl(txStore, nftStore, userStore)

	if err := logCrawler.Start(); err != nil {
		log.Errorln(err)
	}

	go logHandler.Run(context.Background(), logCrawler.GetLogsChan())

	r := gin.Default()
	r.Use(middleware.Recover(appCtx))
	r.Use(middleware.AllowCORS())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Static("/assets", "./assets")

	v1 := r.Group("/v1")

	mainRoute(v1, appCtx)

	if err := r.Run(); err != nil {
		log.Fatalln(err)
	}
}
