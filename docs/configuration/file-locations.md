# File Locations

Where Envy stores its files on different platforms.

## Overview

Envy uses two main directories:

1. **Data Directory** — Encrypted vault and lock file
2. **Config Directory** — Lua configuration file

## Default Locations

### Linux

| Type | Path |
|------|------|
| Data | `~/.envy/` |
| Config | `~/.config/envy/` |
| Vault | `~/.envy/keys.json` |
| Lock | `~/.envy/.lock` |
| Config File | `~/.config/envy/config.lua` |

### macOS

| Type | Path |
|------|------|
| Data | `~/.envy/` |
| Config | `~/Library/Application Support/envy/` |
| Vault | `~/.envy/keys.json` |
| Lock | `~/.envy/.lock` |
| Config File | `~/Library/Application Support/envy/config.lua` |

**Note:** macOS uses both standard Unix paths (`~/.envy/`) and macOS-specific paths (`~/Library/`).

### Windows

| Type | Path |
|------|------|
| Data | `%APPDATA%\envy\` |
| Config | `%APPDATA%\envy\` |
| Vault | `%APPDATA%\envy\keys.json` |
| Lock | `%APPDATA%\envy\.lock` |
| Config File | `%APPDATA%\envy\config.lua` |

**Windows paths expanded:**
- Typically: `C:\Users\<username>\AppData\Roaming\envy\`

## XDG Base Directory Support

On Linux, Envy respects XDG environment variables:

| Variable | Used For | Default |
|----------|----------|---------|
| `XDG_DATA_HOME` | Data directory | `~/.local/share` |
| `XDG_CONFIG_HOME` | Config directory | `~/.config` |

**Custom XDG paths:**
```bash
export XDG_DATA_HOME="$HOME/.local/share"
export XDG_CONFIG_HOME="$HOME/.config"

# Envy will use:
# Data: $XDG_DATA_HOME/envy/
# Config: $XDG_CONFIG_HOME/envy/
```

## Customizing Locations

### Via Lua Configuration

Set custom paths in `config.lua`:

```lua
return {
  backend = {
    keys_path = "/path/to/custom/keys.json",
    lock_path = "/path/to/custom/.lock"
  }
}
```

### Path Expansion

Envy supports these shortcuts:

| Shortcut | Expands To |
|----------|------------|
| `~` | Home directory |
| Environment variables | Via `os.getenv()` |

**Examples:**
```lua
backend = {
  keys_path = "~/Dropbox/envy/keys.json",
  lock_path = "~/Dropbox/envy/.lock"
}
```

### Platform-Specific Paths

Use the `envy` Lua module:

```lua
return {
  backend = {
    keys_path = envy.home .. "/secrets/envy/keys.json",
    lock_path = envy.home .. "/secrets/envy/.lock"
  }
}
```

## File Permissions

Envy sets strict permissions on sensitive files:

| File/Directory | Permissions | Notes |
|----------------|-------------|-------|
| `keys.json` | 0600 | Owner read/write only |
| `.lock` | 0600 | Owner read/write only |
| Data directory | 0700 | Owner access only |
| Config directory | 0755 | Standard directory |
| `config.lua` | 0644 | Not sensitive |

### Checking Permissions

```bash
# Linux/macOS
ls -la ~/.envy/
# Should show:
# -rw------- keys.json
# -rw------- .lock
# drwx------ .

