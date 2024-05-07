package handlers

import (
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/logger"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/storage/maps"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string) (*http.Response, string) {

	// Создаем кастомный клиент
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Отключаем автоматический редирект
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)

}

func TestRoute(t *testing.T) {
	db := maps.NewMapStorage()
	storage := usecase.NewUsecaseMemStorage(db)

	logg, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logg.Sync()
	logger.InitLogger(logg)
	sugar := *logg.Sugar()
	sugar.Infow("server started")

	//добавим в бд тестовую запись
	ts := httptest.NewServer(InitRoutes(*storage))
	storage.UsecaseUpdateGauge("Allock", 100)
	defer ts.Close()

	type values struct {
		url    string
		method string
	}

	type want struct {
		status int
		body   string
	}

	var tests = []struct {
		name   string
		values values
		want   want
	}{
		{
			name: "Test #1 Ошибка метод для вывода всех метрик",
			values: values{
				url:    "/",
				method: "POST",
			},
			want: want{
				status: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "Test #2 Метод GET для вывода всех метрик",
			values: values{
				url:    "/",
				method: "GET",
			},
			want: want{
				status: http.StatusOK,
			},
		},
		{
			name: "Test #3 не тот метод handlerPost",
			values: values{
				url:    "/update/gauge/name/112",
				method: "GET",
			},
			want: want{
				status: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "Test #4 не валидный value в запросе handlerPost",
			values: values{
				url:    "/update/gauge/name/value",
				method: "POST",
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "Test #5 пустой name в запросе handlerPost",
			values: values{
				url:    "/update/gauge//112",
				method: "POST",
			},
			want: want{
				status: http.StatusNotFound,
			},
		},
		{
			name: "Test #6 запрос gauge handlerPost",
			values: values{
				url:    "/update/gauge/Allock/112",
				method: "POST",
			},
			want: want{
				status: http.StatusOK,
			},
		},
		{
			name: "Test #7 неизвестный тип handlerPost",
			values: values{
				url:    "/update/gaga/Allock/112",
				method: "POST",
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "Test #8 запрос counter handlerPost",
			values: values{
				url:    "/update/counter/Allock/112",
				method: "POST",
			},
			want: want{
				status: http.StatusOK,
			},
		},
		{
			name: "Test #9 пустой тип  handlerPost",
			values: values{
				url:    "/update//Allock/112",
				method: "POST",
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "Test #10 не тот метод handlerGet",
			values: values{
				url:    "/value/counter/Allock",
				method: "POST",
			},
			want: want{
				status: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "Test #11 пустой тип handlerGet",
			values: values{
				url:    "/value//Allock",
				method: "GET",
			},
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "Test #12 пустой name handlerGet",
			values: values{
				url:    "/value/counter",
				method: "GET",
			},
			want: want{
				status: http.StatusNotFound,
			},
		},
		{
			name: "Test #13 handlerGet",
			values: values{
				url:    "/value/gauge/Allock",
				method: "GET",
			},
			want: want{
				status: http.StatusOK,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, test.values.method, test.values.url)
			defer resp.Body.Close()
			assert.Equal(t, test.want.status, resp.StatusCode)

		})
	}
}
