#!/bin/bash

# Build script for Shelly Prometheus Exporter
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
VERSION=${VERSION:-"dev"}
COMMIT=${COMMIT:-$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")}
BUILD_TIME=${BUILD_TIME:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}
LDFLAGS="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}"

# Build directory
BUILD_DIR="build"

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to build for a specific platform
build_platform() {
    local os=$1
    local arch=$2
    local output_name=$3
    
    print_status "Building for ${os}/${arch}..."
    
    GOOS=${os} GOARCH=${arch} go build \
        -ldflags="${LDFLAGS}" \
        -o "${BUILD_DIR}/${output_name}" \
        ./cmd/shelly-exporter
    
    if [ $? -eq 0 ]; then
        print_status "Successfully built ${output_name}"
    else
        print_error "Failed to build ${output_name}"
        exit 1
    fi
}

# Function to create archive
create_archive() {
    local name=$1
    local os=$2
    local arch=$3
    
    print_status "Creating archive for ${name}..."
    
    cd "${BUILD_DIR}"
    tar -czf "${name}.tar.gz" "${name}" LICENSE README.md CHANGELOG.md
    cd ..
    
    print_status "Created ${name}.tar.gz"
}

# Main build process
main() {
    print_status "Starting build process..."
    print_status "Version: ${VERSION}"
    print_status "Commit: ${COMMIT}"
    print_status "Build Time: ${BUILD_TIME}"
    
    # Clean previous builds
    if [ -d "${BUILD_DIR}" ]; then
        print_status "Cleaning previous build directory..."
        rm -rf "${BUILD_DIR}"
    fi
    
    # Create build directory
    mkdir -p "${BUILD_DIR}"
    
    # Build for different platforms
    build_platform "linux" "amd64" "shelly-exporter-linux-amd64"
    build_platform "linux" "arm64" "shelly-exporter-linux-arm64"
    build_platform "windows" "amd64" "shelly-exporter-windows-amd64.exe"
    build_platform "darwin" "amd64" "shelly-exporter-darwin-amd64"
    build_platform "darwin" "arm64" "shelly-exporter-darwin-arm64"
    
    # Create archives
    create_archive "shelly-exporter-linux-amd64" "linux" "amd64"
    create_archive "shelly-exporter-linux-arm64" "linux" "arm64"
    create_archive "shelly-exporter-windows-amd64.exe" "windows" "amd64"
    create_archive "shelly-exporter-darwin-amd64" "darwin" "amd64"
    create_archive "shelly-exporter-darwin-arm64" "darwin" "arm64"
    
    # Create checksums
    print_status "Creating checksums..."
    cd "${BUILD_DIR}"
    sha256sum *.tar.gz > checksums.txt
    cd ..
    
    print_status "Build completed successfully!"
    print_status "Artifacts are available in the ${BUILD_DIR}/ directory"
    
    # List artifacts
    print_status "Built artifacts:"
    ls -la "${BUILD_DIR}/"
}

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    print_error "go.mod not found. Please run this script from the project root."
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed or not in PATH"
    exit 1
fi

# Run main function
main "$@"
