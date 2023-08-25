package servers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/supakornn/hexagonal-go/config"
)

type IServer interface {
	Start()
	GetServer() *server
}

type server struct {
	app *fiber.App
	cfg config.Iconfig
	db  *sqlx.DB
}

func NewServer(cfg config.Iconfig, db *sqlx.DB) IServer {
	return &server{
		cfg: cfg,
		db:  db,
		app: fiber.New(fiber.Config{
			AppName:      cfg.App().Name(),
			BodyLimit:    cfg.App().BodyLimit(),
			ReadTimeout:  cfg.App().ReadTimeout(),
			WriteTimeout: cfg.App().WriteTimeout(),
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		}),
	}
}

func (s *server) GetServer() *server {
	return s
}

func (s *server) Start() {

	m := InitMiddlewares(s)
	s.app.Use(m.Logger())
	s.app.Use(m.Cors())

	v1 := s.app.Group("v1")
	modules := InitModule(v1, s, m)
	modules.MonitorModule()
	modules.UsersModule()
	modules.AppinfoModule()
	modules.FilesModule()
	modules.ProductsModule()
	modules.OrdersModule()

	s.app.Use(m.RouterCheck())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		log.Println("server is shutting down...")
		_ = s.app.Shutdown()
	}()

	log.Printf("Sever is running on %v", s.cfg.App().Url())
	s.app.Listen(s.cfg.App().Url())
}
