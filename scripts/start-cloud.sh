#!/usr/bin/env bash
set -euo pipefail
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RUN_DIR="${ROOT_DIR}/.run"
LOG_DIR="${RUN_DIR}/logs"
mkdir -p "${LOG_DIR}"
export NIMBUS_NACOS_HOST="${NIMBUS_NACOS_HOST:-127.0.0.1}"
export NIMBUS_NACOS_PORT="${NIMBUS_NACOS_PORT:-28858}"

services=("system:58081" "infra:58082" "pay:58085" "member:58087" "business:58090" "gateway:58080")
for entry in "${services[@]}"; do
  name="${entry%%:*}"
  port="${entry##*:}"
  binary="${ROOT_DIR}/backend/bin/nimbus-${name}"
  [[ -x "${binary}" ]] || { echo "Missing ${binary}; run make build first" >&2; exit 1; }
  nohup "${binary}" >"${LOG_DIR}/${name}.log" 2>&1 &
  echo "$!" >"${RUN_DIR}/${name}.pid"
  echo "[starting] ${name} pid=$! port=${port}"
done
