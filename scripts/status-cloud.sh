#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

for entry in "gateway:58080" "system:58081" "infra:58082" "pay:58085" "member:58087" "business:58090"; do
  name="${entry%%:*}"
  port="${entry##*:}"
  state="DOWN"
  pid="-"
  pid_file="${ROOT_DIR}/.run/${name}.pid"
  if [[ -f "${pid_file}" ]]; then
    pid="$(<"${pid_file}")"
  fi
  if curl -fsS "http://127.0.0.1:${port}/health" >/dev/null 2>&1; then
    state="UP"
  fi
  printf '%-10s port=%-5s pid=%-7s %s\n' "${name}" "${port}" "${pid}" "${state}"
done

mysql_state="DOWN"
if docker compose -f "${ROOT_DIR}/compose.yaml" exec -T mysql mysqladmin ping -h 127.0.0.1 -u nimbus -pnimbus_dev --silent >/dev/null 2>&1; then
  mysql_state="UP"
fi
printf '%-10s port=%-5s pid=%-7s %s\n' "mysql" "23326" "docker" "${mysql_state}"

redis_state="DOWN"
if docker compose -f "${ROOT_DIR}/compose.yaml" exec -T redis redis-cli ping 2>/dev/null | grep -q PONG; then
  redis_state="UP"
fi
printf '%-10s port=%-5s pid=%-7s %s\n' "redis" "27326" "docker" "${redis_state}"

nacos_state="DOWN"
if curl -fsS "http://127.0.0.1:28090/v3/console/health/readiness" >/dev/null 2>&1; then
  nacos_state="UP"
fi
printf '%-10s port=%-5s pid=%-7s %s\n' "nacos" "28858" "docker" "${nacos_state}"
