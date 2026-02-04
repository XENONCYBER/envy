# Quick Start Guide

Let's get Envy running in the next few minutes.

## First Launch

Type this in your terminal:

```bash
envy
```

Since this is your first time, Envy asks you to create a master password. Choose something strong — this password encrypts everything in your vault. We don't store it anywhere, so if you forget it, your secrets are gone.

## Import Your Existing Secrets

Already have a `.env` file? Bring it in:

```bash
envy -i .env  #or envy --import .env
```

Envy asks for a project name (something like "myapp" or "api-service") and which environment this is for. Most people start with `dev`.

Done. Your secrets are now encrypted and stored safely.

## Using the TUI

Launch the interface anytime:

```bash
envy
```

You'll see a grid of your projects. Here's how to move around:

| Key | What it does |
|-----|--------------|
| Arrow keys or `h/j/k/l` | Move between projects |
| `Enter` | Open a project |
| `i` | Search through everything |
| `ctrl+n` | Create a new project |
| `q` | Quit |

Once you're inside a project (press Enter on one), you can:

| Key | What it does |
|-----|--------------|
| Arrows or `j/k` | Move between secrets |
| `Enter` or `Space` | Show/hide the secret value |
| `y` | Copy to clipboard (clears in 30 seconds) |
| `e` | Edit a secret |
| `E` | Add or remove keys from this project |
| `H` | See old versions of this secret |
| `esc` | Go back to the grid |

### Searching

Press `i` to search. Start typing and projects filter instantly. Press `Tab` to switch between searching everything, just project names, or just key names.

## Using the CLI

For quick tasks, the command line is faster:

### Add a Secret

```bash
# Add to dev environment (default)
envy set myproject DATABASE_URL=postgresql://localhost/mydb

# Add to production
envy set myproject API_KEY=secret123 -e prod

# Shorthand version
envy -s myproject API_KEY=secret123 -e prod
```

If the project doesn't exist, Envy creates it. If the key already exists, the old value gets saved to history.

### Run Commands

```bash
envy run myproject -- npm start
```

This loads all secrets from "myproject", makes them available as environment variables, then runs your command. When your command finishes, those environment variables disappear — they never leak into your shell.

### Export to .env

Sometimes you need a regular .env file:

```bash
envy --export myproject #or envy -t myproject
```

This creates a `.env` file in your current directory. **Be careful** — it's plaintext. Secure it immediately:

```bash
chmod 600 .env
```

And delete it when you're done. Never commit .env files to git.

## Your Daily Workflow

Most people use Envy like this:

```bash
# Morning — start your dev server
envy run myproject -- npm run dev

# Need to add a new API key?
envy set myproject NEW_API_KEY=whatever

# Want to browse or copy something?
envy
# (use the TUI to find and copy)
```

## Where Everything Lives

- **Your encrypted vault:** `~/.envy/keys.json` (Linux/macOS) or `%APPDATA%\envy\keys.json` (Windows)
- **Your settings:** `~/.config/envy/config.lua`
- **Lock file:** `~/.envy/.lock` (prevents two Envy instances from running at once)

The vault is encrypted with AES-256-GCM. Without your master password, it's just random noise.

## What's Next?

- Learn the TUI inside out with the [TUI Guide](../usage/tui-guide.md)
- See all the [CLI commands](../usage/cli-commands.md)
- Customize colors and keys with [Lua configuration](../configuration/lua-config.md)
- Understand how the [encryption works](../security/encryption.md)

---

Hit a snag? Check [Common Issues](../troubleshooting/common-issues.md) for solutions.
