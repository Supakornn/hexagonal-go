package productsHandlers

import (
	"fmt"
	"strings"

	"github.com/Supakornn/hexagonal-go/config"
	"github.com/Supakornn/hexagonal-go/modules/appinfo"
	"github.com/Supakornn/hexagonal-go/modules/entities"
	"github.com/Supakornn/hexagonal-go/modules/files"
	"github.com/Supakornn/hexagonal-go/modules/files/filesUsecases"
	"github.com/Supakornn/hexagonal-go/modules/products"
	"github.com/Supakornn/hexagonal-go/modules/products/productsUsecases"
	"github.com/gofiber/fiber/v2"
)

type productsHandlersErrCode string

const (
	findOneProductErr productsHandlersErrCode = "products-001"
	findProductsErr   productsHandlersErrCode = "products-002"
	insertProductErr  productsHandlersErrCode = "products-003"
	updateProductErr  productsHandlersErrCode = "products-004"
	deleteProductErr  productsHandlersErrCode = "products-005"
)

type IProductsHandler interface {
	FindOneProduct(c *fiber.Ctx) error
	FindProducts(c *fiber.Ctx) error
	InsertProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
}

type productsHandler struct {
	cfg             config.IConfig
	productsUsecase productsUsecases.IProductsUsecase
	filesUsecase    filesUsecases.IFilesUsecase
}

func ProductsHandler(cfg config.IConfig, productsUsecase productsUsecases.IProductsUsecase, filesUsecase filesUsecases.IFilesUsecase) IProductsHandler {
	return &productsHandler{
		cfg:             cfg,
		productsUsecase: productsUsecase,
		filesUsecase:    filesUsecase,
	}
}

func (h *productsHandler) FindOneProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), "")

	product, err := h.productsUsecase.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}

func (h *productsHandler) FindProducts(c *fiber.Ctx) error {
	req := &products.ProductFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findProductsErr),
			err.Error(),
		).Res()
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
		req.Sort = "asc"
	}

	products := h.productsUsecase.FindProducts(req)

	return entities.NewResponse(c).Success(fiber.StatusOK, products).Res()
}

func (h *productsHandler) InsertProduct(c *fiber.Ctx) error {
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

	product, err := h.productsUsecase.InsertProduct(req)
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

	if err := h.filesUsecase.DeleteFromGCP(deleteFileReq); err != nil {
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

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}
