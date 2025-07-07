#!/bin/bash

# CodeQuest CLI installer script
set -e

REPO="crisecheverria/codequest"
BINARY_NAME="codequest"
INSTALL_DIR="/usr/local/bin"

# Detect platform
OS=$(uname -s)
ARCH=$(uname -m)

# Map to GitHub release naming
case $OS in
    "Linux")
        OS_NAME="linux"
        ;;
    "Darwin")
        OS_NAME="darwin"
        ;;
    *)
        echo "❌ Unsupported operating system: $OS"
        echo "Please download manually from: https://github.com/$REPO/releases/latest"
        exit 1
        ;;
esac

case $ARCH in
    "x86_64" | "amd64")
        ARCH_NAME="amd64"
        ;;
    "arm64" | "aarch64")
        ARCH_NAME="arm64"
        ;;
    *)
        echo "❌ Unsupported architecture: $ARCH"
        echo "Please download manually from: https://github.com/$REPO/releases/latest"
        exit 1
        ;;
esac

BINARY_URL="https://github.com/$REPO/releases/latest/download/$BINARY_NAME-$OS_NAME-$ARCH_NAME"

echo "🚀 Installing CodeQuest CLI..."
echo "Platform: $OS_NAME-$ARCH_NAME"
echo "URL: $BINARY_URL"
echo ""

# Download binary
echo "📥 Downloading..."
if command -v curl >/dev/null 2>&1; then
    curl -L "$BINARY_URL" -o "$BINARY_NAME"
elif command -v wget >/dev/null 2>&1; then
    wget "$BINARY_URL" -O "$BINARY_NAME"
else
    echo "❌ Neither curl nor wget found. Please install one of them."
    exit 1
fi

# Make executable
chmod +x "$BINARY_NAME"

# Install to system PATH
echo "📦 Installing to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY_NAME" "$INSTALL_DIR/"
else
    echo "🔐 Installing requires sudo permissions..."
    sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
fi

# Verify installation
if command -v codequest >/dev/null 2>&1; then
    echo "✅ Installation successful!"
    echo ""
    echo "Version installed:"
    codequest --version
    echo ""
    echo "Try it out:"
    echo "  codequest list"
    echo "  codequest fetch <challenge>"
else
    echo "❌ Installation failed. Binary not found in PATH."
    echo "You may need to add $INSTALL_DIR to your PATH."
    exit 1
fi