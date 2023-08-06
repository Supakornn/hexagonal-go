package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/supakornn/hexagonal-go/modules/middlewares/middlewaresHandlers"
	"github.com/supakornn/hexagonal-go/modules/middlewares/middlewaresRepositories"
	"github.com/supakornn/hexagonal-go/modules/middlewares/middlewaresUsecases"
	"github.com/supakornn/hexagonal-go/modules/monitor/monitorHandlers"
	"github.com/supakornn/hexagonal-go/modules/users/usersHandlers"
	"github.com/supakornn/hexagonal-go/modules/users/usersRepositories"
	"github.com/supakornn/hexagonal-go/modules/users/usersUsecases"
)

type IModuleFactory interface {
	MonitorModule()
	UserModule()
}

type moduleFactory struct {
	router     fiber.Router
	server     *server
	middleware middlewaresHandlers.IMiddlewareHandler
}

func InitModule(r fiber.Router, s *server, m middlewaresHandlers.IMiddlewareHandler) IModuleFactory {
	return &moduleFactory{
		router:     r,
		server:     s,
		middleware: m,
	}
}

func InitMiddleware(s *server) middlewaresHandlers.IMiddlewareHandler {
	repository := middlewaresRepositories.MiddlewaresRepo(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	return middlewaresHandlers.MiddlewarsHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.server.cfg)
	m.router.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UserModule() {
	repository := usersRepositories.UserRepository(m.server.db)
	usecase := usersUsecases.UserUsecase(m.server.cfg, repository)
	handler := usersHandlers.UserHandler(m.server.cfg, usecase)

	router := m.router.Group("/users")
	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
	router.Post("/refresh", handler.RefreshPassport)
	router.Post("/signout", handler.SignOut)
	router.Post("/signupadmin", handler.SignUpAdmin)
	router.Get("/:userid", m.middleware.JwtAuth(), m.middleware.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.GenerateAdminToken)
}
