package txmodel

import (
	"github.com/shopspring/decimal"
	"service-nft-marketplace-200lab/common"
)

const EntityName = "Transaction"

type Transaction struct {
	common.SQLModel
	TxHash        string              `json:"tx_hash" gorm:"column:tx_hash;"`
	BlockNumber   uint64              `json:"block_number" gorm:"column:block_number;"`
	NFTId         string              `json:"nft_id" gorm:"column:nft_id;"`
	EventName     string              `json:"event_name" gorm:"column:event_name;"`
	OrderId       string              `json:"order_id" gorm:"column:order_id;"`
	SellerAddress string              `json:"seller_address" gorm:"column:seller_address;"`
	BuyerAddress  string              `json:"buyer_address" gorm:"column:buyer_address;"`
	PayableToken  string              `json:"payable_token" gorm:"column:payable_token;"`
	Price         decimal.NullDecimal `json:"price" gorm:"column:price;"`
}

func (Transaction) TableName() string { return "transactions" }
