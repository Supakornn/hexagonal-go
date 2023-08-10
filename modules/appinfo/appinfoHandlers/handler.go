package appinfoHandlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/supakornn/hexagonal-go/config"
	"github.com/supakornn/hexagonal-go/modules/appinfo"
	"github.com/supakornn/hexagonal-go/modules/appinfo/appinfoUsecases"
	"github.com/supakornn/hexagonal-go/modules/entities"
	"github.com/supakornn/hexagonal-go/pkg/auth"
)

type appinfoHandlersErrCode string

const (
	GenerateApiKeyErr appinfoHandlersErrCode = "appinfoHandlers-001"
	FindCategoryErr   appinfoHandlersErrCode = "appinfoHandlers-002"
)

type IAppinfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error
	FindCategory(c *fiber.Ctx) error
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
		return entities.NewResponse(c).Error(fiber.ErrInternalServerError.Code, string(GenerateApiKeyErr), err.Error()).Res()
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
		return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(FindCategoryErr), err.Error()).Res()
	}

	category, err := h.appinfoUsecase.FindCategory(req)
	if err != nil {
		return entities.NewResponse(c).Error(fiber.ErrInternalServerError.Code, string(FindCategoryErr), err.Error()).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, category).Res()
}
