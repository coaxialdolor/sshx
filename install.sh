#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

echo "SSHX Installer"
echo "=============="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed or not in PATH${NC}"
    echo "Please install Go 1.22+ from https://go.dev/dl/"
    exit 1
fi

# Build binaries
echo "Building binaries..."
cd "$SCRIPT_DIR"
if ! go build -o sshx ./cmd/sshx; then
    echo -e "${RED}Error: Failed to build sshx${NC}"
    exit 1
fi

if ! go build -o sshx-agent ./cmd/sshx-agent; then
    echo -e "${RED}Error: Failed to build sshx-agent${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Binaries built successfully${NC}"
echo ""

# Check if install directory is writable
if [ ! -w "$INSTALL_DIR" ]; then
    echo "Install directory $INSTALL_DIR is not writable."
    echo "You may need to run this script with sudo or set INSTALL_DIR to a writable directory."
    read -p "Continue with sudo? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        SUDO="sudo"
    else
        echo "Installation cancelled."
        exit 1
    fi
else
    SUDO=""
fi

# Install binaries
echo "Installing binaries to $INSTALL_DIR..."
$SUDO cp sshx "$INSTALL_DIR/sshx"
$SUDO cp sshx-agent "$INSTALL_DIR/sshx-agent"
$SUDO chmod +x "$INSTALL_DIR/sshx"
$SUDO chmod +x "$INSTALL_DIR/sshx-agent"

echo -e "${GREEN}✓ Binaries installed${NC}"
echo ""

# Ask about alias
echo "SSHX can enhance your SSH experience by making the 'ssh' command automatically use sshx."
echo ""
echo "This does NOT replace or modify your system SSH binary."
echo "It simply adds a safe alias in your shell profile."
echo ""
read -p "Enable this feature? (1) Yes — make 'ssh' launch sshx (2) No — keep ssh unchanged [1/2]: " choice

if [ "$choice" = "1" ]; then
    # Detect shell
    SHELL_NAME=$(basename "$SHELL")

    case "$SHELL_NAME" in
        zsh)
            PROFILE="$HOME/.zshrc"
            ;;
        bash)
            if [ -f "$HOME/.bash_profile" ]; then
                PROFILE="$HOME/.bash_profile"
            else
                PROFILE="$HOME/.bashrc"
            fi
            ;;
        fish)
            PROFILE="$HOME/.config/fish/config.fish"
            mkdir -p "$HOME/.config/fish"
            ;;
        *)
            echo -e "${YELLOW}Warning: Unsupported shell ($SHELL_NAME). Skipping alias setup.${NC}"
            PROFILE=""
            ;;
    esac

    if [ -n "$PROFILE" ]; then
        # Check if alias already exists
        if grep -q 'alias ssh="sshx"' "$PROFILE" 2>/dev/null; then
            echo -e "${YELLOW}Alias already exists in $PROFILE${NC}"
        else
            echo "" >> "$PROFILE"
            echo "# SSHX alias" >> "$PROFILE"
            echo 'alias ssh="sshx"' >> "$PROFILE"
            echo -e "${GREEN}✓ Alias added to $PROFILE${NC}"
            echo "Run 'source $PROFILE' or restart your shell to use the alias."
        fi
    fi
else
    echo "Alias not added. You can add it manually later or run the installer again."
fi

echo ""
echo "=========================================="
echo "Installation complete!"
echo ""
echo "To start using SSHX:"
echo "  sshx user@host"
echo ""
if [ "$choice" = "1" ] && [ -n "$PROFILE" ]; then
    echo "If you enabled the alias:"
    echo "  ssh user@host"
    echo ""
fi
echo "To undo the alias:"
echo "  sshx uninstall-alias"
echo ""
echo "To start the agent:"
echo "  sshx-agent"
echo ""

