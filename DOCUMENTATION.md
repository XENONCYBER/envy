# Envy Documentation

Welcome to the Envy docs. Whether you're just getting started or looking to customize everything, you'll find what you need here.

## Quick Links

**New to Envy?** Start with [Installation](docs/getting-started/installation.md), then work through the [Quick Start](docs/getting-started/quick-start.md) guide.

**Using Envy daily?** Bookmark the [CLI Commands](docs/usage/cli-commands.md) reference and [TUI Guide](docs/usage/tui-guide.md).

**Setting up for a team?** Read [Team Setup](docs/examples/team-setup.md) and [Security Best Practices](docs/security/best-practices.md).

**Want to customize?** Check out [Lua Configuration](docs/configuration/lua-config.md) and [Themes](docs/configuration/themes.md).

## What's Inside

The documentation is organized into sections that match how you'll actually use Envy:

**Getting Started** — Installation, creating your first vault, and understanding the basics

**Usage** — Detailed guides for the TUI, CLI commands, organizing projects, and managing environments

**Configuration** — Customize everything: keybindings, colors, storage locations, and more via Lua

**Security** — How the encryption works, password recommendations, and keeping your secrets safe

**Reference** — Quick lookups for commands and configuration options

**Troubleshooting** — Solutions to common problems and recovery procedures

**Examples** — Real-world workflows, team configurations, and CI/CD integration

## The Big Ideas

Envy revolves around a few core concepts that are worth understanding upfront:

**Projects** are how you group related secrets. Think of them as folders — you might have one project for your web app, another for your API service, and so on.

**Environments** let you separate credentials by stage. Each project can exist in dev, stage, and prod, with completely different values in each.

**Version History** keeps previous values when you update secrets. Accidentally overwrote a working API key? Pull the old one from history.

**The Vault** is an encrypted JSON file stored on your machine. It's encrypted with AES-256-GCM using a key derived from your master password via Argon2id.

## Common Commands

```bash
# Start the TUI
envy

# Add a secret
envy set myproject DATABASE_URL=postgres://localhost/db

# Run something with your secrets
envy run myproject -- npm start

# Import/export .env files
envy --import .env
envy --export myproject
```

See the [full command reference](docs/usage/cli-commands.md) for details.

## Configuration

Envy reads `~/.config/envy/config.lua` for customization. Here's a minimal example:

```lua
return {
  keys = {
    quit = "q",
    yank = "y"
  },
  theme = {
    accent = "#ff6b6b"
  }
}
```

See the [complete configuration guide](docs/configuration/lua-config.md) for all options.

## Getting Help

Stuck? Try these in order:

1. Check the [FAQ](docs/troubleshooting/faq.md) — common questions answered
2. Look at [Common Issues](docs/troubleshooting/common-issues.md) — solutions to frequent problems
3. Run `envy --help` — quick command reference
4. Open a [GitHub Issue](https://github.com/XENONCYBER/envy/issues) — bug reports and feature requests

---

Ready to dive in? Start with [Installation](docs/getting-started/installation.md) or go straight to the [Quick Start](docs/getting-started/quick-start.md) if you already have Envy installed.
