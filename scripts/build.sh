#!/usr/bin/env bash
################################################################################
# Build Script
################################################################################
# PURPOSE: Build the pr-faq-validator binary with linting checks
# USAGE: ./scripts/build.sh [--skip-lint]
################################################################################

set -euo pipefail

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Source common functions
source "$SCRIPT_DIR/lib/common.sh"

cd "$PROJECT_ROOT"

# Parse arguments
SKIP_LINT=false
if [[ "${1:-}" == "--skip-lint" ]]; then
    SKIP_LINT=true
fi

log_step "Building pr-faq-validator"

# Check required tools
log_info "Checking required tools..."
require_commands go || die "Go is not installed"

# Run linting checks first (unless skipped)
if [ "$SKIP_LINT" = false ]; then
    log_step "Running linting checks..."
    if [ -f "$SCRIPT_DIR/test.sh" ]; then
        "$SCRIPT_DIR/test.sh" --quick
    else
        log_warning "test.sh not found, skipping linting"
    fi
fi

# Clean previous builds
log_step "Cleaning previous builds..."
rm -f pr-faq-validator
log_success "Cleaned"

# Get version info
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

log_info "Version: $VERSION"
log_info "Commit: $GIT_COMMIT"
log_info "Build time: $BUILD_TIME"

# Build binary
log_step "Building binary..."
go build \
    -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT" \
    -o pr-faq-validator \
    .

if [ -f pr-faq-validator ]; then
    log_success "Build successful: pr-faq-validator"
    
    # Show binary info
    SIZE=$(du -h pr-faq-validator | cut -f1)
    log_info "Binary size: $SIZE"
    
    # Make executable
    chmod +x pr-faq-validator
    
    log_success "Build complete! ✨"
    echo ""
    log_info "Run with: ./pr-faq-validator -file testdata/example_prfaq_1.md"
else
    log_error "Build failed"
    exit 1
fi

