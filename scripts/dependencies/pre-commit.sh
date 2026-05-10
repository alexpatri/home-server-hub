#!/usr/bin/env bash
set -e

if command -v pre-commit &>/dev/null; then
  echo "    pre-commit já instalado: $(pre-commit --version)"
  exit 0
fi

case "$(uname -s)" in
  Linux)
    # Instalado via pipx para evitar conflito com ambientes gerenciados pelo SO (PEP 668)
    if ! command -v pipx &>/dev/null; then
      echo "--> pipx não encontrado. Instalando via apt..."
      sudo apt-get install -y pipx -qq
      pipx ensurepath
      export PATH="$PATH:$HOME/.local/bin"
    fi
    echo "--> Instalando pre-commit via pipx..."
    pipx install pre-commit
    ;;
  Darwin)
    if ! command -v brew &>/dev/null; then
      echo "ERRO: Homebrew não encontrado. Instale em https://brew.sh e tente novamente."
      exit 1
    fi
    echo "--> Instalando pre-commit via Homebrew..."
    brew install pre-commit
    ;;
  *)
    echo "ERRO: Sistema não suportado por este script."
    exit 1
    ;;
esac

echo "    pre-commit instalado."
