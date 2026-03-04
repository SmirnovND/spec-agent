#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/../.." && pwd)"
SRC_FILE="$ROOT_DIR/eval/fixtures/draft/source.yaml"
TARGET_DIR="$ROOT_DIR/eval/fixtures/draft/repo"

repo=$(awk -F': ' '/^repo:/{print $2}' "$SRC_FILE")
ref=$(awk -F': ' '/^ref:/{print $2}' "$SRC_FILE")

rm -rf "$TARGET_DIR"
git clone "$repo" "$TARGET_DIR"
(
  cd "$TARGET_DIR"
  git checkout "$ref"
)

echo "draft fixture synced to $TARGET_DIR (ref=$ref)"
