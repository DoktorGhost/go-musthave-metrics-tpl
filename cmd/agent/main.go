package main

import (
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/metrics"
	"net/http"
	"sync"
	"time"
)

func main() {
	client := &http.Client{}
	host := "http://localhost:8080/"

	m := metrics.NewMetrics()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			m.CollectMetrics()
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		defer wg.Done()
		for {
			time.Sleep(10 * time.Second)
			m.UpdateMetrics(client, host)
		}
	}()
	wg.Wait()
}
