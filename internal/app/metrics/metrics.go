package metrics

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
)

type Metrics struct {
	Gauge   map[string]float64
	Counter int64
}

func NewMetrics() *Metrics {
	return &Metrics{
		Gauge:   make(map[string]float64),
		Counter: int64(0),
	}
}

func (m *Metrics) CollectMetrics() {
	var randomValue float64

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.Counter++

	randomValue = rand.Float64() * 100

	m.Gauge["RandomValue"] = randomValue
	m.Gauge["Alloc"] = float64(memStats.Alloc)
	m.Gauge["BuckHashSys"] = float64(memStats.BuckHashSys)
	m.Gauge["Frees"] = float64(memStats.Frees)
	m.Gauge["GCCPUFraction"] = float64(memStats.GCCPUFraction)
	m.Gauge["GCSys"] = float64(memStats.GCSys)
	m.Gauge["HeapAlloc"] = float64(memStats.HeapAlloc)
	m.Gauge["HeapIdle"] = float64(memStats.HeapIdle)
	m.Gauge["HeapInuse"] = float64(memStats.HeapInuse)
	m.Gauge["HeapObjects"] = float64(memStats.HeapObjects)
	m.Gauge["HeapReleased"] = float64(memStats.HeapReleased)
	m.Gauge["HeapSys"] = float64(memStats.HeapSys)
	m.Gauge["LastGC"] = float64(memStats.LastGC)
	m.Gauge["Lookups"] = float64(memStats.Lookups)
	m.Gauge["MCacheInuse"] = float64(memStats.MCacheInuse)
	m.Gauge["MCacheSys"] = float64(memStats.MCacheSys)
	m.Gauge["MSpanInuse"] = float64(memStats.MSpanInuse)
	m.Gauge["MSpanSys"] = float64(memStats.MSpanSys)
	m.Gauge["Mallocs"] = float64(memStats.Mallocs)
	m.Gauge["NextGC"] = float64(memStats.NextGC)
	m.Gauge["NumForcedGC"] = float64(memStats.NumForcedGC)
	m.Gauge["NumGC"] = float64(memStats.NumGC)
	m.Gauge["OtherSys"] = float64(memStats.OtherSys)
	m.Gauge["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	m.Gauge["StackInuse"] = float64(memStats.StackInuse)
	m.Gauge["StackSys"] = float64(memStats.StackSys)
	m.Gauge["Sys"] = float64(memStats.Sys)
	m.Gauge["TotalAlloc"] = float64(memStats.TotalAlloc)

}

func (m *Metrics) UpdateMetrics(client *http.Client, serverAddress string) {

	var endpoints []string
	endpoints = append(endpoints, serverAddress+"update/counter/PollCount/"+strconv.FormatInt(m.Counter, 10))
	for key, value := range m.Gauge {
		endpoint := serverAddress + "update/gauge/" + key + "/" + strconv.FormatFloat(value, 'f', -1, 64)
		endpoints = append(endpoints, endpoint)
	}
	for _, endpoint := range endpoints {
		request, err := http.NewRequest(http.MethodPost, endpoint, nil)
		if err != nil {
			fmt.Println(err)
			break
		}
		request.Header.Add("Content-Type", "text/plain")
		response, err := client.Do(request)
		if err != nil {
			fmt.Println(err)
			break
		}

		defer response.Body.Close()
	}
}
