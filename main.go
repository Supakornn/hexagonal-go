package main

import (
	"os"

	"github.com/Supakornn/go-api/config"
	"github.com/Supakornn/go-api/pkg/databaes"
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
	db := databaes.DbConnect(cfg.Db())
	defer db.Close()
}
