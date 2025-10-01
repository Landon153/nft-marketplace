package mkpmodel

type Filter struct {
	Status    string `json:"status"` // do not allow change status from frontend
	SortedBy  string `json:"sorted_by" form:"sorted_by"`
	Element   string `json:"element" form:"element"`
	IsSelling bool   `json:"is_selling" form:"is_selling"` // for filter NFTs of user
	OwnerId   int    `json:"-"`
}
