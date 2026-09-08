# kryvora-node

Official client daemon for Kryvora Network nodes.

## Overview

Kryvora is a decentralized physical infrastructure network designed for distributed verification and worker coordination. The `kryvora-node` daemon connects worker hardware to the Kryvora Hub, executes telemetry probes, verifies peer state, and reports proof of availability.

## System Requirements

Hardware specs for running a standard verification worker:

* CPU: 2 cores minimum (x86_64 or arm64)
* Memory: 2 GB RAM minimum, 4 GB recommended
* Storage: 20 GB available SSD storage
* Network: Stable broadband connection with inbound port 4177 open for local telemetry probes

## Quick Start

### Option 1: Automated Script (Linux systemd)

Run the installation script to fetch the binary and register the systemd service:

```bash
curl -sSL https://raw.githubusercontent.com/kryvora-network/kryvora-node/main/install.sh | bash
```

Check service status:

```bash
sudo systemctl status kryvora-node
```

### Option 2: Docker Compose

1. Clone the repository:

```bash
git clone https://github.com/kryvora-network/kryvora-node.git
cd kryvora-node
```

2. Copy the default configuration:

```bash
cp config.example.yaml config.yaml
```

3. Start the node container:

```bash
docker compose up -d
```

4. Verify logs:

```bash
docker compose logs -f kryvora-node
```

### Option 3: Build from Source

Requirements: Go 1.22 or higher.

```bash
git clone https://github.com/kryvora-network/kryvora-node.git
cd kryvora-node
make build
./bin/kryvora-node --config config.example.yaml
```

## Configuration

Configuration is loaded from `config.yaml` or set via environment variables.

| Key | Type | Default | Description |
|---|---|---|---|
| `node.id` | string | auto-generated | Unique identifier for the worker node |
| `node.listen_addr` | string | `0.0.0.0:4177` | Local address for telemetry and metrics |
| `node.data_dir` | string | `/var/lib/kryvora` | Local state and task database directory |
| `hub.endpoint` | string | `https://hub.kryvora.network:4188` | Upstream coordination hub address |
| `hub.heartbeat_interval` | int | `30` | Interval in seconds between hub pings |
| `telemetry.enabled` | bool | `true` | Expose Prometheus metrics on `/metrics` |

## Local Endpoints

When running, the daemon exposes the following endpoints on port 4177:

* `GET /health`: JSON response indicating daemon lifecycle status
* `GET /status`: Peer connectivity, current block height, and sync state
* `GET /metrics`: Prometheus-compatible telemetry counters

Example query:

```bash
curl -s http://127.0.0.1:4177/status
```

## Security

Please report vulnerabilities directly to `security@kryvora.network` rather than opening public issues.

## License

Apache License 2.0. See LICENSE for details.
