#!/bin/sh
# Envy Installation Script
# Usage: curl -fsSL https://raw.githubusercontent.com/XENONCYBER/envy/main/install.sh | sh
#
# This script detects OS/CPU, downloads the correct binary from GitHub Releases,
# verifies it, and installs to /usr/local/bin
#
# Supports Linux (x86_64, arm64, arm), macOS (x86_64, arm64), and Windows (x86_64)

set -e

# Configuration
REPO="XENONCYBER/envy"
BINARY_NAME="envy"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print functions
info() {
    printf "${BLUE}[INFO]${NC} %s\n" "$1" >&2
}

success() {
    printf "${GREEN}[OK]${NC} %s\n" "$1" >&2
}

warn() {
    printf "${YELLOW}[WARN]${NC} %s\n" "$1" >&2
}

error() {
    printf "${RED}[ERROR]${NC} %s\n" "$1" >&2
    exit 1
}

# Detect OS
detect_os() {
    OS="$(uname -s)"
    case "${OS}" in
    Linux*) OS="linux" ;;
    Darwin*) OS="darwin" ;;
    MINGW* | MSYS* | CYGWIN*) OS="windows" ;;
    *) error "Unsupported operating system: ${OS}" ;;
    esac
    echo "${OS}"
}

# Detect CPU architecture
detect_arch() {
    ARCH="$(uname -m)"
    case "${ARCH}" in
    x86_64 | amd64) ARCH="amd64" ;;
    aarch64 | arm64) ARCH="arm64" ;;
    armv7l | armv6l) ARCH="arm" ;;
    i386 | i686) ARCH="386" ;;
    *) error "Unsupported architecture: ${ARCH}" ;;
    esac
    echo "${ARCH}"
}

# Detect Linux distribution
detect_distro() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        echo "${ID}"
    elif [ -f /etc/arch-release ]; then
        echo "arch"
    elif [ -f /etc/fedora-release ]; then
        echo "fedora"
    elif [ -f /etc/debian_version ]; then
        echo "debian"
    else
        echo "unknown"
    fi
}

# Check for required commands
check_dependencies() {
    for cmd in curl tar; do
        if ! command -v "${cmd}" >/dev/null 2>&1; then
            error "Required command '${cmd}' not found. Please install it first."
        fi
    done

    # Check for git (needed for building from source)
    if ! command -v git >/dev/null 2>&1; then
        warn "git not found - building from source will not be available"
    fi

    # Check for sha256sum or shasum
    if command -v sha256sum >/dev/null 2>&1; then
        SHA_CMD="sha256sum"
    elif command -v shasum >/dev/null 2>&1; then
        SHA_CMD="shasum -a 256"
    else
        warn "sha256sum/shasum not found - skipping checksum verification"
        SHA_CMD=""
    fi
}

