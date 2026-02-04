# Envy Documentation

Welcome to the comprehensive Envy documentation. This is the official documentation for **Envy** - a secure, encrypted terminal-based secret manager for API keys, passwords, and environment variables.

## Quick Links

**New Users:** [Installation](./getting-started/installation.md) → [Quick Start](./getting-started/quick-start.md) → [TUI Guide](./usage/tui-guide.md)

**Daily Use:** [CLI Commands](./usage/cli-commands.md) | [Projects](./usage/projects.md) | [Environments](./usage/environments.md)

**Configuration:** [Lua Config](./configuration/lua-config.md) | [Themes](./configuration/themes.md) | [Keybindings](./configuration/keybindings.md)

**Security:** [Encryption](./security/encryption.md) | [Master Password](./security/master-password.md) | [Best Practices](./security/best-practices.md)

**Reference:** [Commands](./reference/commands.md)

---

## What is Envy?

Envy is a secure, encrypted secret management tool designed for developers and DevOps professionals. It combines a beautiful terminal user interface (TUI) with powerful command-line capabilities, making it easy to manage API keys, database credentials, and environment variables while maintaining the highest security standards.

## Documentation Sections

| Section | Description | For |
|---------|-------------|-----|
| [Getting Started](./getting-started/) | Installation, first vault, basic concepts | New users |
| [Usage Guide](./usage/) | TUI navigation, CLI commands, projects, environments | Daily users |
| [Configuration](./configuration/) | Lua config, themes, keybindings, file locations | Customizers |
| [Security](./security/) | Encryption details, password best practices, security guide | Security-conscious |
| [Reference](./reference/) | Quick command reference | Everyone |
| [Troubleshooting](./troubleshooting/) | Common issues, FAQ, recovery | Problem solvers |
| [Examples](./examples/) | Real-world workflows, team setup, CI/CD | Implementers |

## Documentation Structure

```
docs/
├── README.md                    # This file
├── getting-started/
│   ├── installation.md          # Installation guide
│   ├── quick-start.md           # Get started in 5 minutes
│   └── first-vault.md           # Understanding your vault
├── usage/
│   ├── tui-guide.md             # Complete TUI navigation guide
│   ├── cli-commands.md          # All CLI commands explained
│   ├── projects.md              # Project organization
│   └── environments.md          # Environment management
├── configuration/
│   ├── lua-config.md            # Lua configuration with examples
│   ├── themes.md                # Theme gallery and customization
│   ├── keybindings.md           # Keybinding customization
│   └── file-locations.md        # Where files are stored
├── security/
│   ├── encryption.md            # Technical encryption details
│   ├── master-password.md       # Password best practices
│   └── best-practices.md        # Security recommendations
├── reference/
│   └── commands.md              # Command quick reference
├── troubleshooting/
│   ├── common-issues.md         # Solutions to common problems
│   ├── faq.md                   # Frequently asked questions
│   └── recovery.md              # Disaster recovery guide
├── examples/
│   ├── basic-workflow.md        # Complete workflow example
│   ├── team-setup.md            # Team configuration guide
│   └── ci-cd-integration.md     # CI/CD pipeline integration
└── development.md               # Development/contribution guide
```

## Core Concepts

### Projects
Projects are logical groupings of related secrets. For example, you might have a "webapp" project containing `DATABASE_URL`, `API_KEY`, and `SECRET_TOKEN`.

### Environments
Each project can have multiple environments:
- **dev** — Local development (green badge)
- **stage** — Staging/testing (yellow badge)
- **prod** — Production (red badge)

This lets you maintain separate credentials for different deployment stages.

### Version History
Every secret update preserves the previous value in history. You can view historical values to track changes or recover old credentials.

### Vault
Your encrypted data store using AES-256-GCM with Argon2id key derivation:
- Linux/macOS: `~/.envy/keys.json`
- Windows: `%APPDATA%\envy\keys.json`

## Quick Start

```bash
# Install
curl -fsSL https://raw.githubusercontent.com/XENONCYBER/envy/main/install.sh | sh

# First launch (creates vault)
envy

# Import existing secrets
envy --import .env

# Set a secret
envy set myproject API_KEY=secret123

# Run with secrets
envy run myproject -- npm start
```

See [Quick Start Guide](./getting-started/quick-start.md) for detailed walkthrough.

## Getting Help

- **Browse docs**: You're already here! Explore the sections above.
- **Quick help**: Run `envy --help` in your terminal
- **GitHub Issues**: [github.com/XENONCYBER/envy/issues](https://github.com/XENONCYBER/envy/issues)
- **FAQ**: Check [troubleshooting/faq.md](./troubleshooting/faq.md)

---

**Ready?** Start with [Installation](./getting-started/installation.md) or [Quick Start](./getting-started/quick-start.md).
