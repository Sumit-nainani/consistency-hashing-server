# Consistency Hashing

Consistency Hashing is a real-time distributed system simulation built with Kubernetes, Go, Python, gRPC, and WebSockets. It implements the High-Level Design (HLD) concept of consistent hashing, where keys are distributed among servers on a virtual ring, with efficient rebalancing (only $k/n$ keys move when a node changes).

---

## 🚀 Summary

This project demonstrates a real-time dashboard and backend system for consistent hashing:

* Keys are hashed and assigned to servers on a virtual ring.
* When servers are added or removed, minimal keys are rebalanced.
* Real-time pod creation/deletion and request handling are shown on a dashboard.

---

## 🛠 Tech Stack

* **Backend**: Go (Golang), gRPC, WebSocket (Gorilla)
* **Middleware**: gRPC server + Kubernetes PodWatcher + WebSocket event emitter
* **Frontend**: Django (serves static files), JS (renders real-time ring)
* **WebSocket Client**: Python Tornado
* **Visualization**: HTML, CSS, JavaScript
* **Orchestration**: Kubernetes (Kind), HPA, Prometheus, Metrics Server
* **Docker**: Distroless static binary images
* **Observability**: Prometheus custom metrics, /metrics endpoint

---

## 🌟 Key Features

* 🌀 **Consistent Hash Ring** logic with AddNode, RemoveNode, and GetNode
* 📈 **Prometheus Metrics Endpoint** (`/metrics`) to expose custom metrics:

  * Unique IP requests in the last 1 minute
  * CPU usage > 80%
* ⚖️ **Horizontal Pod Autoscaling (HPA)** based on metrics
* 🔄 **Real-time Updates** with WebSocket using Protocol Buffers → JSON → JavaScript
* ⚙️ **Kubernetes Watcher** to track pod changes & update hash ring accordingly
* 🎯 **Direct Pod-to-Pod Communication** via Pod IPs
* 💻 **gRPC Interface** to fetch initial data (instead of REST)
* 📊 **UI Dashboard** shows servers & clients on a ring with Docker-style logos
* 📦 **Distroless Docker Image** using Google static build
* 🔧 **Shell Scripts**:

  * `create.sh`: Sets up Kind cluster, venv, installs dependencies
  * `start.sh`: Opens terminal sessions and starts services sequentially

---

## 🧠 Architecture

![Architecture](path/to/your/architecture-image.png)

---

## 🧪 Setup Instructions

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/consistency-hashing.git
cd consistency-hashing
```

### 2. Install Kind

Ensure you have [Kind](https://kind.sigs.k8s.io/) installed locally.

### 3. Create Cluster

```bash
sh create.sh <cluster-name> <config.yaml>
```

This installs dependencies, creates a virtual environment, sets up the cluster, Prometheus, metrics-server, and all Kubernetes components.

### 4. Start Backend and UI

```bash
sh start.sh
```

This launches iTerm2 with three split panes:

* Go backend (`go run main.go`)
* Django + WebSocket client (`python server.py`)
* HPA trigger via `kubectl scale`

**Note**: Modify `start.sh` AppleScript for your OS/terminal if not using macOS/iTerm2.

### 5. Start Prometheus

```bash
kubectl port-forward -n monitoring svc/prometheus-operated 9090
```

Access: [http://localhost:9090](http://localhost:9090)

### 6. Visualize

* Open [http://localhost:8085](http://localhost:8085)
* Send curl requests from different IPs to test autoscaling
* Run:

```bash
kubectl get hpa -n demo
kubectl get pods -n demo
```

To observe scaling behavior

---

## 📂 Directory Structure

```
├── create.sh
├── start.sh
├── backend/
│   ├── main.go
│   └── websocket.go
├── middleware/
│   └── grpc_server.go
├── frontend/
│   ├── server.py
│   ├── templates/index.html
│   └── static/
├── proto/
│   └── message.proto
├── k8s/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── hpa.yaml
```

---

## 🤝 Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

---

## 📄 License

[MIT](LICENSE)

---

## 📬 Contact

For any questions or feedback, feel free to reach out!

> Built with ❤️ by Sumit Nainani
