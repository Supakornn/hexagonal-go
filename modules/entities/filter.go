package entities

type PaginationReq struct {
	Page      int `query:"page"`
	Limit     int `query:"limit"`
	TotalPage int `query:"total_page"`
	TotalItem int `query:"total_item"`
}

type SortReq struct {
	OrderBy string `query:"order_by"`
	Sort    string `query:"sort"`
}
