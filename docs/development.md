# Development Guide

How to contribute to Envy development.

## Prerequisites

- **Go 1.25.4** or later
- **Git** for version control
- **Make** (optional, for build automation)

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/XENONCYBER/envy.git
cd envy
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Build

```bash
# Development build
go build -o envy ./cmd/main.go

# Production build (optimized)
go build -ldflags="-s -w" -o envy ./cmd/main.go

# Run without building
go run ./cmd/main.go
```

### 4. Test

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test ./tests/encryption_test.go

# Verbose output
go test -v ./...
```

## Project Structure

```
envy/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── auth/                # Password prompts
│   │   └── password.go
│   ├── commands/            # CLI commands
│   │   ├── root.go          # Root command setup
│   │   ├── run.go           # 'envy run' command
│   │   ├── set.go           # 'envy set' command
│   │   ├── export.go        # Export functionality
│   │   └── import.go        # Import functionality
│   ├── config/              # Configuration
│   │   ├── lua.go           # Lua config loader
│   │   ├── keymap.go        # Keybindings
│   │   ├── theme.go         # Theme/styling
│   │   ├── paths.go         # File paths
│   │   └── version.go       # Version info
│   ├── crypto/              # Encryption
│   │   └── encryption.go
│   ├── domain/              # Data models
│   │   └── models.go
│   ├── service/             # Business logic
│   │   └── vault.go
│   ├── storage/             # File storage
│   │   ├── store.go
│   │   ├── lock_unix.go
│   │   └── lock_windows.go
│   └── tui/                 # Terminal UI
│       ├── model.go         # TUI state
│       ├── update.go        # Event handling
│       ├── view.go          # Rendering
│       ├── forms.go         # Form views
│       ├── bottombar.go     # Status bar
│       └── styles.go        # Style helpers
├── tests/                   # Integration tests
├── scripts/                 # Build scripts
├── docs/                    # Documentation
├── go.mod                   # Go module
├── go.sum                   # Dependencies
├── install.sh               # Installation script
└── README.md                # Project readme
```

## Development Workflow

### Making Changes

1. **Create a branch**
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make your changes**
   - Follow Go coding conventions
   - Add tests for new functionality
   - Update documentation

3. **Test your changes**
   ```bash
   go test ./...
   go run ./cmd/main.go --help
   ```

4. **Commit**
   ```bash
   git add .
   git commit -m "feat: add new feature description"
   ```

5. **Push and create PR**
   ```bash
   git push origin feature/my-feature
   ```

### Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` — New feature
- `fix:` — Bug fix
- `docs:` — Documentation
- `style:` — Formatting (no code change)
- `refactor:` — Code restructuring
- `test:` — Tests
- `chore:` — Maintenance

Examples:
```
feat: add project search by key name
fix: resolve clipboard clearing issue
docs: update installation instructions
refactor: simplify vault service interface
```

## Code Style

### Go Conventions

Follow standard Go style:

```bash
# Format code
go fmt ./...

# Lint
golangci-lint run

# Vet
go vet ./...
```

### Key Principles

1. **Error handling** — Always check errors
   ```go
   result, err := someFunction()
   if err != nil {
       return fmt.Errorf("context: %w", err)
   }
   ```

2. **Comments** — Document exported functions
   ```go
   // VaultService defines the interface for vault operations.
   type VaultService interface {
       // GetProjects returns all projects in the vault.
       GetProjects() []domain.Project
   }
   ```

3. **Testing** — Write tests for critical paths
   ```go
   func TestEncryption(t *testing.T) {
       // Test cases here
   }
   ```

## Testing

### Unit Tests

```go
// internal/crypto/encryption_test.go
package crypto

import "testing"

func TestEncryptDecrypt(t *testing.T) {
    key := make([]byte, 32)
    plaintext := []byte("secret data")
    
    encrypted, err := Encrypt(plaintext, key)
    if err != nil {
        t.Fatalf("encrypt failed: %v", err)
    }
    
    decrypted, err := Decrypt(encrypted, key)
    if err != nil {
        t.Fatalf("decrypt failed: %v", err)
    }
    
    if string(decrypted) != string(plaintext) {
        t.Error("decrypted text doesn't match original")
    }
}
```

