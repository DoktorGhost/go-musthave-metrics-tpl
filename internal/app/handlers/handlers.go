package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/compressor"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/config"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/crypto"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/logger"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/models"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/usecase"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func InitRoutes(useCase usecase.UsecaseMemStorage, conf *config.Config) chi.Router {
	r := chi.NewRouter()

	r.Use(logger.WithLogging)
	r.Use(compressor.GzipMiddleware)
	r.Use(compressor.DecompressMiddleware)
	r.Use(crypto.CryptoMiddleware(conf))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handlerAllMetrics(w, r, useCase)
	})
	r.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		handlerPost(w, r, useCase)
	})
	r.Post("/update/", func(w http.ResponseWriter, r *http.Request) {
		handlerJSONUpdate(w, r, useCase)
	})
	r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
		handlerJSONUpdate(w, r, useCase)
	})
	r.Get("/value/{type}/{name}", func(w http.ResponseWriter, r *http.Request) {
		handlerGet(w, r, useCase)
	})
	r.Post("/value/", func(w http.ResponseWriter, r *http.Request) {
		handlerJSONValue(w, r, useCase)
	})
	r.Post("/value", func(w http.ResponseWriter, r *http.Request) {
		handlerJSONValue(w, r, useCase)
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		handlerPing(w, r, conf)
	})
	r.Post("/updates/", func(w http.ResponseWriter, r *http.Request) {
		handlerUpdates(w, r, useCase)
	})

	return r
}

func handlerPost(res http.ResponseWriter, req *http.Request, useCase usecase.UsecaseMemStorage) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	typeMetric := chi.URLParam(req, "type")
	nameMetric := chi.URLParam(req, "name")
	valueMetric := chi.URLParam(req, "value")

	if nameMetric == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	if typeMetric == "" || valueMetric == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	if typeMetric == "gauge" {
		valueFloat, err := strconv.ParseFloat(valueMetric, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("ошибка конвертации"))
			return
		}
		useCase.UsecaseUpdateGauge(nameMetric, valueFloat)
		res.WriteHeader(http.StatusOK)
	} else if typeMetric == "counter" {
		valueInt, err := strconv.ParseInt(valueMetric, 10, 64)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("ошибка конвертации"))
			return
		}
		useCase.UsecaseUpdateCounter(nameMetric, valueInt)
		res.WriteHeader(http.StatusOK)
		return
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}

func handlerGet(res http.ResponseWriter, req *http.Request, useCase usecase.UsecaseMemStorage) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	typeMetric := chi.URLParam(req, "type")
	nameMetric := chi.URLParam(req, "name")

	if nameMetric == "" {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	if typeMetric == "" {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	value := useCase.UsecaseRead(typeMetric, nameMetric)
	if value != nil {
		stringValue := fmt.Sprintf("%v", value)
		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(stringValue))
	} else {
		res.WriteHeader(http.StatusNotFound)
	}

}

func handlerAllMetrics(res http.ResponseWriter, req *http.Request, useCase usecase.UsecaseMemStorage) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	metrics := useCase.UsecaseReadAll()

	// Создаем шаблон HTML-страницы
	htmlTemplate := `
        <html>
        <head><title>Metrics</title></head>
        <body>
            <h1>Список метрик</h1>
            <ul>
                {{range $key, $value := .}}
                <li>{{$key}}: {{$value}}</li>
                {{end}}
            </ul>
        </body>
        </html>
    `

	// Создаем HTML-страницу на основе шаблона
	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		http.Error(res, "Ошибка создания HTML-страницы", http.StatusInternalServerError)
		return
	}

	// Отображаем HTML-страницу
	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusOK)
	err = tmpl.Execute(res, metrics)
	if err != nil {
		http.Error(res, "Ошибка отображения HTML-страницы", http.StatusInternalServerError)
		return
	}

}

func handlerJSONUpdate(w http.ResponseWriter, r *http.Request, useCase usecase.UsecaseMemStorage) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if req.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if req.MType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.MType == "gauge" {
		if req.Value != nil {
			useCase.UsecaseUpdateGauge(req.ID, *req.Value)
			key := useCase.UsecaseRead(req.MType, req.ID)
			*req.Value = key.(float64)
		} else {
			return
		}
	} else if req.MType == "counter" {
		if req.Delta != nil {
			useCase.UsecaseUpdateCounter(req.ID, *req.Delta)
			key := useCase.UsecaseRead(req.MType, req.ID)
			*req.Delta = key.(int64)
		} else {
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	if err := enc.Encode(req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handlerJSONValue(w http.ResponseWriter, r *http.Request, useCase usecase.UsecaseMemStorage) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.Metrics
	var res models.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		log.Println("ошибка декодирования JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if req.ID == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if req.MType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	value := useCase.UsecaseRead(req.MType, req.ID)

	if value == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		if req.MType == "gauge" {
			res.ID = req.ID
			res.MType = req.MType
			vv := value.(float64)
			res.Value = &vv
		} else if req.MType == "counter" {
			res.ID = req.ID
			res.MType = req.MType
			vv := value.(int64)
			res.Delta = &vv
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	if err := enc.Encode(res); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handlerPing(res http.ResponseWriter, req *http.Request, conf *config.Config) {
	if req.Method != http.MethodGet {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ps := conf.DatabaseDSN

	db, err := sql.Open("pgx", ps)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func handlerUpdates(w http.ResponseWriter, r *http.Request, useCase usecase.UsecaseMemStorage) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req []models.Metrics
	var res models.Metrics
	var responses []models.Metrics
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&req); err != nil {
		log.Println("ошибка декодирования JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	for _, m := range req {
		if m.ID == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if m.MType == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if m.MType == "gauge" {
			if m.Value != nil {
				useCase.UsecaseUpdateGauge(m.ID, *m.Value)
				key := useCase.UsecaseRead(m.MType, m.ID)
				*m.Value = key.(float64)
			} else {
				return
			}
		} else if m.MType == "counter" {
			if m.Delta != nil {
				useCase.UsecaseUpdateCounter(m.ID, *m.Delta)
				key := useCase.UsecaseRead(m.MType, m.ID)
				*m.Delta = key.(int64)
			} else {
				return
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

		responses = append(responses, res)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	if err := enc.Encode(responses); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
