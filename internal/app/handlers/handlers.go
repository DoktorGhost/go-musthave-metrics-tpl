package handlers

import (
	"fmt"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/usecase"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func InitRoutes(useCase usecase.UsecaseMemStorage) chi.Router {
	r := chi.NewRouter()

	r.Post("/update/{type}/{name}/{value}", func(w http.ResponseWriter, r *http.Request) {
		HandlerPost(w, r, useCase)
	})
	return r
}

func HandlerPost(res http.ResponseWriter, req *http.Request, useCase usecase.UsecaseMemStorage) {
	if req.Method != http.MethodPost {
		res.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if req.Header.Get("Content-Type") != "text/plain" {
		res.WriteHeader(http.StatusBadRequest)
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
