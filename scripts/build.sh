#!/bin/bash
# Consolidated build script for CanIFly
# Handles platform-specific builds with proper error handling

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="canifly-backend"
DIST_DIR="dist"

# Helper functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed. Please install Go first."
        exit 1
    fi
    log_info "Go version: $(go version)"
}

# Build for specific platform
build_platform() {
    local platform=$1
    local goos=$2
    local goarch=$3
    local output=$4
    
    log_info "Building for $platform..."
    
    # Create directory if it doesn't exist
    mkdir -p "$(dirname "$output")"
    
    # Build
    GOOS=$goos GOARCH=$goarch go build -o "$output" .
    
    if [ $? -eq 0 ]; then
        log_info "✓ Built $platform binary: $output"
    else
        log_error "Failed to build for $platform"
        exit 1
    fi
}

# Build React frontend
build_frontend() {
    log_info "Building React frontend..."
    
    if [ ! -d "renderer" ]; then
        log_error "renderer directory not found"
        exit 1
    fi
    
    cd renderer
    
    # Install dependencies if node_modules doesn't exist
    if [ ! -d "node_modules" ]; then
        log_info "Installing frontend dependencies..."
        npm install
    fi
    
    # Build
    npm run build
    
    if [ $? -eq 0 ]; then
        log_info "✓ Frontend built successfully"
    else
        log_error "Failed to build frontend"
        exit 1
    fi
    
    cd ..
}

# Main build function
main() {
    local target=${1:-all}
    
    log_info "Starting CanIFly build process..."
    check_go
    
    case $target in
        "mac"|"darwin")
            build_platform "macOS" "darwin" "amd64" "$DIST_DIR/mac/$BINARY_NAME"
            ;;
        "win"|"windows")
            build_platform "Windows" "windows" "amd64" "$DIST_DIR/win/$BINARY_NAME.exe"
            ;;
        "linux")
            build_platform "Linux" "linux" "amd64" "$DIST_DIR/linux/$BINARY_NAME"
            ;;
        "frontend"|"react")
            build_frontend
            ;;
        "backend"|"go")
            build_platform "macOS" "darwin" "amd64" "$DIST_DIR/mac/$BINARY_NAME"
            build_platform "Windows" "windows" "amd64" "$DIST_DIR/win/$BINARY_NAME.exe"
            build_platform "Linux" "linux" "amd64" "$DIST_DIR/linux/$BINARY_NAME"
            ;;
        "all"|"")
            # Build everything
            log_info "Building all platforms..."
            
            # Backend
            build_platform "macOS" "darwin" "amd64" "$DIST_DIR/mac/$BINARY_NAME"
            build_platform "Windows" "windows" "amd64" "$DIST_DIR/win/$BINARY_NAME.exe"
            build_platform "Linux" "linux" "amd64" "$DIST_DIR/linux/$BINARY_NAME"
            
            # Frontend
            build_frontend
            
            log_info "✓ All builds completed successfully!"
            ;;
        *)
            log_error "Unknown target: $target"
            echo "Usage: $0 [mac|win|linux|frontend|backend|all]"
            exit 1
            ;;
    esac
}

# Run main function
main "$@"