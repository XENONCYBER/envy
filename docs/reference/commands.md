# Command Reference

Quick reference for all Envy commands and options.

## Global Flags

| Flag | Shorthand | Description | Example |
|------|-----------|-------------|---------|
| `--import <file>` | `-i` | Import .env file | `envy -i .env` |
| `--export <project>` | `-t` | Export to .env | `envy -t myapp` |
| `--version` | — | Show version | `envy --version` |
| `--help` | `-h` | Show help | `envy --help` |

## Commands

### envy

Launch TUI (no arguments).

```bash
envy
```

**First run:** Creates vault, prompts for master password.
**Subsequent runs:** Prompts for master password, opens TUI.

---

### envy set

Set or update a secret.

```bash
envy set <project> <KEY=VALUE> [flags]
```

**Arguments:**
- `project` — Project name (creates if missing)
- `KEY=VALUE` — Secret to store

**Flags:**
- `-e, --env <env>` — Environment: `dev` (default), `stage`, `prod`

**Examples:**
```bash
# Set in dev (default)
envy set myapp API_KEY=secret123

# Set in production
envy set myapp API_KEY=prod-key -e prod

# Project with spaces
envy set "My App" DATABASE_URL=postgres://localhost
```

**Shorthand:**
```bash
envy -s myapp API_KEY=secret123 -e prod
```

---

### envy run

Run command with secrets as environment variables.

```bash
envy run <project> -- <command> [args...]
```

**Arguments:**
- `project` — Project name (case-insensitive)
- `--` — Required separator
- `command` — Command to execute
- `args` — Arguments for command

**Examples:**
```bash
# Run npm
envy run myapp -- npm start

# Run Python
envy run production -- python app.py

# Run with args
envy run myapp -- docker run -p 3000:3000 myimage
```

**Important:**
- `--` separator is required
- Secrets only available to child process
- Returns exit code of child command

---

### envy --import

Import .env file into vault.

```bash
envy --import <filepath>
envy -i <filepath>
```

**Arguments:**
- `filepath` — Path to .env file

**Interactive:**
- Prompts for project name
- Prompts for environment (dev/stage/prod)

**Examples:**
```bash
envy --import .env
envy --import ~/projects/myapp/.env
envy -i ./config/production.env
```

---

### envy --export

Export project to .env file.

```bash
envy --export <project>
envy -t <project>
```

**Arguments:**
- `project` — Project name (case-insensitive)

**Output:** Creates `.env` in current directory.

**Examples:**
```bash
envy --export myapp
envy -t production
```

**Warning:** Creates plaintext file with secrets. Secure immediately.

---

### envy --version

Show version information.

```bash
envy --version
```

**Output:**
```
Envy v1.1.2
```

---

### envy --help

Show help.

```bash
envy --help              # General help
envy set --help          # Command-specific help
envy run --help          # Command-specific help
```

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Error (wrong password, invalid args, etc.) |
| N | Exit code from command (with `envy run`) |

## Command Comparison

| Task | Best Command | Notes |
|------|--------------|-------|
| Browse secrets | `envy` | TUI is fastest |
| Quick copy | `envy` → `y` | TUI copy to clipboard |
| Set one secret | `envy set p K=V` | Fast CLI operation |
| Set many secrets | `envy --import file` | Bulk import |
| Run app | `envy run p -- cmd` | Injects env vars |
| Export for deploy | `envy --export p` | Creates .env file |
| Edit secret | `envy` → `e` | TUI edit mode |
| View history | `envy` → `H` | TUI history sidebar |

## Environment Variables

Envy reads:

| Variable | Purpose | Default |
|----------|---------|---------|
| `HOME` | Find home directory | System default |
| `APPDATA` | Windows data path | `%USERPROFILE%\AppData\Roaming` |
| `XDG_DATA_HOME` | Linux data path | `~/.local/share` |
| `XDG_CONFIG_HOME` | Linux config path | `~/.config` |

## File Locations

| OS | Vault | Config |
|----|-------|--------|
| Linux | `~/.envy/keys.json` | `~/.config/envy/config.lua` |
| macOS | `~/.envy/keys.json` | `~/Library/Application Support/envy/config.lua` |
| Windows | `%APPDATA%\envy\keys.json` | `%APPDATA%\envy\config.lua` |

## Keybindings (TUI)

### Grid View

| Key | Action |
|-----|--------|
| `↑↓←→` or `hjkl` | Navigate |
| `Enter` | Open project |
| `i` | Search |
| `Ctrl+n` | New project |
| `d` | Delete project |
| `q` | Quit |

### Detail View

| Key | Action |
|-----|--------|
| `↑↓` or `kj` | Navigate keys |
| `Enter/Space` | Reveal/hide |
| `y` | Copy to clipboard |
| `e` | Edit value |
| `E` | Edit project |
| `H` | View history |
| `d` | Delete key |
| `Esc` | Back |

## Validation Rules

### Project Names

- Required (cannot be empty)
- Max 256 characters
- Any characters allowed except system restrictions

### Key Names

- Required (cannot be empty)
- Max 256 characters
- Cannot contain: `=`, newline, carriage return

### Values

- Can be empty
- No length limit (within reason)
- Any characters allowed

### Environments

Must be exactly one of:
- `dev`
- `stage`
- `prod`

## Limitations

| Aspect | Limit |
|--------|-------|
| Projects | Unlimited (practical limit: thousands) |
| Keys per project | Unlimited (practical limit: hundreds) |
| Key name length | 256 characters |
| Value length | Unlimited (tested to MB range) |
| History depth | Unlimited (all versions kept) |
| Vault file size | Limited by disk space |
| Concurrent access | One user at a time (file locking) |

## Quick Examples

### Daily Workflow

```bash
# Morning: Start TUI
envy

# Copy a secret: Select project → y

# Set new secret
envy set myapp NEW_KEY=value

# Run dev server
envy run myapp-dev -- npm run dev

# Deploy with prod secrets
envy run myapp-prod -- ./deploy.sh
```

### Setup Workflow

```bash
# New project setup
cd ~/projects/newapp

# Import existing .env
envy --import .env

# Add more secrets
envy set newapp DATABASE_URL=postgres://localhost/newapp
envy set newapp API_KEY=generated-key-123

# Launch to verify
envy
```

---

See [CLI Commands](../usage/cli-commands.md) for detailed command documentation.
