package server

import (
	"fmt"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/config"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/handlers"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/logger"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/models"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/osfile"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/storage/maps"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/storage/postgres"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/usecase"
	"go.uber.org/zap"
	"io"
	"net/http"
	"sync"
	"time"
)

func StartServer(conf *config.Config) error {
	//логирование
	logg, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logg.Sync()
	logger.InitLogger(logg)
	sugar := *logg.Sugar()
	sugar.Infow("server started", "addr", conf.Host+":"+conf.Port)

	//инициализация бд
	var useCase *usecase.UsecaseMemStorage

	if conf.DatabaseDSN != "" {
		db, err := postgres.NewPostgresStorage(conf.DatabaseDSN)
		if err != nil {
			sugar.Fatalw("Ошибка при подключении к БД", "error", err)
		}
		useCase = usecase.NewUsecaseMemStorage(db)
		sugar.Infow("Успешное подключение к БД")
	} else {
		db := maps.NewMapStorage()
		useCase = usecase.NewUsecaseMemStorage(db)
		sugar.Infow("Использование оперативной памяти вместо БД")
	}

	//загрузка ранее сохраненных метрик
	if conf.Restore {
		cons, err := osfile.NewConsumer(conf.FileStoragePath)
		if err != nil {
			sugar.Infow("ошибка чтения конфигурациооного файла", err)
		} else {
			for {
				metric, err := cons.ReadEvent()
				if err != nil {
					if err == io.EOF {
						break
					}
					sugar.Infow("ошибка чтения события", err)
					continue
				}
				if metric == nil {
					break
				}
				if metric.MType == "gauge" {
					useCase.UsecaseUpdateGauge(metric.ID, *metric.Value)
					sugar.Infow("Запись из файла восстановлена в БД")
				} else if metric.MType == "counter" {
					useCase.UsecaseUpdateCounter(metric.ID, *metric.Delta)
					sugar.Infow("Запись из файла восстановлена в БД")
				}
			}
		}
		defer cons.Close()
	}

	r := handlers.InitRoutes(*useCase, conf)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			time.Sleep(time.Duration(conf.StoreInterval) * time.Second)
			mapsMetrics := useCase.UsecaseReadAll()
			prod, err := osfile.NewProducer(conf.FileStoragePath)
			if err != nil {
				sugar.Infow("Ошибка создания Producer:", err)
				continue
			}

			for key, value := range mapsMetrics {
				var metr models.Metrics
				switch v := value.(type) {
				case int64:
					metr = models.Metrics{
						ID:    key,
						MType: "counter",
						Delta: &v,
					}
				case float64:
					metr = models.Metrics{
						ID:    key,
						MType: "gauge",
						Value: &v,
					}
				default:
					sugar.Infow("Неизвестный тип метрики:", "key", key, "type", fmt.Sprintf("%T", value))
					continue
				}

				if err := prod.WriteEvent(&metr); err != nil {
					sugar.Infow("Ошибка записи события:", "error", err)
					continue
				} else {
					sugar.Infow("Успешная запись метрик")
				}
			}
			prod.Close()
		}
	}()

	err = http.ListenAndServe(":"+conf.Port, r)
	wg.Wait()
	if err != nil {
		return err
	}
	return nil
}
