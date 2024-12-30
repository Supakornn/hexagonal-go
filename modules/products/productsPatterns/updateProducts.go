package productsPatterns

import "github.com/Supakornn/hexagonal-go/modules/entities"

type IUpdateProductBuilder interface {
	initTransaction() error
	initQuery()
	updateTitleQuery()
	updateDescriptionQuery()
	updatePriceQuery()
	updateCategory() error
	insertImages() error
	getOldImages() []*entities.Image
	deleteOldImages() error
	closeQuery()
	updateProduct() error
	getQueryFields() []string
	getValues() []any
	getQuery() string
	setQuery(query string)
	getImagesLen() int
	commit() error
}

type updateProductBuilder struct{}

func UpdateProductBuilder() IUpdateProductBuilder {
	return &updateProductBuilder{}
}

func (b *updateProductBuilder) initTransaction() error {
	return nil
}

func (b *updateProductBuilder) initQuery() {

}

func (b *updateProductBuilder) updateTitleQuery() {

}

func (b *updateProductBuilder) updateDescriptionQuery() {

}

func (b *updateProductBuilder) updatePriceQuery() {

}

func (b *updateProductBuilder) updateCategory() error {
	return nil
}

func (b *updateProductBuilder) insertImages() error {
	return nil
}

func (b *updateProductBuilder) getOldImages() []*entities.Image {
	return nil
}

func (b *updateProductBuilder) deleteOldImages() error {
	return nil
}

func (b *updateProductBuilder) closeQuery() {
}

func (b *updateProductBuilder) updateProduct() error {
	return nil
}

func (b *updateProductBuilder) getQueryFields() []string {
	return nil
}

func (b *updateProductBuilder) getValues() []any {
	return nil
}

func (b *updateProductBuilder) getQuery() string {
	return ""
}

func (b *updateProductBuilder) setQuery(query string) {
}

func (b *updateProductBuilder) getImagesLen() int {
	return 0
}

func (b *updateProductBuilder) commit() error {
	return nil
}

type updateProductEngineer struct {
	builder IFindProductBuilder
}

func UpdateProductEngineer(b IFindProductBuilder) *updateProductEngineer {
	return &updateProductEngineer{
		builder: b,
	}
}
