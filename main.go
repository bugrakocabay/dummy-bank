package main

import (
	"database/sql"
	"github.com/bugrakocabay/dummy-bank/api"
	db "github.com/bugrakocabay/dummy-bank/db/sqlc"
	"log"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	if err = server.Start(serverAddress); err != nil {
		log.Fatal("Server start failed:", err)
	}
}
