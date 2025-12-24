#!/usr/bin/env bash
set -euo pipefail

APP_NAME="pi-as-a-sensor"
BIN_NAME="pi-as-a-sensor"
INSTALL_DIR="/usr/local/bin"
SERVICE_DIR="/etc/systemd/system"
SERVICE_FILE="${SERVICE_DIR}/${APP_NAME}.service"

# Ensure root
if [[ $EUID -ne 0 ]]; then
  echo "Please run as root (use sudo)"
  exit 1
fi

# Build (assumes Go is installed)
GO_BIN="${GO_BIN:-}"
if [[ -z "${GO_BIN}" ]]; then
  if command -v go >/dev/null 2>&1; then
    GO_BIN="$(command -v go)"
  elif [[ -x /usr/local/go/bin/go ]]; then
    GO_BIN="/usr/local/go/bin/go"
  else
    echo "go not found. Ensure Go is installed and in PATH, or set GO_BIN=/path/to/go"
    exit 1
  fi
fi

echo "==> Building binary"
"${GO_BIN}" build -o "${BIN_NAME}" ./cmd/agent

echo "==> Installing ${APP_NAME}"

# Ask user for node identifier
read -rp "Enter server domain (e.g. pizero-1.local): " SERVER
if ! [[ "$SERVER" =~ ^[a-zA-Z0-9.-]+(:[0-9]+)?$ ]]; then
  echo "Invalid server address"
  exit 1
fi

read -rp "Use HTTPS? [y/N]: " USE_HTTPS
PROTO="http"
[[ "${USE_HTTPS}" =~ ^[Yy]$ ]] && PROTO="https"

# Install binary
echo "==> Installing binary to ${INSTALL_DIR}"
install -m 0755 "${BIN_NAME}" "${INSTALL_DIR}/${BIN_NAME}"

# Install systemd service
echo "==> Installing systemd service"
cat > "${SERVICE_FILE}" <<EOF
[Unit]
Description=Raspberry Pi Metrics Agent
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=${INSTALL_DIR}/${BIN_NAME}
Restart=always
RestartSec=5

Environment="ENDPOINT=${PROTO}://${SERVER}/api/measurements"
Environment="INTERVAL=1s"
Environment="TIMEOUT=5s"
Environment="SENSOR_NAME=system"

NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=full
ProtectHome=true

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd and enable service
echo "==> Enabling service"
systemctl daemon-reload
systemctl enable "${APP_NAME}.service"

echo "==> Starting service"
systemctl restart "${APP_NAME}.service"

echo "==> Checking status with: systemctl status ${APP_NAME}.service"
systemctl --no-pager status "${APP_NAME}.service" || true

echo "==> Done"
