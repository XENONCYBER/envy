# Recovery Guide

How to recover from various failure scenarios.

## Scenario: Forgotten Master Password

### Reality Check

**If you forgot your master password and have no backup, your secrets are permanently lost.**

Envy uses encryption that is designed to be unbreakable. We cannot:
- Reset your password
- Recover your secrets
- Bypass the encryption

### Prevention (Before It Happens)

1. **Write it down:**
   ```
   Paper → Envelope → Safe deposit box
   ```

2. **Use a passphrase you'll remember:**
   ```
   Correct-Horse-Battery-Staple-47
   ```

3. **Store in password manager:**
   - 1Password, Bitwarden, KeePass, etc.

4. **Regular use:**
   - Use Envy regularly to reinforce memory

### What To Try (If You Think You Remember)

1. **Common variations:**
   - Different capitalization
   - With/without numbers at end
   - Old passwords you might have reused

2. **Check everywhere:**
   - Password manager
   - Browser saved passwords
   - Notes apps
   - Physical notes
   - Safe deposit box
   - Cloud documents

3. **Think about hints:**
   - What was important to you at the time?
   - Any personal references?

### Recovery Steps

If truly forgotten:

```bash
# 1. Accept data loss
# 2. Remove old vault
rm ~/.envy/keys.json

# 3. Create new vault
envy
# Create new master password

# 4. Rebuild from any sources:
# - .env files you exported
# - Documentation
# - Team members (if shared)
# - Cloud provider dashboards (rotate keys)
```

## Scenario: Corrupted Vault

### Symptoms

- "Failed to parse storage file"
- Vault file is empty (0 bytes)
- JSON syntax errors
- Strange characters in file

### Recovery Steps

#### Step 1: Don't Panic

Stop using the vault immediately to prevent further corruption.

#### Step 2: Assess the Damage

```bash
# Check file size
ls -la ~/.envy/keys.json

# Check if it's valid JSON
python3 -m json.tool ~/.envy/keys.json > /dev/null && echo "Valid JSON"

# Try to read structure
head -c 1000 ~/.envy/keys.json
```

#### Step 3: Find Backups

```bash
# List all backup files
ls -la ~/.envy/
ls -la ~/Backups/envy/
ls -la ~/Dropbox/
find ~ -name "*keys*" -o -name "*envy*" 2>/dev/null

# Check timestamps
ls -lt ~/.envy/keys.json*
```

#### Step 4: Restore from Backup

```bash
# Identify best backup (before corruption)
# Usually the most recent backup before the issue

# Restore
cp ~/.envy/keys.json.backup.20240115 ~/.envy/keys.json

# Verify permissions
chmod 600 ~/.envy/keys.json

# Try to unlock
envy
```

#### Step 5: If No Backup Exists

Check for exported .env files:

```bash
# Search for any .env exports
find ~ -name ".env" -type f 2>/dev/null
find ~ -name "*.env" -type f 2>/dev/null

# Check common locations
ls -la ~/projects/*/.env
ls -la ~/workspace/*/.env
```

If you find exports:
```bash
# Recreate vault from exports
envy --import ~/projects/myapp/.env
# Repeat for each project
```

#### Step 6: Rebuild from Scratch

If no backups or exports:

1. Document what secrets you had
2. Rotate all API keys (generate new ones)
3. Update database passwords
4. Reconfigure services
5. Create new Envy vault
6. Store new secrets properly
7. **Set up backups immediately!**

## Scenario: Accidental Deletion

### Deleted Project

**If in TUI:**
1. You probably saw a confirmation
2. Recovery only possible from backup

**Restore from backup:**
```bash
# Find backup from before deletion
cp ~/.envy/keys.json.backup.20240115 ~/.envy/keys.json
```

**Or recreate:**
```bash
# Add secrets back manually
envy set myapp KEY1=value1
envy set myapp KEY2=value2
```

### Deleted Key

**Restore:**
1. If deleted recently, check if you have the value elsewhere
2. Or restore from backup
3. Or regenerate the key (API keys, passwords)

**Note:** Unlike password updates, deletions are not saved to history.

## Scenario: Lost or Stolen Device

### Immediate Actions

1. **Don't panic**

2. **Change all passwords immediately** (assume vault will be cracked):
   - Database passwords
   - API keys
   - Service credentials
   - Cloud provider access
   - SSH keys

3. **Revoke access:**
   - Invalidate API tokens
   - Rotate SSH keys
   - Change database credentials
   - Update cloud access keys

4. **Monitor for usage:**
   - Check access logs
   - Look for unusual activity
   - Set up alerts

5. **File reports:**
   - Police report (for insurance)
   - Employer (if work device)

### If Vault Was Synced

If you used cloud sync and device was compromised:

1. **Assume vault is compromised**
2. **Rotate ALL secrets immediately**
3. **Remove synced vault from cloud:**
   ```bash
   # Remove from Dropbox
   rm ~/Dropbox/envy/keys.json
   
   # Remove from iCloud
   rm ~/Library/Mobile\ Documents/com~apple~CloudDocs/envy/keys.json
   ```

