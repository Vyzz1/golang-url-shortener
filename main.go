package main

import (
	"url-shortener/api"
	db "url-shortener/db/sqlc"
	"url-shortener/utils"
)

func main() {

	config, err := utils.LoadConfig(".")

	if err != nil {
		panic(err)
	}

	store, err := db.NewStore(config.DBSOURCE)
	if err != nil {
		panic(err)
	}

	server := api.NewServer(&config, store)

	server.Start(config.HttpServerAddress)
}
