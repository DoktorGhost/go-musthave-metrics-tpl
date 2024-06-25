package crypto

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/config"
	"io"
	"net/http"
)

func CryptoMiddleware(conf *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := conf.SecretKey

			if key != "" {
				// Выполняем подпись
				body, err := io.ReadAll(r.Body)
				if err != nil {
					http.Error(w, "Unable to read body", http.StatusInternalServerError)
					return
				}
				// Восстанавливаем тело запроса
				r.Body = io.NopCloser(bytes.NewBuffer(body))

				clientHash := r.Header.Get("HashSHA256")
				serverHash := hashData(body, key)

				if clientHash != serverHash {
					http.Error(w, "Invalid hash", http.StatusBadRequest)
					return
				}

				// Создаем записывающий прокси для тела ответа
				rec := newResponseRecorder(w)
				next.ServeHTTP(rec, r)

				// Вычисляем и добавляем хеш в заголовок ответа
				responseHash := hashData(rec.body.Bytes(), key)
				w.Header().Set("HashSHA256", responseHash)

				// Копируем статус и заголовки ответа
				w.WriteHeader(rec.statusCode)
				rec.body.WriteTo(w)
			} else {
				// Пропускаем проверку и подпись, если ключ не задан
				next.ServeHTTP(w, r)
			}
		})
	}
}

func hashData(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		body:           new(bytes.Buffer),
	}
}

func (rec *responseRecorder) WriteHeader(code int) {
	rec.statusCode = code
}

func (rec *responseRecorder) Write(b []byte) (int, error) {
	rec.body.Write(b)
	return rec.ResponseWriter.Write(b)
}
