# CI/CD Integration

Using Envy in continuous integration and deployment pipelines.

## Overview

Envy can securely provide secrets to CI/CD pipelines without exposing them in:
- Repository code
- Build logs
- Environment variable UIs
- Docker images

## Security Principles

1. **Never commit secrets** — Not even encrypted
2. **No secrets in build logs** — Mask or prevent
3. **Rotate CI secrets** — Different from production
4. **Least privilege** — CI only gets needed secrets
5. **Audit access** — Track what secrets CI uses

## Approach 1: Direct Injection (Recommended)

Use `envy run` to inject secrets directly into build commands.

### GitHub Actions

```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install Envy
        run: |
          curl -fsSL https://raw.githubusercontent.com/XENONCYBER/envy/main/install.sh | sh
      
      - name: Setup Envy Vault
        env:
          ENVY_MASTER_PASSWORD: ${{ secrets.ENVY_MASTER_PASSWORD }}
        run: |
          mkdir -p ~/.envy
          echo "${{ secrets.ENVY_VAULT }}" | base64 -d > ~/.envy/keys.json
          chmod 600 ~/.envy/keys.json
      
      - name: Deploy with Secrets
        env:
          ENVY_MASTER_PASSWORD: ${{ secrets.ENVY_MASTER_PASSWORD }}
        run: |
          envy run myapp-prod -- ./deploy.sh
```

**Secrets to configure in GitHub:**
- `ENVY_MASTER_PASSWORD` — Your vault master password
- `ENVY_VAULT` — Base64-encoded vault file

### GitLab CI

```yaml
# .gitlab-ci.yml
stages:
  - deploy

deploy:
  stage: deploy
  image: ubuntu:latest
  before_script:
    - curl -fsSL https://raw.githubusercontent.com/XENONCYBER/envy/main/install.sh | sh
    - mkdir -p ~/.envy
    - echo "$ENVY_VAULT" | base64 -d > ~/.envy/keys.json
    - chmod 600 ~/.envy/keys.json
  script:
    - envy run myapp-prod -- ./deploy.sh
  only:
    - main
```

**Variables to configure in GitLab:**
- `ENVY_MASTER_PASSWORD` — Protected, masked
- `ENVY_VAULT` — Base64-encoded vault, protected

## Approach 2: Export to Temporary File

When you need a .env file for tools that require it:

```yaml
# GitHub Actions example with .env
- name: Create temporary .env
  env:
    ENVY_MASTER_PASSWORD: ${{ secrets.ENVY_MASTER_PASSWORD }}
  run: |
    mkdir -p ~/.envy
    echo "${{ secrets.ENVY_VAULT }}" | base64 -d > ~/.envy/keys.json
    chmod 600 ~/.envy/keys.json
    
    # Export to temp .env
    cd /tmp
    envy --export myapp-prod
    chmod 600 .env
    
    # Run command that needs .env
    docker run --env-file /tmp/.env myimage
    
    # Immediately delete
    rm -f /tmp/.env
```

**Note:** This is less secure than direct injection. Use only when necessary.

## Approach 3: CI-Specific Vault

Create a separate vault with only the secrets CI needs:

### Setup

```bash
# 1. Create CI vault locally
mkdir -p ~/.envy-ci
export ENVY_DATA_DIR=~/.envy-ci

# 2. Create new vault
envy
# Set different master password

# 3. Add limited secrets
envy set webapp DATABASE_URL=postgres://ci-readonly@prod-db
envy set webapp API_KEY=ci-limited-key-123
# Do NOT add admin credentials

# 4. Export CI vault for CI system
base64 ~/.envy-ci/keys.json
# Store this as ENVY_VAULT_CI in CI system
```

### CI Configuration

```yaml
# Use CI-specific vault
- name: Deploy
  env:
    ENVY_MASTER_PASSWORD: ${{ secrets.ENVY_MASTER_PASSWORD_CI }}
  run: |
    mkdir -p ~/.envy
    echo "${{ secrets.ENVY_VAULT_CI }}" | base64 -d > ~/.envy/keys.json
    chmod 600 ~/.envy/keys.json
    envy run webapp-prod -- ./deploy.sh
```

**Benefits:**
- CI only gets minimal required access
- Can rotate CI secrets independently
- Breach of CI vault ≠ breach of main vault
- Easy to revoke CI access

## Approach 4: Environment-Specific Projects

Separate projects for CI vs human use:

```bash
# Human access
envy set webapp DATABASE_URL=postgres://admin@prod -e prod

# CI access (different credentials)
envy set webapp-ci DATABASE_URL=postgres://ci@prod -e prod
envy set webapp-ci DEPLOY_TOKEN=ci-token-123 -e prod
```

**CI uses:**
```bash
envy run webapp-ci -- ./deploy.sh
```

## Docker Integration

### Multi-stage Build with Envy

```dockerfile
# Dockerfile
FROM node:18 AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .

# Don't build here - no secrets yet

# Production image
FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
COPY package.json .

# Install Envy in container
RUN apk add --no-cache curl && \
    curl -fsSL https://raw.githubusercontent.com/XENONCYBER/envy/main/install.sh | sh && \
    apk del curl

# Entrypoint uses envy
ENTRYPOINT ["envy", "run", "myapp-prod", "--"]
CMD ["node", "dist/server.js"]
```

**Usage:**
```bash
# Mount vault into container
docker run \
  -v ~/.envy:/root/.envy:ro \
  -e ENVY_MASTER_PASSWORD \
  myapp:latest
```

### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'

