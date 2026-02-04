# Your First Vault

Learn how Envy's vault system works and how to manage it effectively.

## What is the Vault?

The vault is your encrypted secrets database. It's a single JSON file containing:

- **Projects** — Logical groupings of secrets
- **Keys** — Individual secrets with metadata
- **History** — Previous values of each key
- **Encryption metadata** — Salt and authentication hash

## Vault File Location

The vault is stored at:

| OS | Path |
|----|------|
| Linux | `~/.envy/keys.json` |
| macOS | `~/.envy/keys.json` |
| Windows | `%APPDATA%\envy\keys.json` |

You can customize this path in your [Lua configuration](../configuration/lua-config.md).

## Vault Structure

```json
{
  "version": 1,
  "salt": "base64_encoded_salt",
  "auth_hash": "base64_encoded_hash",
  "projects": [
    {
      "name": "myproject",
      "environment": "prod",
      "keys": [
        {
          "title": "API_KEY",
          "key": "API_KEY",
          "current": {
            "value": "encrypted_value",
            "created_at": "2024-01-15T10:30:00Z",
            "created_by": "cli-set"
          },
          "history": [
            {
              "value": "encrypted_old_value",
              "created_at": "2024-01-14T09:00:00Z",
              "created_by": "tui"
            }
          ]
        }
      ]
    }
  ]
}
```

**Note:** All secret values are encrypted with AES-256-GCM. The salt and auth_hash are used for password verification and key derivation.

## Creating Your First Project

### Via TUI

1. Launch Envy: `envy`
2. Press `ctrl+n` to create a new project
3. Enter project name
4. Select environment (dev/stage/prod)
5. Add at least one key-value pair
6. Press `Save` (or `ctrl+s`)

### Via CLI

```bash
# Create project and add first secret
envy set myproject API_KEY=mysecretvalue

# You'll be prompted:
# - Master password (to unlock vault)
# - Confirm creating new project 'myproject' (dev)
```

## Understanding Environments

Each project exists in one of three environments:

- **`dev`** — Local development (green badge in TUI)
- **`stage`** — Staging/testing (yellow badge)
- **`prod`** — Production (red badge)

You can have multiple projects with the same name in different environments:

```bash
envy set webapp DATABASE_URL=localhost -e dev
envy set webapp DATABASE_URL=staging-db -e stage
envy set webapp DATABASE_URL=prod-db -e prod
```

This creates three separate projects, each with its own set of secrets.

## Version History

Every time you update a secret, the old value is preserved:

```bash
# Initial value
envy set myapp API_KEY=old_key_123

# Update - old value moves to history
envy set myapp API_KEY=new_key_456

# View history in TUI: Select project → Choose key → Press 'H'
```

History shows:
- Current value (green badge)
- Previous values (yellow badges)
- Timestamp of each change
- Who/what made the change (`cli-set`, `tui`, `tui-edit`, `cli-import`)

## Vault Security

### File Permissions

Envy sets strict permissions on vault files:

- **Vault file:** 0600 (readable/writable by owner only)
- **Data directory:** 0700 (accessible by owner only)

### Encryption Details

- **Algorithm:** AES-256-GCM (authenticated encryption)
- **Key Derivation:** Argon2id with:
  - Time: 1 iteration
  - Memory: 64 MB
  - Threads: 4
  - Salt: 16 random bytes
- **Key Size:** 256 bits (32 bytes)

### Master Password

Your master password is never stored. Instead:
1. Argon2id derives a 256-bit key from your password + salt
2. A SHA-256 hash of this key is stored as `auth_hash`
3. During unlock, the hash is recomputed and compared

This means:
- We can't recover your password if you forget it
- We can't access your secrets without your password
- If you forget your password, your secrets are lost forever

## Backup and Recovery

### Manual Backup

```bash
# Simple backup with timestamp
cp ~/.envy/keys.json ~/.envy/keys.json.backup.$(date +%Y%m%d)

# Or backup to different location
cp ~/.envy/keys.json ~/Dropbox/envy-backup.json
```

### Automated Backup Script

```bash
#!/bin/bash
# Add to crontab: 0 2 * * * /path/to/backup-envy.sh

BACKUP_DIR="$HOME/envy-backups"
mkdir -p "$BACKUP_DIR"

# Create timestamped backup
cp ~/.envy/keys.json "$BACKUP_DIR/keys-$(date +%Y%m%d-%H%M%S).json"

# Keep only last 30 backups
ls -t "$BACKUP_DIR"/keys-*.json | tail -n +31 | xargs rm -f
```

### Recovery

If your vault becomes corrupted or you need to restore:

```bash
# Restore from backup
cp ~/.envy/keys.json.backup.20240115 ~/.envy/keys.json

# Launch Envy - you'll need your master password
envy
```

## Migration

### Moving to a New Machine

1. Copy the vault file to the new machine:
   ```bash
   scp ~/.envy/keys.json newmachine:~/.envy/
   ```

2. Set correct permissions:
   ```bash
   chmod 600 ~/.envy/keys.json
   chmod 700 ~/.envy/
   ```

3. Copy your config (optional):
   ```bash
   scp ~/.config/envy/config.lua newmachine:~/.config/envy/
   ```

### Export/Import All Projects

```bash
# On old machine - export each project
for project in $(envy list-projects); do
  envy --export "$project"
  mv .env "${project}.env"
done

# On new machine - import each
for envfile in *.env; do
  envy --import "$envfile"
done
```

## Vault Maintenance

### Checking Vault Health

```bash
# Verify JSON is valid
python3 -m json.tool ~/.envy/keys.json > /dev/null && echo "Valid JSON"

# Check file permissions
ls -la ~/.envy/keys.json
# Should show: -rw------- (0600)
```

### Compacting History

Old history entries accumulate over time. While Envy doesn't have a built-in compact command, you can:

1. Export projects
2. Create new vault
3. Re-import (only current values, not history)

## Troubleshooting Vault Issues

### "Authentication failed: incorrect password"

- Double-check your password (caps lock?)
- If you recently restored from backup, ensure the backup was from after your last password change

### "Failed to parse storage file (corrupted?)"

- Restore from backup immediately
- Check disk space and file system health
- If no backup exists, the vault is unrecoverable

### Lock File Issues

If Envy crashes, it may leave a stale lock file:

```bash
# Remove stale lock (only if Envy is not running!)
rm ~/.envy/.lock
```

---

**Next:** Learn about [Projects](../usage/projects.md) and how to organize your secrets effectively.
