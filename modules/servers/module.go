package servers

import (
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
	router.Post("/refresh", handler.RefreshPassport)
	router.Post("/signout", handler.SignOut)
	router.Post("/signup-admin", handler.SignUpAdmin)
	router.Get("/secret", m.middlewares.JwtAuth(), handler.GenerateAdminToken)
}
