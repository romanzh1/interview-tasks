package handler

import (
	"log/slog"
	"net/http"
)

type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rec *statusRecorder) WriteHeader(status int) {
	if !rec.wroteHeader {
		rec.status = status
		rec.ResponseWriter.WriteHeader(status)
		rec.wroteHeader = true
	}
}

func (rec *statusRecorder) Write(b []byte) (int, error) {
	if !rec.wroteHeader {
		rec.WriteHeader(http.StatusOK)
	}
	return rec.ResponseWriter.Write(b)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w}

		defer func() {
			if rec.status >= 400 {
				slog.Error("HTTP request failed", "method", r.Method, "path", r.URL.Path, "status", rec.status)
			}
		}()

		next.ServeHTTP(rec, r)
	})
}
