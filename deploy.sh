#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$ROOT_DIR"

echo "[1/3] Pull latest base images..."
docker compose pull || true

echo "[2/3] Build images..."
docker compose build

echo "[3/3] Start containers..."
docker compose up -d

echo "Containers are up:"
docker compose ps
