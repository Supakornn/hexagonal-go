package servers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/Supakornn/hexagonal-go/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type IServer interface {
	Start()
}

type server struct {
	app *fiber.App
	cfg config.IConfig
	db  *sqlx.DB
}

func NewServer(cfg config.IConfig, db *sqlx.DB) IServer {
	return &server{
		cfg: cfg,
		db:  db,
		app: fiber.New(fiber.Config{
			AppName:      cfg.App().Name(),
			BodyLimit:    cfg.App().BodyLimit(),
			ReadTimeout:  cfg.App().ReadTimeout(),
			WriteTimeout: cfg.App().ReadTimeout(),
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		}),
	}
}

func (s *server) Start() {

	middlewares := InitMiddlewares(s)
	s.app.Use(middlewares.Logger())
	s.app.Use(middlewares.Cors())

	v1 := s.app.Group("/v1")
	modules := InitModule(v1, s, middlewares)
	modules.MonitorModule()
	modules.UsersModule()
	modules.AppinfoModule()
	modules.FilesModule()
	modules.FilesModule()
	modules.ProductsModule()
	modules.OrdersModule()

	s.app.Use(middlewares.RouterCheck())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		log.Println("Shutting down server...")
		_ = s.app.Shutdown()
	}()

	log.Printf("Server is running on url:%s", s.cfg.App().Url())
	s.app.Listen(s.cfg.App().Url())
}
