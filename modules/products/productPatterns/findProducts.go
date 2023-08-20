package productPatterns

import (
	"github.com/jmoiron/sqlx"
	"github.com/supakornn/hexagonal-go/modules/products"
)

type IFindProductBuilder interface {
	openJsonQuery()
	initQuery()
	countQuery()
	whereQuery()
	sort()
	paginate()
	closeJsonQuery()
	resetQuery()
	Result() []*products.Product
	Count() int
	PrintQuery()
}

type findProductBuilder struct {
	db             *sqlx.DB
	req            *products.ProductFilter
	query          string
	lastStackIndex int
	values         []any
}

func FindProductBuilder(db *sqlx.DB, req *products.ProductFilter) IFindProductBuilder {
	return &findProductBuilder{
		db:  db,
		req: req,
	}
}

func (b *findProductBuilder) openJsonQuery() {
	b.query += `
	SELECT
		array_to_json(array_agg("t"))
	FROM (`
}

func (b *findProductBuilder) initQuery() {}

func (b *findProductBuilder) countQuery() {}

func (b *findProductBuilder) whereQuery() {}

func (b *findProductBuilder) sort() {}

func (b *findProductBuilder) paginate() {}

func (b *findProductBuilder) closeJsonQuery() {}

func (b *findProductBuilder) resetQuery() {}

func (b *findProductBuilder) Result() []*products.Product { return nil }

func (b *findProductBuilder) Count() int { return 0 }

func (b *findProductBuilder) PrintQuery() {}

type findProductEngineer struct {
	builder IFindProductBuilder
}

func FindProductEngineer(builder IFindProductBuilder) *findProductEngineer {
	return &findProductEngineer{
		builder: builder,
	}
}

func (en *findProductEngineer) FindProduct() IFindProductBuilder {
	return en.builder
}

func (en *findProductEngineer) CountProduct() IFindProductBuilder {
	return en.builder
}
