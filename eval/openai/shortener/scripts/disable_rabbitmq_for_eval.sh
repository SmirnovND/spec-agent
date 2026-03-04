#!/usr/bin/env bash
set -euo pipefail

: "${WORKSPACE:?WORKSPACE is required}"

TARGET="$WORKSPACE/cmd/server/main.go"
if [ ! -f "$TARGET" ]; then
  echo "skip: $TARGET not found"
  exit 0
fi

# Remove direct RabbitMQ import (unused after patch).
perl -0777 -i -pe 's@\n\s*"github\.com/SmirnovND/toolbox/pkg/rabbitmq"@@g' "$TARGET"

# Remove mandatory RabbitMQ init block for eval runs.
perl -0777 -i -pe 's@\n\s*// Инициализация RabbitMQ компонентов через контейнер \(управление жизненным циклом\).*?\n\s*log\.Println\("RabbitMQ initialized successfully"\)\n@\n\t// RabbitMQ init is disabled in local eval mode to speed up isolated checks.\n@sg' "$TARGET"

