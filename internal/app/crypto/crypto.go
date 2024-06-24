package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/DoktorGhost/go-musthave-metrics-tpl/internal/app/config"
	"io"
	"net/http"
)

func СryptoMiddleware(conf *config.Config) func(http.Handler) http.Handler {
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

				clientHash := r.Header.Get("HashSHA256")
				serverHash := hashData(body, key)

				if clientHash != serverHash {
					http.Error(w, "Invalid hash", http.StatusBadRequest)
					return
				}

				// Передаем управление следующему обработчику
				next.ServeHTTP(w, r)

				// Вычисляем и добавляем хеш в заголовок ответа
				responseHash := hashData([]byte("Response data"), key)
				w.Header().Set("HashSHA256", responseHash)
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
