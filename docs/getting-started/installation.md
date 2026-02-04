# Installation

This guide covers installing Envy on Linux, macOS, and Windows.

## Requirements

- **Linux**: x86_64, ARM64, or ARM processor
- **macOS**: Intel (x86_64) or Apple Silicon (ARM64)
- **Windows**: x86_64 processor
- **Optional**: Go 1.25.4+ (only if building from source)

## Quick Install (Recommended)

Install the latest release with one command:

```bash
curl -fsSL https://raw.githubusercontent.com/XENONCYBER/envy/main/install.sh | sh
```

This script automatically:
- Detects your OS and architecture
- Downloads the correct binary
- Verifies the checksum
- Installs to `/usr/local/bin`

### Install Script Options

```bash
# Show help
curl -fsSL .../install.sh | sh -s -- --help

# Check version without installing
curl -fsSL .../install.sh | sh -s -- --version

# Install to custom directory
INSTALL_DIR=~/bin curl -fsSL .../install.sh | sh

# Uninstall Envy
curl -fsSL .../install.sh | sh -s -- --uninstall
```

## Manual Installation

### From GitHub Releases

1. Download the appropriate archive from [Releases](https://github.com/XENONCYBER/envy/releases):
   - Linux: `envy_vX.X.X_linux_{amd64,arm64,arm}.tar.gz`
   - macOS: `envy_vX.X.X_darwin_{amd64,arm64}.tar.gz`
   - Windows: `envy_vX.X.X_windows_amd64.tar.gz`

2. Extract and install:

**Linux/macOS:**
```bash
tar -xzf envy_vX.X.X_OS_ARCH.tar.gz
sudo mv envy /usr/local/bin/
chmod +x /usr/local/bin/envy
```

**Windows (PowerShell as Administrator):**
```powershell
Expand-Archive -Path envy_vX.X.X_windows_amd64.tar.gz -DestinationPath C:\Windows\System32
```

## Build from Source

Requires Go 1.25.4 or later:

```bash
# Clone repository
git clone https://github.com/XENONCYBER/envy.git
cd envy

# Build
go build -ldflags="-s -w" -o envy ./cmd/main.go

# Install to system (optional)
sudo mv envy /usr/local/bin/
```

### Development Build

Run without installing:
```bash
go run ./cmd/main.go --help
```

## Verify Installation

```bash
# Check version
envy --version
# Output: Envy v1.1.2

# Show help
envy --help
```

## Troubleshooting

### "command not found"

Ensure `/usr/local/bin` is in your PATH:

```bash
# Add to shell profile (~/.bashrc, ~/.zshrc, etc.)
export PATH="$PATH:/usr/local/bin"

# Reload configuration
source ~/.bashrc  # or ~/.zshrc
```

### macOS "developer cannot be verified"

Remove the quarantine attribute:
```bash
xattr -d com.apple.quarantine /usr/local/bin/envy
```

### Permission Denied

```bash
sudo chmod +x /usr/local/bin/envy
```

## Upgrading

Run the install script again—it automatically replaces the existing binary. Your vault (`~/.envy/keys.json`) remains untouched.

## Uninstallation

```bash
# Using install script
curl -fsSL .../install.sh | sh -s -- --uninstall

# Or manually
sudo rm /usr/local/bin/envy

# Remove vault (WARNING: deletes all secrets!)
rm -rf ~/.envy
```

---

**Next:** [Quick Start Guide](./quick-start.md)
