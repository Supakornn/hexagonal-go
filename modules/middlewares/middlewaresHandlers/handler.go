package middlewaresHandlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/supakornn/hexagonal-go/config"
	"github.com/supakornn/hexagonal-go/modules/entities"
	"github.com/supakornn/hexagonal-go/modules/middlewares/middlewaresUsecases"
	"github.com/supakornn/hexagonal-go/pkg/auth"
)

type middlewareHanlderErrCode string

const (
	routerCheckErr middlewareHanlderErrCode = "router-001"
	jwtAuthErr     middlewareHanlderErrCode = "router-002"
	paramsCheckErr middlewareHanlderErrCode = "router-003"
)

type IMiddlewareHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
}

type middlewarehanlder struct {
	cfg                 config.Iconfig
	middlerwaresUsecase middlewaresUsecases.IMiddlewaresUsecase
}

func MiddlewarsHandler(cfg config.Iconfig, middlerwaresUsecase middlewaresUsecases.IMiddlewaresUsecase) IMiddlewareHandler {
	return &middlewarehanlder{
		cfg:                 cfg,
		middlerwaresUsecase: middlerwaresUsecase,
	}
}

func (h *middlewarehanlder) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (h *middlewarehanlder) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"router not found",
		).Res()
	}
}

func (h *middlewarehanlder) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} ${ip} ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Bangkok/Asia",
	})
}

func (h *middlewarehanlder) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := auth.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(fiber.ErrUnauthorized.Code, string(jwtAuthErr), err.Error()).Res()
		}

		claims := result.Claims
		if !h.middlerwaresUsecase.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(fiber.ErrUnauthorized.Code, string(jwtAuthErr), "no permission").Res()
		}

		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)
		return c.Next()
	}
}

func (h *middlewarehanlder) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		if c.Params("userId") != userId {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(paramsCheckErr),
				"never gonna give you",
			).Res()
		}
		return c.Next()
	}
}