services:
  app:
    build: .
    volumes:
      - ~/.envy:/root/.envy:ro
    environment:
      - ENVY_MASTER_PASSWORD
    command: ["envy", "run", "myapp-prod", "--", "node", "server.js"]
```

## Kubernetes Integration

### Secret Management

Option 1: Init Container with Envy

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: myapp
spec:
  initContainers:
  - name: envy-setup
    image: alpine
    command:
    - sh
    - -c
    - |
      apk add curl
      curl -fsSL https://raw.githubusercontent.com/XENONCYBER/envy/main/install.sh | sh
      mkdir -p /shared/envy
      echo "$VAULT_DATA" | base64 -d > /shared/envy/keys.json
    env:
    - name: VAULT_DATA
      valueFrom:
        secretKeyRef:
          name: envy-vault
          key: vault
    volumeMounts:
    - name: shared
      mountPath: /shared
  
  containers:
  - name: app
    image: myapp:latest
    command:
    - envy
    - run
    - myapp-prod
    - --
    - node
    - server.js
    env:
    - name: ENVY_MASTER_PASSWORD
      valueFrom:
        secretKeyRef:
          name: envy-vault
          key: password
    volumeMounts:
    - name: shared
      mountPath: /root/.envy
  
  volumes:
  - name: shared
    emptyDir: {}
```

### Kubernetes Secret

```bash
# Create secret from vault
kubectl create secret generic envy-vault \
  --from-literal=password="$(cat ~/envy-password.txt)" \
  --from-literal=vault="$(base64 ~/.envy/keys.json)"
```

## Best Practices

### 1. Mask Secrets in Logs

```yaml
# GitHub Actions
- name: Deploy
  env:
    ENVY_MASTER_PASSWORD: ${{ secrets.ENVY_MASTER_PASSWORD }}
  run: |
    # Envy will prompt for password, mask it
    echo "::add-mask::$ENVY_MASTER_PASSWORD"
    envy run myapp-prod -- ./deploy.sh
```

### 2. Use Branch Protection

Protect branches that deploy to production:
- Require PR reviews
- Require status checks
- Restrict who can push

### 3. Separate CI Secrets

```bash
# Rotate CI secrets monthly
ci-rotate-secrets.sh
```

### 4. Audit CI Access

```yaml
# Log all CI deployments
- name: Deploy
  run: |
    echo "$(date): Deployment by $GITHUB_ACTOR" >> deployments.log
    envy run myapp-prod -- ./deploy.sh
```

### 5. Short-lived CI Vaults

Generate temporary vaults for each deployment:

```bash
#!/bin/bash
# create-ci-vault.sh

# Create temp vault with just this deployment's secrets
mkdir -p /tmp/envy-ci
export HOME=/tmp/envy-ci

# Generate random password
PASSWORD=$(openssl rand -base64 32)

# Create vault non-interactively (if supported)
# Or use pre-created minimal vault

# Export for CI
echo "ENVY_VAULT=$(base64 /tmp/envy-ci/.envy/keys.json)"
echo "ENVY_PASSWORD=$PASSWORD"

# Cleanup after deployment
rm -rf /tmp/envy-ci
```

## CI/CD Security Checklist

- CI vault separate from production vault
- Master password stored as encrypted secret
- Vault file base64-encoded in CI secret
- No secrets in build logs
- .env files immediately deleted after use
- Branch protection on deploy branches
- Rotation schedule for CI credentials
- Audit logging of deployments
- Emergency access procedure documented
- Test disaster recovery monthly

## Example: Complete GitHub Actions Workflow

```yaml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Install Envy
        run: curl -fsSL https://raw.githubusercontent.com/XENONCYBER/envy/main/install.sh | sh
      
      - name: Setup Test Secrets
        env:
          ENVY_MASTER_PASSWORD: ${{ secrets.ENVY_CI_PASSWORD }}
        run: |
          mkdir -p ~/.envy
          echo "${{ secrets.ENVY_CI_VAULT }}" | base64 -d > ~/.envy/keys.json
          chmod 600 ~/.envy/keys.json
      
      - name: Run Tests
        env:
          ENVY_MASTER_PASSWORD: ${{ secrets.ENVY_CI_PASSWORD }}
        run: envy run myapp-test -- npm test

  deploy-staging:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/develop'
    environment: staging
    steps:
      - uses: actions/checkout@v3
      
      - name: Install Envy
        run: curl -fsSL https://raw.githubusercontent.com/XENONCYBER/envy/main/install.sh | sh
      
      - name: Deploy to Staging
        env:
          ENVY_MASTER_PASSWORD: ${{ secrets.ENVY_STAGING_PASSWORD }}
        run: |
          mkdir -p ~/.envy
          echo "${{ secrets.ENVY_STAGING_VAULT }}" | base64 -d > ~/.envy/keys.json
          chmod 600 ~/.envy/keys.json
          envy run myapp-stage -- ./deploy-staging.sh

  deploy-production:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment: production
    steps:
      - uses: actions/checkout@v3
      
      - name: Install Envy
        run: curl -fsSL https://raw.githubusercontent.com/XENONCYBER/envy/main/install.sh | sh
      
      - name: Deploy to Production
        env:
          ENVY_MASTER_PASSWORD: ${{ secrets.ENVY_PROD_PASSWORD }}
        run: |
          mkdir -p ~/.envy
          echo "${{ secrets.ENVY_PROD_VAULT }}" | base64 -d > ~/.envy/keys.json
          chmod 600 ~/.envy/keys.json
          envy run myapp-prod -- ./deploy-production.sh
```

---

**Next:** Review [Security Best Practices](../security/best-practices.md) for production deployments.
