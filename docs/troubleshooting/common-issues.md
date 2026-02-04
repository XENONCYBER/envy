# Common Issues

Solutions to frequently encountered problems.

## Installation Issues

### "Command not found" after installation

**Problem:** `envy` command not recognized after install.

**Solutions:**

1. Check if `/usr/local/bin` is in PATH:
```bash
echo $PATH | grep /usr/local/bin
```

2. Add to PATH:
```bash
# Bash
echo 'export PATH="$PATH:/usr/local/bin"' >> ~/.bashrc
source ~/.bashrc

# Zsh
echo 'export PATH="$PATH:/usr/local/bin"' >> ~/.zshrc
source ~/.zshrc

# Fish
set -U fish_user_paths /usr/local/bin $fish_user_paths
```

3. Verify installation:
```bash
ls -la /usr/local/bin/envy
which envy
```

### Permission denied

**Problem:** Cannot run `envy` due to permissions.

**Solution:**
```bash
sudo chmod +x /usr/local/bin/envy
```

### macOS "developer cannot be verified"

**Problem:** macOS blocks the binary.

**Solutions:**

1. Remove quarantine attribute:
```bash
xattr -d com.apple.quarantine /usr/local/bin/envy
```

2. Or allow in System Preferences:
   - Apple menu → System Preferences → Security & Privacy
   - Click "Allow Anyway" for Envy

### Build fails

**Problem:** `go build` fails.

**Check:**
```bash
# Go version (need 1.25.4+)
go version

# If too old, update from https://go.dev/dl/

# Try clean build
go clean
go build -o envy ./cmd/main.go
```

## Vault Issues

### "Authentication failed: incorrect password"

**Problem:** Cannot unlock vault.

**Causes & Solutions:**

1. **Wrong password:**
   - Check Caps Lock
   - Try common variations you might have used

2. **Wrong vault file:**
   - Check `~/.envy/keys.json` exists
   - Verify you're using correct machine/account

3. **Corrupted vault:**
   - Restore from backup:
   ```bash
   cp ~/.envy/keys.json.backup.20240115 ~/.envy/keys.json
   ```

4. **Last resort:**
   - If no backup and password truly forgotten, vault is unrecoverable
   - Remove and recreate: `rm ~/.envy/keys.json`

### "Failed to parse storage file"

**Problem:** Vault file is corrupted.

**Solutions:**

1. Restore from backup:
```bash
cp ~/.envy/keys.json.latest ~/.envy/keys.json
```

2. Check file isn't empty:
```bash
ls -la ~/.envy/keys.json
# Should show non-zero size
```

3. Verify JSON structure:
```bash
python3 -m json.tool ~/.envy/keys.json > /dev/null && echo "Valid JSON"
```

4. If corrupted beyond repair:
   - Export what you can from memory/notes
   - Create new vault
   - Restore from any .env backups

### Lock file issues

**Problem:** "Failed to acquire lock" or stale lock file.

**Symptoms:**
- "Vault is locked by another process"
- Envy crashed and won't restart

**Solutions:**

1. Check if Envy is actually running:
```bash
ps aux | grep envy
```

2. If not running, remove stale lock:
```bash
rm ~/.envy/.lock
```

3. If running, wait or kill it:
```bash
killall envy
rm ~/.envy/.lock
```

**Note:** Only remove lock if you're certain no other Envy instance is running!

### Vault file permissions

**Problem:** Permission errors accessing vault.

**Check:**
```bash
ls -la ~/.envy/
```

**Should show:**
```
drwx------  2 user user 4096 .        (0700)
-rw-------  1 user user 1234 keys.json (0600)
-rw-------  1 user user    5 .lock     (0600)
```

**Fix:**
```bash
chmod 700 ~/.envy/
chmod 600 ~/.envy/keys.json
chmod 600 ~/.envy/.lock
```

## TUI Issues

### Display garbled / artifacts

**Problem:** Interface doesn't render correctly.

**Causes:**
1. Terminal doesn't support Unicode
2. Terminal has limited colors
3. Font issues

**Solutions:**

1. Use a modern terminal:
   - iTerm2 (macOS)
   - Windows Terminal (Windows)
   - Alacritty, Kitty, WezTerm (cross-platform)

2. Check terminal capabilities:
```bash
echo $TERM
# Should show: xterm-256color or similar
```

3. Set TERM if needed:
```bash
export TERM=xterm-256color
```

### Keys not working

**Problem:** Keyboard shortcuts don't respond.

**Common causes:**

1. **Terminal intercepts keys:**
   - `Ctrl+s` freezes terminal (press `Ctrl+q`)
   - TMUX prefix conflicts (`Ctrl+b`)
   - Screen prefix conflicts (`Ctrl+a`)

2. **SSH issues:**
   - Terminal type not set correctly
   - Key codes not transmitted

**Solutions:**

1. Remap conflicting keys in `config.lua`:
```lua
keys = {
  save = "ctrl+x",      -- Instead of ctrl+s
  create = "ctrl+n",    -- Keep if not conflicting
  add = "ctrl+e"        -- Instead of ctrl+a
}
```

2. For TMUX, add to `~/.tmux.conf`:
```bash
set -g prefix C-a        # Change from C-b
unbind C-b
```

3. For SSH, set TERM:
```bash
export TERM=xterm-256color
```

### TUI won't open

**Problem:** Running `envy` shows error or exits immediately.

