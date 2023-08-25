package productsHandlers

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/supakornn/hexagonal-go/config"
	"github.com/supakornn/hexagonal-go/modules/appinfo"
	"github.com/supakornn/hexagonal-go/modules/entities"
	"github.com/supakornn/hexagonal-go/modules/files"
	"github.com/supakornn/hexagonal-go/modules/files/filesUsecases"
	"github.com/supakornn/hexagonal-go/modules/products"
	"github.com/supakornn/hexagonal-go/modules/products/productsUsecases"
)

type productsHandlerErrCode string

const (
	findOneProductErr productsHandlerErrCode = "products-001"
	findProductErr    productsHandlerErrCode = "products-002"
	insertProductErr  productsHandlerErrCode = "products-003"
	deleteProductErr  productsHandlerErrCode = "products-004"
	updateProductErr  productsHandlerErrCode = "products-005"
)

type IProductsHandler interface {
	FindOneProduct(c *fiber.Ctx) error
	FindProduct(c *fiber.Ctx) error
	AddProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
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

func (h *productsHandler) FindProduct(c *fiber.Ctx) error {
	req := &products.ProductFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(findProductErr), err.Error()).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 5 {
		req.Limit = 5
	}

	if req.OrderBy == "" {
		req.OrderBy = "title"
	}

	if req.Sort == "" {
		req.Sort = "ASC"
	}

	products := h.productsUsecase.FindProduct(req)

	return entities.NewResponse(c).Success(fiber.StatusOK, products).Res()
}

func (h *productsHandler) AddProduct(c *fiber.Ctx) error {
	req := &products.Product{
		Category: &appinfo.Category{},
		Images:   make([]*entities.Image, 0),
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProductErr),
			err.Error(),
		).Res()
	}
	if req.Category.Id <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProductErr),
			"category id is invalid",
		).Res()
	}

	product, err := h.productsUsecase.AddProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertProductErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, product).Res()
}

func (h *productsHandler) UpdateProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")
	req := &products.Product{
		Images:   make([]*entities.Image, 0),
		Category: &appinfo.Category{},
	}
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateProductErr),
			err.Error(),
		).Res()
	}
	req.Id = productId

	product, err := h.productsUsecase.UpdateProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateProductErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}

func (h *productsHandler) DeleteProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	product, err := h.productsUsecase.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}

	deleteFileReq := make([]*files.DeleteFileReq, 0)
	for _, p := range product.Images {
		deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
			Destination: fmt.Sprintf("images/products/%s", p.FileName),
		})
	}
	if err := h.fileUsecase.DeletefileOnGCP(deleteFileReq); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}

	if err := h.productsUsecase.DeleteProduct(productId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusNoContent, nil).Res()
}
