# Encryption

Technical deep-dive into Envy's encryption implementation.

## Overview

Envy uses industry-standard encryption algorithms designed to protect your secrets even if the vault file is compromised.

**Security Stack:**
- **AES-256-GCM** — Symmetric encryption with authentication
- **Argon2id** — Modern password-based key derivation
- **SHA-256** — Password verification hash
- **Cryptographically secure RNG** — For salts and nonces

## Encryption Flow

### When Saving (Encrypt)

```
1. User enters master password
2. Argon2id(password, salt) → 256-bit encryption key
3. For each secret:
   a. Generate random nonce (12 bytes)
   b. AES-256-GCM encrypt(secret, key, nonce) → ciphertext
   c. Store: nonce + ciphertext (base64 encoded)
4. Store auth_hash = SHA256(key) for verification
```

### When Loading (Decrypt)

```
1. User enters master password
2. Load salt from vault file
3. Argon2id(password, salt) → 256-bit key
4. Verify: SHA256(key) == auth_hash?
   - No → "Incorrect password"
   - Yes → Continue
5. For each encrypted secret:
   a. Extract nonce from beginning
   b. AES-256-GCM decrypt(ciphertext, key, nonce)
   c. Return plaintext secret
```

## Algorithm Details

### AES-256-GCM

**What it is:**
Advanced Encryption Standard with Galois/Counter Mode and 256-bit keys.

**Why we use it:**
- NIST approved, widely audited
- Authenticated encryption (detects tampering)
- Fast in hardware (AES-NI)
- No known practical attacks

**Parameters:**
- Key size: 256 bits (32 bytes)
- Block size: 128 bits
- Nonce size: 96 bits (12 bytes)
- Tag size: 128 bits (authentication tag)

**Implementation:**
```go
// Using Go's crypto/aes and crypto/cipher
block, _ := aes.NewCipher(key)        // AES-256
gcm, _ := cipher.NewGCM(block)        // GCM mode
nonce := make([]byte, gcm.NonceSize()) // 12 bytes
rand.Read(nonce)                       // Random nonce
ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
```

### Argon2id

**What it is:**
Winner of the Password Hashing Competition (2015), modern successor to bcrypt/scrypt/PBKDF2.

**Why we use it:**
- Resistant to GPU attacks
- Resistant to ASIC attacks
- Memory-hard (requires RAM)
- Adjustable parameters

**Parameters:**
```go
const (
    time    = 1          // Iterations
    memory  = 64 * 1024  // 64 MB
    threads = 4          // Parallelism
    keyLen  = 32         // 256-bit output
)

key := argon2.IDKey(password, salt, time, memory, threads, keyLen)
```

**Security analysis:**
- 64 MB memory usage makes GPU attacks expensive
- Single iteration with memory hardness is sufficient
- 4 threads utilize modern CPUs efficiently

### Salt Generation

**Purpose:**
Ensures unique keys even if two users have the same password.

**Generation:**
```go
salt := make([]byte, 16)  // 128-bit salt
rand.Read(salt)           // Cryptographically secure RNG
```

**Storage:**
- Stored as base64 in vault file
- Not secret (can be public)
- Must be unique per vault

### Password Verification

**Why not store the key?**
- Key is derived from password
- Storing key = storing password equivalent
- We only store a hash of the key

**Process:**
```go
// On vault creation:
authHash := sha256(key)
store(authHash)

// On unlock:
derivedKey := argon2id(password, salt)
if sha256(derivedKey) == storedAuthHash {
    password is correct
}
```

**Benefits:**
- Fast verification (single SHA-256)
- No password equivalent stored
- Still requires full Argon2id to derive key

## Threat Model

### What Envy Protects Against

| Threat | Protection |
|--------|------------|
| **Vault theft** | AES-256 encryption requires password |
| **Brute force** | Argon2id makes attempts expensive |
| **GPU cracking** | Memory-hard KDF |
| **Tampering** | GCM authentication detects modification |
| **Known plaintext** | Unique nonce per encryption |

### What Envy Cannot Protect Against

| Threat | Reason |
|--------|--------|
| **Weak passwords** | If password is "123456", vault can be cracked |
| **Keyloggers** | Password captured during entry |
| **Memory dumps** | Secrets in RAM while unlocked |
| **Shoulder surfing** | Password visible during entry |
| **Backup theft** | Same protection as original vault |

