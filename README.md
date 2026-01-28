# Envy

A secure encrypted vault for managing API keys, secrets, and environment variables with CLI and TUI interfaces. Envy implements a "fuzzy" matching algorithm for quick secret retrieval and provides military-grade encryption to protect your sensitive data.

## Highlights

- **Secure** — AES-256-GCM encryption with Argon2id key derivation
- **Fast** — Optimized for instant access to thousands of secrets
- **Portable** — Single binary installation with no external dependencies
- **Versatile** — CLI and TUI interfaces with shell integration support
- **All-inclusive** — Import/export .env files, version history, and multi-environment support

## Table of Contents

- [Installation](#installation)
  - [Quick Install (curl)](#quick-install-curl)
  - [Build from Source](#build-from-source)
  - [Package Managers](#package-managers)
  - [Upgrading Envy](#upgrading-envy)
- [Quick Start](#quick-start)
- [Usage](#usage)
  - [TUI Mode](#tui-mode)
  - [CLI Commands](#cli-commands)
  - [Import/Export](#importexport)
- [Configuration](#configuration)
  - [Environment Variables](#environment-variables)
  - [Key Bindings](#key-bindings)
- [Security](#security)
- [Advanced Topics](#advanced-topics)
  - [Shell Integration](#shell-integration)
  - [Backup and Recovery](#backup-and-recovery)
  - [Performance Tips](#performance-tips)
- [Development](#development)
- [License](#license)

## Installation

### Quick Install (recommended)

Install the latest binary with one command:

```bash
curl -fsSL https://raw.githubusercontent.com/XENONCYBER/envy/main/install.sh | sh
```

### Build from Source

If you have Go 1.25+ installed:

```bash
git clone https://github.com/XENONCYBER/envy.git
cd envy
go build -o envy ./cmd/main.go
sudo mv envy /usr/local/bin/
```

Or build and run locally:

```bash
go run ./cmd/main.go --help
```

### Upgrading Envy

To upgrade to the latest version:

```bash
curl -fsSL https://raw.githubusercontent.com/XENONCYBER/envy/main/install.sh | sh
```

The installer will automatically replace your existing installation.

## Quick Start

First run creates a new encrypted vault:

```bash
envy
```

You will be prompted to create a master password. This password encrypts all your secrets and is never stored - only a salted hash is kept for verification.

**Get started immediately:**

1. **Import existing secrets:** `envy --import path/to/.env` or `envy -i path/to/.env`
2. **Launch TUI:** `envy`
3. **Search and copy secrets:** Type to search, press `y` to copy

Your encrypted vault is stored at `~/.envy/.envy.json` with secure file permissions.

## Usage

### TUI Mode

Launch the interactive interface:

```bash
envy
```

**Search Syntax:**
- Fuzzy matching: type partial names like `api` for `API_KEY`
- Exact match: wrap in quotes `"API_KEY"`
- Filter by project: `project:prod api`

### Import/Export

**Import .env file:**
```bash
envy --import path/to/.env
# or specify project name
envy --import path/to/.env --project myproject
```

**Export to .env file:**
```bash
envy --export project-name
# or specify output file
envy --export project-name --output path/to/.env
```

## Configuration

### Key Bindings

Customize key bindings in `~/.config/envy/config.lua`:

```lua
return {
  keys = {
    quit = "q",
    copy = "y",
    edit = "e",
    delete = "d",
    search = "/",
    create = "n"
  }
}
```

## Security

- **Encryption:** AES-256-GCM with authenticated encryption
- **Key Derivation:** Argon2id with configurable parameters
- **Master Password:** Never stored, only salted hash for verification
- **File Permissions:** Vault stored with 0600 permissions
- **Atomic Writes:** Prevents data corruption during saves
- **Clipboard:** Auto-clears after 30 seconds by default
- **Memory:** Secrets cleared from memory after use

## Advanced Topics

### Backup and Recovery

**Manual backup:**
```bash
cp ~/.envy.json ~/.envy.backup.$(date +%Y%m%d)
```

**Automated backup script:**
```bash
#!/bin/bash
BACKUP_DIR="$HOME/envy-backups"
mkdir -p "$BACKUP_DIR"
cp ~/.envy.json "$BACKUP_DIR/envy-$(date +%Y%m%d-%H%M%S).json"
ls -t "$BACKUP_DIR"/envy-*.json | tail -n +31 | xargs rm -f
```

**Recovery:**
```bash
cp ~/.envy.backup.20240101 ~/.envy.json
```

### Performance Tips

- **Large vaults:** Use search filters to narrow results
- **Frequent access:** Pin commonly used projects
- **Memory usage:** Vault is loaded on-demand, not kept in memory
- **Network drives:** Avoid storing vault on network filesystems

## Development

### Prerequisites

- Go 1.25.4 or later
- Make (optional, for build scripts)

### Building

```bash
# Build for current platform
go build -o envy ./cmd/main.go

# Build for all platforms
make build-all

# Build with version info
go build -ldflags "-X main.version=$(git describe --tags)" -o envy ./cmd/main.go
```

### Running

```bash
# Run from source
go run ./cmd/main.go
```

### Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test ./tests/encryption_test.go
```

### Development Workflow

```bash
# Install development dependencies
go mod download

# Run linter
go vet ./...

# Format code
go fmt ./...

# Build and test
make test
```

### Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes and add tests
4. Ensure all tests pass: `make test`
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Related Projects

- [pass](https://www.passwordstore.org/) - Standard unix password manager
- [gopass](https://www.gopass.pw/) - The slightly more awesome standard unix password manager
- [vault](https://www.vaultproject.io/) - HashiCorp's secrets management tool

---

**Envy** - Secure secret management made simple.  
[GitHub](https://github.com/XENONCYBER/envy) | [Issues](https://github.com/XENONCYBER/envy/issues) | [Releases](https://github.com/XENONCYBER/envy/releases)
