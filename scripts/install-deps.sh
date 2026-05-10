#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "==> Verificando dependências de desenvolvimento..."

for dep in "$SCRIPT_DIR/dependencies"/*.sh; do
  bash "$dep"
done

echo "==> Dependências OK."
