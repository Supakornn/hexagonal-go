package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/supakornn/hexagonal-go/config"
	"github.com/supakornn/hexagonal-go/modules/entities"
	"github.com/supakornn/hexagonal-go/modules/middlewares/usecases"
)

type middlewareHanlderErrCode string

const (
	routerCheckErr middlewareHanlderErrCode = "router-001"
)

type IHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
}

type hanlder struct {
	cfg                 config.Iconfig
	middlerwaresUsecase usecases.IUsecase
}

func MiddlewarsHandler(cfg config.Iconfig, middlerwaresUsecase usecases.IUsecase) IHandler {
	return &hanlder{
		cfg:                 cfg,
		middlerwaresUsecase: middlerwaresUsecase,
	}
}

func (h *hanlder) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (h *hanlder) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"router not found",
		).Res()
	}
}

func (h *hanlder) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} ${ip} ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Bangkok/Asia",
	})
}
