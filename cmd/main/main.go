package main

import (
	"log"
	"newsapi/internal/di"
)

func main() {
	server := di.InitHTTPServer()
	server.ConnectCoreWithEcho()

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("failed to serve http:", err.Error())
	}
}
