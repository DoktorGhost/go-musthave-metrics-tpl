package main

import (
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/config"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/server"
	"strconv"
)

func main() {
	hostPort := config.ParseConfigServer()

	err := server.StartServer(strconv.Itoa(hostPort.Port))
	if err != nil {
		panic(err)
	}
}
