#!/usr/bin/env bash
set -e

GO_VERSION="1.25.3"
GO_INSTALL_DIR="/usr/local"

_version_gte() {
  [ "$(printf '%s\n' "$1" "$2" | sort -V | head -1)" = "$2" ]
}

if command -v go &>/dev/null; then
  INSTALLED_VERSION="$(go version | awk '{print $3}' | sed 's/go//')"
  if _version_gte "$INSTALLED_VERSION" "$GO_VERSION"; then
    echo "    Go já instalado: go$INSTALLED_VERSION"
    exit 0
  fi
  echo "--> Versão instalada (go$INSTALLED_VERSION) é menor que a requerida (go$GO_VERSION). Atualizando..."
fi

OS_NAME="$(uname -s)"
case "$OS_NAME" in
  Linux)  TARBALL_OS="linux"  ;;
  Darwin) TARBALL_OS="darwin" ;;
  *)
    echo "ERRO: Sistema '$OS_NAME' não suportado por este script."
    exit 1
    ;;
esac

ARCH="$(uname -m)"
case "$ARCH" in
  x86_64)        ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *)
    echo "ERRO: Arquitetura '$ARCH' não suportada por este script."
    exit 1
    ;;
esac

TARBALL="go${GO_VERSION}.${TARBALL_OS}-${ARCH}.tar.gz"
TMP_DIR="$(mktemp -d)"

echo "--> Baixando Go $GO_VERSION..."
curl -fsSL "https://go.dev/dl/${TARBALL}" -o "$TMP_DIR/$TARBALL"

echo "--> Instalando em $GO_INSTALL_DIR..."
sudo rm -rf "$GO_INSTALL_DIR/go"
sudo tar -C "$GO_INSTALL_DIR" -xzf "$TMP_DIR/$TARBALL"
rm -rf "$TMP_DIR"

GO_PATH_LINE='export PATH=$PATH:/usr/local/go/bin'
if [ "$OS_NAME" = "Darwin" ]; then
  PROFILE_FILE="${ZDOTDIR:-$HOME}/.zshrc"
  [ -f "$HOME/.bash_profile" ] && [ ! -f "$HOME/.zshrc" ] && PROFILE_FILE="$HOME/.bash_profile"
else
  PROFILE_FILE="$HOME/.profile"
fi

if ! grep -qF "$GO_PATH_LINE" "$PROFILE_FILE" 2>/dev/null; then
  echo "$GO_PATH_LINE" >> "$PROFILE_FILE"
fi

export PATH="$PATH:/usr/local/go/bin"

echo "    Go $GO_VERSION instalado. Reinicie o terminal ou execute: source $PROFILE_FILE"
