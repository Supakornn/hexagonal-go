package servers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/supakornn/hexagonal-go/modules/appinfo/appinfoHandlers"
	"github.com/supakornn/hexagonal-go/modules/appinfo/appinfoRepositories"
	"github.com/supakornn/hexagonal-go/modules/appinfo/appinfoUsecases"
	"github.com/supakornn/hexagonal-go/modules/files/filesHandlers"
	"github.com/supakornn/hexagonal-go/modules/files/filesUsecases"
	"github.com/supakornn/hexagonal-go/modules/middlewares/middlewaresHandlers"
	"github.com/supakornn/hexagonal-go/modules/middlewares/middlewaresRepositories"
	"github.com/supakornn/hexagonal-go/modules/middlewares/middlewaresUsecases"
	"github.com/supakornn/hexagonal-go/modules/monitor/monitorHandlers"
	"github.com/supakornn/hexagonal-go/modules/products/productsHandlers"
	"github.com/supakornn/hexagonal-go/modules/products/productsRepositories"
	"github.com/supakornn/hexagonal-go/modules/products/productsUsecases"
	"github.com/supakornn/hexagonal-go/modules/users/usersHandlers"
	"github.com/supakornn/hexagonal-go/modules/users/usersRepositories"
	"github.com/supakornn/hexagonal-go/modules/users/usersUsecases"
)

type IModuleFactory interface {
	MonitorModule()
	UserModule()
	AppinfoModule()
	FileModule()
	ProductsModule()
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
	router.Post("/signup", m.middleware.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", m.middleware.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", m.middleware.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", m.middleware.ApiKeyAuth(), handler.SignOut)

	router.Post("/signupadmin", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.SignUpAdmin)
	router.Get("/:userid", m.middleware.JwtAuth(), m.middleware.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.GenerateAdminToken)
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.server.db)
	usecase := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.server.cfg, usecase)

	router := m.router.Group("/appinfo")

	router.Get("/apikey", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.GenerateApiKey)
	router.Get("/categories", m.middleware.ApiKeyAuth(), handler.FindCategory)
	router.Post("/categories", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.AddCategory)
	router.Delete("/:category_id/categories", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.RemoveCategory)
}

func (m *moduleFactory) FileModule() {
	usecase := filesUsecases.FilesUsecase(m.server.cfg)
	handler := filesHandlers.FilesHandler(m.server.cfg, usecase)

	router := m.router.Group("/files")

	router.Post("/upload", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.UploadFiles)
	router.Patch("/delete", m.middleware.JwtAuth(), m.middleware.Authorize(2), handler.DeleteFile)
}
func (m *moduleFactory) ProductsModule() {
	filesUsecase := filesUsecases.FilesUsecase(m.server.cfg)
	repository := productsRepositories.ProductsRepository(m.server.db, m.server.cfg, filesUsecase)
	usecase := productsUsecases.ProductsUsecase(repository)
	handler := productsHandlers.ProductsHandler(m.server.cfg, usecase, filesUsecase)

	router := m.router.Group("/products")
	_ = handler
	_ = router
}
