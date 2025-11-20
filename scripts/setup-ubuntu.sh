#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib/common.sh"
init_script

log_header "Ubuntu Development Environment Setup"

sudo apt-get update -qq
sudo apt-get install -y golang-go curl wget unzip build-essential
log_success "Ubuntu setup complete!"
