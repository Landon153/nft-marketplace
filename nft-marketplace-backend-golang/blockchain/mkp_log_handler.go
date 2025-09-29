package blockchain

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"math/big"
	"math/rand"
	"os"
	appCommon "service-nft-marketplace-200lab/common"
	"service-nft-marketplace-200lab/component/asyncjob"
	authmodel "service-nft-marketplace-200lab/modules/auth/model"
	txmodel "service-nft-marketplace-200lab/modules/transaction/model"
	"strconv"
	"strings"
	"time"
)

const (
	EvtOrderAdded      = "0x2a40b778ad2a0e665f93ad917831cd03068b829743c6f6a350518360e2dd97df"
	EvtOrderCancel     = "0x61b9399f2f0f32ca39ce8d7be32caed5ec22fe07a6daba3a467ed479ec606582"
	EvtOrderMatched    = "0xd92904f74973e583bd5940ae2a3a1a28fa8a41e9e9233efe6eb60fc6d0536522"
	DefaultABIFilePath = "./marketplace.abi"
)

type TxStorage interface {
	GetDataWithCondition(
		ctx context.Context,
		condition map[string]interface{},
		moreKeys ...string,
	) (*txmodel.Transaction, error)

	CreateData(
		ctx context.Context,
		tx *txmodel.Transaction,
	) error
}

type PetNFTStorage interface {
	UpdateDataWithCondition(
		ctx context.Context,
		condition map[string]interface{},
		data map[string]interface{},
	) error
}

type UserStore interface {
	CreateUser(
		ctx context.Context,
		data *authmodel.AuthDataCreation,
	) error
	FindUserWithCondition(
		ctx context.Context,
		condition map[string]interface{},
	) (*authmodel.AuthData, error)
}

type mkpHdl struct {
	txStore   TxStorage
	petStore  PetNFTStorage
	userStore UserStore
	abiObj    abi.ABI
}

func NewMkpHdl(txStore TxStorage, petStore PetNFTStorage, userStore UserStore) *mkpHdl {
	// Open & Parse file ABI
	f, err := os.Open(DefaultABIFilePath)

	mkpABI, err := abi.JSON(f)

	if err != nil {
		log.Fatalln(err)
	}

	return &mkpHdl{
		txStore:   txStore,
		petStore:  petStore,
		userStore: userStore,
		abiObj:    mkpABI,
	}
}

func (p *mkpHdl) Run(ctx context.Context, queue <-chan types.Log) {
	for l := range queue {
		functionHash := strings.ToLower(l.Topics[0].Hex())
		var job asyncjob.Job

		switch functionHash {
		case EvtOrderAdded:
			job = asyncjob.NewJob(func(ctx context.Context) error {
				return p.handleOrderAdded(ctx, l)
			})
		case EvtOrderCancel:
			job = asyncjob.NewJob(func(ctx context.Context) error {
				return p.handleOrderCancelled(ctx, l)
			})

		case EvtOrderMatched:
			job = asyncjob.NewJob(func(ctx context.Context) error {
				return p.handleOrderMatched(ctx, l)
			})
		default:
			log.Printf("Block %d - Tx %s - Event %s \n", l.BlockNumber, l.TxHash.Hex(), functionHash)
			continue
		}

		job.SetRetryDurations(time.Second, time.Second, time.Second, time.Second) // 4 times (1s each)

		if err := job.Execute(ctx); err != nil {
			log.Errorln(err)
		}
	}
}

func (p *mkpHdl) handleOrderAdded(ctx context.Context, l types.Log) error {
	log.Infof("Block %d - Tx %s - Event ORDER_ADDED \n", l.BlockNumber, l.TxHash.Hex())

	tx, err := p.txStore.GetDataWithCondition(ctx, map[string]interface{}{"tx_hash": l.TxHash.Hex()})

	if tx != nil {
		return nil
	}

	if err != appCommon.ErrRecordNotFound {
		return err
	}

	// No tx in DB

	event := struct {
		PaymentToken common.Address
		Price        *big.Int
	}{}

	err = p.abiObj.UnpackIntoInterface(&event, "OrderAdded", l.Data)

	if err != nil {
		log.Errorln(err)
		return err
	}

	orderId, _ := strconv.ParseInt(strings.ReplaceAll(l.Topics[1].Hex(), "0x", ""), 16, 64)
	sellerAddress := strings.ToLower(common.HexToAddress(l.Topics[2].Hex()).Hex())
	tokenId, _ := strconv.ParseInt(strings.ReplaceAll(l.Topics[3].Hex(), "0x", ""), 16, 64)

	paymentToken := strings.ToLower(event.PaymentToken.Hex())
	price := decimal.NewNullDecimal(decimal.NewFromBigInt(event.Price, 0))

	newTx := txmodel.Transaction{
		SQLModel:      appCommon.NewSQLModel(),
		TxHash:        l.TxHash.Hex(),
		BlockNumber:   l.BlockNumber,
		NFTId:         fmt.Sprintf("%d", tokenId),
		EventName:     "order_added",
		OrderId:       fmt.Sprintf("%d", orderId),
		SellerAddress: sellerAddress,
		BuyerAddress:  "",
		PayableToken:  paymentToken,
		Price:         price,
	}

	newTx.Status = "success"

	if err := p.txStore.CreateData(ctx, &newTx); err != nil {
		return err
	}

	// Update table nft_pets
	// 1. Find user to update owner id
	// 2. Set NFT Pets to listing status
	user, err := p.findOrCreateUser(ctx, sellerAddress)

	if err != nil {
		log.Errorln(err)
		return err
	}

	_ = p.petStore.UpdateDataWithCondition(
		ctx,
		map[string]interface{}{"nft_id": newTx.NFTId},
		map[string]interface{}{
			"owner_id":      user.Id,
			"status":        "selling",
			"listing_price": price,
			"listed_at":     time.Now().UTC(),
			"order_id":      fmt.Sprintf("%d", orderId),
			"payable_token": paymentToken,
		},
	)

	return nil
}

