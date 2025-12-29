package main

import (
	"log"
	"url-shortener/api"
	db "url-shortener/db/sqlc"
	"url-shortener/utils"
)

func main() {

	config, err := utils.LoadConfig(".")

	if err != nil {
		log.Fatal("cannot load config:", err)
		panic(err)
	}

	store, err := db.NewStore(config.DbSource)
	if err != nil {
		log.Fatal("cannot create store:", err)
		panic(err)
	}

	server := api.NewServer(&config, store)

	var ServerAddress = config.HttpServerAddress

	if ServerAddress == "" {
		ServerAddress = ":8080"
	}

	server.Start(ServerAddress)
}
