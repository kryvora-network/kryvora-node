#!/usr/bin/env bash
set -euo pipefail

ARCH="$(uname -m)"
case "${ARCH}" in
  x86_64) TARGET_ARCH="amd64" ;;
  aarch64|arm64) TARGET_ARCH="arm64" ;;
  *) echo "Unsupported architecture: ${ARCH}"; exit 1 ;;
esac

INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/kryvora"
DATA_DIR="/var/lib/kryvora"

echo "Installing kryvora-node for ${TARGET_ARCH}..."
mkdir -p "${CONFIG_DIR}" "${DATA_DIR}"

if [ -f "./bin/kryvora-node" ]; then
  install -m 0755 ./bin/kryvora-node "${INSTALL_DIR}/kryvora-node"
elif command -v go >/dev/null 2>&1; then
  go build -o "${INSTALL_DIR}/kryvora-node" ./cmd/kryvora-node
else
  echo "Please build the binary or use Docker."
  exit 1
fi

if [ ! -f "${CONFIG_DIR}/config.yaml" ]; then
  cp config.example.yaml "${CONFIG_DIR}/config.yaml"
fi

if [ -f "systemd/kryvora-node.service" ] && command -v systemctl >/dev/null 2>&1; then
  cp systemd/kryvora-node.service /etc/systemd/system/kryvora-node.service
  systemctl daemon-reload
  systemctl enable kryvora-node
  echo "Service installed and enabled. Start with: sudo systemctl start kryvora-node"
fi

echo "kryvora-node installation complete."
