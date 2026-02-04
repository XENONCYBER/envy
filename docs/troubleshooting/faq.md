# Frequently Asked Questions

Common questions about Envy.

## General

### What is Envy?

Envy is a secure, encrypted secret manager for developers. It stores API keys, passwords, and environment variables in an encrypted vault with both a terminal UI (TUI) and command-line interface (CLI).

### Who is Envy for?

Developers and DevOps professionals who need to:
- Manage multiple sets of credentials
- Keep secrets secure and organized
- Access secrets quickly from the terminal
- Inject secrets into applications

### Is Envy free?

Yes. Envy is open source and free under the MIT License.

### What platforms does Envy support?

Linux (x86_64, ARM64, ARM), macOS (Intel and Apple Silicon), and Windows (x86_64).

## Security

### How secure is Envy?

Envy uses military-grade encryption:
- **AES-256-GCM** for encryption
- **Argon2id** for key derivation (resistant to GPU attacks)
- Your master password is never stored

### Can Envy be hacked?

With a strong master password, the vault is practically unbreakable. However:
- Weak passwords can be brute-forced
- If your system is compromised while vault is unlocked, secrets may be exposed
- Physical access + unlocked computer = access to secrets

### What if I forget my master password?

**Your secrets are lost forever.** We don't store your password and cannot recover it. This is by design for security.

**Prevention:**
- Use a memorable passphrase
- Write it down and store securely
- Use a password manager

### Is the vault file safe to back up?

Yes. The vault file (`keys.json`) is encrypted. You can safely back it up to cloud storage, USB drives, etc. Just ensure:
- Your master password is strong
- Backup locations are secure
- You have multiple backups

### Can I share my vault with teammates?

You can share the vault file, but:
- They'll need your master password
- Everyone has access to all secrets (VERY IMPORTANT keep in mind)
- Consider separate vaults for different access levels

### Is Envy audited?

Envy uses well-established, audited cryptographic libraries (Go standard library and golang.org/x/crypto). The application code itself is open source for community review.

## Usage

### How do I start Envy?

```bash
# First time: creates vault
envy

# Subsequent times: unlocks vault
envy

# CLI commands
envy set myapp KEY=value
envy run myapp -- npm start
```

### How do I copy a secret?

**TUI:**
1. Launch: `envy`
2. Select project with `Enter`
3. Navigate to secret with arrow keys
4. Press `y` to copy
5. Clipboard auto-clears after 30 seconds

**CLI:**
Envy doesn't support direct copy via CLI (security feature). Use TUI or pipe to clipboard tool:
```bash
# Not recommended for security
envy --export myapp && cat .env | grep KEY | xclip
```

### Can I change my master password?

Currently, no. You must:
1. Export all projects
2. Remove vault (`rm ~/.envy/keys.json`)
3. Create new vault with new password
4. Re-import projects

This is a planned feature for future releases.

### How do I delete a project?

**TUI:**
1. Navigate to project in grid
2. Press `d`
3. Confirm with `y`

**CLI:**
No direct delete command. Use TUI or:
1. Export projects you want to keep
2. Delete vault
3. Recreate and re-import

### Can I rename a project?

Yes, you can:
1. Go into project deatails view by pressing enter on the project
2. Press `Shift+E` to enter edit mode for project.
3. Enter new name

### How do I move secrets between environments?

```bash
# Export from source
envy --export myapp  # Exports dev

# Import to target environment
envy --import .env
# When prompted, specify different environment
```

Or use TUI to manually copy values between projects.

## Technical

### Where is my data stored?

- **Linux:** `~/.envy/keys.json`
- **macOS:** `~/.envy/keys.json`
- **Windows:** `%APPDATA%\envy\keys.json`

### What format is the vault?

JSON with encrypted values. Structure:
```json
{
  "version": 1,
  "salt": "...",
  "auth_hash": "...",
  "projects": [...]
}
```

All secret values are AES-256-GCM encrypted.

### Can I use Envy in CI/CD?

Yes:
```bash
# Inject secrets during deployment
envy run production -- ./deploy.sh
```

For fully automated CI/CD, you'll need to:
- Provide master password non-interactively (feature planned)
- Or use environment-specific vaults with different access

