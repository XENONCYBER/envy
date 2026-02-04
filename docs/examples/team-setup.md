# Team Setup

Best practices for using Envy in team environments.

## Overview

Teams have unique requirements:
- Multiple developers need access
- Different access levels per environment
- Onboarding/offboarding procedures
- Secret rotation policies

## Team Structure Options

### Option 1: Shared Vault

**Best for:** Small, trusted teams (2-5 people)

All team members share one vault with the same master password.

**Setup:**
```bash
# 1. One person creates vault
envy
# Set strong team master password

# 2. Add all projects
envy set shared-app DATABASE_URL=...
envy set shared-app API_KEY=...

# 3. Share vault file securely
# - Encrypted USB
# - Secure file share
# - Password manager attachment
# - GPG encrypted email
```

**Distribution:**
```bash
# Team member setup:
curl -fsSL ... | sh  # Install Envy
mkdir -p ~/.envy
chmod 700 ~/.envy

# Copy vault (secure method)
cp /secure/path/keys.json ~/.envy/
chmod 600 ~/.envy/keys.json

# Launch and enter team password
envy
```

**Pros:**
- Simple setup
- Everyone has access to everything
- Easy secret updates

**Cons:**
- No access control
- Hard to revoke access
- Everyone knows master password
- No audit trail

### Option 2: Environment-Based Vaults

**Best for:** Teams with clear environment separation

Separate vaults per environment:
- `~/.envy/dev/keys.json` — Development (all devs)
- `~/.envy/stage/keys.json` — Staging (senior devs)
- `~/.envy/prod/keys.json` — Production (ops only)

**Setup:**

Create configs:
```lua
-- ~/.config/envy/dev.lua
return {
  backend = {
    keys_path = "~/.envy/dev/keys.json",
    lock_path = "~/.envy/dev/.lock"
  }
}

-- ~/.config/envy/stage.lua
return {
  backend = {
    keys_path = "~/.envy/stage/keys.json",
    lock_path = "~/.envy/stage/.lock"
  }
}

-- ~/.config/envy/prod.lua
return {
  backend = {
    keys_path = "~/.envy/prod/keys.json",
    lock_path = "~/.envy/prod/.lock"
  }
}
```

Wrapper script:
```bash
#!/bin/bash
# /usr/local/bin/envy-team

case "$1" in
  dev)
    cp ~/.config/envy/dev.lua ~/.config/envy/config.lua
    shift
    envy "$@"
    ;;
  stage)
    cp ~/.config/envy/stage.lua ~/.config/envy/config.lua
    shift
    envy "$@"
    ;;
  prod)
    cp ~/.config/envy/prod.lua ~/.config/envy/config.lua
    shift
    envy "$@"
    ;;
  *)
    echo "Usage: envy-team {dev|stage|prod} [envy-args]"
    exit 1
    ;;
esac
```

**Usage:**
```bash
envy-team dev          # Use dev vault
envy-team stage        # Use stage vault
envy-team prod         # Use prod vault
```

**Pros:**
- Clear access boundaries
- Different passwords per environment
- Easier to revoke access to prod

**Cons:**
- Multiple vaults to manage
- More complex setup
- Context switching overhead

### Option 3: Role-Based Vaults

**Best for:** Larger teams with role separation

Vaults per role:
- `~/.envy/developers/keys.json` — Dev secrets
- `~/.envy/operations/keys.json` — Ops secrets
- `~/.envy/admin/keys.json` — Admin secrets

Similar setup to Option 2, but organized by role instead of environment.

### Option 4: Individual Vaults with Sharing

**Best for:** Distributed teams, contractors

Each person has their own vault. Secrets shared via secure export.

**New developer onboarding:**
```bash
# Senior dev exports needed secrets
envy --export myapp-dev
# Securely sends .env file to new dev

# New dev imports
envy --import myapp-dev.env
```

