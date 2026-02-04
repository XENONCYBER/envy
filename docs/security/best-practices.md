# Security Best Practices

Comprehensive security guidance for using Envy effectively.

## Quick Reference

| Do | Don't |
|-------|----------|
| Use a strong master password | Use weak or common passwords |
| Backup your vault regularly | Store vault in public/cloud without encryption |
| Set correct file permissions | Share your vault file |
| Use different passwords per environment | Reuse passwords across services |
| Lock your computer when away | Leave secrets visible on screen |
| Clear clipboard after copying | Paste secrets into untrusted apps |
| Keep Envy updated | Use outdated versions |

## Master Password

### Create a Strong Password

**Use a passphrase:**
```
Correct-Horse-Battery-Staple-47
Five-Random-Words-Are-Secure!
```

**Or random characters (16+):**
```
xK9#mP2$vL5@nQ8*
```

**Requirements:**
- 12+ characters minimum
- Unique to Envy (don't reuse)
- Not based on personal information
- Memorable or stored securely

### Store Password Securely

**Option 1: Memorize**
- Use a passphrase you'll remember

**Option 2: Physical backup**
```
Write on paper → Seal in envelope → Store in safe
```

**Option 3: Password manager**
- Store in your primary password manager

**Not recommended:**
- Plain text files
- Unencrypted notes apps
- Browser password managers
- Email drafts

## Vault Protection

### File System Security

**Verify permissions:**
```bash
ls -la ~/.envy/
# Should show:
# drwx------  2 user user 4096 .        (0700)
# -rw-------  1 user user 1234 keys.json (0600)
# -rw-------  1 user user    5 .lock     (0600)
```

**Fix if wrong:**
```bash
chmod 700 ~/.envy/
chmod 600 ~/.envy/keys.json
chmod 600 ~/.envy/.lock
```

### Backup Strategy

**3-2-1 Rule:**
- **3** copies of data
- **2** different media types
- **1** offsite

**Implementation:**
```bash
# 1. Primary: ~/.envy/keys.json

# 2. Local backup (encrypted USB)
cp ~/.envy/keys.json /mnt/encrypted-usb/envy-backup-$(date +%Y%m%d).json

# 3. Offsite backup (encrypted cloud)
gpg -c ~/.envy/keys.json
# Upload keys.json.gpg to secure cloud storage
```

**Backup automation:**
```bash
#!/bin/bash
# ~/.local/bin/backup-envy.sh
DATE=$(date +%Y%m%d)
cp ~/.envy/keys.json ~/Backups/envy/keys-$DATE.json
ls -t ~/Backups/envy/keys-*.json | tail -n +31 | xargs rm -f
```

Add to crontab:
```bash
0 2 * * * /home/user/.local/bin/backup-envy.sh
```

### Full Disk Encryption

**Essential for vault security:**

| OS | Method | Setup |
|----|--------|-------|
| Linux | LUKS | During installation or `cryptsetup` |
| macOS | FileVault | System Preferences → Security |
| Windows | BitLocker | Control Panel → BitLocker |

**Why:** If your laptop is stolen, thief can't read vault without login password.

## Operational Security

### Clipboard Management

**Envy's protection:**
- Auto-clears clipboard after 30 seconds. (Currently only works when envy is running acitvely as a process)
- Shows warning: "Copied (will clear in 30s)"

**Your responsibility:**
- Don't paste into untrusted applications
- Clear clipboard manually if needed: `echo "" | xclip -selection clipboard`
- Be aware of clipboard history tools (Ditto, CopyQ) - they may retain secrets

### Shoulder Surfing Prevention

**When entering password:**
- Check surroundings
- Use body to shield keyboard
- Consider privacy screens

**When secrets are visible:**
- Reveal secrets only when needed
- Hide immediately after use (press Enter/Space again)
- Don't leave detail view open

## Environment Security

### CI/CD Security

**In automated pipelines:**
```bash
# Good: Direct injection
envy run production -- ./deploy.sh

# Risky: Export to file
envy --export production
deploy.sh
rm .env  # Hope this runs!
```

**Best practices:**
- Use `envy run` for direct injection
- If exporting, use secure temporary storage
- Delete exported files immediately
- Don't log secrets

## Network Security

### Sync and Cloud Storage

**Vault file in cloud:**
```bash
# Encrypt before uploading
gpg -c ~/.envy/keys.json
# Upload keys.json.gpg

# Download and decrypt
gpg keys.json.gpg
# Move to ~/.envy/
```

**Sync services:**
- Keep a backup file saved locally somewhere(recommended)
- iCloud (encrypted)
- Dropbox (encrypted at rest)
- Google Drive (encrypted at rest)
- Self-hosted Nextcloud (ensure HTTPS)

### Remote Access

**SSH sessions:**
- Use `envy run` to inject secrets remotely
- Don't transfer vault file over unencrypted channels
- Verify host keys

```bash
# Good: Run remote command with local secrets
ssh server 'envy run production -- ./restart.sh'

# Risky: Copying vault to server
scp ~/.envy/keys.json server:~/.envy/
```

## Physical Security

### Device Protection

**Laptop/desktop:**
- Full disk encryption (essential)
- Screen lock on sleep
- BIOS password
- TPM/Secure Boot if available

**Mobile devices:**
- If accessing via SSH, use secure apps (Termius, Blink)

## Threat Response

### Suspected Compromise

**If you think your vault is compromised:**

1. **Immediately:**
   - Don't panic
   - Don't access the vault
   - Isolate the system if possible

2. **Assess:**
   - Who had access?
   - What systems use these secrets?
   - When did compromise occur?

3. **Contain:**
   - Rotate all secrets in vault
   - Check access logs
   - Scan for malware

4. **Recover:**
   - Create new vault with new password
   - Re-add all secrets (rotated values)
   - Improve security practices

### Weak Password Realization

**If you realize your password is weak:**

1. **Don't panic**
2. **Export all projects immediately**
3. **Create new vault with strong password**
4. **Re-import all projects**
5. **Delete old vault**
6. **Monitor for suspicious activity**

## Advanced Security

### Separate Vaults

For high-security separation:

```bash
# Personal vault
~/.envy/personal/keys.json

# Work vault
~/.envy/work/keys.json

# Client A vault
~/.envy/client-a/keys.json

# Client B vault
~/.envy/client-b/keys.json
```

Switch via config files or scripts.

### Air-Gapped Usage

For maximum security:
- Use Envy on air-gapped machine
- Export to encrypted USB
- Transfer to production manually

**Workflow:**
```bash
# On secure machine
envy --export production
gpg -c .env
# Copy to USB

# On production
# Decrypt and use
```

### Hardware Security Keys

Future feature: YubiKey integration
- Store encryption key on hardware
- Two-factor authentication
- Physical presence required

## Security Resources

- [OWASP Secrets Management](https://cheatsheetseries.owasp.org/cheatsheets/Secrets_Management_Cheat_Sheet.html)
- [NIST Password Guidelines](https://pages.nist.gov/800-63-3/sp800-63b.html#sec5)
- [CISA Cybersecurity](https://www.cisa.gov/cybersecurity)
- [Have I Been Pwned](https://haveibeenpwned.com/) - Check if credentials leaked

---

**Questions?** Check [FAQ](../troubleshooting/faq.md) or [Common Issues](../troubleshooting/common-issues.md).
