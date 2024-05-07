package server

import (
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/config"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/handlers"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/logger"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/storage/maps"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/usecase"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func StartServer(hostPort *config.Config) error {
	db := maps.NewMapStorage()
	useCase := usecase.NewUsecaseMemStorage(db)

	//логирование
	logg, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logg.Sync()
	logger.InitLogger(logg)
	sugar := *logg.Sugar()
	sugar.Infow("server started", "addr", hostPort.String())

	r := handlers.InitRoutes(*useCase)

	err = http.ListenAndServe(":"+strconv.Itoa(hostPort.Port), r)

	if err != nil {
		return err
	}
	return nil
}
