# Bifrost

**Bifrost** is a Go-based service designed around the OpenTelemetry SDK to be a drop in replacement for OTEL Collector for processing metrics, with a focus on extensibility and reliability. It includes support for configuration via YAML, Docker-based deployment, and a modular architecture for metric handling. It provides an added custom MetricProcessor which core function is to be able to track Metric Names at ingestion and persist them to long term storage DB (default: Redis)

## 🧩 Features

- Configurable via `config.yaml`
- Metric ingestion and processing logic
- Dockerized for easy deployment

## 📦 Tech Stack

- **Language:** Go
- **Config:** YAML
- **Containerization:** Docker

## 🚀 Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) 1.20+
- [Docker](https://www.docker.com/) (for containerized usage)

### Installation

```bash
git clone https://github.com/rodolfo-mora/bifrost.git
cd bifrost
go build -o bifrost main.go
./bifrost --config=/path/to/config.yaml

