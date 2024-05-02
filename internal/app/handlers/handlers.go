package handlers

import (
	"fmt"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/usecase"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
	"strconv"
)

func InitRoutes(useCase usecase.UsecaseMemStorage) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		handlerAllMetrics(w, r, useCase)
	})
	r.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		handlerPost(w, r, useCase)
	})
	r.Get("/value/{type}/{name}", func(w http.ResponseWriter, r *http.Request) {
		handlerGet(w, r, useCase)
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
	if typeMetric == "guage" {
		valueFloat, err := strconv.ParseFloat(valueMetric, 64)
		if err != nil {
			fmt.Println("Ошибка конвертации значения")
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		useCase.UsecaseUpdateGuage(nameMetric, valueFloat)
		res.WriteHeader(http.StatusOK)
	} else if typeMetric == "counter" {
		valueInt, err := strconv.ParseInt(valueMetric, 10, 64)
		if err != nil {
			fmt.Println("Ошибка конвертации значения")
			res.WriteHeader(http.StatusInternalServerError)
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
	stringValue := fmt.Sprintf("%v", value)
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte(stringValue))

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
}
