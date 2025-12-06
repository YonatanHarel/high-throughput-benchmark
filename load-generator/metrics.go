package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	loadgenRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "loadgen_requests_total",
			Help: "Total number of requests sent by the load generator, labeled by result.",
		},
		[]string{"result"}, // "success" | "error"
	)

	loadgenRequestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "loadgen_request_duration_seconds",
			Help:    "Latency of requests from the load generator point of view.",
			Buckets: prometheus.DefBuckets,
		},
	)

	loadgenConfiguredRate = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "loadgen_rate_configured",
			Help: "Configured total requests per second.",
		},
	)

	loadgenWorkers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "loadgen_workers",
			Help: "Number of worker goroutines.",
		},
	)
)

func init() {
	prometheus.MustRegister(
		loadgenRequestsTotal,
		loadgenRequestDuration,
		loadgenConfiguredRate,
		loadgenWorkers,
	)
}

// metricsHandler returns the /metrics HTTP handler.
func metricsHandler() http.Handler {
	return promhttp.Handler()
}
