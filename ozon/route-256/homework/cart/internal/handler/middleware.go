package handler

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/trace"

	"route256/libs/metrics"
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

func (h *Handler) loggingAndObserveMiddleware(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx, span := h.tracer.Start(r.Context(), r.URL.Path)
		defer span.End()

		traceID := trace.SpanContextFromContext(ctx).TraceID().String()
		spanID := trace.SpanContextFromContext(ctx).SpanID().String()

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		defer func() {
			duration := time.Since(start).Seconds()
			statusCode := strconv.Itoa(rec.status)

			metrics.RequestDuration.WithLabelValues(r.URL.Path, statusCode).Observe(duration)
			metrics.TotalRequests.WithLabelValues(r.URL.Path, statusCode).Inc()

			logTime := time.Now().Format(time.RFC3339)

			if rec.status >= 400 {
				slog.Error("HTTP request failed",
					"service", "cart",
					"time", logTime,
					"method", r.Method,
					"path", r.URL.Path,
					"status", rec.status,
					"trace_id", traceID,
					"span_id", spanID)
			} else {
				slog.Info("HTTP request succeeded",
					"service", "cart",
					"time", logTime,
					"method", r.Method,
					"path", r.URL.Path,
					"status", rec.status,
					"duration", duration,
					"trace_id", traceID,
					"span_id", spanID)
			}
		}()

		handlerFunc(rec, r.WithContext(ctx))
	}
}
