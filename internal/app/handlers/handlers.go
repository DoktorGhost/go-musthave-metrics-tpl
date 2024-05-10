package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/logger"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/models"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/usecase"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func InitRoutes(useCase usecase.UsecaseMemStorage) chi.Router {
	r := chi.NewRouter()

	r.Use(logger.WithLogging)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handlerAllMetrics(w, r, useCase)
	})
	r.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		handlerPost(w, r, useCase)
	})
	r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
		handlerJSONUpdate(w, r, useCase)
	})
	r.Get("/value/{type}/{name}", func(w http.ResponseWriter, r *http.Request) {
		handlerGet(w, r, useCase)
	})
	r.Post("/value", func(w http.ResponseWriter, r *http.Request) {
		handlerJSONValue(w, r, useCase)
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
	err = tmpl.Execute(res, metrics)
	if err != nil {
		http.Error(res, "Ошибка отображения HTML-страницы", http.StatusInternalServerError)
		return
	}
	res.WriteHeader(http.StatusOK)

}

func handlerJSONUpdate(w http.ResponseWriter, r *http.Request, useCase usecase.UsecaseMemStorage) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req models.Metrics
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

	if req.MType == "gauge" {
		useCase.UsecaseUpdateGauge(req.ID, *req.Value)
		key := useCase.UsecaseRead(req.MType, req.ID)
		*req.Value = key.(float64)
	} else if req.MType == "counter" {
		useCase.UsecaseUpdateCounter(req.ID, *req.Delta)
		key := useCase.UsecaseRead(req.MType, req.ID)
		*req.Delta += key.(int64)
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
	} else {
		if req.MType == "gauge" {
			*req.Value = value.(float64)
		} else if req.MType == "counter" {
			*req.Delta = value.(int64)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	enc := json.NewEncoder(w)
	if err := enc.Encode(req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
