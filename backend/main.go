package main

import (
	"MORE.Tech/backend/db"
	"MORE.Tech/backend/server"
	"log"
)

func main() {
	err := db.GetDB().DB().Ping()
	if err != nil {
		log.Fatal(err)
		return
	}

	server.RunForever()
}
