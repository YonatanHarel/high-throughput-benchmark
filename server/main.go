package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type ProcessRequest struct {
	ID int64 `json:"status"`
	Value interface{} `json:"value"`
}

type ProcessResponse struct {
	Status string `json:"status"`
	RecievedAt int64 `json:"recieved_at"`
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contectt-Type", "application/json")
	resp := ProcessResponse{
		Status: "ok",
		RecievedAt: time.Now().Unix(),
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(([]byte("ok")))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/process", processHandler)
	mux.HandleFunc("/health", healthHandler)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Http server is starting on port %s", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}