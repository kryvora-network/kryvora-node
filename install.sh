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

echo "=========================================================="
echo " Kryvora Node Installer (${TARGET_ARCH})"
echo "=========================================================="

# Check for bootstrap token
BOOTSTRAP_TOKEN="${KRYVORA_BOOTSTRAP_TOKEN:-}"

if [ -z "${BOOTSTRAP_TOKEN}" ]; then
  if [ -t 0 ]; then
    echo ""
    echo "Operating a Kryvora node requires an active Genesis Node Key."
    echo "Contract: 0xBCf2cD12D1D37578fA7C69805fE77788e866BE1a (Arbitrum One)"
    echo "Activate your key at: https://node.kryvora.network"
    echo ""
    read -r -p "Enter your Node Bootstrap Token (or press Enter to configure later): " INPUT_TOKEN || true
    BOOTSTRAP_TOKEN="${INPUT_TOKEN:-}"
  fi
fi

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
  if [ -n "${BOOTSTRAP_TOKEN}" ]; then
    sed -i "s|bootstrap_token: \"\"|bootstrap_token: \"${BOOTSTRAP_TOKEN}\"|g" "${CONFIG_DIR}/config.yaml"
    echo "Bootstrap token written to ${CONFIG_DIR}/config.yaml"
  fi
fi

if [ -z "${BOOTSTRAP_TOKEN}" ]; then
  echo ""
  echo "=========================================================="
  echo " ACTION REQUIRED BEFORE STARTING NODE"
  echo "=========================================================="
  echo "Node requires a Genesis Node Key bootstrap token to operate."
  echo "1. Obtain your key at https://node.kryvora.network"
  echo "2. Edit ${CONFIG_DIR}/config.yaml and set 'auth.bootstrap_token'"
  echo "=========================================================="
  echo ""
fi

if [ -f "systemd/kryvora-node.service" ] && command -v systemctl >/dev/null 2>&1; then
  cp systemd/kryvora-node.service /etc/systemd/system/kryvora-node.service
  systemctl daemon-reload
  systemctl enable kryvora-node
  echo "Service installed and enabled. Start with: sudo systemctl start kryvora-node"
fi

echo "kryvora-node installation complete."
