#!/usr/bin/env bash
set -e

echo "==> Configurando git hooks..."

pre-commit install --hook-type pre-commit --hook-type pre-push

echo "==> Hooks ativos: pre-commit (lint) e pre-push (testes)."