## Security Recommendations

### Master Password

**Requirements enforced:**
- Minimum 8 characters (configurable in code)

**Recommendations:**
- Use 12+ characters
- Mix upper, lower, numbers, symbols
- Use a passphrase (4-5 random words)
- Unique to Envy (don't reuse)

**Examples:**
```
Good: Correct-Horse-Battery-Staple!47
Good: Tr0ub4dor&3 (but prefer passphrases)
Bad: password123
Bad: 12345678
```

### Vault Backup

**Best practices:**
1. Keep 3 copies: primary, local backup, offsite backup
2. Encrypt backups (they're already encrypted!)
3. Test restore procedure
4. Version backups (keep last 30 days)

**Script:**
```bash
#!/bin/bash
# Daily backup
DATE=$(date +%Y%m%d)
cp ~/.envy/keys.json ~/Backups/envy-keys-$DATE.json
# Keep only last 30
ls -t ~/Backups/envy-keys-*.json | tail -n +31 | xargs rm -f
```

### File Permissions

**Verify regularly:**
```bash
ls -la ~/.envy/
# Should show:
# drwx------  2 user user  4096 Jan 15 10:00 .
# -rw-------  1 user user 12345 Jan 15 10:00 keys.json
# -rw-------  1 user user     5 Jan 15 10:00 .lock
```

### Operating System Security

**Recommendations:**
- Use full-disk encryption (LUKS, BitLocker, FileVault)
- Lock screen when away
- Keep OS and software updated
- Use security-focused OS (Qubes, Tails for extreme cases)

## Cryptographic Audit

### External Libraries

Envy uses battle-tested libraries:

| Library | Purpose | Status |
|---------|---------|--------|
| Go crypto/aes | AES implementation | Standard library, audited |
| Go crypto/cipher | GCM mode | Standard library, audited |
| golang.org/x/crypto/argon2 | Argon2id | Official extension, audited |
| Go crypto/sha256 | Hashing | Standard library, audited |
| Go crypto/rand | RNG | Standard library, audited |

### Implementation Review

**Attack surfaces:**
1. **Key derivation** — Argon2id with safe parameters
2. **Encryption** — AES-256-GCM with random nonces
3. **Password verification** — SHA-256 of derived key
4. **RNG** — OS-provided CSPRNG

**No custom crypto:**
- All primitives from standard libraries
- No home-grown algorithms
- Standard construction (KDF → AEAD)

## Comparison with Alternatives

| Tool | Encryption | KDF | Notes |
|------|-----------|-----|-------|
| **Envy** | AES-256-GCM | Argon2id | Modern, authenticated |
| pass | GPG (varies) | N/A | Depends on GPG settings |
| gopass | GPG (varies) | N/A | Depends on GPG settings |
| Vault | AES-256-GCM | N/A | Server-side, different model |
| 1Password | AES-256 | PBKDF2 | Commercial, cloud-based |
| Bitwarden | AES-256 | PBKDF2 | Open source, cloud option |

**Envy advantages:**
- Modern Argon2id (vs older PBKDF2)
- Local-only (vs cloud sync attack surface)
- Transparent (open source, auditable)

## Future Security Improvements

Potential enhancements under consideration:

1. **Time-based Unlock**
   - Auto-lock after inactivity
   - Require password re-entry

2. **Audit Logging**
   - Log access attempts
   - Alert on suspicious activity

3. **Plausible Deniability**
   - Hidden vault within vault
   - Different passwords → different data

## For Security Researchers

### Code Review

Security-critical files:
- `internal/crypto/encryption.go`
- `internal/storage/store.go`
- `internal/auth/password.go`

We welcome security-focused code reviews and PRs.

## Further Reading

- [Argon2 RFC](https://tools.ietf.org/html/rfc9106)
- [AES-GCM Paper](https://csrc.nist.gov/publications/detail/sp/800-38d/final)
- [OWASP Password Storage](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
- [Go Cryptography Principles](https://go.dev/blog/cryptography-principles)

---

**Next:** Learn about [Master Password Best Practices](./master-password.md).
