package ordersRepositories

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Supakornn/hexagonal-go/modules/orders"
	"github.com/Supakornn/hexagonal-go/modules/orders/ordersPatterns"
	"github.com/jmoiron/sqlx"
)

type IOrdersRepository interface {
	FindOneOrder(orderId string) (*orders.Order, error)
	FindOrders(req *orders.OrderFilter) ([]*orders.Order, int)
	InsertOrder(req *orders.Order) (string, error)
	UpdateOrder(req *orders.Order) error
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
				) AS "pt"
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
		Products: make([]*orders.ProductsOrder, 0),
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

func (r *ordersRepository) FindOrders(req *orders.OrderFilter) ([]*orders.Order, int) {
	builder := ordersPatterns.FindOrdersBuilder(r.db, req)
	engineer := ordersPatterns.FindOrdersEngineer(builder)

	return engineer.FindOrders(), engineer.CountOrders()
}

func (r *ordersRepository) InsertOrder(req *orders.Order) (string, error) {
	builder := ordersPatterns.InsertOrderBuilder(r.db, req)
	orderId, err := ordersPatterns.InsertOrderEngineer(builder).InsertOrder()
	if err != nil {
		return "", err
	}

	return orderId, nil
}

func (r *ordersRepository) UpdateOrder(req *orders.Order) error {
	query := `
	UPDATE "orders" SET `

	queryWhereStack := make([]string, 0)
	values := make([]any, 0)
	lastIndex := 1

	if req.Status != "" {
		values = append(values, req.Status)
		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`"status" = $%d`, lastIndex))

		lastIndex++
	}

	if req.TransferSlip != nil {
		values = append(values, req.TransferSlip)
		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`"transfer_slip" = $%d`, lastIndex))

		lastIndex++
	}

	values = append(values, req.Id)

	queryClose := fmt.Sprintf(` WHERE "id" = $%d;`, lastIndex)

	for i := range queryWhereStack {
		if i != len(queryWhereStack)-1 {
			query += queryWhereStack[i] + ","
		} else {
			query += queryWhereStack[i]
		}
	}

	query += queryClose

	if _, err := r.db.ExecContext(context.Background(), query, values...); err != nil {
		return fmt.Errorf("update order failed: %w", err)
	}

	return nil
}
