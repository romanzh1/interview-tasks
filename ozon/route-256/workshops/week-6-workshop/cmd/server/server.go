package main

import (
	"context"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.ozon.dev/12/week-6-workshop/internal/pkg/cache"
	logger_custom "gitlab.ozon.dev/12/week-6-workshop/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	resource2 "go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	trace2 "go.opentelemetry.io/otel/trace"
	"io"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var (
	requestCounter = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "app",
			Name:      "request_total_counter",
			Help:      "Total amount of request ",
		},
		[]string{"handler"},
	)

	requestHistogram = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "app",
			Name:      "request_duration_histogram",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"handler"})
)

func main() {
	rootCtx := context.Background()
	_, err := logger_custom.New()
	if err != nil {
		panic(err)
	}

	logger_custom.Infow(rootCtx, "server starting")

	exporter, err := otlptracehttp.New(rootCtx, otlptracehttp.WithEndpointURL("http://jaeger:4318"))
	if err != nil {
		logger_custom.Panicw(rootCtx, "otel exporter error")
	}

	resource, err := resource2.Merge(
		resource2.Default(),
		resource2.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("workshop-6"),
			semconv.DeploymentEnvironment("development"),
		),
	)
	if err != nil {
		logger_custom.Panicw(rootCtx, "creating resource return error")
	}

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource),
	)

	defer func() {
		traceProvider.Shutdown(rootCtx)
	}()

	otel.SetTracerProvider(traceProvider)

	tracer := otel.GetTracerProvider().Tracer("workshop-6")

	cacheService := cache.NewCache(1 << 10)

	http.Handle("GET /metrics", promhttp.Handler())

	http.HandleFunc("POST /set", func(w http.ResponseWriter, r *http.Request) {
		defer func(createdAt time.Time) {
			requestHistogram.WithLabelValues("set").Observe(time.Since(createdAt).Seconds())
		}(time.Now())

		ctx := r.Context()

		ctx, span := tracer.Start(ctx, "handler_set", trace2.WithAttributes(attribute.String("one", "two")))
		defer span.End()

		requestCounter.WithLabelValues("set").Inc()

		var request struct {
			Key   string
			Value string
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		err = json.Unmarshal(body, &request)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		if request.Key == "" || request.Value == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)

			return
		}

		logger_custom.Infow(ctx, "set-cache is used", "key", request.Key, "value", request.Value)

		err = cacheService.Set(ctx, request.Key, []byte(request.Value))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	})

	http.HandleFunc("GET /get", func(w http.ResponseWriter, r *http.Request) {
		defer func(createdAt time.Time) {
			requestHistogram.WithLabelValues("get").Observe(time.Since(createdAt).Seconds())
		}(time.Now())

		ctx := r.Context()

		ctx, span := tracer.Start(ctx, "handler_get", trace2.WithAttributes(attribute.Int("1", 2)))
		defer span.End()

		randomMillisecond := rand.Intn(2000)
		time.Sleep(time.Duration(randomMillisecond) * time.Millisecond)

		span.AddEvent("metric call")

		requestCounter.WithLabelValues("get").Inc()

		key := r.URL.Query().Get("key")
		if key == "" {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		l := logger_custom.With("key", key)
		ctx = logger_custom.ToContext(ctx, l)

		logger_custom.Infow(ctx, "get-cache is used")

		value, err := cacheService.Get(ctx, key)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = w.Write(value)
		if err != nil {
			logger_custom.Errorw(ctx, "get-cache is used", "error", err)
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger_custom.Panicw(rootCtx, "server can't start", "error", err)
	}
}
