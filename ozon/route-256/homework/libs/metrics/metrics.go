package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of requests",
		},
		[]string{"method", "status"},
	)
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of requests in seconds",
		},
		[]string{"method", "status"},
	)
	DBRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "db_requests_total",
			Help: "Total number of DB requests",
		},
		[]string{"query_type", "status"},
	)
	DBRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "db_request_duration_seconds",
			Help: "Duration of DB requests in seconds",
		},
		[]string{"query_type", "status"},
	)
	ExternalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "external_requests_total",
			Help: "Total number of external requests",
		},
		[]string{"url", "status"},
	)
	ExternalRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "external_request_duration_seconds",
			Help: "Duration of external requests in seconds",
		},
		[]string{"url", "status"},
	)
)

func init() {
	prometheus.MustRegister(TotalRequests, RequestDuration, DBRequests, DBRequestDuration, ExternalRequests, ExternalRequestDuration)
}

func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
