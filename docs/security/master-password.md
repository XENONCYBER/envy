# Master Password

Everything you need to know about Envy's master password.

## What is the Master Password?

The master password is the single password that protects all your secrets in Envy. It's used to:

1. **Derive encryption key** — Your password is transformed into a 256-bit encryption key
2. **Verify identity** — Prove you are authorized to access the vault
3. **Protect secrets** — Without it, vault contents are unreadable

## Password Requirements

### Minimum Requirements

Envy enforces:
- **Minimum length:** 8 characters
- **Cannot be empty**
- **Confirmation required** on creation

### Recommended Requirements

For strong security:
- **12+ characters**
- **Mix of character types:**
  - Uppercase (A-Z)
  - Lowercase (a-z)
  - Numbers (0-9)
  - Symbols (!@#$%^&*)
- **Not easily guessable:**
  - Not dictionary words
  - Not personal info (birthdays, names)
  - Not keyboard patterns (qwerty, 123456)

## Password Strength

### Strong Password Examples

**Passphrases (Recommended):**
```
Correct-Horse-Battery-Staple-47
Purple-Monkey-Dishwasher!2024
Seven-Angry-Dolphins-Spin-Fast
```

**Complex passwords:**
```
Tr0ub4dor&3_Complexity!
xK9#mP2$vL5@nQ8*
J3s8#kL0!mN5$pQ2
```

### Weak Password Examples

**Never use:**
```
password123
12345678
qwerty
letmein
admin
envy123
yourname2024
```

**Also avoid:**
- Single dictionary words
- Personal information
- Keyboard patterns
- Reused passwords from other services

## How It Works

### Password to Key Derivation

```
Your Password
      ↓
Argon2id(password, salt)
      ↓
256-bit Encryption Key
      ↓
AES-256-GCM (encrypt/decrypt secrets)
```

**Key points:**
- Your password is never stored
- Only a hash of the derived key is stored (for verification)
- Same password + same salt = same key
- Different vaults have different salts

### Password Verification

```
Login:
1. Enter password
2. Derive key with stored salt
3. Hash the key: SHA256(key)
4. Compare to stored auth_hash
5. Match? Decrypt vault
   No match? "Incorrect password"
```

## Creating a Strong Password

### Method 1: Passphrase (Recommended)

Use 4-5 random words with separators:

```
correct-horse-battery-staple
green-elephants-dance-quickly
pizza-delivery-unicorn-wizard
```

**Why this works:**
- Easy to remember
- High entropy (randomness)
- Resistant to brute force
- Can add numbers/symbols for extra strength

**Generating:**
- Use a passphrase generator (e.g., Diceware)
- Pick random words from a book
- Combine unrelated concepts

### Method 2: Random Characters

Use a password manager to generate:

```
Length: 16-20 characters
Include: Upper, lower, numbers, symbols
Example: xK9#mP2$vL5@nQ8*rT3
```

**Pros:** Maximum entropy
**Cons:** Hard to remember, requires password manager

### Method 3: Modified Passphrase

Combine passphrase with symbols/numbers:

```
Correct-Horse-47!
Battery&Staple-2024
Dancing-Elephants@Night
```

## Changing Your Password

**Current limitation:** Envy does not support changing the master password without creating a new vault.

**Workaround:**
```bash
# 1. Export all projects
mkdir temp_export
cd temp_export

# Export each project (repeat for each)
envy --export project1
mv .env project1.env

# 2. Remove old vault
rm ~/.envy/keys.json

# 3. Create new vault with new password
envy --import project1.env
# Enter new password when prompted

# 4. Import all projects
# Repeat for each .env file

# 5. Clean up
cd ..
rm -rf temp_export
```

**Future feature:** In-place password change.

## Password Recovery

### The Reality

**If you forget your master password, your secrets are lost forever.**

**Why:**
- We don't store your password
- We can't derive the key without it
- Encryption is designed to be unbreakable
- No "reset password" option (by design)

### Prevention Strategies

1. **Write it down**
   - Physical paper in secure location
   - Safe deposit box
   - Trusted family member (partial)

2. **Use a memorable passphrase**
   - Something you'll never forget
   - But others can't guess

3. **Password manager**
   - Store in your primary password manager
   - Only if you trust the password manager

4. **Hint system**
   - Create a personal hint
   - Don't make it obvious to others

5. **Regular use**
   - Use Envy regularly to reinforce memory
   - Muscle memory helps

## Password Security

### Where NOT to Store

**Never store in:**
- Plain text files on your computer
- Notes apps without encryption

**Safe storage options:**
- Physical paper (secure location)
- Hardware security key (YubiKey)
- Dedicated password manager (1Password, Bitwarden)
- Encrypted file (separate from vault)

### Who to Tell

**Principle of least sharing:**

- **Don't share** with coworkers (use separate vaults)
- **Don't share** with friends
- **Consider sharing** with spouse/partner (emergency access)

### Emergency Access

Plan for emergencies:

1. **Create encrypted backup**
   ```bash
   gpg -c ~/.envy/keys.json
   # Store keys.json.gpg securely
   ```

2. **Share decryption method**
   - GPG passphrase in safe deposit box
   - Or share with trusted person

3. **Document recovery process**
   - Write down steps
   - Store with backup

## Password Hygiene

### Unique to Envy

Use a password **only** for Envy:

Bad:
```
Envy password: MyPassword123!
Also used for: Email, Bank, Facebook
```

Good:
```
Envy password: Unique-Envy-Only-Passphrase-47!
Not used anywhere else
```

**Why:** If another service is breached, your vault remains secure.

### Regular Use

- Enter password regularly to maintain memory
- Don't rely solely on muscle memory
- Practice from memory occasionally

### Rotation

**Not required for master passwords** (unlike service passwords):

- Only change if compromised
- Or if you've used it elsewhere
- Or if it's weak and you want to strengthen

**Remember:** Changing master password requires vault recreation (currently).

## Forgotten Password Recovery

### Before You Panic

1. **Try variations**
   - Different capitalization
   - With/without numbers
   - Old passwords you might have used

2. **Check password manager**
   - Maybe you stored it there

3. **Check written records**
   - Physical notes
   - Safe deposit box

4. **Think about hints**
   - What was your thought process?

### If Truly Forgotten

Unfortunately, there's no recovery option. You must:

1. Accept data loss
2. Remove old vault: `rm ~/.envy/keys.json`
3. Create new vault: `envy`
4. Re-add all secrets

**Prevention for next time:**
- Use a passphrase you'll remember
- Write it down and secure it
- Store in password manager

## Advanced: Password Analysis

### Entropy Calculation

**Rough estimates:**

| Type | Example | Entropy |
|------|---------|---------|
| 8 chars random | `xK9#mP2$` | ~52 bits |
| 12 chars random | `xK9#mP2$vL5@` | ~78 bits |
| 4-word passphrase | `correct-horse-battery-staple` | ~88 bits |
| 5-word passphrase | `purple-monkey-dishwasher-spinning-wheel` | ~110 bits |

**Recommendation:** 80+ bits of entropy (4+ word passphrase or 12+ random chars)

### Brute Force Estimates

Time to crack with specialized hardware:

| Password | Time to Crack |
|----------|---------------|
| `password123` | < 1 second |
| `Tr0ub4dor&3` | ~3 days |
| `correct-horse-battery-staple` | ~500 years |
| 16 random chars | ~1 trillion years |

**Note:** These assume offline attack on stolen vault with specialized hardware.

---

**Next:** Learn [Security Best Practices](./best-practices.md) for using Envy securely.