**Pros:**
- Maximum isolation
- Easy to revoke (don't send new secrets)
- No shared master password

**Cons:**
- Manual secret distribution
- Secrets can diverge
- Hard to update shared secrets

## Recommended: Hybrid Approach

Combine options based on sensitivity:

```
┌─────────────────────────────────────────┐
│  Development Vault (Shared)             │
│  - All developers have access           │
│  - Shared master password               │
│  - Contains dev/stage secrets           │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│  Production Vault (Restricted)          │
│  - Only senior devs + ops               │
│  - Different master password            │
│  - Stored separately                    │
└─────────────────────────────────────────┘
```

## Onboarding Process

### New Developer Checklist

**Before first day:**
- [ ] Create development vault access
- [ ] Prepare onboarding secrets package
- [ ] Send master password via secure channel

**First day:**
```bash
# 1. Install Envy
curl -fsSL https://raw.githubusercontent.com/XENONCYBER/envy/main/install.sh | sh

# 2. Create directories
mkdir -p ~/.envy
chmod 700 ~/.envy

# 3. IT provides vault file
# Copy to ~/.envy/keys.json
chmod 600 ~/.envy/keys.json

# 4. Test access
envy
# Enter master password

# 5. Verify can run app
envy run myapp-dev -- npm run dev
```

**Documentation to provide:**
- Master password (securely)
- Vault backup procedure
- Secret naming conventions
- Environment guidelines
- Emergency contacts

## Offboarding Process

### Developer Leaving

**Immediate:**
- Remove from vault distribution
- Rotate all secrets they had access to
- Revoke any personal API keys

**Steps:**
```bash
# 1. Generate new master password
# 2. Create new vault with rotated secrets
# 3. Distribute to remaining team

# Rotate all secrets they knew:
# - Database passwords
# - API keys
# - Service credentials
# - SSH keys
```

**Communication:**
- Notify team of new vault/password
- Document rotated secrets
- Update access documentation

### Environment Isolation

**Rule:** Never use production secrets in development.

```bash
# Create separate projects per environment
envy set myapp DATABASE_URL=postgres://dev -e dev
envy set myapp DATABASE_URL=postgres://prod -e prod

# Or separate vaults (recommended for teams)
# dev-vault: myapp (dev)
# prod-vault: myapp (prod)
```

### Rotation Schedule

| Secret Type | Rotation Frequency | Owner |
|-------------|-------------------|-------|
| Database passwords | Quarterly | DBA/Operations |
| API keys (external) | Monthly | Security team |
| JWT secrets | Semi-annually | Backend team |
| Service accounts | Quarterly | DevOps |
| TLS certificates | Before expiry | Security team |

**Rotation procedure:**
```bash
# 1. Generate new secret
NEW_KEY=$(openssl rand -hex 32)

# 2. Update in Envy
envy set myapp API_SECRET_KEY="$NEW_KEY" -e prod

# 3. Update in service (no restart needed if hot-reload)
envy run myapp-prod -- ./reload-config.sh

# 4. Verify service works

# 5. Old value in history if rollback needed
```

## Access Control Documentation

### Vault Access Matrix

| Role | Dev Vault | Stage Vault | Prod Vault |
|------|-----------|-------------|------------|
| Junior Dev | Read/Write | — | — |
| Senior Dev | Read/Write | Read/Write | Read |
| DevOps | Read/Write | Read/Write | Read/Write |
| Contractor | Read (specific projects) | — | — |

## Team Automation

### Git Hook for Secret Prevention

Prevent accidental .env commits:

```bash
#!/bin/bash
# .git/hooks/pre-commit

if git diff --cached --name-only | grep -E '\.env($|\.)'; then
  echo "ERROR: Attempting to commit .env file!"
  echo "Use Envy instead: envy --export <project> if needed"
  exit 1
fi
```

### Makefile Integration

```makefile
# Makefile

.PHONY: dev prod secrets-check

dev:
	envy run myapp-dev -- npm run dev

prod-deploy:
	@echo "Deploying to production..."
	envy run myapp-prod -- ./deploy.sh

secrets-check:
	@echo "Checking for exposed secrets..."
	@grep -r "sk_live_" . --include="*.js" --include="*.ts" 2>/dev/null && echo "WARNING: Live keys found!" || echo "OK"
	@grep -r "postgres://" . --include="*.js" --include="*.ts" 2>/dev/null && echo "WARNING: DB URLs found!" || echo "OK"
```

### Docker Integration

```dockerfile
# Dockerfile

# Don't copy .env files!
# Use envy run in container

FROM node:18
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .

# Expose port
EXPOSE 3000

# Run with envy (requires vault in container - mount or copy)
CMD ["envy", "run", "myapp", "--", "npm", "start"]
```

**Next:** Learn about [CI/CD Integration](./ci-cd-integration.md).
