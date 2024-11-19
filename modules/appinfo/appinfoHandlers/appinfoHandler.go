package appinfoHandlers

import (
	"github.com/Supakornn/hexagonal-go/config"
	"github.com/Supakornn/hexagonal-go/modules/appinfo"
	"github.com/Supakornn/hexagonal-go/modules/appinfo/appinfoUsecases"
	"github.com/Supakornn/hexagonal-go/modules/entities"
	"github.com/Supakornn/hexagonal-go/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

type appinfoHandlerErrCode string

const (
	GenerateApiKeyErr  appinfoHandlerErrCode = "appinfo-001"
	FindCategoryErrErr appinfoHandlerErrCode = "appinfo-002"
	AddCategoryErr     appinfoHandlerErrCode = "appinfo-003"
)

type IAppinfoHandler interface {
	GenerateApiKey(ctx *fiber.Ctx) error
	FindCategory(c *fiber.Ctx) error
	AddCategory(c *fiber.Ctx) error
}

type appinfoHandler struct {
	cfg            config.IConfig
	appinfoUsecase appinfoUsecases.IAppinfoUsecase
}

func AppinfoHandler(cfg config.IConfig, appinfoUsecase appinfoUsecases.IAppinfoUsecase) IAppinfoHandler {
	return &appinfoHandler{
		cfg:            cfg,
		appinfoUsecase: appinfoUsecase,
	}
}

func (h *appinfoHandler) GenerateApiKey(c *fiber.Ctx) error {
	apikey, err := auth.NewAuth(auth.Apikey, h.cfg.Jwt(), nil)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(GenerateApiKeyErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Key string `json:"key"`
		}{
			Key: apikey.SignToken(),
		},
	).Res()
}

func (h *appinfoHandler) FindCategory(c *fiber.Ctx) error {
	req := new(appinfo.CategoryFilter)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(FindCategoryErrErr),
			"invalid request",
		).Res()
	}

	category, err := h.appinfoUsecase.FindCategory(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(FindCategoryErrErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, category).Res()
}

func (h *appinfoHandler) AddCategory(c *fiber.Ctx) error {
	req := make([]*appinfo.Category, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(AddCategoryErr),
			"invalid request",
		).Res()
	}

	if err := h.appinfoUsecase.InsertCategory(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(AddCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, req).Res()
}
