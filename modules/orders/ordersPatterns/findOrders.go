package ordersPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Supakornn/hexagonal-go/modules/orders"
	"github.com/jmoiron/sqlx"
)

type IFindOrdersBuilder interface {
	initQuery()
	initCountQuery()
	buildWhereSearch()
	buildWhereStatus()
	buildWhereDate()
	buildSort()
	buildPaginate()
	closeQuery()
	getQuery() string
	setQuery(query string)
	getValues() []any
	setValues(values []any)
	setLastIndex(index int)
	getDb() *sqlx.DB
	reset()
}

type findOrdersBuilder struct {
	db        *sqlx.DB
	req       *orders.OrderFilter
	query     string
	values    []any
	lastIndex int
}

func FindOrdersBuilder(db *sqlx.DB, req *orders.OrderFilter) IFindOrdersBuilder {
	return &findOrdersBuilder{
		db:     db,
		req:    req,
		values: make([]any, 0),
	}
}

func (b *findOrdersBuilder) initQuery() {
	b.query += `
	SELECT 
		array_to_json(array_agg("at"))
	FROM(
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
        WHERE 1 = 1`
}

func (b *findOrdersBuilder) initCountQuery() {
	b.query += `
        SELECT
            COUNT(*) AS "count"
        FROM "orders" "o"
        WHERE 1 = 1`
}

func (b *findOrdersBuilder) buildWhereSearch() {
	if b.req.Search != "" {
		b.values = append(b.values,
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
		)

		query := fmt.Sprintf(`
		AND (
			LOWER("o"."user_id") LIKE $%d OR
			LOWER("o"."address") LIKE $%d OR
			LOWER("o"."contact") LIKE $%d 
		)`,
			b.lastIndex+1,
			b.lastIndex+2,
			b.lastIndex+3,
		)

		temp := b.getQuery()
		temp += query
		b.setQuery(temp)
		b.lastIndex = len(b.values)
	}
}

func (b *findOrdersBuilder) buildWhereStatus() {
	if b.req.Status != "" {
		b.values = append(b.values,
			strings.ToLower(b.req.Search),
		)

		query := fmt.Sprintf(`
		AND "o"."status" = $%d`,
			b.lastIndex+1,
		)

		temp := b.getQuery()
		temp += query
		b.setQuery(temp)
		b.lastIndex = len(b.values)
	}
}

func (b *findOrdersBuilder) buildWhereDate() {
	if b.req.StartDate != "" && b.req.EndDate != "" {
		b.values = append(b.values,
			b.req.StartDate,
			b.req.EndDate,
		)

		query := fmt.Sprintf(`
		AND "o"."created_at" BETWEEN DATE($%d) AND ($%d)::DATE + 1`,
			b.lastIndex+1,
			b.lastIndex+2,
		)

		temp := b.getQuery()
		temp += query
		b.setQuery(temp)
		b.lastIndex = len(b.values)
	}
}

func (b *findOrdersBuilder) buildSort() {
	b.values = append(b.values, b.req.OrderBy)

	b.query += fmt.Sprintf(`
    ORDER BY $%d %s`,
		b.lastIndex+1,
		b.req.Sort,
	)
	b.lastIndex = len(b.values)
}

func (b *findOrdersBuilder) buildPaginate() {
	b.values = append(b.values,
		(b.req.Page-1)*b.req.Limit,
		b.req.Limit,
	)

	b.query += fmt.Sprintf(`
    OFFSET $%d LIMIT $%d`,
		b.lastIndex+1,
		b.lastIndex+2,
	)
	b.lastIndex = len(b.values)
}

func (b *findOrdersBuilder) closeQuery() {
	b.query += `
	) AS "at"
	`
}

func (b *findOrdersBuilder) getQuery() string {
	return b.query
}

func (b *findOrdersBuilder) setQuery(query string) {
	b.query = query
}

func (b *findOrdersBuilder) getValues() []any {
	return b.values
}

func (b *findOrdersBuilder) setValues(values []any) {
	b.values = values
}

func (b *findOrdersBuilder) setLastIndex(index int) {
	b.lastIndex = index
}

func (b *findOrdersBuilder) getDb() *sqlx.DB {
	return b.db
}

func (b *findOrdersBuilder) reset() {
	b.query = ""
	b.values = make([]any, 0)
	b.lastIndex = 0
}

type findOrdersEngineer struct {
	builder IFindOrdersBuilder
}

func FindOrdersEngineer(builder IFindOrdersBuilder) *findOrdersEngineer {
	return &findOrdersEngineer{builder: builder}
}

func (en *findOrdersEngineer) FindOrders() []*orders.Order {
	_, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	en.builder.initQuery()
	en.builder.buildWhereSearch()
	en.builder.buildWhereStatus()
	en.builder.buildWhereDate()
	en.builder.buildSort()
	en.builder.buildPaginate()
	en.builder.closeQuery()

	raw := make([]byte, 0)
	if err := en.builder.getDb().Get(&raw, en.builder.getQuery(), en.builder.getValues()...); err != nil {
		log.Printf("get orders error: %v", err)
		return make([]*orders.Order, 0)
	}

	ordersData := make([]*orders.Order, 0)
	if err := json.Unmarshal(raw, &ordersData); err != nil {
		log.Printf("unmarshal orders error: %v", err)
	}

	en.builder.reset()
	return ordersData
}

func (en *findOrdersEngineer) CountOrders() int {
	_, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	en.builder.initCountQuery()
	en.builder.buildWhereSearch()
	en.builder.buildWhereStatus()
	en.builder.buildWhereDate()

	var count int
	if err := en.builder.getDb().Get(&count, en.builder.getQuery(), en.builder.getValues()...); err != nil {
		log.Printf("get count orders error: %v", err)
		return 0
	}

	en.builder.reset()
	return count
}
