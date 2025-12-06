
# High-Throughput Request Benchmarking System  
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
![Status](https://img.shields.io/badge/status-active-success.svg)
![Contributions welcome](https://img.shields.io/badge/contributions-welcome-blue.svg)
![Requests](https://img.shields.io/badge/Requests-1M%2B%2Fsec-orange.svg)

A high-performance benchmarking framework written in **Go**, demonstrating how to **generate, send, and process up to one million HTTP REST requests per second (1M+ RPS)**.

The project includes:

- A **high-throughput HTTP server** in Go.
- A **load generator** in Go.
- Optional **Prometheus + Grafana** monitoring.
- Docker-based setup for reproducible experiments.

---

## 📚 Table of Contents

1. [Overview](#overview)  
2. [Features](#features)  
3. [Architecture](#architecture)  
4. [System Flow](#system-flow)  
5. [Tech Stack](#tech-stack)  
6. [Project Structure](#project-structure)  
7. [Setup & Installation](#setup--installation)  
8. [Running the Server](#running-the-server)  
9. [Running the Load Generator](#running-the-load-generator)  
10. [Configuration](#configuration)  
11. [Metrics & Dashboards](#metrics--dashboards)  
12. [Example REST Endpoint](#example-rest-endpoint)  
13. [Performance Notes](#performance-notes)  
14. [Contributing](#contributing)  
15. [License](#license)

---

## 🔭 Overview

This repository demonstrates:

- How to implement a **minimal, fast HTTP server** in Go that can handle massive concurrency.
- How to implement a **load generator** that pushes extremely high RPS using goroutines.
- How to **measure and visualize** throughput, latency and error rates.
- How to run all components using **Docker Compose**.

> ⚠️ This is a benchmarking / demo project. Do **not** point the load generator at services you do not own.

---

## ✨ Features

- 🚀 Generate up to **1,000,000+ HTTP requests per second** (hardware and OS dependent).
- ⚡ Go-based **HTTP server** using `net/http`.
- 🧵 Highly concurrent **load generator** using goroutines and a tuned HTTP client.
- 📊 Optional **Prometheus metrics**.
- 📈 Optional **Grafana dashboards**.
- 🐳 **Docker Compose** environment for server + monitoring.

---

## 🏗 Architecture

### High-Level Diagram

```mermaid
flowchart LR
    LG[Load Generator (Go)] -->|1M RPS| SRV(HTTP Server (Go))
    SRV --> H[Request Handler]
    H --> M[Metrics / Instrumentation]
    M --> P[Prometheus]
    P --> G[Grafana Dashboard]
```
</br>

## How to run locally
1. from project root:
```
go run ./server
```

Test if server running:
```
curl http://localhost:8080/health
```

2. Test ```/process``` manually
```
curl -X POST http://localhost:8080/process -d '{"id":1,"value":"test"}' -H "Content-Type: application/json"
```

3. Test the ```/metrics``` endpoint
```
curl http://localhost:9090/metrics
```

Now, let's run the Load Generator locally
in new terminal, from the root:
```
go run ./load-generator \
  -target http://localhost:8080/process \
  -rate 50000 \
  -connections 100 \
  -duration 10s
  ```
</br>
💡Note: </br>
Running 1M RPS locally on a laptop <b>will not work</b> — start with <b>5k</b> → <b>20k</b> → <b>50k</b> RPS.
</br>
</br>

### Verify the server is recieving load
#### Health check:
```
curl http://localhost:9090/metrics | grep http_requests_total
```


## Want to run both in Docker instead?
