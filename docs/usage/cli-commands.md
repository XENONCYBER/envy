# CLI Commands

Complete reference for Envy's command-line interface.

## Command Overview

```
envy [flags]                    # Launch TUI
envy [command] [args] [flags]   # Run CLI command
```

## Global Flags

These flags work at the top level (without subcommands):

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--import <file>` | `-i` | Import .env file into vault |
| `--export <project>` | `-t` | Export project to .env file |
| `--version` | — | Show version information |
| `--help` | `-h` | Show help |

## Commands

### `envy` (No Arguments)

Launch the interactive TUI.

```bash
envy
```

**Behavior:**
- If vault doesn't exist: Prompts to create master password
- If vault exists: Prompts for master password, then opens TUI

**When to use:**
- Browsing and managing secrets interactively
- Creating new projects
- Viewing version history
- Copying secrets to clipboard

---

### `envy set`

Set or update a secret in a project.

```bash
envy set [project] [KEY=VALUE] [flags]
```

**Arguments:**
- `project` — Project name (creates if doesn't exist)
- `KEY=VALUE` — Secret in KEY=VALUE format

**Flags:**
- `-e, --env <environment>` — Target environment: `dev` (default), `stage`, or `prod`

**Examples:**

```bash
# Set in default environment (dev)
envy set myapp DATABASE_URL=postgresql://localhost/db

# Set in production
envy set myapp API_KEY=sk-abc123 -e prod
#or
envy -s myproj -e prod KEY=VALUE

# Project names with spaces (quote them)
envy set "My App" SECRET= value -e dev
```

**What it does:**
1. Prompts for master password
2. Creates project if it doesn't exist
3. If key exists: Updates it and moves old value to history
4. If key doesn't exist: Creates it
5. Saves vault

**Output:**
```
Enter master password: ********
Added 'API_KEY' to project 'myapp' (prod)
Vault saved successfully
```

**Limitations:**
- Key name cannot contain `=`, newlines, or carriage returns
- Key and value cannot be empty
- Key name max length: 256 characters
- Environment must be exactly `dev`, `stage`, or `prod`

**Shorthand Syntax:**

You can also use `-s` flag directly:

```bash
envy -s myapp API_KEY=secret123 -e prod
#or
envy -s myproj -e prod KEY=VALUE
```

This is equivalent to `envy set myapp API_KEY=secret123 -e prod`.

---

### `envy run`

Run a command with project secrets injected as environment variables.

```bash
envy run [project] -- [command] [args...]
```

**Arguments:**
- `project` — Project name (case-insensitive)
- `--` — Separator between Envy args and your command
- `command` — Command to execute
- `args...` — Arguments for your command

**Examples:**

```bash
# Run npm with secrets
envy run myapp -- npm start

# Run Python script
envy run production -- python app.py

# Run with arguments
envy run myapp -- docker run -e NODE_ENV=production myimage

# Make build
envy run staging -- make build

# Start development server
envy run dev -- ./start-server.sh
```

**What it does:**
1. Prompts for master password
2. Loads all secrets from the specified project
3. Adds them to the environment
4. Executes your command with that environment
5. Returns when your command exits

**Important characteristics:**
- Secrets are **only** available to the child process
- Your shell environment is **not** modified
- Secrets are cleaned up when the process exits
- Case-insensitive project name matching
- Returns your command's exit code

**Output:**
```
Enter master password: ********
Loaded 5 secrets from 'myapp' (prod)
Running: npm start

> myapp@1.0.0 start
> node server.js

Server running on port 3000...
```

**Error cases:**

```bash
# Missing -- separator
envy run myapp npm start
# Error: missing '--' separator

# Project not found
envy run nonexistent -- npm start
# Error: project 'nonexistent' not found

# Missing command
envy run myapp --
# Error: missing command after '--'
```

**Security Note:**
The `--` separator is required to prevent command injection. Envy will not execute the command without it.

---

### `envy --import`

Import a .env file into the vault.

```bash
envy --import <filepath>
envy -i <filepath>
```

**Arguments:**
- `filepath` — Path to .env file

**Interactive prompts:**
1. Project name
2. Environment (dev/stage/prod, default: dev)

**Examples:**

```bash
# Import current directory .env
envy --import .env

