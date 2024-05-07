package main

import (
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/config"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/server"
)

func main() {
	hostPort := config.ParseConfigServer()

	err := server.StartServer(hostPort)
	if err != nil {
		panic(err)
	}
}