# Get latest release version from GitHub
get_latest_version() {
    LATEST_URL="https://api.github.com/repos/${REPO}/releases/latest"

    if command -v curl >/dev/null 2>&1; then
        # Check if releases exist first
        HTTP_CODE=$(curl -sSL -o /dev/null -w "%{http_code}" "${LATEST_URL}")
        if [ "${HTTP_CODE}" = "404" ]; then
            echo "none"
            return
        fi
        
        # Get release info and check for assets
        RELEASE_INFO=$(curl -fsSL "${LATEST_URL}")
        VERSION=$(echo "${RELEASE_INFO}" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
        ASSETS_COUNT=$(echo "${RELEASE_INFO}" | grep '"assets"' -A 10 | grep '\[' -A 10 | grep -c 'name' || echo "0")
        
        # If no assets available, return "none" to trigger build from source
        if [ "${ASSETS_COUNT}" = "0" ]; then
            echo "none"
            return
        fi
    else
        error "curl is required to fetch the latest version"
    fi

    if [ -z "${VERSION}" ]; then
        echo "none"
        return
    fi

    echo "${VERSION}"
}

# Build from source fallback
build_from_source() {
    TMP_DIR="$1"
    info "No pre-built binaries available. Building from source..."
    
    # Check if Go is installed
    if ! command -v go >/dev/null 2>&1; then
        error "Go is required to build from source. Please install Go first:"
        info "  https://go.dev/doc/install"
    fi
    
    # Check Go version
    GO_VERSION=$(go version 2>/dev/null | awk '{print $3}' | sed 's/go//')
    if [ -n "${GO_VERSION}" ]; then
        info "Using Go ${GO_VERSION}"
    fi
    
    # Clone repository
    info "Cloning source code..."
    if ! git clone --depth 1 https://github.com/${REPO}.git "${TMP_DIR}"; then
        error "Failed to clone repository"
    fi
    
    # Build
    info "Building binary..."
    cd "${TMP_DIR}"
    
    # Run the build
    if ! go build -ldflags="-s -w" -o "${BINARY_NAME}" ./cmd/main.go; then
        error "Build failed"
    fi
    
    # Make executable
    chmod +x "${BINARY_NAME}"
    
    cd - >/dev/null
}

# Download and verify binary
download_binary() {
    VERSION="$1"
    OS="$2"
    ARCH="$3"
    TMP_DIR="$4"

    # Construct download URL
    # Format: envy_VERSION_OS_ARCH.tar.gz
    FILENAME="${BINARY_NAME}_${VERSION#v}_${OS}_${ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${FILENAME}"
    CHECKSUM_URL="https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt"

    info "Downloading ${BINARY_NAME} ${VERSION} for ${OS}/${ARCH}..."

    # Download binary
    if ! curl -fsSL "${DOWNLOAD_URL}" -o "${TMP_DIR}/${FILENAME}"; then
        error "Failed to download ${DOWNLOAD_URL}"
    fi

    success "Downloaded ${FILENAME}"

    # Download and verify checksum if available
    if [ -n "${SHA_CMD}" ]; then
        info "Verifying checksum..."
        if curl -fsSL "${CHECKSUM_URL}" -o "${TMP_DIR}/checksums.txt" 2>/dev/null; then
            cd "${TMP_DIR}"
            EXPECTED_SUM=$(grep "${FILENAME}" checksums.txt | awk '{print $1}')
            if [ -n "${EXPECTED_SUM}" ]; then
                ACTUAL_SUM=$(${SHA_CMD} "${FILENAME}" | awk '{print $1}')
                if [ "${EXPECTED_SUM}" != "${ACTUAL_SUM}" ]; then
                    error "Checksum verification failed!"
                fi
                success "Checksum verified"
            else
                warn "Checksum for ${FILENAME} not found in checksums.txt"
            fi
            cd - >/dev/null
        else
            warn "Checksums file not available - skipping verification"
        fi
    fi

    # Extract binary
    info "Extracting binary..."
    tar -xzf "${TMP_DIR}/${FILENAME}" -C "${TMP_DIR}"

    # Find the binary (handling potential subdirectories from tar)
    # We look for the binary name, ensuring it is a file
    FOUND_BINARY=$(find "${TMP_DIR}" -type f -name "${BINARY_NAME}" | head -n 1)

    if [ -z "${FOUND_BINARY}" ]; then
        error "Binary '${BINARY_NAME}' not found in downloaded archive."
    fi

    # If the binary is not in the root of TMP_DIR, move it there
    if [ "${FOUND_BINARY}" != "${TMP_DIR}/${BINARY_NAME}" ]; then
        mv "${FOUND_BINARY}" "${TMP_DIR}/${BINARY_NAME}"
    fi

    # Make executable
    chmod +x "${TMP_DIR}/${BINARY_NAME}"
}

# Install binary to system
install_binary() {
    TMP_DIR="$1"

    info "Installing to ${INSTALL_DIR}..."

    # Check if we need sudo
    if [ -w "${INSTALL_DIR}" ]; then
        mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        if command -v sudo >/dev/null 2>&1; then
            sudo mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
        elif command -v doas >/dev/null 2>&1; then
            doas mv "${TMP_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
        else
            error "Cannot write to ${INSTALL_DIR}. Please run with sudo or as root."
        fi
    fi

    success "Installed ${BINARY_NAME} to ${INSTALL_DIR}/${BINARY_NAME}"
}

# Verify installation
verify_installation() {
    if command -v "${BINARY_NAME}" >/dev/null 2>&1; then
        INSTALLED_VERSION=$("${BINARY_NAME}" version 2>/dev/null || echo "unknown")
        success "Installation complete!"
        info "Version: ${INSTALLED_VERSION}"
        info "Run '${BINARY_NAME}' to get started"
    else
        warn "Installation complete, but '${BINARY_NAME}' not found in PATH"
        info "You may need to add ${INSTALL_DIR} to your PATH"
        info "Add this to your shell profile:"
        info "  export PATH=\"\$PATH:${INSTALL_DIR}\""
    fi
}

# Print system info
print_system_info() {
    OS=$(detect_os)
    ARCH=$(detect_arch)

    if [ "${OS}" = "linux" ]; then
        DISTRO=$(detect_distro)
        info "Detected: ${OS}/${ARCH} (${DISTRO})"
    else
        info "Detected: ${OS}/${ARCH}"
    fi
}

# Main installation flow
main() {
    printf "\n"
    printf "${GREEN}  ███████╗███╗   ██╗██╗   ██╗██╗   ██╗${NC}\n"
    printf "${GREEN}  ██╔════╝████╗  ██║██║   ██║╚██╗ ██╔╝${NC}\n"
    printf "${GREEN}  █████╗  ██╔██╗ ██║██║   ██║ ╚████╔╝ ${NC}\n"
    printf "${GREEN}  ██╔══╝  ██║╚██╗██║╚██╗ ██╔╝  ╚██╔╝  ${NC}\n"
    printf "${GREEN}  ███████╗██║ ╚████║ ╚████╔╝    ██║   ${NC}\n"
    printf "${GREEN}  ╚══════╝╚═╝  ╚═══╝  ╚═══╝     ╚═╝   ${NC}\n"
    printf "\n"
    printf "  ${BLUE}Secure Secret Manager${NC}\n"
    printf "\n"

    # Check dependencies
    check_dependencies

    # Detect system
    print_system_info
    OS=$(detect_os)
    ARCH=$(detect_arch)

    # Get latest version
    info "Fetching latest version..."
    VERSION=$(get_latest_version)
    
    # Create a single temporary directory for either build or download
    TMP_DIR=$(mktemp -d)
    trap "rm -rf ${TMP_DIR}" EXIT
    
    if [ "${VERSION}" = "none" ]; then
        warn "No pre-built binaries found"
        # Build from source instead
        build_from_source "${TMP_DIR}"
    else
        success "Latest version: ${VERSION}"
        # Download binary
        download_binary "${VERSION}" "${OS}" "${ARCH}" "${TMP_DIR}"
    fi

    # Install
    install_binary "${TMP_DIR}"

    # Verify
    verify_installation

    printf "\n"
    info "Quick start:"
    info "  1. Import existing secrets: ${BINARY_NAME} --import .env"
    info "  2. Launch the TUI: ${BINARY_NAME}"
    info "  3. Search and copy secrets: Type to search, press 'y' to copy"
    printf "\n"
    info "Documentation: https://github.com/${REPO}"
    printf "\n"
}

# Handle uninstall flag
if [ "$1" = "uninstall" ] || [ "$1" = "--uninstall" ]; then
    info "Uninstalling ${BINARY_NAME}..."
    if [ -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
        if [ -w "${INSTALL_DIR}" ]; then
            rm -f "${INSTALL_DIR}/${BINARY_NAME}"
        else
            sudo rm -f "${INSTALL_DIR}/${BINARY_NAME}"
        fi
        success "Uninstalled ${BINARY_NAME}"
        warn "Note: Your encrypted vault at ~/.envy.json remains untouched"
        info "To remove it completely: rm ~/.envy.json"
    else
        warn "${BINARY_NAME} not found in ${INSTALL_DIR}"
    fi
    exit 0
fi

# Handle help flag
if [ "$1" = "help" ] || [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    echo "Envy Installation Script"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --version    Show latest version available"
    echo "  --uninstall   Remove Envy from system"
    echo "  --help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  curl -fsSL https://raw.githubusercontent.com/XENONCYBER/envy/main/install.sh | sh"
    echo "  curl -fsSL https://raw.githubusercontent.com/XENONCYBER/envy/main/install.sh | sh -s -- --version"
    echo "  curl -fsSL https://raw.githubusercontent.com/XENONCYBER/envy/main/install.sh | sh -s -- --uninstall"
    exit 0
fi

# Handle version flag
if [ "$1" = "version" ] || [ "$1" = "--version" ]; then
    echo "Fetching latest version..."
    VERSION=$(get_latest_version)
    echo "Latest Envy version: ${VERSION}"
    exit 0
fi

# Run main installation
main
