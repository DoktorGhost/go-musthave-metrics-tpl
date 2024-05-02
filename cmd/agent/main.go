package main

import (
	"fmt"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/metrics"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func main() {
	client := &http.Client{}

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
			var endpoints []string
			endpoints = append(endpoints, "http://localhost:8080/update/counter/PollCount/"+strconv.FormatInt(m.Counter, 10))
			for key, value := range m.Guage {
				endpoint := "http://localhost:8080/update/guage/" + key + "/" + strconv.FormatFloat(value, 'f', -1, 64)
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
				// выводим код ответа
				//fmt.Println("Статус-код ", response.Status)
				defer response.Body.Close()
			}
		}
	}()

	wg.Wait()
}
