package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Total number of HTTP requests processed, labeled by path, method, status.",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
)

func init() {
	// Register the metrics with the default registry
	prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
}

func Handler() http.Handler {
	// Handler exposes the Prometheus metrics HTTP handler.
	return promhttp.Handler()
}

// statusRecorder wraps http.ResponseWriter so we can capture the status code.
type statusRecorder struct {
	http.ResponseWriter 
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func InstrunebtHandler(path string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		start := time.Now()
		next.ServeHTTP(rec, r)
		elapsed := time.Since(start).Seconds()

		statusCode := rec.status
		method := r.Method

		httpRequestsTotal.WithLabelValues(path, method, strconv.Itoa(statusCode)).Inc()
		httpRequestDuration.WithLabelValues(path, method).Observe(elapsed)
	})
}