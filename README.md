# ⚡ Smart Adaptive Load Balancer

**Smart Adaptive Load Balancer** is a Go-based load balancing system that supports **Layer 7 (HTTP)** and **Layer 4 (TCP)** traffic distribution using **traditional routing algorithms**. It includes **health checks, live metrics**, and a **real-time visualization dashboard** to observe routing behavior under load.

---

## 🧰 Built With

* **Go (Golang)** – core backend & networking
* **net/http & TCP sockets** – L7/L4 proxies
* **HTML, CSS, JavaScript (Canvas API)** – live dashboard
* **Traditional Algorithms** – Round Robin, Least Connections, Random, Weighted Round Robin

---

## ✨ Features

* ✅ **L7 HTTP Reverse Proxy**
* ✅ **L4 TCP Load Balancer**
* 🔁 **Routing Algorithms**

  * Round Robin
  * Least Connections
  * Random
* ❤️ **Active Health Checks** with auto-failover
* 📊 **Metrics API** for backend stats
* 🖥️ **Live Traffic Visualization Dashboard**
* ⚡ **Traffic Simulation** (animated requests)
* 🛠️ **Admin API** to dynamically add backends

---

## 🧠 Architecture Overview

```
Client
  ↓
Smart Load Balancer (L4 / L7)
  ↓
Routing Algorithm
  ↓
Healthy Backend Servers
```

The system continuously monitors backend health, tracks active connections, and adapts routing decisions in real time.

---

## 🚀 Getting Started

Follow these steps to set up and run the Smart Adaptive Load Balancer on your local machine.

---

## 📦 Prerequisites

Before you begin, ensure you have the following installed:

* **Go 1.20+**
* **Git**
* A modern web browser (Chrome / Firefox)

Verify Go installation:

```bash
go version
```

---

## 📥 Cloning the Repository

Open a terminal and run:

```bash
git clone https://github.com/sbirmecha99/smart-adaptive-load-balancer.git
cd smart-adaptive-load-balancer
```

---

## ▶️ Running the Load Balancer

### 🔹 L7 Mode (Default – HTTP)

```bash
go run cmd/balancer/main.go
```

Expected output:

```text
Starting L7 HTTP Load Balancer on :8080
```

---

### 🔹 L4 Mode (TCP)

```bash
LB_MODE=L4 go run cmd/main.go
```

---

## 🖥️ Running the Visualization Dashboard

1. Navigate to the dashboard folder:

```
dashboard/
```

2. Open `index.html` and `simulation.html` in your browser

The dashboard will:

* Fetch live metrics from `http://localhost:8080/metrics`
* Fetch live status from `http://localhost:8080/status`
* Display backend health (Alive / Down)
* Animate request routing based on the selected algorithm

---

## 📊 API Endpoints

### 📈 Metrics Endpoint

```
GET /metrics
```

Returns real-time backend metrics:

* Alive status
* Active connections
* Latency
* Error count

---

### ➕ Add Backend Server

```
POST /admin/add
```

Dynamically adds a backend server to the pool.

---


## 📁 Project Structure

```
smart-adaptive-load-balancer/
│
├── cmd/balancer
│   └── main.go
│
├── internal/
│   ├── api/          # Metrics & admin handlers
│   ├── core/         # Backend model
│   ├── health/       # Health checker
│   ├── proxy/
│   │   ├── l4/       # TCP proxy
│   │   └── l7/       # HTTP proxy
│   └── routing/     # Routing algorithms
│
├── dashboard/
│   ├── index.html
│   ├── style.css
│   └── script.js
│
└── README.md
```


