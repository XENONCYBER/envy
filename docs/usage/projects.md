# Projects

Understanding project-based organization in Envy.

## What is a Project?

A **project** in Envy is a logical grouping of related secrets. Think of it as a container for all the API keys, database credentials, and configuration values that belong to a single application or service.

**Example projects:**
- `webapp` — Web application secrets (DATABASE_URL, API_KEY, SECRET_TOKEN)
- `api-service` — Backend API secrets
- `mobile-app` — Mobile application config
- `infrastructure` — Infrastructure credentials (AWS keys, SSH keys)

## Project Structure

Each project contains:

- **Name** — Identifier (e.g., "myapp")
- **Environment** — dev, stage, or prod
- **Keys** — List of secret key-value pairs
- **History** — Previous values for each key

## Creating Projects

### Via TUI

1. Launch Envy: `envy`
2. Press `shift+n` for "New Project"
3. Enter project name
4. Select environment (dev/stage/prod)
5. Add at least one key-value pair
6. Press `Save` (`shift+s`)

### Via CLI

```bash
# Simple creation
envy set myproject API_KEY=mysecret

# This will:
# 1. Create project 'myproject' (dev) if it doesn't exist
# 2. Add or update the API_KEY
# 3. Save to vault
```

## Naming Conventions

### Valid Names

- Alphanumeric characters: `myapp`, `api-v2`, `service_1`
- Spaces allowed (quote in CLI): `"My App"`, `'another app'`
- Max length: 256 characters

### Invalid Names

- Cannot be empty
- Max 256 characters

### Best Practices

```bash
# Good names
webapp              # Simple, clear
api-service         # Descriptive
mobile-app-v2       # Version included
mycompany-web       # Namespace prefix

# Avoid
app                 # Too generic
project1            # Not descriptive
temp                # Temporary becomes permanent
```

## Project Organization Strategies

### By Application

Organize by the application or service:

```
├── webapp (dev/stage/prod)
├── api-service (dev/stage/prod)
├── worker-service (dev/stage/prod)
└── mobile-app (dev/stage/prod)
```

**Pros:**
- Logical grouping
- Easy to find related secrets
- Matches code repositories

**Cons:**
- Can lead to duplication across environments

### By Team/Department

Organize by ownership:

```
├── backend-team (dev/stage/prod)
├── frontend-team (dev/stage/prod)
├── devops-team (dev/stage/prod)
└── data-team (dev/stage/prod)
```

**Pros:**
- Clear ownership
- Easy team-based access control (in future)

**Cons:**
- Less intuitive for cross-team projects

### By Environment Only

Flat structure by environment:

```
├── dev (all dev secrets)
├── stage (all staging secrets)
└── prod (all production secrets)
```

**Pros:**
- Simple structure
- Easy to export entire environment

**Cons:**
- Harder to find specific application secrets
- Can become unwieldy with many secrets

### Recommended: Hybrid Approach

Combine application and environment:

```
├── webapp-dev
├── webapp-stage
├── webapp-prod
├── api-dev
├── api-stage
└── api-prod
```

Or use Envy's built-in environment system:

```
├── webapp (with dev/stage/prod environments)
└── api (with dev/stage/prod environments)
```

## Managing Project Secrets

### Adding Secrets

**TUI:**
1. Select project with `Enter`
2. Press `shift+e` for "Edit Project"
3. Navigate to "New Key Name/Value" fields
4. Enter details and press `Add`
5. Press `Save`

**CLI:**
```bash
# Add multiple secrets to same project
envy set myapp DATABASE_URL=postgres://localhost/myapp
envy set myapp API_KEY=sk-abc123
envy set myapp SECRET_TOKEN=xyz789
envy set myapp REDIS_URL=redis://localhost:6379
```

### Updating Secrets

When you update a secret, the old value is automatically saved to history.

```bash
# First set
envy set myapp API_KEY=old-key-123

# Later update - old value moves to history
envy set myapp API_KEY=new-key-456
```

### Viewing History

**TUI:**
1. Select project with `Enter`
2. Navigate to key with `↑/↓`
3. Press `shift+h` for history

Shows:
- Current value (green badge)
- All previous values with timestamps
- Up to 5 most recent previous values

### Deleting Secrets

**TUI:**
1. Select project with `Enter`
2. Navigate to key with `↑/↓`
3. Press `shift+d` to delete (with confirmation)

**Note:** Deleted secrets are **not** recoverable. History is preserved for updated keys, but deletion is permanent.

## Project Workflows

### Development Workflow

```bash
# 1. Create dev project
envy set myapp-dev DATABASE_URL=postgres://localhost/dev

# 2. Run development server
envy run myapp-dev -- npm run dev

# 3. Add more secrets as needed
envy set myapp-dev API_KEY=dev-key-123
envy set myapp-dev DEBUG=true
```

### Promotion Workflow

Moving secrets from dev → stage → prod:

```bash
# Method 1: Export and re-import
envy --export myapp-dev
envy --import .env
# When prompted: name=myapp, environment=stage

# Method 2: Set individually
envy set myapp DATABASE_URL=postgres://stage-db -e stage
envy set myapp API_KEY=stage-key-123 -e stage

# Method 3: Copy from another environment (manual)
# 1. View in TUI
# 2. Copy values
# 3. Create in target environment
```

### Team Onboarding

New team member setup:

```bash
# 1. Share vault file securely (encrypted)
# Team member copies to ~/.envy/keys.json

# 2. Set permissions
chmod 600 ~/.envy/keys.json

# 3. Verify access
envy
```

## Searching Projects

### TUI Search

Press `i` in grid view to search:

- **All mode** (default): Searches both project names AND key names
- **Projects mode**: Searches project names only
- **Keys mode**: Searches key names only

Press `Tab` to cycle modes.

**Search examples:**
- Type `api` → Matches projects/keys containing "api"
- Type `prod` → Finds production projects
- Type `database` → Finds database-related keys

### Filtering Tips

- Search is case-insensitive
- Partial matching works (type `db` to find `DATABASE_URL`)
- Real-time filtering as you type

## Project Metadata

Each project tracks:

- **Creation** — Implicit (first secret added)
- **Last modified** — Not explicitly tracked (use key timestamps)
- **Environment** — dev/stage/prod
- **Key count** — Visible in TUI grid

## Best Practices

### Naming

- Use lowercase with hyphens: `my-service`, `api-gateway`
- Include version if relevant: `api-v1`, `api-v2`
- Be descriptive but concise
- Use consistent naming across team

### Organization

- One project per application/service
- Use environments for deployment stages
- Group related keys together
- Don't overload a project with unrelated secrets

### Security

- Separate projects for different trust levels
- Use different master passwords for different vaults (if needed)
- Export only when necessary, delete .env files immediately
- Review and rotate secrets regularly

### Maintenance

- Delete unused projects to reduce clutter
- Update outdated secrets
- Review history periodically
- Backup vault regularly

## Common Patterns

### Microservices Setup

```
├── gateway (dev/stage/prod)
├── auth-service (dev/stage/prod)
├── user-service (dev/stage/prod)
├── payment-service (dev/stage/prod)
└── notification-service (dev/stage/prod)
```

### Monolith Setup

```
├── main-app (dev/stage/prod)
├── background-jobs (dev/stage/prod)
└── third-party-apis (dev/stage/prod)
```

### Infrastructure Setup

```
├── aws-credentials (prod)
├── ssh-keys (prod)
├── database-cluster (prod)
└── monitoring (prod)
```

---

**Next:** Learn about [Environments](./environments.md) for managing deployment stages.
