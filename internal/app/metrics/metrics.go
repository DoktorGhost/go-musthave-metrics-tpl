package metrics

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
)

type Metrics struct {
	Guage   map[string]float64
	Counter int64
}

func NewMetrics() *Metrics {
	return &Metrics{
		Guage:   make(map[string]float64),
		Counter: int64(0),
	}
}

func (m *Metrics) CollectMetrics() {
	var randomValue float64

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	m.Counter++

	randomValue = rand.Float64() * 100

	m.Guage["RandomValue"] = randomValue
	m.Guage["Alloc"] = float64(memStats.Alloc)
	m.Guage["BuckHashSys"] = float64(memStats.BuckHashSys)
	m.Guage["Frees"] = float64(memStats.Frees)
	m.Guage["GCCPUFraction"] = float64(memStats.GCCPUFraction)
	m.Guage["GCSys"] = float64(memStats.GCSys)
	m.Guage["HeapAlloc"] = float64(memStats.HeapAlloc)
	m.Guage["HeapIdle"] = float64(memStats.HeapIdle)
	m.Guage["HeapInuse"] = float64(memStats.HeapInuse)
	m.Guage["HeapObjects"] = float64(memStats.HeapObjects)
	m.Guage["HeapReleased"] = float64(memStats.HeapReleased)
	m.Guage["HeapSys"] = float64(memStats.HeapSys)
	m.Guage["LastGC"] = float64(memStats.LastGC)
	m.Guage["Lookups"] = float64(memStats.Lookups)
	m.Guage["MCacheInuse"] = float64(memStats.MCacheInuse)
	m.Guage["MCacheSys"] = float64(memStats.MCacheSys)
	m.Guage["MSpanInuse"] = float64(memStats.MSpanInuse)
	m.Guage["MSpanSys"] = float64(memStats.MSpanSys)
	m.Guage["Mallocs"] = float64(memStats.Mallocs)
	m.Guage["NextGC"] = float64(memStats.NextGC)
	m.Guage["NumForcedGC"] = float64(memStats.NumForcedGC)
	m.Guage["NumGC"] = float64(memStats.NumGC)
	m.Guage["OtherSys"] = float64(memStats.OtherSys)
	m.Guage["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	m.Guage["StackInuse"] = float64(memStats.StackInuse)
	m.Guage["StackSys"] = float64(memStats.StackSys)
	m.Guage["Sys"] = float64(memStats.Sys)
	m.Guage["TotalAlloc"] = float64(memStats.TotalAlloc)

}

func (m *Metrics) UpdateMetrics(client *http.Client, serverAddress string) {

	var endpoints []string
	endpoints = append(endpoints, serverAddress+"update/counter/PollCount/"+strconv.FormatInt(m.Counter, 10))
	for key, value := range m.Guage {
		endpoint := serverAddress + "update/guage/" + key + "/" + strconv.FormatFloat(value, 'f', -1, 64)
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
