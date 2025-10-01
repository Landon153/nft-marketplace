package mkpmodel

import (
	"github.com/shopspring/decimal"
	"service-nft-marketplace-200lab/common"
	"time"
)

const EntityName = "Marketplace NFT item"

type NFTPet struct {
	common.SQLModel
	PetId        int                 `json:"-" gorm:"column:pet_id;"`
	NFTId        string              `json:"nft_id" gorm:"column:nft_id;"`
	OwnerId      int                 `json:"-" gorm:"column:owner_id;"`
	Name         string              `json:"name" gorm:"column:name;"`
	Gender       string              `json:"gender" gorm:"column:gender;"`
	StatsAttack  int                 `json:"stats_attack" gorm:"column:stats_attack"`
	StatsDef     int                 `json:"stats_def" gorm:"column:stats_def"`
	StatsHp      int                 `json:"stats_hp" gorm:"column:stats_hp"`
	StatsSpeed   int                 `json:"stats_speed" gorm:"column:stats_speed"`
	ListingPrice decimal.NullDecimal `json:"listing_price" gorm:"column:listing_price;"`
	PayableToken string              `json:"payable_token" gorm:"column:payable_token;"`
	OrderId      string              `json:"order_id" gorm:"column:order_id;"`
	ListedAt     *time.Time          `json:"listed_at" gorm:"column:listed_at;"`
	Pet          *Pet                `json:"pet" gorm:"PRELOAD:false;foreignKey:PetId;"`
	User         *User               `json:"user" gorm:"PRELOAD:false;foreignKey:OwnerId;"`
}

func (NFTPet) TableName() string { return "nft_pets" }

func (p *NFTPet) Mask(assetDomain string) {
	p.SQLModel.Mask(common.DbTypeNFTPet)

	if v := p.Pet; v != nil {
		v.Mask(assetDomain)
	}

	if u := p.User; u != nil {
		u.Mask(common.DbTypeUser)
	}
}

type Pet struct {
	common.SQLModel
	Element string        `json:"element" gorm:"column:element;"`
	Image   *common.Image `json:"image" gorm:"column:image;"`
}

func (Pet) TableName() string { return "pets" }

func (p *Pet) Mask(assetDomain string) {
	p.SQLModel.Mask(common.DbTypePet)
	p.Image.Fulfill(assetDomain)
}

type User struct {
	common.SQLModel
	DisplayName   string `json:"display_name" gorm:"column:display_name;"`
	WalletAddress string `json:"wallet_address" gorm:"column:wallet_address;"`
}

func (User) TableName() string { return "users" }
