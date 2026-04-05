#!/bin/sh
set -eu

mkdir -p /app/storage/avatar /app/storage/cover /app/storage/video

if [ ! -f /app/storage/avatar/default.svg ] && [ -f /app/defaults/avatar/default.svg ]; then
  cp /app/defaults/avatar/default.svg /app/storage/avatar/default.svg
fi

exec /app/http
