package main

import (
	"os"

	"github.com/supakornn/hexagonal-go/config"
	"github.com/supakornn/hexagonal-go/modules/servers"
	"github.com/supakornn/hexagonal-go/pkg/databases"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := databases.DbConnect(cfg.Db())
	defer db.Close()

	servers.NewServer(cfg, db).Start()
}
