package servers

import (
	"github.com/Supakornn/go-api/modules/middlewares/middlewareHandlers"
	"github.com/Supakornn/go-api/modules/middlewares/middlewareUsecases"
	"github.com/Supakornn/go-api/modules/middlewares/middlewaresRepositories"
	"github.com/Supakornn/go-api/modules/monitor/monitorHandlers"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
}

type moduleFactory struct {
	router      fiber.Router
	server      *server
	middlewares middlewareHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, m middlewareHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		router:      r,
		server:      s,
		middlewares: m,
	}
}

func InitMiddlewares(s *server) middlewareHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewareUsecases.MiddlewaresUsecase(repository)
	return middlewareHandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandlers(m.server.cfg)
	m.router.Get("/", handler.HealthCheck)
}
