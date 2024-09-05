package internalhttp

import (
	"log/slog"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		// Сохраним оригинальный ResponseWriter для получения статуса ответа
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		latency := time.Since(startTime)
		clientIP := r.RemoteAddr
		method := r.Method
		path := r.URL.Path
		proto := r.Proto
		status := ww.statusCode
		userAgent := r.UserAgent()

		logger := slog.Default()
		logger.Info("Request",
			"client_ip", clientIP,
			"time", startTime.Format(time.RFC1123),
			"method", method,
			"path", path,
			"proto", proto,
			"status", status,
			"latency", latency,
			"user_agent", userAgent,
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