### Integration Tests

```bash
# Run integration tests
go test ./tests/...

# Test with real vault operations
./scripts/integration_test.go
```

### Manual Testing

Test changes manually:

```bash
# Build and test locally
go build -o envy-dev ./cmd/main.go

# Test TUI
./envy-dev

# Test CLI commands
./envy-dev set testproject KEY=value
./envy-dev run testproject -- env | grep KEY
./envy-dev --export testproject
```

## Adding Features

### New CLI Command

1. Create file in `internal/commands/`:
   ```go
   // internal/commands/mycmd.go
   package commands
   
   import "github.com/spf13/cobra"
   
   var myCmd = &cobra.Command{
       Use:   "mycmd [args]",
       Short: "Brief description",
       RunE:  runMyCmd,
   }
   
   func init() {
       RootCmd.AddCommand(myCmd)
   }
   
   func runMyCmd(cmd *cobra.Command, args []string) error {
       // Implementation
       return nil
   }
   ```

2. Add tests
3. Update documentation

### New TUI Feature

1. Update model in `internal/tui/model.go`:
   ```go
   type Model struct {
       // Add new fields
       myNewField string
   }
   ```

2. Handle in `internal/tui/update.go`:
   ```go
   func (m Model) updateMyView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
       // Handle input
   }
   ```

3. Render in `internal/tui/view.go`:
   ```go
   func (m Model) viewMyView() string {
       // Return view string
   }
   ```

### New Configuration Option

1. Add to config struct in `internal/config/lua.go`:
   ```go
   type AppConfig struct {
       MyNewOption string
   }
   ```

2. Extract from Lua in `extractConfig()` function
3. Add default value in `DefaultAppConfig()`
4. Document in config docs

## Debugging

### Using Delve

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug
dlv debug ./cmd/main.go

# Or attach to running process
dlv attach <pid>
```

### Logging

Add temporary logging:

```go
import "log"

func someFunction() {
    log.Printf("Debug: value=%v", someValue)
    // ...
}
```

### Environment Variables

```bash
# Debug mode (if implemented)
ENVY_DEBUG=1 go run ./cmd/main.go
```

## Building Releases

### Manual Build

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o envy-linux-amd64 ./cmd/main.go

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o envy-darwin-amd64 ./cmd/main.go

# macOS ARM
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o envy-darwin-arm64 ./cmd/main.go

# Windows
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o envy-windows-amd64.exe ./cmd/main.go
```

### Using GoReleaser

```bash
# Install goreleaser
go install github.com/goreleaser/goreleaser@latest

# Build snapshot (no release)
goreleaser build --snapshot --rm-dist

# Full release (requires GITHUB_TOKEN)
goreleaser release
```

## Documentation

### Code Documentation

- Document all exported functions, types, and packages
- Use complete sentences
- Provide usage examples for complex functions

### User Documentation

Update relevant docs in `docs/`:
- `docs/usage/` for new features
- `docs/configuration/` for new options
- `docs/troubleshooting/` for known issues

## Submitting Changes

### Pull Request Process

1. **Update README** if needed
2. **Add tests** for new functionality
3. **Update docs** for user-facing changes
4. **Ensure tests pass**
5. **Fill out PR template**

### PR Checklist

- Code follows style guidelines
- Tests added/updated
- Documentation updated
- CHANGELOG.md updated
- Commits are meaningful
- No breaking changes (or clearly documented)

### Review Process

1. Maintainers review within 7 days
2. Address feedback
3. Merge when approved

## Getting Help

- **GitHub Discussions:** For questions
- **GitHub Issues:** For bugs and features

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Envy!
