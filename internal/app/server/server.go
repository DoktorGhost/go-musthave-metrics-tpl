package server

import (
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/handlers"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/storage/maps"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/usecase"
	"net/http"
)

func StartServer(port string) error {
	db := maps.NewMapStorage()
	useCase := usecase.NewUsecaseMemStorage(db)
	r := handlers.InitRoutes(*useCase)

	err := http.ListenAndServe(":"+port, r)

	if err != nil {
		return err
	}
	return nil
}
