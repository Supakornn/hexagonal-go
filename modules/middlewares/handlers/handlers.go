package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/supakornn/hexagonal-go/config"
	"github.com/supakornn/hexagonal-go/modules/middlewares/usecases"
)

type IHandler interface {
	Cors() fiber.Handler
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
