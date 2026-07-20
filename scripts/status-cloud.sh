#!/usr/bin/env bash
set -euo pipefail
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
for entry in "gateway:58080" "system:58081" "infra:58082" "pay:58085" "member:58087" "business:58090"; do
  name="${entry%%:*}"; port="${entry##*:}"; state="DOWN"
  if [[ -f "${ROOT_DIR}/.run/${name}.pid" ]] && kill -0 "$(<"${ROOT_DIR}/.run/${name}.pid")" 2>/dev/null && nc -z 127.0.0.1 "${port}" 2>/dev/null; then state="UP"; fi
  printf '%-10s port=%s %s\n' "${name}" "${port}" "${state}"
done

