# Basic Workflow Example

A complete walkthrough of typical Envy usage.

## Scenario

You're starting a new web application called "taskapp" with:
- Node.js backend
- PostgreSQL database
- Redis cache
- External API integration

You need to manage secrets for development and production.

## Step 1: Initial Setup

### Install Envy

```bash
curl -fsSL https://raw.githubusercontent.com/XENONCYBER/envy/main/install.sh | sh
```

### First Launch

```bash
envy
```

Create your master password:
```
Welcome to Envy - Secure Secret Manager
No vault found. Let's create one!

Create master password: YourSecurePassphrase123!
Confirm master password: YourSecurePassphrase123!

Vault created successfully!
```

The TUI opens with an empty grid.

## Step 2: Development Environment Setup

### Create Development Project

In the TUI:
1. Press `ctrl+n` to create new project
2. Enter name: `taskapp`
3. Select environment: `DEV` (use arrow keys or space)
4. Enter first key name: `DATABASE_URL`
5. Enter value: `postgresql://localhost:5432/taskapp_dev`
6. Press `+ Add` (or `ctrl+a`)
7. Add more keys:

Continue adding:
```
Key Name: REDIS_URL
Value: redis://localhost:6379/0
[+ Add]

Key Name: API_SECRET_KEY
Value: dev-secret-key-not-for-production
[+ Add]

Key Name: JWT_SECRET
Value: dev-jwt-secret-change-in-prod
[+ Add]

Key Name: PORT
Value: 3000
```

Finally, press `Save` (`ctrl+s`).

### Verify in Grid

Back in the grid view, you'll see a new card labeled "taskapp" with a green [DEV] badge. The card displays a preview of the first few keys you've added and indicates there are more. This visual confirmation shows your project is properly stored and ready to use.

### Run Development Server

```bash
cd ~/projects/taskapp

# Run with secrets injected
envy run taskapp -- npm run dev

Output:
Enter master password: ********
Loaded 5 secrets from 'taskapp' (dev)
Running: npm run dev

> taskapp@1.0.0 dev
> nodemon server.js

Connected to PostgreSQL at localhost:5432
Redis connected
Server running on port 3000
```

Your app now has all environment variables set.

## Step 3: Adding More Secrets

### Adding Database Migration Secret

```bash
# Add via CLI
envy set taskapp MIGRATE_DB=true
envy set taskapp MIGRATE_FORCE=false
```

### Adding External API Keys

```bash
# Stripe (test key)
envy set taskapp STRIPE_SECRET_KEY=sk_test_abc123
envy set taskapp STRIPE_PUBLISHABLE_KEY=pk_test_xyz789

# SendGrid (email)
envy set taskapp SENDGRID_API_KEY=SG.test.key123
```

### Verify in TUI

Launch TUI and open taskapp to see all 9 keys.

## Step 4: Production Environment

### Create Production Project

```bash
# Add production secrets via CLI

# Database (production cluster)
envy set taskapp DATABASE_URL=postgres://prod-user:prod-pass@prod-db-host:5432/taskapp -e prod

# Redis (production)
envy set taskapp REDIS_URL=redis://prod-redis:6379/0 -e prod

# Secrets (strong, unique values)
envy set taskapp API_SECRET_KEY=$(openssl rand -hex 32) -e prod
envy set taskapp JWT_SECRET=$(openssl rand -hex 32) -e prod

# Stripe (live keys)
envy set taskapp STRIPE_SECRET_KEY=sk_live_xyz789 -e prod
envy set taskapp STRIPE_PUBLISHABLE_KEY=pk_live_abc123 -e prod

# SendGrid (production)
envy set taskapp SENDGRID_API_KEY=SG.production.key456 -e prod

# Port
envy set taskapp PORT=8080 -e prod
```

### View Both Environments

In TUI:
1. Grid now shows two `taskapp` cards:
   - `taskapp` with [DEV] badge
   - `taskapp` with [PROD] badge

2. Navigate between them with arrow keys
3. Open each to verify different values

## Step 5: Daily Development Workflow

### Morning Routine

```bash
cd ~/projects/taskapp

# Start development server
envy run taskapp -- npm run dev

# Or with debug mode
envy run taskapp -- npm run dev:debug
```

### Copying a Secret

Need to paste a secret into a config file?

1. Launch TUI: `envy`
2. Select `taskapp` (dev) with `Enter`
3. Navigate to key (e.g., `SENDGRID_API_KEY`)
4. Press `y` to copy
5. Paste where needed
6. Clipboard auto-clears in 30 seconds

### Adding a New Secret

Your team adds a new service:

```bash
# Add to both environments

# Dev
envy set taskapp NEW_SERVICE_API_KEY=dev-key-123

# Prod
envy set taskapp NEW_SERVICE_API_KEY=prod-key-456 -e prod
```

