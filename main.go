package main

import (
	"database/sql"
	"github.com/bugrakocabay/dummy-bank/api"
	db "github.com/bugrakocabay/dummy-bank/db/sqlc"
	"github.com/bugrakocabay/dummy-bank/util"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Error with loading env: ", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Cannot start server:", err)
	}

	if err = server.Start(config.ServerAddress); err != nil {
		log.Fatal("Server start failed:", err)
	}
}
