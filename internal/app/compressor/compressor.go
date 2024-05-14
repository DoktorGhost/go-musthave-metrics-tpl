package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// GzipMiddleware сжимает запросы и ответы с поддержкой gzip для определенных типов контента.
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем поддержку клиентом сжатия ответов
		acceptEncoding := r.Header.Get("Accept-Encoding")
		if strings.Contains(acceptEncoding, "gzip") {
			// Проверяем тип контента запроса
			contentType := r.Header.Get("Content-Type")
			if contentType == "application/json" || contentType == "text/html" {
				// Добавляем заголовок Content-Encoding для указания, что ответ будет сжат
				w.Header().Set("Content-Encoding", "gzip")

				// Создаем сжатый Writer для записи ответа
				gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
				if err != nil {
					io.WriteString(w, err.Error())
					return
				}
				defer gz.Close()
				// Создаем ResponseWriter с оберткой для сжатия данных
				gzWriter := &gzipWriter{Writer: gz, ResponseWriter: w}
				w = gzWriter
			}
		}

		// Продолжаем обработку запроса в следующем обработчике
		next.ServeHTTP(w, r)
	})
}

// gzipResponseWriter обертывает http.ResponseWriter для поддержки сжатия данных.
type gzipWriter struct {
	io.Writer
	http.ResponseWriter
}

// Write перенаправляет запись в сжатый Writer.
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func DecompressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, содержит ли заголовок Content-Encoding информацию о сжатии
		contentEncoding := r.Header.Get("Content-Encoding")
		if contentEncoding == "gzip" {
			// Декомпрессия тела запроса
			zr, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress request body", http.StatusBadRequest)
				return
			}
			defer zr.Close()

			// Создаем новый запрос с декомпрессированными данными
			r.Body = http.MaxBytesReader(w, zr, 10<<20) // Устанавливаем максимальный размер тела запроса
		}

		// Продолжаем обработку запроса в следующем обработчике
		next.ServeHTTP(w, r)
	})
}
