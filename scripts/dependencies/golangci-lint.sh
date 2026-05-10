#!/usr/bin/env bash
set -e

if ! command -v golangci-lint &>/dev/null; then
  echo "--> Instalando golangci-lint..."
  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
  echo "    golangci-lint instalado."
else
  echo "    golangci-lint já instalado: $(golangci-lint --version 2>&1 | head -1)"
fi
