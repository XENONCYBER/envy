<p align="center">
  <img src="docs/logo.png" alt="Envy Logo" width="400">
</p>

<p align="center">
  <a href="#quick-start">Quick Start</a> •
  <a href="docs/configuration/lua-config.md">Configuration</a> •
  <a href="docs/usage/cli-commands.md">Commands</a> •
  <a href="docs/configuration/keybindings.md">Keybindings</a> •
  <a href="docs/troubleshooting/faq.md">FAQ</a> •
  <a href="DOCUMENTATION.md">Documentation</a>
</p>

<p align="center">
  A secure encrypted vault for managing API keys, secrets, and environment variables. Built for developers who live in the terminal.
</p>

## What is Envy?

I built Envy because I was tired of juggling .env files and keeping secrets scattered across password managers, sticky notes, and Slack threads. It's a single binary that gives you both a slick TUI for browsing secrets and a CLI for automation.

Here's what you get:

- **Real encryption** — AES-256-GCM with Argon2id, not just base64 obfuscation
- **Fast access** — Fuzzy search finds secrets in milliseconds
- **Organized by project** — Group related secrets together
- **Multiple environments** — Keep dev, stage, and prod separate
- **Version history** — See what changed and when
- **Import/export** — Bring in your existing .env files

## Quick Install

One command:

```bash
curl -fsSL https://raw.githubusercontent.com/XENONCYBER/envy/main/install.sh | sh
```

Or build from source:
```bash
git clone https://github.com/XENONCYBER/envy.git && cd envy
go build -o envy ./cmd/main.go
```

## Quick Start

```bash
# Create your vault (first run only)
envy

# Import existing .env files
envy --import .env

# Set secrets via CLI
envy set myproject API_KEY=secret123

# Run commands with secrets injected
envy run myproject -- npm start

# Export project to .env file
envy --export myproject
```

## Security

Your secrets are encrypted with AES-256-GCM. The key is derived from your master password using Argon2id — a memory-hard function designed to resist GPU cracking attempts.

- Your password is never stored anywhere
- We only keep a hash to verify it's correct
- The vault file is useless without your password
- If you forget your password, your secrets are gone. Period.

Read more about [how the encryption works](docs/security/encryption.md).

## Documentation

The full docs live in the [docs/](docs/) folder:

- **[Getting Started](docs/getting-started/)** — Installation, your first vault, basic concepts
- **[Usage Guide](docs/usage/)** — TUI navigation, CLI commands, workflows
- **[Configuration](docs/configuration/)** — Lua configuration, themes, keybindings
- **[Security](docs/security/)** — Encryption details, best practices
- **[Reference](docs/reference/)** — Command reference
- **[Troubleshooting](docs/troubleshooting/)** — Common issues, FAQ
- **[Examples](docs/examples/)** — Real-world workflows

## Development

```bash
# Prerequisites: Go 1.25.4+

# Clone and build
git clone https://github.com/XENONCYBER/envy.git
cd envy
go build -o envy ./cmd/main.go

# Run tests
go test ./...

# Run locally
go run ./cmd/main.go
```

See [Development Guide](docs/development.md) for contribution guidelines.

## License

MIT License — see [LICENSE](LICENSE) for details.

## Similar Projects

- [pass](https://www.passwordstore.org/) — The standard Unix password manager
- [gopass](https://www.gopass.pw/) — A more feature-rich pass implementation
- [vault](https://www.vaultproject.io/) — HashiCorp's enterprise secrets solution

---

**[GitHub](https://github.com/XENONCYBER/envy)** · **[Issues](https://github.com/XENONCYBER/envy/issues)** · **[Releases](https://github.com/XENONCYBER/envy/releases)**
