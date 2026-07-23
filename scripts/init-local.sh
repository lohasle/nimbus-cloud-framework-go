#!/usr/bin/env bash
set -euo pipefail
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${ROOT_DIR}"
docker compose up -d mysql redis nacos
until docker compose exec -T mysql mysqladmin ping -h 127.0.0.1 -u nimbus -pnimbus_dev --silent >/dev/null 2>&1; do sleep 1; done
until docker compose exec -T redis redis-cli ping 2>/dev/null | grep -q PONG; do sleep 1; done
until curl -fsS http://127.0.0.1:28090/v3/console/health/readiness >/dev/null 2>&1; do sleep 2; done
echo "MySQL 8.4, Redis and Nacos are ready"