# Verify config permissions
ls -la ~/.config/envy/
```

### Fixing Permissions

```bash
# If permissions are wrong
chmod 600 ~/.envy/keys.json
chmod 600 ~/.envy/.lock
chmod 700 ~/.envy/
```

## File Descriptions

### keys.json

The encrypted vault containing all your secrets.

**Format:** JSON with encrypted values
**Size:** Depends on number of secrets (typically 1-100 KB)
**Critical:** Yes - backup regularly!

Structure:
```json
{
  "version": 1,
  "salt": "base64_salt",
  "auth_hash": "base64_hash",
  "projects": [...]
}
```

### .lock

Prevents concurrent access to the vault.

**Format:** Empty file or process ID
**Critical:** No - can be safely deleted if stale
**Automatic:** Created on vault open, deleted on close

### config.lua

User configuration in Lua format.

**Format:** Lua script
**Critical:** No - Envy works without it
**Optional:** Yes

## Migration Scenarios

### Moving to New Machine

1. Copy vault file:
   ```bash
   scp ~/.envy/keys.json newmachine:~/.envy/
   ```

2. Copy config (optional):
   ```bash
   scp ~/.config/envy/config.lua newmachine:~/.config/envy/
   ```

3. Set permissions:
   ```bash
   chmod 600 ~/.envy/keys.json
   chmod 700 ~/.envy/
   ```

### Using Cloud Storage

Store vault in synced folder:

```lua
-- config.lua
return {
  backend = {
    keys_path = "~/Dropbox/envy/keys.json",
    lock_path = "~/Dropbox/envy/.lock"
  }
}
```

**Caution:** 
- Lock file may cause issues with sync
- Ensure cloud storage is encrypted
- Don't sync to public/shared folders

### Multiple Vaults

Use different configs for different contexts:

```bash
# Work vault
ENVY_CONFIG=~/.config/envy/work.lua envy

# Personal vault  
ENVY_CONFIG=~/.config/envy/personal.lua envy
```

**Note:** Future feature - currently requires manual config switching.

### Separate Config Per Project

```lua
-- ~/.config/envy/team-a.lua
return {
  backend = {
    keys_path = "~/.envy/team-a/keys.json",
    lock_path = "~/.envy/team-a/.lock"
  }
}

-- ~/.config/envy/team-b.lua
return {
  backend = {
    keys_path = "~/.envy/team-b/keys.json",
    lock_path = "~/.envy/team-b/.lock"
  }
}
```

## Backup Locations

Recommended backup strategy:

```bash
# Primary location
~/.envy/keys.json

# Backup locations
~/Backups/envy/keys.json.$(date +%Y%m%d)
/path/to/encrypted/backup/keys.json
```

## Security Considerations

### What to Protect

| File | Sensitivity | Action |
|------|-------------|--------|
| `keys.json` | **Critical** | Secure backups, never share |
| `.lock` | Low | Can be deleted |
| `config.lua` | Low | May contain paths (mild sensitivity) |

### What NOT to Do

Don't store vault in:
- Public Git repositories
- Unencrypted cloud storage
- Network shares without encryption
- World-readable directories

Don't:
- Share your vault file with others
- Copy vault to untrusted machines
- Upload to Pastebin/services

Do:
- Keep regular encrypted backups
- Store on encrypted filesystems
- Use strong master passwords
- Audit file permissions

## Environment Variables

Future versions may support:

```bash
# Override data directory
export ENVY_DATA_DIR="/secure/storage/envy"

# Override config directory
export ENVY_CONFIG_DIR="/etc/envy"

# Override specific files
export ENVY_KEYS_PATH="/secure/keys.json"
export ENVY_LOCK_PATH="/tmp/envy.lock"
```

**Current status:** Not implemented - use `config.lua` instead.

## Troubleshooting File Issues

### "Permission denied"

```bash
# Check permissions
ls -la ~/.envy/

# Fix
chmod 600 ~/.envy/keys.json
chmod 700 ~/.envy/
```

### "No such file or directory"

```bash
# Create directories
mkdir -p ~/.envy
mkdir -p ~/.config/envy
```

### Lock file stuck

```bash
# Remove stale lock (only if Envy is not running!)
rm ~/.envy/.lock
```

### Corrupted vault

Restore from backup:
```bash
cp ~/.envy/keys.json.backup.20240115 ~/.envy/keys.json
```

---

**Next:** Learn about [Encryption](../security/encryption.md) and how your secrets are protected.
