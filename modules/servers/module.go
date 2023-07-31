package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/supakornn/hexagonal-go/modules/middlewares/middlewaresHandlers"
	"github.com/supakornn/hexagonal-go/modules/middlewares/middlewaresRepositories"
	"github.com/supakornn/hexagonal-go/modules/middlewares/middlewaresUsecases"
	"github.com/supakornn/hexagonal-go/modules/monitor/monitorHandlers"
)

type IModuleFactory interface {
	MonitorModule()
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