# Import specific file
envy --import ~/projects/myapp/.env

# Import production config
envy --import ./.env.production
```

**What it does:**
1. Parses the .env file (handles comments, empty lines, quotes)
2. If first run: Creates vault and prompts for master password
3. Prompts for project name and environment
4. Creates project with all key-value pairs
5. Saves vault

**Supported .env formats:**

```bash
# Comments and empty lines are ignored
KEY=value
KEY2="quoted value"
KEY3='single quoted'
KEY4=value with spaces  # This works too
# KEY5=ignored (commented)
```

**Conflict handling:**
If a project with the same name and environment exists:
```
Project 'myapp' (prod) already exists. Overwrite? [y/N]:
```

- Press `y` to replace the existing project
- Press `n` or `Enter` to cancel

**Limitations:**
- Invalid key names (containing `=`, newlines) are skipped with a warning
- Empty values are allowed
- File must be readable text (UTF-8)

---

### `envy --export`

Export a project to a .env file.

```bash
envy --export <project>
envy -t <project>
```

**Arguments:**
- `project` — Project name (case-insensitive matching)

**Output file:**
Creates `.env` in the current directory.

**Examples:**

```bash
# Export to current directory
envy --export myapp
# Creates: ./.env

# Export then move
envy --export myapp && mv .env ~/projects/myapp/
```

**What it does:**
1. Prompts for master password
2. Finds project (case-insensitive)
3. Creates `.env` file with all secrets
4. Warns about plaintext secrets

**Output file format:**

```bash
# Exported from Envy - Project: myapp (prod)
DATABASE_URL=postgresql://user:pass@localhost/db
API_KEY=sk-abc123xyz789
SECRET_TOKEN=supersecrettoken123
```

**Security warnings:**
- Exported file contains **plaintext secrets**
- Set restrictive permissions: `chmod 600 .env`
- Never commit .env files to version control
- Delete after use: `rm .env`

**Limitations:**
- Only exports current values (history not included)
- Overwrites existing `.env` file without warning
- Case-insensitive project matching (first match wins)

---

### `envy --version`

Show version information.

```bash
envy --version
```

**Output:**
```
Envy v1.1.2
```

---

### `envy --help`

Show help for Envy and its commands.

```bash
envy --help              # General help
envy set --help          # Help for set command
envy run --help          # Help for run command
```

---

## Command Comparison Table

| Task | Command | Notes |
|------|---------|-------|
| Browse secrets | `envy` | Opens TUI |
| Set single secret | `envy set p K=V` | Quick CLI operation |
| Run with secrets | `envy run p -- cmd` | Injects env vars |
| Import .env | `envy --import file` | Creates project |
| Export to .env | `envy --export p` | Plaintext output |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error (invalid args, auth failed, etc.) |
| Exit code from command | When using `envy run`, returns the child's exit code |

## Environment Variables

Envy reads these environment variables:

| Variable | Purpose |
|----------|---------|
| `ENVY_DATA_DIR` | Override vault storage location |
| `ENVY_CONFIG_DIR` | Override config file location |

See [Environment Variables](../reference/environment-variables.md) for details.

## Scripting Examples

### CI/CD Pipeline

```bash
#!/bin/bash
set -e

# Export secrets for build
envy --export myapp

# Source them (in sub-shell to avoid polluting)
(
  source .env
  npm run build
)

# Clean up
rm -f .env
```

### Backup Script

```bash
#!/bin/bash
# Backup all projects

for project in $(envy list-projects 2>/dev/null); do
  mkdir -p backups
  envy --export "$project"
  mv .env "backups/${project}-$(date +%Y%m%d).env"
done
```

### Development Workflow

```bash
#!/bin/bash
# start-dev.sh

# Check if we have the dev project
if ! envy run myapp-dev -- echo "test" >/dev/null 2>&1; then
  echo "Creating dev project..."
  envy set myapp-dev DATABASE_URL=postgres://localhost/dev
  envy set myapp-dev API_KEY=dev-key-123
fi

# Start development server
envy run myapp-dev -- npm run dev
```

---

**Next:** Learn about organizing secrets with [Projects](./projects.md).
