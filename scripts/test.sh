#!/usr/bin/env bash
################################################################################
# Test Script
################################################################################
# PURPOSE: Run all tests and linting checks
# USAGE: ./scripts/test.sh [--quick|--full]
################################################################################

set -euo pipefail

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Source common functions
source "$SCRIPT_DIR/lib/common.sh"

cd "$PROJECT_ROOT"

# Parse arguments
QUICK_MODE=false
if [[ "${1:-}" == "--quick" ]]; then
    QUICK_MODE=true
fi

log_step "Running Go tests and linting checks"

# Check required tools
log_info "Checking required tools..."
require_commands go golangci-lint || die "Missing required tools"

# Run gofmt check
log_step "Checking code formatting (gofmt)..."
UNFORMATTED=$(gofmt -l . 2>&1 | grep -v vendor || true)
if [ -n "$UNFORMATTED" ]; then
    log_error "The following files are not formatted:"
    echo "$UNFORMATTED"
    log_info "Run: gofmt -w ."
    exit 1
fi
log_success "Code formatting check passed"

# Run goimports check (if available)
if command_exists goimports; then
    log_step "Checking import formatting (goimports)..."
    UNFORMATTED_IMPORTS=$(goimports -l . 2>&1 | grep -v vendor || true)
    if [ -n "$UNFORMATTED_IMPORTS" ]; then
        log_warning "The following files have unformatted imports:"
        echo "$UNFORMATTED_IMPORTS"
        log_info "Run: goimports -w ."
    else
        log_success "Import formatting check passed"
    fi
fi

# Run go vet
log_step "Running go vet..."
if go vet ./...; then
    log_success "go vet passed"
else
    log_error "go vet failed"
    exit 1
fi

# Run golangci-lint
log_step "Running golangci-lint..."
if golangci-lint run; then
    log_success "golangci-lint passed"
else
    log_error "golangci-lint failed"
    exit 1
fi

# Run Go tests
log_step "Running Go tests..."
if go test -v -race -coverprofile=coverage.out ./...; then
    log_success "All tests passed"
    
    # Show coverage
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}')
    log_info "Test coverage: $COVERAGE"
else
    log_error "Tests failed"
    exit 1
fi

log_success "All checks passed! ✨"

