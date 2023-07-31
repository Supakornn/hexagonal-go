package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/supakornn/hexagonal-go/modules/middlewares/handlers"
	"github.com/supakornn/hexagonal-go/modules/middlewares/repositories"
	"github.com/supakornn/hexagonal-go/modules/middlewares/usecases"
	"github.com/supakornn/hexagonal-go/modules/monitor/monitorHandlers"
)

type IModuleFactory interface {
	MonitorModule()
}

type moduleFactory struct {
	router     fiber.Router
	server     *server
	middleware handlers.IHandler
}

func InitModule(r fiber.Router, s *server, m handlers.IHandler) IModuleFactory {
	return &moduleFactory{
		router:     r,
		server:     s,
		middleware: m,
	}
}

func InitMiddleware(s *server) handlers.IHandler {
	repository := repositories.MiddlewarsRepo(s.db)
	usecase := usecases.MiddlewareUsecase(repository)
	return handlers.MiddlewarsHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.server.cfg)

	m.router.Get("/", handler.HealthCheck)

}
