package ordersRepositories

import (
	"encoding/json"
	"fmt"

	"github.com/Supakornn/hexagonal-go/modules/orders"
	"github.com/jmoiron/sqlx"
)

type IOrdersRepository interface {
	FindOneOrder(orderId string) (*orders.Order, error)
}

type ordersRepository struct {
	db *sqlx.DB
}

func OrdersRepository(db *sqlx.DB) IOrdersRepository {
	return &ordersRepository{db: db}
}

func (r *ordersRepository) FindOneOrder(orderId string) (*orders.Order, error) {
	query := `
	SELECT 
		to_jsonb("t")
	FROM (
		SELECT 
			"o"."id",
			"o"."user_id",
			"o"."transfer_slip",
			"o"."status",
			(
				SELECT
					array_to_json(array_agg("pt")) 
				FROM (
					SELECT
						"po"."id",
						"po"."qty",
						"po"."product"
					FROM "products_orders" "po"
					WHERE "po"."order_id" = "o"."id"
				)AS "pt"
			) AS "products",
			"o"."address",
			"o"."contact",
			(
				SELECT
					SUM(COALESCE(("po"."product"->>'price')::FLOAT * ("po"."qty")::FLOAT, 0))
				FROM "products_orders" "po"
				WHERE "po"."order_id" = "o"."id"
			)AS "total_paid",
			"o"."created_at",
			"o"."updated_at"
		FROM "orders" "o"
		WHERE "o"."id" = $1
	)AS "t";`

	orderData := &orders.Order{
		TransferSlip: &orders.TransferSlip{},
		Products:     make([]*orders.ProductsOrder, 0),
	}

	raw := make([]byte, 0)
	if err := r.db.Get(&raw, query, orderId); err != nil {
		return nil, fmt.Errorf("get order failed: %w", err)
	}

	if err := json.Unmarshal(raw, &orderData); err != nil {
		return nil, fmt.Errorf("unmarshal order failed: %w", err)
	}

	return orderData, nil
}
