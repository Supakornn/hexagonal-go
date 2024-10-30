package servers

import (
	"github.com/Supakornn/hexagonal-go/modules/appinfo/appinfoHandlers"
	"github.com/Supakornn/hexagonal-go/modules/appinfo/appinfoRepositories"
	"github.com/Supakornn/hexagonal-go/modules/appinfo/appinfoUsecases"
	"github.com/Supakornn/hexagonal-go/modules/middlewares/middlewaresHandlers"
	"github.com/Supakornn/hexagonal-go/modules/middlewares/middlewaresRepositories"
	"github.com/Supakornn/hexagonal-go/modules/middlewares/middlewaresUsecases"
	"github.com/Supakornn/hexagonal-go/modules/monitor/monitorHandlers"
	"github.com/Supakornn/hexagonal-go/modules/users/usersHandlers"
	"github.com/Supakornn/hexagonal-go/modules/users/usersRepositories"
	"github.com/Supakornn/hexagonal-go/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
}

type moduleFactory struct {
	router      fiber.Router
	server      *server
	middlewares middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, m middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		router:      r,
		server:      s,
		middlewares: m,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	return middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandlers(m.server.cfg)
	m.router.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.server.db)
	usecase := usersUsecases.UsersUsecase(m.server.cfg, repository)
	handler := usersHandlers.UsersHandler(m.server.cfg, usecase)

	router := m.router.Group("/users")

	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
	router.Post("/signout", handler.SignOut)
	router.Get("/:user_id", m.middlewares.JwtAuth(), m.middlewares.ParamsCheck(), handler.GetProfile)

	router.Post("/signup-admin", handler.SignUpAdmin)
	router.Get("/admin/secret", m.middlewares.JwtAuth(), m.middlewares.Authorize(2), handler.GenerateAdminToken)

	router.Post("/refresh", handler.RefreshPassport)
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.server.db)
	usecase := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.server.cfg, usecase)

	router := m.router.Group("/appinfo")
	router.Get("/apikey", m.middlewares.JwtAuth(), m.middlewares.Authorize(2), handler.GenerateApiKey)
}
