package middlewareHandlers

import (
	"github.com/Supakornn/go-api/config"
	"github.com/Supakornn/go-api/modules/middlewares/middlewareUsecases"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler
}

type middlewaresHandler struct {
	cfg               config.IConfig
	middlewareUsecase middlewareUsecases.IMiddlewaresUsecase
}

func MiddlewaresHandler(cfg config.IConfig, middlewareUsecase middlewareUsecases.IMiddlewaresUsecase) IMiddlewaresHandler {
	return &middlewaresHandler{
		cfg:               cfg,
		middlewareUsecase: middlewareUsecase,
	}
}

func (h *middlewaresHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,PATCH,DELETE",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}
