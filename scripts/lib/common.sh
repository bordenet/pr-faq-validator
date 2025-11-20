#!/usr/bin/env bash

init_script() {
    echo "Initializing script..."
}

log_header() {
    echo "=== $1 ==="
}

log_success() {
    echo "✅ $1"
}

log_info() {
    echo "ℹ️  $1"
}

log_error() {
    echo "❌ $1"
}

die() {
    log_error "$1"
    exit 1
}
