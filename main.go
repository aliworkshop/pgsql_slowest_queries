package main

import (
	"github.com/aliworkshop/pgsql_slowest_queries/application"
	"github.com/aliworkshop/pgsql_slowest_queries/handler"
	"log"
)

func main() {

	app := application.NewApp()
	handlers := handler.NewHandler(app.GetDB())

	app.RegisterRoutes(handlers.Hello, "/", application.GET)
	app.RegisterRoutes(handlers.SlowestConnectionHandler, "/slowest_connection", application.GET)

	if err := app.Start(); err != nil {
		log.Fatal("error happened on app starting... : " + err.Error())
	}
}
