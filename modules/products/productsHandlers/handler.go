package productsHandlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/supakornn/hexagonal-go/config"
	"github.com/supakornn/hexagonal-go/modules/entities"
	"github.com/supakornn/hexagonal-go/modules/files/filesUsecases"
	"github.com/supakornn/hexagonal-go/modules/products/productsUsecases"
)

type productsHandlerErrCode string

const (
	findOneProductErr productsHandlerErrCode = "products-001"
)

type IProductsHandler interface {
	FindOneProduct(c *fiber.Ctx) error
}

type productsHandler struct {
	cfg             config.Iconfig
	productsUsecase productsUsecases.IProductsUsecase
	fileUsecase     filesUsecases.IFilesUsecase
}

func ProductsHandler(cfg config.Iconfig, productUsecase productsUsecases.IProductsUsecase, fileUsecase filesUsecases.IFilesUsecase) IProductsHandler {
	return &productsHandler{
		productsUsecase: productUsecase,
		cfg:             cfg,
		fileUsecase:     fileUsecase,
	}
}

func (h *productsHandler) FindOneProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	product, err := h.productsUsecase.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.ErrInternalServerError.Code, string(findOneProductErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}
