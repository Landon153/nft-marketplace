package petmodel

import "service-nft-marketplace-200lab/common"

const EntityName = "Pet"

type Pet struct {
	common.SQLModel
	Element string        `json:"element" gorm:"column:element;"`
	Image   *common.Image `json:"image" gorm:"column:image;"`
}

func (Pet) TableName() string { return "pets" }
