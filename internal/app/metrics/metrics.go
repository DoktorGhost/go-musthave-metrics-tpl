package metrics

import (
	"math/rand"
	"runtime"
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

	gauge := m.Guage
	m.Counter++

	randomValue = rand.Float64() * 100
	gauge["RandomValue"] = randomValue
	gauge["Alloc"] = float64(memStats.Alloc)
	gauge["BuckHashSys"] = float64(memStats.BuckHashSys)
	gauge["Frees"] = float64(memStats.Frees)
	gauge["GCCPUFraction"] = float64(memStats.GCCPUFraction)
	gauge["GCSys"] = float64(memStats.GCSys)
	gauge["HeapAlloc"] = float64(memStats.HeapAlloc)
	gauge["HeapIdle"] = float64(memStats.HeapIdle)
	gauge["HeapInuse"] = float64(memStats.HeapInuse)
	gauge["HeapObjects"] = float64(memStats.HeapObjects)
	gauge["HeapReleased"] = float64(memStats.HeapReleased)
	gauge["HeapSys"] = float64(memStats.HeapSys)
	gauge["LastGC"] = float64(memStats.LastGC)
	gauge["Lookups"] = float64(memStats.Lookups)
	gauge["MCacheInuse"] = float64(memStats.MCacheInuse)
	gauge["MCacheSys"] = float64(memStats.MCacheSys)
	gauge["MSpanInuse"] = float64(memStats.MSpanInuse)
	gauge["MSpanSys"] = float64(memStats.MSpanSys)
	gauge["Mallocs"] = float64(memStats.Mallocs)
	gauge["NextGC"] = float64(memStats.NextGC)
	gauge["NumForcedGC"] = float64(memStats.NumForcedGC)
	gauge["NumGC"] = float64(memStats.NumGC)
	gauge["OtherSys"] = float64(memStats.OtherSys)
	gauge["PauseTotalNs"] = float64(memStats.PauseTotalNs)
	gauge["StackInuse"] = float64(memStats.StackInuse)
	gauge["StackSys"] = float64(memStats.StackSys)
	gauge["Sys"] = float64(memStats.Sys)
	gauge["TotalAlloc"] = float64(memStats.TotalAlloc)

}

func (m *Metrics) UpdateMetrics() {

}
