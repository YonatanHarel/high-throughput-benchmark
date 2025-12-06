package main

import (
	"bytes"
	"flag"
	"log"
	"math"
	"net"
	"net/http"
	"sync"
	"time"
)

func main() {
	target := flag.String("target", "http://localhost:8080/process", "Target URL to send requests to")
	rate := flag.Int("rate", 1000000, "Number of requests per second")
	connections := flag.Int("connections", 2000, "Number of concurrent connections")
	duration := flag.Duration("duration", 30*time.Second, "Test duration")
	payload := flag.String("payload", "{}", "JSON payload for POST requests")

	flag.Parse()

	if *connections <= 0 {
		log.Fatalf("connections muxt be > 0")
	}
	if *rate <= 0 {
		log.Fatalf("rate must be > 0")
	}

	perWorkerRate := float64(*rate) / float64(*connections)
	log.Printf("Target: %s\n", *target)
	log.Printf("Total rate: %d RPS, connections: %d, per-worker: ~%.2f RPS, duration: %s\n", *rate, *connections, perWorkerRate, duration.String())

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:		*connections*2,
			MaxConnsPerHost:	*connections*2,
			IdleConnTimeout:	30 * time.Second,
			DisableCompression:	true,
			DisableKeepAlives:	false,
			ForceAttemptHTTP2: 	false,
			DialContext: (&net.Dialer{
				Timeout:   2 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		},
		Timeout: 5 * time.Second,
	}

	var wg sync.WaitGroup
	stop := time.Now().Add(*duration)

	for i := 0; i < *connections; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(client, *target, perWorkerRate, stop, []byte(*payload), id)	
		}(i)
	}
	wg.Wait()
	log.Println("Load test finished")
}

func worker(client *http.Client, target string, perWorkerRate float64, stop time.Time, payload []byte, workerID int) {
	// interval between requests for this worker
	interval := time.Duration(float64(time.Second) / perWorkerRate)
	if interval <= 0 {
		interval = 0
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for time.Now().Before(stop) {
		<-ticker.C

		req, err := http.NewRequest(http.MethodPost, target, bytes.NewReader(payload))
		if err != nil {
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		_ = resp.Body.Close()
	}
}

// helper to avoid NaNs if someone passes weird flags
func safeInterval(rate float64) time.Duration {
	if rate <= 0 {
		return time.Second
	}
	secPerReq := 1.0 / rate
	if math.IsInf(secPerReq, 0) || math.IsNaN(secPerReq) {
		return time.Second
	}
	return time.Duration(secPerReq * float64(time.Second))
}
