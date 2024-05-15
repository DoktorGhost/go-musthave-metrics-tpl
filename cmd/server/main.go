package main

import (
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/config"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/server"
)

func main() {
	conf := config.ParseConfigServer()

	err := server.StartServer(conf)
	if err != nil {
		panic(err)
	}
}