### Does Envy support multiple vaults?

Yes, but indirectly. Use different config files with different `backend.keys_path`:

```lua
-- ~/.config/envy/personal.lua
return {
  backend = {
    keys_path = "~/.envy/personal/keys.json"
  }
}

-- ~/.config/envy/work.lua
return {
  backend = {
    keys_path = "~/.envy/work/keys.json"
  }
}
```

Switch by swapping config files or using a wrapper script.

### Can I sync my vault between devices?

Yes. The vault file can be synced like any file:
1. Store in Dropbox/iCloud/etc
2. Or use `rsync`, `scp`, etc

**Important:** The file is encrypted, so safe to sync. Just ensure:
- You use the same master password on all devices
- Sync doesn't corrupt the file
- Handle conflicts carefully

### Does Envy work offline?

Yes. Envy is entirely local. No internet connection required.

## Comparison

### How is Envy different from 1Password/Bitwarden?

| Feature | Envy | 1Password/Bitwarden |
|---------|------|---------------------|
| Local/Cloud | Local only | Cloud sync |
| Target users | Developers | General users |
| CLI/TUI | Yes | Limited CLI |
| Browser integration | No | Yes |
| Mobile apps | No | Yes |
| Price | Free | Subscription |

### How is Envy different from pass/gopass?

| Feature | Envy | pass/gopass |
|---------|------|-------------|
| Encryption | AES-256-GCM | GPG |
| KDF | Argon2id | N/A (GPG) |
| Structure | Projects + Environments | Folders |
| TUI | Built-in | External tools |
| Git integration | No | Yes |

### When should I use Envy vs HashiCorp Vault?

**Use Envy when:**
- Personal or small team use
- Need quick CLI/TUI access
- Want simple, local storage
- Don't need complex access policies

**Use HashiCorp Vault when:**
- Enterprise environment
- Need dynamic secrets
- Require complex ACLs
- Centralized secret management

## Troubleshooting

### Envy won't start

1. Check if vault exists: `ls ~/.envy/keys.json`
2. Check permissions: `ls -la ~/.envy/`
3. Try moving config temporarily:
   ```bash
   mv ~/.config/envy/config.lua ~/.config/envy/config.lua.bak
   envy
   ```

### Forgot master password

Sorry, but your secrets are unrecoverable. This is by design.

Prevention:
- Use memorable passphrase
- Write down and store securely
- Use password manager

### Clipboard not working

Envy uses system clipboard. If `y` (yank) doesn't work:
1. Check you have `xclip` or `xsel` installed (Linux)
2. Try running from terminal (not launcher)
3. Check if clipboard manager is interfering

### Search not working in TUI

1. Ensure you're in Grid view (not Detail)
2. Press `i` to activate search
3. Check if custom keybinding changed it

## Future Features

### What's planned?

- Master password change
- Project renaming
- YubiKey/hardware key support
- Biometric unlock
- Plugin system
- Additional import formats
- Secret sharing
- Audit logging

### How can I request features?

Open an issue on [GitHub](https://github.com/XENONCYBER/envy/issues) with:
- Clear use case
- Proposed solution
- Why it helps

## Contributing

### Can I contribute?

Yes! See [Development Guide](../development.md) for:
- Setting up development environment
- Coding standards
- Submitting PRs

### How do I report bugs?

[GitHub Issues](https://github.com/XENONCYBER/envy/issues) with:
- Envy version
- OS and terminal
- Steps to reproduce
- Expected vs actual behavior

## Miscellaneous

### Why is it called Envy?

Envy = **ENV**ironment + **Y** (because Y not?)

Also, it sounds cool and the logo looks nice in ASCII art.

### How do I pronounce Envy?

Seriously,
/ˈɛnvi/ — like the word "envy"

### Is there a mobile app?

No. Envy is terminal-based and designed for desktop/server use.

### Can I use Envy for passwords?

Yes, but it's designed for API keys and environment variables. For general password management, dedicated password managers may be more convenient.

### Does Envy auto-lock?

Not currently. The vault locks when you quit Envy. Future versions may add inactivity timeout.

---

**Still have questions?** Open an issue on GitHub or check [Common Issues](./common-issues.md).
