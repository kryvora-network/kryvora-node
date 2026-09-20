# kryvora-node

[![Go Version](https://img.shields.io/badge/go-%3E%3D%201.22-30363d.svg?style=flat-square&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-Apache--2.0-30363d.svg?style=flat-square)](LICENSE)
[![Release](https://img.shields.io/github/v/release/kryvora-network/kryvora-node?style=flat-square&color=30363d&label=release)](https://github.com/kryvora-network/kryvora-node/releases)
[![CI Status](https://img.shields.io/github/actions/workflow/status/kryvora-network/kryvora-node/ci.yml?branch=main&style=flat-square&label=ci&color=30363d)](https://github.com/kryvora-network/kryvora-node/actions)
[![Network](https://img.shields.io/badge/network-Kryvora%20DePIN-30363d.svg?style=flat-square)](https://kryvora.network)

Official client daemon for Kryvora Network nodes.

## Overview

Kryvora is a decentralized physical infrastructure network designed for distributed verification and worker coordination. The `kryvora-node` daemon connects worker hardware to the Kryvora Hub, executes telemetry probes, verifies peer state, and reports proof of availability.

## System Topology

```
+-------------------------------------------------------------+
|                     Kryvora Node Host                       |
|                                                             |
|  +--------------------+             +--------------------+  |
|  |   Worker Engine    |             |  Telemetry Daemon  |  |
|  |  (Task Execution)  |<----------->|    (Port 4177)     |  |
|  +---------+----------+   IPC/Bus   +---------+----------+  |
|            |                                  |             |
|            | Unix Socket (/var/run/node.sock) | HTTP API    |
|            v                                  v             |
|  +-------------------------------------------------------+  |
|  |                 Kryvora Storage Store                 |  |
|  |            (State, Peer Nonces, Telemetry)            |  |
|  +-------------------------------------------------------+  |
+------------------------------+------------------------------+
                               |
                   Mutual TLS  |  Heartbeat (30s)
                   gRPC / HTTP |  Verification Probes
                               v
+-------------------------------------------------------------+
|                 Kryvora Coordination Hub                    |
|                        (Port 4188)                          |
+-------------------------------------------------------------+
```

## Node Licensing & Genesis Key Requirement

Operating a Kryvora verification node requires an active **Genesis Node Key**.
The Genesis Node Key is an on-chain soulbound verification credential deployed on Arbitrum One:

* **Contract**: `0xBCf2cD12D1D37578fA7C69805fE77788e866BE1a`
* **Network**: Arbitrum One (Chain ID: `42161`)
* **Portal**: [https://node.kryvora.network](https://node.kryvora.network)

Before launching the daemon, obtain your node bootstrap token by connecting your wallet at [node.kryvora.network](https://node.kryvora.network). Unauthenticated nodes without a valid bootstrap token will be rejected by the Coordination Hub (`hub.kryvora.network:4188`).

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

The script will prompt for your Genesis Node Key bootstrap token and register the systemd service.

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

2. Copy the configuration template and set your token:

```bash
cp config.example.yaml config.yaml
# Edit config.yaml and insert your bootstrap_token under auth:
# Or set KRYVORA_BOOTSTRAP_TOKEN in your environment
```

3. Start the node container:

```bash
KRYVORA_BOOTSTRAP_TOKEN="your_token_here" docker compose up -d
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
./bin/kryvora-node --config config.yaml --token "your_token_here"
```

## Configuration

Configuration is loaded from `config.yaml` or set via environment variables.

| Key | Type | Default | Description |
|---|---|---|---|
| `node.id` | string | auto-generated | Unique identifier for worker node |
| `node.listen_addr` | string | `0.0.0.0:4177` | Local address for telemetry and metrics |
| `node.data_dir` | string | `/var/lib/kryvora` | Local state and task database directory |
| `auth.bootstrap_token` | string | `""` | Mandatory: Genesis Node Key activation token |
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