### Updating a Secret

API key was rotated:

```bash
# Update - old value saved to history
envy set taskapp SENDGRID_API_KEY=SG.new.rotated.key789
```

View history in TUI:
1. Open taskapp
2. Navigate to `SENDGRID_API_KEY`
3. Press `H`
4. See current and previous values

## Step 6: Deployment Workflow

### Pre-Deployment

```bash
# Export production secrets for deployment script
envy --export taskapp -e prod

# Warning appears:
# "Exported .env file contains secrets in plain text. Keep it secure!"

# Move to deployment directory
mv .env ~/deployments/taskapp/

# Set restrictive permissions
chmod 600 ~/deployments/taskapp/.env

# Deploy
cd ~/deployments/taskapp
./deploy.sh

# Immediately delete
rm -f .env
```

### Better: Direct Injection

```bash
# On production server

# Pull latest code
git pull origin main

# Install dependencies
npm ci

# Run deployment with secrets (no .env file!)
envy run taskapp -- ./deploy.sh

# Or specifically for prod
envy run taskapp -- npm run migrate
envy run taskapp -- npm start
```

## Step 7: Team Onboarding

### New Developer Joins

**Option 1: Share Vault File**

```bash
# Secure transfer
scp ~/.envy/keys.json newdev:~/.envy/

# Or encrypted

gpg -e -r newdev@company.com ~/.envy/keys.json
# Send keys.json.gpg
```

New developer:
```bash
# 1. Install Envy
curl -fsSL ... | sh

# 2. Create directories
mkdir -p ~/.envy
chmod 700 ~/.envy

# 3. Copy vault
cp keys.json ~/.envy/
chmod 600 ~/.envy/keys.json

# 4. Launch
envy
# Enter shared master password
```

**Option 2: Export/Import**

```bash
# Export each project
envy --export taskapp
mv .env taskapp.env

# Send .env files (securely!)
# New developer imports
envy --import taskapp.env
```

## Step 8: Backup and Maintenance

### Automated Backups

```bash
#!/bin/bash
# ~/.local/bin/backup-envy.sh

DATE=$(date +%Y%m%d)
BACKUP_DIR="$HOME/Backups/envy"

mkdir -p "$BACKUP_DIR"

# Backup vault
cp ~/.envy/keys.json "$BACKUP_DIR/keys-$DATE.json"

# Keep only last 30 days
ls -t "$BACKUP_DIR"/keys-*.json | tail -n +31 | xargs rm -f

echo "Backup completed: $DATE"
```

Add to crontab:
```bash
0 2 * * * /home/user/.local/bin/backup-envy.sh
```

### Monthly Maintenance

```bash
# 1. Test backup restore
cp ~/Backups/envy/keys-$(date +%Y%m%d).json /tmp/test-keys.json
# Verify can unlock

# 2. Clean up old history (if needed)
# Export projects, re-import to clear history

# 3. Rotate critical secrets
envy
# View production keys
# Rotate API keys older than 90 days
```

## Step 9: Adding Staging Environment

```bash
# Create staging secrets
envy set taskapp DATABASE_URL=postgres://stage-db/taskapp -e stage
envy set taskapp REDIS_URL=redis://stage-redis:6379/0 -e stage
envy set taskapp API_SECRET_KEY=$(openssl rand -hex 32) -e stage
envy set taskapp JWT_SECRET=$(openssl rand -hex 32) -e stage
envy set taskapp STRIPE_SECRET_KEY=sk_test_staging -e stage
envy set taskapp SENDGRID_API_KEY=SG.staging.key -e stage
envy set taskapp PORT=8080 -e stage
```

Now you have three environments:
- `taskapp` (dev) - Green badge
- `taskapp` (stage) - Yellow badge
- `taskapp` (prod) - Red badge

## Complete Workflow Summary

```bash
# Daily development
envy run taskapp -- npm run dev

# Add secret
envy set taskapp NEW_KEY=value

# Update secret
envy set taskapp EXISTING_KEY=newvalue

# Copy secret (TUI)
envy
# Navigate → y to copy

# Deploy to production
envy run taskapp -- ./deploy.sh

# Backup
~/.local/bin/backup-envy.sh
```

## Tips

1. **Use TUI for browsing/copying** — Fastest way to find and copy secrets
2. **Use CLI for adding/updating** — Quick one-liners
3. **Always use `envy run`** — Avoid creating .env files
4. **Backup regularly** — Automate daily backups
5. **Separate environments** — Never mix dev/prod secrets
6. **Clean up exports** — Delete .env files immediately after use

---

**Next:** Learn about [Team Setup](./team-setup.md) for larger organizations.
