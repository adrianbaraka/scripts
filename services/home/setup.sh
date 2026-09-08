#!/usr/bin/env bash


# Resolve the absolute path of the directory this script lives
SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

DNS_SCRIPT="$SCRIPT_DIR/dns-update.sh"
SERVICE_NAME="dns-update.service"
SERVICE_PATH="/etc/systemd/system/$SERVICE_NAME"

if [[ $EUID -ne 0 ]]; then
    echo "This script needs root to write to /etc/systemd/system and reload systemd." >&2
    echo "Re-run with: sudo $0" >&2
    exit 1
fi

if [[ ! -f "$DNS_SCRIPT" ]]; then
    echo "Expected to find dns-update.sh at: $DNS_SCRIPT" >&2
    echo "Make sure this setup script lives alongside dns-update.sh." >&2
    exit 1
fi

if [[ ! -x "$DNS_SCRIPT" ]]; then
    echo "Making $DNS_SCRIPT executable..."
    chmod u+x "$DNS_SCRIPT"
fi

echo "Writing $SERVICE_PATH..."
cat > "$SERVICE_PATH" <<EOF
[Unit]
Description=Update DNS-Exit record with Tailscale IP
After=network-online.target
Wants=network-online.target

[Service]
Type=oneshot
WorkingDirectory=$SCRIPT_DIR
ExecStart=$DNS_SCRIPT

[Install]
WantedBy=multi-user.target
EOF

echo "Reloading systemd daemon..."
systemctl daemon-reload

echo "Enabling $SERVICE_NAME to run at boot..."
systemctl enable "$SERVICE_NAME"

echo ""
echo "Done. $SERVICE_NAME is enabled and will run on next boot."
