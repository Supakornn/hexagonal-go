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
	"github.com/supakornn/hexagonal-go/pkg/utils"
)

type middlewareHanlderErrCode string

const (
	routerCheckErr middlewareHanlderErrCode = "router-001"
	jwtAuthErr     middlewareHanlderErrCode = "router-002"
	paramsCheckErr middlewareHanlderErrCode = "router-003"
	AuthorizeErr   middlewareHanlderErrCode = "router-004"
	ApiKeyErr      middlewareHanlderErrCode = "router-005"
)

type IMiddlewareHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
	Authorize(expectRoleId ...int) fiber.Handler
	ApiKeyAuth() fiber.Handler
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
		if c.Locals("userRoleId").(int) == 2 {
			return c.Next()
		}
		if c.Params("user_id") != userId {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(paramsCheckErr),
				"never gonna give you up",
			).Res()
		}
		return c.Next()
	}
}

func (h *middlewarehanlder) Authorize(expectRoleId ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleId, ok := c.Locals("userRoleId").(int)
		if !ok {
			return entities.NewResponse(c).Error(fiber.ErrBadRequest.Code, string(AuthorizeErr), "role id is not int").Res()
		}

		roles, err := h.middlerwaresUsecase.FindRole()
		if err != nil {
			return entities.NewResponse(c).Error(fiber.ErrInternalServerError.Code, string(AuthorizeErr), err.Error()).Res()
		}

		sum := 0

		for _, v := range expectRoleId {
			sum += v
		}

		expectValueBinary := utils.BinaryConverter(sum, len(roles))
		userValueBinary := utils.BinaryConverter(roleId, len(roles))

		for i := range userValueBinary {
			if userValueBinary[i]&expectValueBinary[i] == 1 {
				return c.Next()
			}
		}
		return entities.NewResponse(c).Error(fiber.ErrUnauthorized.Code, string(AuthorizeErr), "no permission").Res()
	}
}

func (m *middlewarehanlder) ApiKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-API-KEY")
		if _, err := auth.ParseApiKey(m.cfg.Jwt(), key); err != nil {
			return entities.NewResponse(c).Error(fiber.ErrUnauthorized.Code, "api key is invalid", err.Error()).Res()
		}
		return c.Next()
	}
}
