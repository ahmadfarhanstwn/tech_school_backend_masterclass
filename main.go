package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ahmadfarhanstwn/simple_bank/api"
	db "github.com/ahmadfarhanstwn/simple_bank/db/sqlc"
	"github.com/ahmadfarhanstwn/simple_bank/util"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("starting server ...")
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config", err)
	}
	conn, err := sql.Open(config.DB_Driver, config.DB_Source)
	if err != nil {
		log.Fatal("can't connect to database:", err)
	}
	store := db.NewStore(conn)
	server,err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server")
	}
	err = server.Start(config.Server_Address)
	if err != nil {
		log.Fatal("cannot start server")
	}
}