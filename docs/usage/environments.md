# Environments

Understanding and managing deployment environments in Envy.

## What are Environments?

Environments in Envy represent deployment stages. Each project can exist in three environments:

- **`dev`** — Local development
- **`stage`** — Staging/testing/pre-production
- **`prod`** — Production/live

This allows you to maintain separate credentials for different stages of your deployment pipeline.

## Environment Isolation

Each environment is completely isolated:

```
Project: myapp
├── Environment: dev
│   ├── DATABASE_URL=postgres://localhost/devdb
│   └── API_KEY=dev-key-123
├── Environment: stage
│   ├── DATABASE_URL=postgres://stage-db/myapp
│   └── API_KEY=stage-key-456
└── Environment: prod
    ├── DATABASE_URL=postgres://prod-db/myapp
    └── API_KEY=prod-key-789
```

**Key points:**
- Same project name, different environments = different projects
- Secrets in `dev` don't affect `stage` or `prod`
- You can have the same key names with different values

## Visual Indicators

In the TUI, environments are color-coded:

| Environment | Badge | Color | Usage |
|-------------|-------|-------|-------|
| `dev` | [DEV] | Green | Local development, testing |
| `stage` | [STAGE] | Yellow | Staging, QA, pre-production |
| `prod` | [PROD] | Red | Production, live systems |

These colors help you quickly identify the context of secrets.

## Specifying Environments

### CLI Commands

Use the `-e` or `--env` flag:

```bash
# Set secret in specific environment
envy set myapp DATABASE_URL=postgres://localhost -e dev
envy set myapp DATABASE_URL=postgres://stage-server -e stage
envy set myapp DATABASE_URL=postgres://prod-cluster -e prod

# If not specified, defaults to 'dev'
envy set myapp DEBUG=true
# Same as: envy set myapp DEBUG=true -e dev
```

### TUI

When creating a project:

1. Press `ctrl+n` to create
2. Enter project name
3. Navigate to environment selector
4. Use `←/→` or `Space` to select
5. Press `Enter` to confirm

## Environment Workflows

### Development → Staging → Production

Typical promotion flow:

```bash
# 1. Development
envy set myapp API_KEY=dev-testing-key -e dev
envy run myapp -- npm run dev

# 2. Promote to staging
# Option A: Re-set manually
envy set myapp API_KEY=staging-key-123 -e stage

# Option B: Export from dev and import to stage
envy --export myapp  # Exports dev (default)
# Then manually update values and import to stage

# 3. Promote to production
envy set myapp API_KEY=live-production-key -e prod
```

### Testing Changes

```bash
# Test in dev first
envy set myapp EXPERIMENTAL_FEATURE=true -e dev
envy run myapp -- ./test-new-feature.sh

# Then promote to stage
envy set myapp EXPERIMENTAL_FEATURE=true -e stage
envy run myapp -- ./staging-tests.sh

# Finally to prod
envy set myapp EXPERIMENTAL_FEATURE=true -e prod
```

### Rollback Scenario

If a production secret needs to be reverted:

```bash
# View history in TUI to find previous value
# Or if you know the old value:

envy set myapp API_KEY=previous-working-key -e prod
# Old (broken) value automatically saved to history
```

## Environment-Specific Secrets

Some secrets only exist in certain environments:

```bash
# Debug mode only in dev
envy set myapp DEBUG=true -e dev

# Real payment processor only in prod
envy set myapp STRIPE_KEY=sk_live_xxx -e prod
envy set myapp STRIPE_KEY=sk_test_xxx -e dev
envy set myapp STRIPE_KEY=sk_test_xxx -e stage

# Local database only in dev
envy set myapp DATABASE_URL=postgres://localhost/myapp -e dev
```

## Common Environment Patterns

### Pattern 1: Full Separation

Every secret is environment-specific:

```bash
# Dev
envy set myapp API_URL=http://localhost:3000 -e dev
envy set myapp DATABASE_URL=postgres://localhost/dev -e dev

# Stage
envy set myapp API_URL=https://stage-api.example.com -e stage
envy set myapp DATABASE_URL=postgres://stage-db/stage -e stage

# Prod
envy set myapp API_URL=https://api.example.com -e prod
envy set myapp DATABASE_URL=postgres://prod-cluster/prod -e prod
```

### Pattern 2: Partial Sharing

Some secrets shared, others environment-specific:

```bash
# Shared across environments (set in all)
envy set myapp APP_NAME=MyApplication -e dev
envy set myapp APP_NAME=MyApplication -e stage
envy set myapp APP_NAME=MyApplication -e prod

# Environment-specific
envy set myapp LOG_LEVEL=debug -e dev
envy set myapp LOG_LEVEL=info -e stage
envy set myapp LOG_LEVEL=warn -e prod
```

### Pattern 3: Environment Variables Only

Use `envy run` to inject environment-specific values:

```bash
# In code, use generic key names
db_url = os.getenv('DATABASE_URL')

# But set different values per environment
envy run myapp-dev -- python app.py    # Uses dev database
envy run myapp-stage -- python app.py  # Uses stage database
envy run myapp-prod -- python app.py   # Uses prod database
```

## Environment Best Practices

### Security

1. **Never share prod secrets**
   - Prod environment should have restricted access
   - Consider separate vault for production

2. **Use different credentials per environment**
   - Different database users
   - Different API keys
   - Isolated infrastructure

3. **Rotate production secrets more frequently**
   - Dev: Rotate occasionally
   - Stage: Rotate before major releases
   - Prod: Rotate regularly (monthly/quarterly)

### Organization

1. **Consistent naming**
   ```bash
   # Good
   myapp (dev)
   myapp (stage)
   myapp (prod)

   # Avoid mixing conventions
   myapp-dev (dev)  # Redundant
   ```

2. **Environment-specific prefixes** (optional)
   ```bash
   # Some teams prefer this
   envy set myapp-dev DATABASE_URL=... -e dev
   envy set myapp-prod DATABASE_URL=... -e prod
   ```

3. **Document environment differences**
   Keep notes on what differs between environments.

### Automation

1. **CI/CD Integration**
   ```bash
   # Deploy script
   case $ENVIRONMENT in
     dev)
       envy run myapp-dev -- ./deploy.sh
       ;;
     stage)
       envy run myapp-stage -- ./deploy.sh
       ;;
     prod)
       envy run myapp-prod -- ./deploy.sh
       ;;
   esac
   ```

2. **Environment-specific configs**
   ```bash
   # .env files per environment
   envy --export myapp
   mv .env .env.production

   # In deployment:
   cp .env.$ENVIRONMENT .env
   ```

## Switching Environments

### Quick Switch in TUI

1. Exit current project view (`Esc`)
2. Search for same project name (`i`)
3. Navigate to desired environment
4. Press `Enter`

### CLI: Always Specify

```bash
# Be explicit about environment
envy run myapp -- ./script.sh          # Uses dev (default)
envy run myapp -- ./script.sh -e prod  # Uses prod

# Or set ENVY_ENVIRONMENT (if supported in future)
```

## Troubleshooting

### "Project not found"

```bash
# Check which environments exist
envy
# Then search for your project - you may have created it in a different environment
```

### Wrong environment selected

```bash
# Verify current environment
envy
# Look at the badge color: [DEV] vs [STAGE] vs [PROD]
```

### Accidentally created in wrong environment

1. Export the project: `envy --export myapp`
2. Delete from wrong environment (in TUI)
3. Import with correct environment

---

**Next:** Learn how to customize Envy with [Lua Configuration](../configuration/lua-config.md).
