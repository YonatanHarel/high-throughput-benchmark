package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/YonatanHarel/high-throughput-benchmark/server/config"
	"github.com/YonatanHarel/high-throughput-benchmark/server/logging"
	"github.com/YonatanHarel/high-throughput-benchmark/server/metrics"
)

type ProcessRequest struct {
	ID int64 `json:"status"`
	Value interface{} `json:"value"`
}

type ProcessResponse struct {
	Status string `json:"status"`
	ReceivedAt int64 `json:"recieved_at"`
}

func processHandler(includeTimestamp bool, minimizeResponse bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		resp := ProcessResponse{
			Status: "ok",
		}

		if includeTimestamp && !minimizeResponse {
			resp.ReceivedAt = time.Now().Unix()
		}

		if minimizeResponse {
			// Just write "ok" to reduce JSON encoding overhead
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`"ok"`))
			return
		}

		_ = json.NewEncoder(w).Encode(resp)
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(([]byte("ok")))
}

func main() {
	cfg, err := config.Load("server/config/server.yaml")
	if err != nil {
		logging.Infof("Could load config, using defaults: %v", err)
		cfg = &config.Config{}
	}

	// Initialize logger
	logging.Init(logging.Config{
		Level: 				cfg.Logging.Level,
		Format: 			cfg.Logging.Format,
		DisableReqLogs: 	cfg.Logging.DisableReqLogs,
	})
	
	port := cfg.Server.Port
	if port == 0 {
		if p := os.Getenv("PORT"); p != "" {
			fmt.Sscanf(p, "%d", &port)
		} else {
			port = 8080
		}
	}

	mux := http.NewServeMux()
	// Handlers with metrics instrumentation
	ph := processHandler(cfg.Performance.IncludeTimestamp, cfg.Performance.MinimizeResponse)
	mux.Handle("/process", metrics.InstrumentHandler("/process", ph))
	mux.HandleFunc("/health", healthHandler)

	if cfg.Metrics.Enabled {
		// addr := fmt.Sprintf(":%d", cfg.Metrics.Port)
		path := cfg.Metrics.Path
		if path == "" {
			path = "/metrics"
		}
		mux.Handle(path, metrics.Handler())
		logging.Infof("Metrics available at %s", path)

		// Serve /metrics on a separate server/port for clarity
		// go func() {
		// 	metricsMux := http.NewServeMux()
		// 	metricsMux.Handle(path, metrics.Handler())
		// 	logging.Infof("Metrics server listening on %s%s", addr, path)
		// 	if err := http.ListenAndServe(addr, metricsMux); err != nil {
		// 		logging.Errorf("metrics server error: %v", err)
		// 	}
		// }()
	}

	addr := fmt.Sprintf(":%d", port)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
		// You can plug in timeouts from cfg.Server.* if you want:
		// ReadTimeout:  time.Duration(cfg.Server.ReadTimeoutMS) * time.Millisecond,
		// WriteTimeout: time.Duration(cfg.Server.WriteTimeoutMS) * time.Millisecond,
	}

	logging.Infof("HTTP server listening on %s", addr)
	if err := server.ListenAndServe(); err != nil {
		logging.Errorf("server error: %v", err)
	}
}