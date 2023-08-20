package products

import (
	"github.com/supakornn/hexagonal-go/modules/appinfo"
	"github.com/supakornn/hexagonal-go/modules/entities"
)

type Product struct {
	Id          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Category    *appinfo.Category `json:"category"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	Price       float32           `json:"price"`
	Images      []*entities.Image `json:"images"`
}

type ProductFilter struct {
	Id     string `query:"id"`
	Search string `query:"search"`
	*entities.PaginationReq
	*entities.SortReq
}
