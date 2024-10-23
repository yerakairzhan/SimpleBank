package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/yerakairzhan/SimpleBank/api"
	db "github.com/yerakairzhan/SimpleBank/db/sqlc"
	"github.com/yerakairzhan/SimpleBank/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	log.Printf("Loaded config: %+v", config)

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to DB:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
