#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Basic validation tiers
case "${1:-med}" in
    --p1)
        echo "Running P1 validation..."
        go mod tidy
        go build
        echo "✅ P1 validation passed"
        ;;
    --med)
        echo "Running medium validation..."
        go mod tidy
        go build
        go test ./...
        echo "✅ Medium validation passed"
        ;;
    --all)
        echo "Running full validation..."
        go mod tidy
        go build
        go test ./...
        echo "✅ Full validation passed"
        ;;
    *)
        echo "Usage: $0 [--p1|--med|--all]"
        exit 1
        ;;
esac
