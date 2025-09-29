package blockchain

import (
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"math/big"
	appCommon "service-nft-marketplace-200lab/common"
	"service-nft-marketplace-200lab/component/asyncjob"
	"strconv"
	"strings"
	"time"
)

const (
	BlockStart            int64 = 16061292
	DefaultRPCURL               = "https://data-seed-prebsc-1-s1.binance.org:8545"
	NFTMarketplaceAddress       = "0x9F53eF4DC95FC1C5411d68f461c403F42A160e93"
)

type BlockTrackerStorage interface {
	GetLatestBlockNumber(ctx context.Context) (int64, error)
	UpdateLatestBlockNumber(ctx context.Context, num int64) error
}

type mkpLogCrawler struct {
	rpcURL              string
	currentBlock        int64
	latestBlock         int64
	contractAddress     string
	contractABIFilePath string
	client              *ethclient.Client
	storage             BlockTrackerStorage
	logChan             chan types.Log
}

func NewMarketplaceLogCrawler(
	rpcURL string,
	blockStart int64,
	contractAddress string,
	store BlockTrackerStorage,
) *mkpLogCrawler {
	mkpCrawler := &mkpLogCrawler{
		rpcURL:          stringOrDefault(rpcURL, DefaultRPCURL),
		currentBlock:    max(blockStart, BlockStart),
		contractAddress: stringOrDefault(contractAddress, NFTMarketplaceAddress),
		client:          nil,
		storage:         store,
		logChan:         make(chan types.Log, 100),
	}

	if blockStart > BlockStart {
		mkpCrawler.currentBlock = blockStart
	}

	return mkpCrawler
}

func (parser *mkpLogCrawler) Start() error {
	// Connect to blockchain node
	client, err := ethclient.Dial(parser.rpcURL)
	if err != nil {
		return err
	}

	parser.client = client

	latestBlockNumber, err := parser.latestBlockNumber()
	parser.latestBlock = latestBlockNumber

	if err != nil {
		return err
	}

	currentBlockNumber, err := parser.latestDbBlockNumber()

	if err != nil {
		return err
	}

	parser.currentBlock = max(currentBlockNumber, parser.currentBlock)

	go func() {
		var stepBlockFastScan int64 = 200
		for {
			time.Sleep(time.Second * 2)
			//log.Println("Update block number in DB here!!", parser.currentBlock)
			if err := parser.storage.UpdateLatestBlockNumber(context.Background(), parser.currentBlock); err != nil {
				log.Errorln(err)
				continue
			}

			if latestBlockNumber > parser.currentBlock {
				if v := latestBlockNumber - parser.currentBlock; v < stepBlockFastScan {
					stepBlockFastScan = v
				}

				if err := parser.scanBlock(parser.currentBlock, stepBlockFastScan); err != nil {
					log.Errorln(err)
					continue
				}

				parser.currentBlock += stepBlockFastScan + 1
				continue
			}

			latestBlockNumber, err = parser.latestBlockNumber()

			if err != nil {
				continue
			}

			parser.latestBlock = latestBlockNumber

			if err != nil {
				log.Errorln("Get latestBlockNumber:", err)
				continue
			}

			if latestBlockNumber <= parser.currentBlock {
				log.Debugln("Still no block available")
				continue
			}

			if err := parser.scanBlock(parser.currentBlock, 0); err != nil {
				log.Errorln(err)
				continue
			}

			parser.currentBlock += 1
		}
	}()

	return nil
}

func (parser *mkpLogCrawler) latestBlockNumber() (int64, error) {
	var result int64 = 1

	job := asyncjob.NewJob(func(ctx context.Context) error {
		latestBlock, err := parser.client.HeaderByNumber(context.Background(), nil)

		if err != nil {
			return err
		}

		result, err = strconv.ParseInt(latestBlock.Number.String(), 10, 64)

		if err != nil {
			return err
		}

		return nil
	})

	job.SetRetryDurations(time.Second, time.Second*2, time.Second*3)

	if err := asyncjob.NewGroup(false, job).Run(context.Background()); err != nil {
		return 0, err
	}

	return result, nil
}

func (parser *mkpLogCrawler) latestDbBlockNumber() (int64, error) {
	var result int64 = 1

	job := asyncjob.NewJob(func(ctx context.Context) error {
		rs, err := parser.storage.GetLatestBlockNumber(ctx)

		if err != nil {
			if err == appCommon.ErrRecordNotFound {
				result = 1
				return nil
			}

			return err
		}

		result = rs
		return nil
	})

	job.SetRetryDurations(time.Second, time.Second*2, time.Second*3)

	if err := asyncjob.NewGroup(false, job).Run(context.Background()); err != nil {
		return 0, err
	}

	if result < parser.currentBlock {
		result = parser.currentBlock
	}

	return result, nil
}

func (parser *mkpLogCrawler) scanBlock(from, step int64) error {
	//log.Printf("Starting scan block from %d to %d. Latest onchain: %d", from, from+step, parser.latestBlock)

	logs, err := parser.client.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: big.NewInt(from),
		ToBlock:   big.NewInt(from + step),
		Addresses: []common.Address{common.HexToAddress(parser.contractAddress)},
	})

	if err != nil {
		return err
	}

	for i, l := range logs {
		functionHash := strings.ToLower(l.Topics[0].Hex())
		log.Printf("Block %d - Tx %s - Event %s \n", l.BlockNumber, l.TxHash.Hex(), functionHash)

		parser.logChan <- logs[i]
	}

	return nil
}

func (parser *mkpLogCrawler) GetLogsChan() <-chan types.Log { return parser.logChan }

func stringOrDefault(s, d string) string {
	if s == "" {
		return d
	}

	return s
}

func max(n1, n2 int64) int64 {
	if n1 > n2 {
		return n1
	}

	return n2
}
