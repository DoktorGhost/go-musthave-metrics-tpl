package main

import (
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/config"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/logger"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/metrics"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func main() {

	client := &http.Client{}
	addr, ReportInterval, PollInterval := config.ParseConfigClient()

	host := "http://" + addr.Host + ":" + strconv.Itoa(addr.Port) + "/"

	logg, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logg.Sync()
	logger.InitLogger(logg)
	sugar := *logg.Sugar()
	sugar.Infow("сlient started", "addr", host)

	m := metrics.NewMetrics()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			m.CollectMetrics()
			time.Sleep(time.Duration(PollInterval) * time.Second)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			time.Sleep(time.Duration(ReportInterval) * time.Second)
			m.UpdateMetrics(client, host)
		}
	}()
	wg.Wait()
}
