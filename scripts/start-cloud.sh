#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
RUN_DIR="${ROOT_DIR}/.run"
LOG_DIR="${RUN_DIR}/logs"
mkdir -p "${LOG_DIR}"

"${ROOT_DIR}/scripts/init-local.sh"

export NIMBUS_NACOS_HOST="${NIMBUS_NACOS_HOST:-127.0.0.1}"
export NIMBUS_NACOS_PORT="${NIMBUS_NACOS_PORT:-28858}"
export NIMBUS_REDIS_ADDR="${NIMBUS_REDIS_ADDR:-127.0.0.1:27326}"
export NIMBUS_DB_DSN="${NIMBUS_DB_DSN:-nimbus:nimbus_dev@tcp(127.0.0.1:23326)/nimbus_cloud_go?charset=utf8mb4&parseTime=True&loc=Local}"

start_process() {
  local name="$1"
  local port="$2"
  local binary="${ROOT_DIR}/backend/bin/nimbus-${name}"
  local pid_file="${RUN_DIR}/${name}.pid"

  [[ -x "${binary}" ]] || {
    echo "Missing ${binary}; run 'cd backend && make build' first" >&2
    exit 1
  }
  if [[ -f "${pid_file}" ]] && kill -0 "$(<"${pid_file}")" 2>/dev/null; then
    echo "[running] ${name} pid=$(<"${pid_file}") port=${port}"
    return
  fi
  rm -f "${pid_file}"
  if nc -z 127.0.0.1 "${port}" 2>/dev/null; then
    echo "Port ${port} is already occupied; stop the existing process first" >&2
    exit 1
  fi
  nohup "${binary}" </dev/null >"${LOG_DIR}/${name}.log" 2>&1 &
  echo "$!" >"${pid_file}"
  echo "[starting] ${name} pid=$! port=${port}"
}

for entry in "infra:58082" "system:58081" "pay:58085" "member:58087" "business:58090" "gateway:58080"; do
  start_process "${entry%%:*}" "${entry##*:}"
done

for entry in "infra:58082" "system:58081" "pay:58085" "member:58087" "business:58090" "gateway:58080"; do
  name="${entry%%:*}"
  port="${entry##*:}"
  for _ in {1..60}; do
    curl -fsS "http://127.0.0.1:${port}/health" >/dev/null 2>&1 && break
    sleep 1
  done
  curl -fsS "http://127.0.0.1:${port}/health" >/dev/null || {
    echo "${name} failed to become healthy; see ${LOG_DIR}/${name}.log" >&2
    exit 1
  }
done

"${ROOT_DIR}/scripts/status-cloud.sh"