4. **Create new vault** with new master password

### Rebuilding

```bash
# 1. On new device, install Envy
curl -fsSL ... | sh

# 2. Create new vault
envy

# 3. Generate new secrets everywhere
# - Cloud provider consoles
# - Database admin panels
# - API key management
# - SSH keygen

# 4. Add new secrets to Envy
envy set myapp NEW_API_KEY=...

# 5. Set up proper backups
```

## Scenario: Lock File Stuck

### Symptoms
- "Failed to acquire lock"
- "Vault is locked by another process"
- Envy crashed and won't restart

### Solution

```bash
# 1. Check if Envy is actually running
ps aux | grep envy
pgrep -f envy

# 2. If running, wait or kill
killall envy
# or
pkill -f envy

# 3. Remove stale lock file
rm ~/.envy/.lock

# 4. Try again
envy
```

### Prevention

- Don't force-quit Envy (use `q` to quit properly)
- If system crashes, manually remove lock on reboot
- Consider alias:
```bash
alias envy-clean='rm -f ~/.envy/.lock; envy'
```

## Scenario: Sync Conflict

### Dropbox/Google Drive Conflict

If vault is synced and conflicts occur:

```bash
# 1. Stop all Envy instances everywhere

# 2. On each machine, rename vault
mv ~/.envy/keys.json ~/.envy/keys.json.machine1

# 3. Let sync resolve

# 4. Compare file sizes and dates
ls -la ~/.envy/keys.json*

# 5. Choose most recent valid vault
cp ~/.envy/keys.json.most_recent ~/.envy/keys.json

# 6. Delete others
rm ~/.envy/keys.json.machine*

# 7. Sync and verify on all machines
```

### Prevention

- Close Envy before system sleep/hibernate
- Avoid using Envy on multiple machines simultaneously
- Use different vaults for different machines (if needed)

## Scenario: Bad Config

### Symptoms
- Envy won't start
- Strange behavior
- Error messages about Lua

### Recovery

```bash
# 1. Backup current config
cp ~/.config/envy/config.lua ~/.config/envy/config.lua.broken

# 2. Remove or reset config
rm ~/.config/envy/config.lua

# 3. Try Envy with defaults
envy

# 4. If works, rebuild config gradually
# Start with minimal config and add options
```

### Minimal Working Config

```lua
-- Start with this, then add more
return {
  keys = {
    quit = "q"
  }
}
```

## Backup Strategy

### 3-2-1 Rule

After any recovery, implement proper backups:

```bash
#!/bin/bash
# ~/.local/bin/backup-envy.sh

# 3 copies
# 2 different media
# 1 offsite

DATE=$(date +%Y%m%d-%H%M%S)

# Copy 1: Local backup
mkdir -p ~/Backups/envy
cp ~/.envy/keys.json ~/Backups/envy/keys-$DATE.json

# Copy 2: External drive (if mounted)
if [ -d /mnt/backup ]; then
  cp ~/.envy/keys.json /mnt/backup/envy/keys-$DATE.json
fi

# Copy 3: Encrypted cloud (manual or scripted)
gpg --batch --yes --passphrase-file ~/.config/envy/backup-passphrase \
  -c ~/.envy/keys.json
# Upload keys.json.gpg to cloud

# Keep only last 30
ls -t ~/Backups/envy/keys-*.json | tail -n +31 | xargs rm -f

# Log
logger "Envy backup completed: $DATE"
```

### Automation

```bash
# Add to crontab - daily at 2 AM
0 2 * * * /home/user/.local/bin/backup-envy.sh

# Or use systemd timer (Linux)
# Or launchd (macOS)
```

## Testing Recovery

### Practice Recovery

Regularly test your recovery process:

```bash
# Monthly test

# 1. Create temp directory
mkdir -p /tmp/envy-test
cd /tmp/envy-test

# 2. Copy latest backup
cp ~/Backups/envy/keys-$(date +%Y%m%d).json ./keys.json

# 3. Try to use it
# Point Envy to this vault (via config)

# 4. Verify you can unlock and read secrets

# 5. Clean up
cd /
rm -rf /tmp/envy-test
```

### Document Your Setup

Create a recovery document:

```markdown
# Envy Recovery Info

## Master Password Hint
[Your hint here - not the actual password]

## Backup Locations
- Local: ~/Backups/envy/
- External: /mnt/backup/envy/
- Cloud: [Service] at [Path]

## Recovery Steps
1. Install Envy
2. Restore from [location]
3. Verify permissions: chmod 600 keys.json
4. Test unlock

## Emergency Contacts
- [Name]: [Contact] - Has backup copy

Last updated: [Date]
```

Store this securely (printed in safe, encrypted file, etc.).

---

**Remember:** Backups are only useful if they work. Test them regularly!
