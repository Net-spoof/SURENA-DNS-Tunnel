#!/usr/bin/env bash
# Guarded DNSFabric subscription portal updater for WhiteDNS v0.8.2 only.
set -Eeuo pipefail
umask 077
export LC_ALL=C

[[ $EUID -eq 0 ]] || { echo "ERROR: run as root" >&2; exit 1; }
for cmd in curl sha256sum tar systemctl mktemp install; do
  command -v "$cmd" >/dev/null || { echo "ERROR: missing $cmd" >&2; exit 1; }
done

BIN=/usr/local/bin/soar-panel
DB=/var/lib/soar-panel/panel.db
SERVICE=soar-panel.service
RELEASE=dnsfabric-subportal-v0.8.2.5
BASE="https://github.com/Net-spoof/dnsfabric/releases/download/$RELEASE"
[[ -x "$BIN" && -s "$DB" ]] || { echo "ERROR: WhiteDNS panel not installed" >&2; exit 1; }
[[ -f /var/lib/soar-panel/access.json ]] || { echo "ERROR: access settings unavailable" >&2; exit 1; }
systemctl is-active --quiet "$SERVICE" || { echo "ERROR: panel is not running" >&2; exit 1; }

VERSION="$("$BIN" version 2>&1 || true)"
if [[ "$VERSION" != *0.8.2* ]]; then
  echo "ERROR: updater requires WhiteDNS 0.8.2; installed: $VERSION" >&2
  exit 1
fi
EXECSTART="$(systemctl show "$SERVICE" -p ExecStart --value)"
[[ "$EXECSTART" == *"$BIN"* ]] || { echo "ERROR: unexpected panel binary path" >&2; exit 1; }

case "$(uname -m)" in
  x86_64) ARCH=amd64 ;;
  aarch64) ARCH=arm64 ;;
  *) echo "ERROR: unsupported architecture" >&2; exit 1 ;;
esac

if [[ $# -eq 1 && $1 == "--check" ]]; then
  echo "PASS: WhiteDNS $VERSION compatible on $ARCH; no changes"
  exit 0
fi
[[ $# -eq 0 ]] || { echo "Usage: bash update-subportal.sh [--check]" >&2; exit 1; }

WORK="$(mktemp -d /tmp/dnsfabric-subportal.XXXXXXXX)"
trap 'rm -rf -- "$WORK"' EXIT
FILE="dnsfabric-subportal-linux-$ARCH.tar.gz"

echo "===== DOWNLOAD RELEASE ====="
curl -fsSL --retry 3 --connect-timeout 20 --max-time 180 -o "$WORK/$FILE" "$BASE/$FILE"
curl -fsSL --retry 3 --connect-timeout 20 --max-time 90 -o "$WORK/SHA256SUMS" "$BASE/SHA256SUMS"
cd "$WORK"
awk -v file="$FILE" '$2==file {print; matched++} END {if(matched!=1) exit 1}' SHA256SUMS > TARGET.sha256
sha256sum --check --strict TARGET.sha256

mapfile -t ENTRIES < <(tar -tzf "$FILE")
if [[ ${#ENTRIES[@]} -ne 2 || ${ENTRIES[0]} != "soar-panel" || ${ENTRIES[1]} != "WHITE-DNS-LICENSE" ]]; then
  echo "ERROR: unexpected archive entries; abort" >&2
  exit 1
fi
tar -xzf "$FILE"
[[ -s soar-panel ]] || { echo "ERROR: patched binary missing" >&2; exit 1; }
chmod 0700 soar-panel
[[ "$(./soar-panel version 2>&1)" == *0.8.2* ]] || { echo "ERROR: wrong release" >&2; exit 1; }

STAMP="$(date -u +%Y%m%dT%H%M%SZ)"
BACKUP="/root/dnsfabric-subportal-backup-$STAMP"
mkdir -m 0700 -- "$BACKUP"
cp -a -- "$BIN" "$BACKUP/soar-panel.before"
cp -a -- /var/lib/soar-panel/access.json "$BACKUP/access.before.json"

STOPPED=0
DONE=0
rollback() {
  code=$?
  if (( STOPPED && ! DONE )); then
    echo "ERROR: update not verified. Restoring original binary..." >&2
    systemctl stop "$SERVICE" || true
    cp -a -- "$BACKUP/soar-panel.before" "$BIN" || true
    systemctl start "$SERVICE" || true
  fi
  rm -rf -- "$WORK"
  exit "$code"
}
trap rollback EXIT

echo "===== UPGRADE (ONE CONTROLLED RESTART) ====="
systemctl stop "$SERVICE"
STOPPED=1
cp -a -- "$DB" "$BACKUP/panel.before.db"
install -m 0755 -o root -g root "$WORK/soar-panel" "$BIN"
systemctl start "$SERVICE"
sleep 2
systemctl is-active --quiet "$SERVICE"
if systemctl is-failed --quiet "$SERVICE"; then
  echo "ERROR: panel failed after upgrade" >&2
  exit 1
fi
DONE=1
STOPPED=0
echo "PASS: DNSFabric subscriber portal updated; users and usage preserved."
echo "Legacy /sub/ remains raw. /s/ shows the portal in browsers and stays raw for apps."
echo "Subscriber portal now offers Add to App (copy subscription + guided in-app Import), alongside profile deep links and QR."
echo "Backup saved: $BACKUP"