func (p *mkpHdl) handleOrderCancelled(ctx context.Context, l types.Log) error {
	log.Infof("Block %d - Tx %s - Event ORDER_CANCELLED \n", l.BlockNumber, l.TxHash.Hex())

	tx, err := p.txStore.GetDataWithCondition(ctx, map[string]interface{}{"tx_hash": l.TxHash.Hex()})

	if tx != nil {
		return nil
	}

	if err != appCommon.ErrRecordNotFound {
		return err
	}

	// No tx in DB
	orderId, _ := strconv.ParseInt(strings.ReplaceAll(l.Topics[1].Hex(), "0x", ""), 16, 64)

	newTx := txmodel.Transaction{
		SQLModel:    appCommon.NewSQLModel(),
		TxHash:      l.TxHash.Hex(),
		BlockNumber: l.BlockNumber,
		EventName:   "order_cancelled",
		OrderId:     fmt.Sprintf("%d", orderId),
	}

	newTx.Status = "success"

	if err := p.txStore.CreateData(ctx, &newTx); err != nil {
		return err
	}

	_ = p.petStore.UpdateDataWithCondition(
		ctx,
		map[string]interface{}{"order_id": orderId},
		map[string]interface{}{
			"status":        "activated",
			"listing_price": 0,
			"listed_at":     time.Now().UTC(),
			"order_id":      0,
		},
	)

	return nil
}

func (p *mkpHdl) handleOrderMatched(ctx context.Context, l types.Log) error {
	log.Infof("Block %d - Tx %s - Event ORDER_MATCHED \n", l.BlockNumber, l.TxHash.Hex())

	tx, err := p.txStore.GetDataWithCondition(ctx, map[string]interface{}{"tx_hash": l.TxHash.Hex()})

	if tx != nil {
		return nil
	}

	if err != appCommon.ErrRecordNotFound {
		return err
	}

	event := struct {
		TokenId      *big.Int
		PaymentToken common.Address
		Price        *big.Int
	}{}

	err = p.abiObj.UnpackIntoInterface(&event, "OrderMatched", l.Data)

	if err != nil {
		log.Errorln(err)
		return err
	}

	orderId, _ := strconv.ParseInt(strings.ReplaceAll(l.Topics[1].Hex(), "0x", ""), 16, 64)
	sellerAddress := strings.ToLower(common.HexToAddress(l.Topics[2].Hex()).Hex())
	buyerAddress := strings.ToLower(common.HexToAddress(l.Topics[3].Hex()).Hex())
	price := decimal.NewNullDecimal(decimal.NewFromBigInt(event.Price, 0))
	paymentToken := strings.ToLower(event.PaymentToken.Hex())
	tokenId, _ := strconv.ParseInt(event.TokenId.String(), 10, 64)

	newTx := txmodel.Transaction{
		SQLModel:      appCommon.NewSQLModel(),
		TxHash:        l.TxHash.Hex(),
		BlockNumber:   l.BlockNumber,
		NFTId:         fmt.Sprintf("%d", tokenId),
		EventName:     "order_matched",
		OrderId:       fmt.Sprintf("%d", orderId),
		SellerAddress: sellerAddress,
		BuyerAddress:  buyerAddress,
		PayableToken:  paymentToken,
		Price:         price,
	}

	newTx.Status = "success"

	if err := p.txStore.CreateData(ctx, &newTx); err != nil {
		return err
	}

	user, err := p.findOrCreateUser(ctx, buyerAddress)

	if err != nil {
		log.Errorln(err)
		return err
	}

	_ = p.petStore.UpdateDataWithCondition(
		ctx,
		map[string]interface{}{"nft_id": newTx.NFTId},
		map[string]interface{}{
			"owner_id":      user.Id,
			"status":        "activated",
			"listing_price": price,
			"order_id":      fmt.Sprintf("%d", orderId),
		},
	)

	return nil
}

func (p *mkpHdl) findOrCreateUser(ctx context.Context, walletAddress string) (*authmodel.AuthData, error) {
	user, err := p.userStore.FindUserWithCondition(ctx, map[string]interface{}{"wallet_address": walletAddress})

	if user != nil {
		return user, nil
	}

	if err != appCommon.ErrRecordNotFound {
		return nil, err
	}

	if err == appCommon.ErrRecordNotFound {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		nonce := r1.Intn(9999) + 10000

		newUser := &authmodel.AuthDataCreation{
			WalletAddress: walletAddress,
			Nonce:         nonce,
		}
		newUser.PrepareForCreating()
		newUser.Status = "not_verified"

		if err := p.userStore.CreateUser(ctx, newUser); err != nil {
			return nil, err
		}

		user = &authmodel.AuthData{
			Id:            newUser.Id,
			WalletAddress: walletAddress,
		}
	}

	return user, nil
}
