package orders

import "github.com/Supakornn/hexagonal-go/modules/products"

type Order struct {
	Id           string           `db:"id" json:"id"`
	UserId       string           `db:"user_id" json:"user_id"`
	TransferSlip *TransferSlip    `json:"transfer_slip"`
	Products     []*ProductsOrder `json:"products"`
	Address      string           `db:"address" json:"address"`
	Contact      string           `db:"contact" json:"contact"`
	Status       string           `db:"status" json:"status"`
	TotalPaid    float64          `db:"total_paid" json:"total_paid"`
	CreatedAt    string           `db:"created_at" json:"created_at"`
	UpdatedAt    string           `db:"updated_at" json:"updated_at"`
}

type TransferSlip struct {
	Id        string `json:"id"`
	Filename  string `json:"filename"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

type ProductsOrder struct {
	Id      string            `db:"id" json:"id"`
	Qty     int               `db:"qty" json:"qty"`
	Product *products.Product `db:"product" json:"product"`
}
