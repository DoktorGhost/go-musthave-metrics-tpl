package metrics

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/config"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/models"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"runtime"
)

var logger *zap.Logger

func init() {
	// Инициализация логгера
	logger, _ = zap.NewProduction()
}

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

func (m *Metrics) UpdateMetrics(client *http.Client, serverAddress string, conf *config.Config) {

	var bodys []models.Metrics
	bodys = append(bodys, models.Metrics{
		ID:    "PollCount",
		MType: "counter",
		Delta: &m.Counter,
	})
	for key, value := range m.Gauge {
		bodys = append(bodys, models.Metrics{
			ID:    key,
			MType: "gauge",
			Value: &value,
		})
	}

	for _, body := range bodys {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			logger.Error("Error occurred", zap.Error(err))
			break
		}
		reader := bytes.NewReader(jsonBody)

		request, err := http.NewRequest(http.MethodPost, serverAddress+"update", reader)
		if err != nil {
			logger.Error("Error occurred", zap.Error(err))
			break
		}
		request.Header.Add("Content-Type", "application/json")

		if conf.SecretKey != "" {
			h := hmac.New(sha256.New, []byte(conf.SecretKey))
			h.Write(jsonBody)
			hash := hex.EncodeToString(h.Sum(nil))
			request.Header.Add("HashSHA256", hash)
		}

		response, err := client.Do(request)
		if err != nil {
			logger.Error("Error occurred", zap.Error(err))
			break
		}

		defer response.Body.Close()
	}
}
