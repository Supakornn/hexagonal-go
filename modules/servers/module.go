package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/supakornn/hexagonal-go/modules/appinfo/appinfoHandlers"
	"github.com/supakornn/hexagonal-go/modules/appinfo/appinfoRepositories"
	"github.com/supakornn/hexagonal-go/modules/appinfo/appinfoUsecases"
	"github.com/supakornn/hexagonal-go/modules/middlewares/middlewaresHandlers"
	"github.com/supakornn/hexagonal-go/modules/middlewares/middlewaresRepositories"
	"github.com/supakornn/hexagonal-go/modules/middlewares/middlewaresUsecases"
	"github.com/supakornn/hexagonal-go/modules/monitor/monitorHandlers"
	"github.com/supakornn/hexagonal-go/modules/orders/ordersHandlers.go"
	"github.com/supakornn/hexagonal-go/modules/orders/ordersRepositories"
	"github.com/supakornn/hexagonal-go/modules/orders/ordersUsecases.go"
	"github.com/supakornn/hexagonal-go/modules/users/usersHandlers"
	"github.com/supakornn/hexagonal-go/modules/users/usersRepositories"
	"github.com/supakornn/hexagonal-go/modules/users/usersUsecases"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
	FilesModule() IFilesModule
	ProductsModule() IProductsModule
	OrdersModule()
}

type moduleFactory struct {
	router     fiber.Router
	server     *server
	middleware middlewaresHandlers.IMiddlewareHandler
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewareHandler) IModuleFactory {
	return &moduleFactory{
		router:     r,
		server:     s,
		middleware: mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewareHandler {
	repository := middlewaresRepositories.MiddlewaresRepo(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	return middlewaresHandlers.MiddlewarsHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.server.cfg)

	m.router.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UserRepository(m.server.db)
	usecase := usersUsecases.UserUsecase(m.server.cfg, repository)
	handler := usersHandlers.UserHandler(m.server.cfg, usecase)

	router := m.router.Group("/users")

	router.Post("/signup", m.middleware.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", m.middleware.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", m.middleware.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", m.middleware.ApiKeyAuth(), handler.SignOut)
	router.Post("/signup-admin", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.SignOut)

	router.Get("/:user_id", m.middleware.JwtAuth(), m.middleware.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.GenerateAdminToken)
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.server.db)
	usecase := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.server.cfg, usecase)

	router := m.router.Group("/appinfo")

	router.Post("/categories", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.AddCategory)

	router.Get("/categories", m.middleware.ApiKeyAuth(), handler.FindCategory)
	router.Get("/apikey", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.GenerateApiKey)

	router.Delete("/:category_id/categories", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.RemoveCategory)
}

func (m *moduleFactory) OrdersModule() {
	ordersRepository := ordersRepositories.OrdersRepository(m.server.db)
	ordersUsecase := ordersUsecases.OrdersUsecase(ordersRepository, m.ProductsModule().Repository())
	ordersHandler := ordersHandlers.OrdersHandler(m.server.cfg, ordersUsecase)

	router := m.router.Group("/orders")

	router.Post("/", m.middleware.JwtAuth(), ordersHandler.InsertOrder)

	router.Get("/", m.middleware.JwtAuth(), m.middleware.Authorize(2), ordersHandler.FindOrder)
	router.Get("/:user_id/:order_id", m.middleware.JwtAuth(), m.middleware.ParamsCheck(), ordersHandler.FindOneOrder)

	router.Patch("/:user_id/:order_id", m.middleware.JwtAuth(), m.middleware.ParamsCheck(), ordersHandler.UpdateOrder)
}
