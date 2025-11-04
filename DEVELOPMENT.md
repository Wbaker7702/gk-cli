# Development Guide

## Quick Start

1. **Install Go** (1.21 or later)
   ```bash
   # Check version
   go version
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Build**
   ```bash
   make build
   # or
   go build -o gk ./...
   ```

4. **Run**
   ```bash
   ./gk --help
   ```

## Project Structure

```
gk-cli/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command setup
│   ├── login.go            # Authentication
│   ├── workspace.go        # Workspace management
│   ├── workspace_list.go   # Workspace listing/showing
│   ├── provider.go         # Provider integration
│   ├── pr.go               # Pull requests
│   ├── patch.go            # Cloud patches
│   ├── launchpad.go        # Launchpad dashboard
│   ├── graph.go            # Visual graph
│   └── setting.go           # Settings
├── internal/
│   ├── auth/               # Authentication logic
│   │   └── oauth.go        # OAuth flow
│   ├── config/             # Configuration
│   │   └── config.go       # Config management
│   ├── workspace/          # Workspace operations
│   │   ├── workspace.go   # Workspace CRUD
│   │   └── git.go          # Git operations
│   └── theme/              # Theme system
│       └── theme.go        # Theme loading
└── pkg/
    └── utils/              # Shared utilities
        ├── errors.go       # Error handling
        ├── input.go        # User input
        └── path.go         # Path utilities
```

## Implemented Features

### ✅ Workspace Management
- Create local/cloud workspaces
- Add repositories (by path or URL)
- List and show workspace details
- Locate repositories in a directory
- Clone all repositories in a workspace

### ✅ Git Operations
- Run git commands (fetch, pull, push, checkout) on all repos in a workspace
- Validates git repositories before operations
- Provides feedback for each repository

### ✅ Authentication
- OAuth flow skeleton (requires GitKraken OAuth credentials)
- Token storage and refresh
- Logout functionality

### ✅ Configuration
- YAML-based configuration
- Theme management
- Workspace persistence

## Usage Examples

### Workspace Management

```bash
# Create a workspace
gk ws create my-workspace --type local

# Add a repository
gk ws add-repo /path/to/repo
# or
gk ws add-repo https://github.com/user/repo.git

# List workspaces
gk ws list

# Show workspace details
gk ws show my-workspace

# Locate repositories in current directory
gk ws locate

# Clone all repositories
gk ws clone ./repos
```

### Git Operations

```bash
# Fetch all repos in workspace
gk ws fetch

# Pull all repos
gk ws pull

# Push all repos
gk ws push

# Checkout a branch in all repos
gk ws checkout main
```

### Authentication

```bash
# Login (requires OAuth credentials)
export GITKRAKEN_CLIENT_ID=your_client_id
export GITKRAKEN_CLIENT_SECRET=your_client_secret
gk login

# Logout
gk logout
```

## Testing

```bash
# Run tests
make test
# or
go test ./...

# Test specific package
go test ./internal/workspace
```

## Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build for specific platform
make build-linux
make build-darwin
make build-windows
```

## Development Workflow

1. **Create a feature branch**
   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make changes**
   - Follow Go conventions
   - Add tests for new functionality
   - Update documentation

3. **Test locally**
   ```bash
   make build
   ./gk --help
   ```

4. **Commit and push**
   ```bash
   git add .
   git commit -m "Add feature X"
   git push origin feature/my-feature
   ```

## TODO / Next Steps

- [ ] Implement API client for GitKraken services
- [ ] Add provider API integrations (GitHub, GitLab, Bitbucket)
- [ ] Implement PR listing and viewing
- [ ] Add Cloud Patch functionality
- [ ] Create Launchpad TUI
- [ ] Add comprehensive tests
- [ ] Implement Code Suggest feature
- [ ] Add Visual Graph integration

## Configuration

Configuration is stored in `~/.config/gk/config.yaml`:

```yaml
theme: default
auth:
  token: "..."
  refresh_token: "..."
  expires_at: "..."
providers: {}
workspaces: {}
```

Workspaces are stored in `~/.config/gk/workspaces/*.json`:

```json
{
  "name": "my-workspace",
  "type": "local",
  "repos": [
    {
      "name": "my-repo",
      "path": "/path/to/repo",
      "remote": "https://github.com/user/repo.git",
      "provider": "github"
    }
  ]
}
```

## Environment Variables

- `GITKRAKEN_CLIENT_ID` - OAuth client ID (for login)
- `GITKRAKEN_CLIENT_SECRET` - OAuth client secret (for login)
- `GK_CONFIG_PATH` - Override config file path

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for contribution guidelines.
