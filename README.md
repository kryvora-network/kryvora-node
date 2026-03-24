# kryvora-node

Official client daemon for Kryvora Network nodes.

## Overview

Kryvora is a decentralized physical infrastructure network designed for distributed verification and worker coordination. The `kryvora-node` daemon connects worker hardware to the Kryvora Hub, executes telemetry probes, verifies peer state, and reports proof of availability.

## System Requirements

Hardware specs for running a standard verification worker:

* CPU: 2 cores minimum (x86_64 or arm64)
* Memory: 2 GB RAM minimum, 4 GB recommended
* Storage: 20 GB available SSD storage
* Network: Stable connection with inbound port 4177 open for telemetry probes

## Quick Start

Requirements: Go 1.22 or higher.

```bash
git clone https://github.com/kryvora-network/kryvora-node.git
cd kryvora-node
make build
./bin/kryvora-node --config config.example.yaml
```

## Configuration

Configuration is loaded from `config.yaml` or set via environment variables.

* `node.listen_addr`: Local address for telemetry and metrics (default: `0.0.0.0:4177`)
* `node.data_dir`: Local state directory (default: `/var/lib/kryvora`)
* `hub.endpoint`: Upstream coordination hub address (default: `https://hub.kryvora.network:4188`)
* `hub.heartbeat_interval_sec`: Interval in seconds between hub pings (default: `30`)

## Local Endpoints

When running, the daemon exposes the following endpoints on port 4177:

* `GET /health`: Basic health probe
* `GET /status`: Node state and uptime summary
* `GET /metrics`: Prometheus telemetry counters

## License

Apache License 2.0. See LICENSE for details.
