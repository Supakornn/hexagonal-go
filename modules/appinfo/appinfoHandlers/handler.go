package appinfoHandlers

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/supakornn/hexagonal-go/config"
	"github.com/supakornn/hexagonal-go/modules/appinfo"
	"github.com/supakornn/hexagonal-go/modules/appinfo/appinfoUsecases"
	"github.com/supakornn/hexagonal-go/modules/entities"
	"github.com/supakornn/hexagonal-go/pkg/auth"
)

type appinfoHandlersErrCode string

const (
	generateApiKeyErr appinfoHandlersErrCode = "appinfoHandlers-001"
	findCategoryErr   appinfoHandlersErrCode = "appinfoHandlers-002"
	addCategoryErr    appinfoHandlersErrCode = "appinfoHandlers-003"
	removeCategoryErr appinfoHandlersErrCode = "appinfoHandlers-004"
)

type IAppinfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error
	FindCategory(c *fiber.Ctx) error
	AddCategory(c *fiber.Ctx) error
	RemoveCategory(c *fiber.Ctx) error
}

type appinfoHandler struct {
	cfg            config.Iconfig
	appinfoUsecase appinfoUsecases.IAppinfoUsecase
}

func AppinfoHandler(cfg config.Iconfig, appinfoUsecase appinfoUsecases.IAppinfoUsecase) IAppinfoHandler {
	return &appinfoHandler{
		cfg:            cfg,
		appinfoUsecase: appinfoUsecase,
	}
}

func (h *appinfoHandler) GenerateApiKey(c *fiber.Ctx) error {
	apiKey, err := auth.NewAuth("apikey", h.cfg.Jwt(), nil)

	if err != nil {
		return entities.NewResponse(c).Error(fiber.ErrInternalServerError.Code, string(generateApiKeyErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, &struct {
		Key string `json:"key"`
	}{
		Key: apiKey.SignToken(),
	},
	).Res()
}

func (h *appinfoHandler) FindCategory(c *fiber.Ctx) error {
	req := new(appinfo.CategoryFilter)
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(findCategoryErr), err.Error()).Res()
	}

	category, err := h.appinfoUsecase.FindCategory(req)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.ErrInternalServerError.Code, string(findCategoryErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, category).Res()
}

func (h *appinfoHandler) AddCategory(c *fiber.Ctx) error {
	req := make([]*appinfo.Category, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(addCategoryErr), err.Error()).Res()
	}

	if len(req) == 0 {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(addCategoryErr), "request body is empty").Res()
	}

	if err := h.appinfoUsecase.InsertCategory(req); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrInternalServerError.Code, string(addCategoryErr), err.Error()).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, req).Res()
}

func (h *appinfoHandler) RemoveCategory(c *fiber.Ctx) error {
	categoryId := strings.Trim(c.Params("category_id"), " ")
	categoryIdInt, err := strconv.Atoi(categoryId)

	if err != nil {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(removeCategoryErr), err.Error()).Res()
	}

	if categoryIdInt <= 0 {
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(removeCategoryErr), "Id must more than 0").Res()
	}

	if err := h.appinfoUsecase.DeleteCategory(categoryIdInt); err != nil {
		return entities.NewResponse(c).Error(fiber.ErrInternalServerError.Code, string(removeCategoryErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK,
		&struct {
			CategoryId int `json:"category_id"`
		}{
			CategoryId: categoryIdInt,
		},
	).Res()
}
