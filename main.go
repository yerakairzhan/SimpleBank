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
		log.Fatal("cannot connect config")
	}

	// Connect to the database
	conn, err := sql.Open(config.DBSource, config.DBDriver)
	if err != nil {
		log.Fatal("cannot connect to DB:", err)
	}

	// Initialize the store
	store := db.NewStore(conn)

	// Create a new API server
	server := api.NewServer(store)

	// Start the server
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
