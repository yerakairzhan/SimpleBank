package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/yerakairzhan/SimpleBank/api"
	db "github.com/yerakairzhan/SimpleBank/db/sqlc"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	// Connect to the database
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to DB:", err)
	}

	// Initialize the store
	store := db.NewStore(conn)

	// Create a new API server
	server := api.NewServer(store)

	// Start the server
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
