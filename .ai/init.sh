#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
AI_DIR="$ROOT_DIR/.ai"

echo "=== AI Dev Team Setup ==="

mkdir -p "$AI_DIR/tickets/shipments"
mkdir -p "$AI_DIR/agents"
mkdir -p "$AI_DIR/templates"

echo "  Directories created."
echo ""
echo "  AI Dev Team is ready."
echo "  Use /ship to start a shipment."
