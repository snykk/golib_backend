package main

import (
	"log"

	"github.com/snykk/golib_backend/config"
	"github.com/snykk/golib_backend/server"
)

func init() {
	if err := config.InitializeAppConfig(); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	app, err := server.NewApp()
	if err != nil {
		panic(err)
	}
	if err := app.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
