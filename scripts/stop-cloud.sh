#!/usr/bin/env bash
set -euo pipefail
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
for file in "${ROOT_DIR}"/.run/*.pid; do
  [[ -f "${file}" ]] || continue
  pid="$(<"${file}")"
  kill "${pid}" 2>/dev/null || true
  rm -f "${file}"
done

