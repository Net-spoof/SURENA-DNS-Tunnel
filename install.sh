#!/usr/bin/env bash
# DNSFabric: guided WhiteDNS installer for fresh VPS installations.
set -Eeuo pipefail
umask 077

if [[ ${EUID} -ne 0 ]]; then
  printf 'ERROR: Run as root (e.g. sudo bash dnsfabric-install.sh).\n' >&2
  exit 1
fi

for cmd in python3 curl mktemp; do
  command -v "$cmd" >/dev/null 2>&1 || {
    printf 'ERROR: Missing prerequisite: %s\n' "$cmd" >&2
    exit 1
  }
done

ROOT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
if [[ -f "$ROOT_DIR/installer.py" ]]; then
  exec python3 "$ROOT_DIR/installer.py" "$@"
fi

# For standalone install.sh downloads, fetch installer.py from the same repository.
WORK_DIR="$(mktemp -d /tmp/dnsfabric-installer.XXXXXXXX)"
trap 'rm -rf -- "$WORK_DIR"' EXIT

INSTALLER_URL="https://raw.githubusercontent.com/Net-spoof/dnsfabric/main/installer.py"
printf 'DNSFabric: fetching installer.py from Net-spoof/dnsfabric...\n'
if ! curl --fail --show-error --silent --location --retry 3 \
    --connect-timeout 20 --max-time 120 \
    --proto '=https' --proto-redir '=https' \
    --output "$WORK_DIR/installer.py" "$INSTALLER_URL"; then
  printf 'ERROR: Could not download DNSFabric installer. Check GitHub access / DNS / HTTPS.\n' >&2
  exit 1
fi

if [[ ! -s "$WORK_DIR/installer.py" ]]; then
  printf 'ERROR: Downloaded installer is empty.\n' >&2
  exit 1
fi

python3 "$WORK_DIR/installer.py" "$@"