**Check:**

1. Vault exists:
```bash
ls ~/.envy/keys.json
```

2. Config is valid (if exists):
```bash
lua -c ~/.config/envy/config.lua 2>&1
```

3. Try without config:
```bash
# Temporarily move config
mv ~/.config/envy/config.lua ~/.config/envy/config.lua.bak
envy
# Restore if needed
mv ~/.config/envy/config.lua.bak ~/.config/envy/config.lua
```

### Search not working

**Problem:** Pressing `i` doesn't activate search.

**Check:**
- Are you in Grid view? (Search only works there)
- Is search mode already active? (Press `Esc` first)
- Try `/` instead of `i`

**Custom keybinding fix:**
```lua
keys = {
  search = "/"  -- Use / instead of i
}
```

## CLI Issues

### "Missing '--' separator"

**Problem:** Running `envy run` without `--`.

**Wrong:**
```bash
envy run myapp npm start
```

**Right:**
```bash
envy run myapp -- npm start
```

### "Project not found"

**Problem:** Project name doesn't exist.

**Solutions:**

1. Check spelling (case-insensitive, but must match)
2. List projects in TUI to see available names
3. Check you're using correct environment:
```bash
# Maybe it's in different environment
envy set myapp KEY=value -e prod
envy run myapp -- ./script.sh  # Looks in dev by default
```

### Invalid KEY=VALUE format

**Problem:** Error setting secret.

**Wrong:**
```bash
envy set myapp API_KEY  # Missing =value
envy set myapp =value   # Missing key
envy set myapp A=B=C    # Extra = in value (actually OK)
```

**Right:**
```bash
envy set myapp API_KEY=secret123
envy set myapp "KEY WITH SPACES=value"
envy set myapp A=B=C    # Value is "B=C"
```

### Import fails

**Problem:** `envy --import .env` fails.

**Check:**

1. File exists and is readable:
```bash
cat .env
```

2. File format is valid:
```bash
# Should show KEY=value pairs
cat .env
```

3. Try absolute path:
```bash
envy --import $(pwd)/.env
```

## Performance Issues

### Slow startup

**Problem:** Envy takes long to open.

**Causes:**
- Very large vault (thousands of projects)
- Slow disk I/O
- Network filesystem (don't store vault on NFS!)

**Solutions:**

1. Ensure vault is on local disk
2. Archive old projects
3. Split into multiple vaults if needed

### Lag in TUI

**Problem:** Interface feels slow.

**Solutions:**

1. Close other heavy applications
2. Use faster terminal emulator
3. Reduce grid size in config:
```lua
theme = {
  grid_cols = 2,          -- Instead of 3
  grid_visible_rows = 2   -- Keep default
}
```

### Large vault handling

**Problem:** Vault is very large (>10MB).

**Solutions:**

1. Archive old projects (export and delete)
2. Split by environment or team
3. Check for unnecessary history accumulation

## Configuration Issues

### Config not loading

**Problem:** Changes to `config.lua` not taking effect.

**Check:**

1. File location:
```bash
# Linux
ls ~/.config/envy/config.lua

# macOS
ls ~/Library/Application\ Support/envy/config.lua

# Windows
dir %APPDATA%\envy\config.lua
```

2. Syntax errors:
```bash
lua -c ~/.config/envy/config.lua
```

3. Must return a table:
```lua
-- WRONG
local config = { keys = { quit = "q" } }

-- RIGHT
return { keys = { quit = "q" } }
```

### Theme not applying

**Problem:** Colors look wrong or default theme shows.

**Check:**

1. Valid hex colors:
```lua
-- WRONG
base = "#1e1e2e  -- Missing closing quote

-- RIGHT
base = "#1e1e2e"
```

2. All required fields:
```lua
theme = {
  base = "#1e1e2e",
  text = "#cdd6f4",
  accent = "#cba6f7",
  -- ... all other fields
}
```

### Keybindings not working

**Problem:** Custom keys don't respond.

**Check:**

1. Valid key names:
```lua
-- WRONG
keys = {
  quit = "ctrl+q"  -- '+' is not valid syntax
}

-- RIGHT
keys = {
  quit = "ctrl+q"  -- Actually fine, just no spaces
}
```

2. No duplicates:
```lua
-- WRONG
keys = {
  quit = "q",
  search = "q"  -- Conflict!
}
```

3. Case sensitivity:
```lua
keys = {
  quit = "q",     -- Lowercase q
  search = "Q"    -- Uppercase Q (Shift+q)
}
```

## Still Having Issues?

### Debug Mode

Run with verbose output (if supported in future):
```bash
ENVY_DEBUG=1 envy
```

### Get Help

1. Check [FAQ](./faq.md)
2. Search [GitHub Issues](https://github.com/XENONCYBER/envy/issues)
3. Create new issue with:
   - Envy version (`envy --version`)
   - Operating system
   - Terminal emulator
   - Exact error message
   - Steps to reproduce

### Emergency Recovery

If vault is corrupted and no backup:

```bash
# 1. Save corrupted file for analysis
mv ~/.envy/keys.json ~/.envy/keys.json.corrupted.$(date +%Y%m%d)

# 2. Check for any .env exports you might have
find ~ -name "*.env" -type f

# 3. Rebuild vault from scratch
envy  # Creates new vault
# Manually re-add secrets
```

---

**Next:** Check [FAQ](./faq.md) for more questions and answers.
