package entities

type PaginationReq struct {
	Page      int `query:"page"`
	Limit     int `query:"limit"`
	TotalPage int `json:"total_page"`
	TotalItem int `json:"total_item"`
}

type SortReq struct {
	OrderBy string `query:"order_by"`
	Sort    string `query:"sort"`
}
