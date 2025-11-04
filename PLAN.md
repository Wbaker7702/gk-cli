# GitKraken CLI Development Plan

## Overview
Build a comprehensive CLI tool (`gk`) for GitKraken that provides workspace management, pull request/issue access, cloud patches, and integration with GitKraken Desktop and GitLens.

## Technology Stack
- **Language**: Go (golang) - Excellent for CLI tools, cross-platform support
- **CLI Framework**: [cobra](https://github.com/spf13/cobra) - Industry standard for Go CLIs
- **Config Management**: [viper](https://github.com/spf13/viper) - Configuration management
- **HTTP Client**: `net/http` or `resty` for API calls
- **JSON**: `encoding/json` for data handling

## Core Features to Implement

### 1. Authentication & Configuration
- `gk login` - Browser-based OAuth flow
- `gk logout` - Clear authentication
- `gk setting` - Configuration management
  - Theme management
  - Provider configuration

### 2. Workspace Management
- `gk ws create` - Create local/cloud workspaces
- `gk ws add-repo` - Add repositories to workspace
- `gk ws locate` - Locate local repos
- `gk ws clone` - Clone all repos in workspace
- `gk ws [action]` - Git operations on multiple repos (fetch, pull, push, checkout)
- `gk ws insights` - PR metrics and insights

### 3. Provider Integration
- `gk provider add` - Connect GitHub/GitLab/Bitbucket
- Provider authentication and token management

### 4. Pull Requests & Issues
- `gk pr list` - List all PRs in workspace
- `gk pr view` - View PR details
- `gk pr suggest` - Code suggestions for PRs

### 5. Cloud Patches
- `gk patch create` - Create cloud patch
- `gk patch apply` - Apply cloud patch
- `gk patch view` - Preview patch changes
- `gk patch list` - List all patches
- `gk patch delete` - Delete patch

### 6. Launchpad
- `gk launchpad` - Interactive dashboard for PRs, Issues, WIPs
- Pin/unpin items
- Snooze/unsnooze items

### 7. Visual Graph
- `gk graph` - Open visual commit graph in GitKraken Desktop/GitLens

## Project Structure

```
gk-cli/
├── cmd/
│   ├── root.go          # Root command and CLI setup
│   ├── login.go         # Authentication commands
│   ├── workspace.go     # Workspace commands
│   ├── provider.go      # Provider commands
│   ├── pr.go            # Pull request commands
│   ├── patch.go         # Cloud patch commands
│   ├── launchpad.go     # Launchpad command
│   └── graph.go         # Graph command
├── internal/
│   ├── api/             # API client for GitKraken services
│   ├── auth/            # Authentication logic
│   ├── config/          # Configuration management
│   ├── workspace/       # Workspace operations
│   ├── provider/        # Provider integrations
│   └── theme/           # Theme system
├── pkg/
│   └── utils/           # Shared utilities
├── themes/              # Theme files
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Implementation Phases

### Phase 1: Foundation (Current)
- [x] Project structure setup
- [ ] Basic CLI framework with cobra
- [ ] Configuration system
- [ ] Authentication flow skeleton

### Phase 2: Core Commands
- [ ] Workspace management
- [ ] Provider integration
- [ ] Basic PR/Issue listing

### Phase 3: Advanced Features
- [ ] Cloud Patches
- [ ] Code Suggest
- [ ] Launchpad UI
- [ ] Theme system

### Phase 4: Polish
- [ ] Error handling
- [ ] Testing
- [ ] Documentation
- [ ] Cross-platform builds

## Configuration Files

### Config Location
- Linux/Mac: `~/.config/gk/config.json` or `~/.gk/config.json`
- Windows: `%APPDATA%\gk\config.json`

### Config Structure
```json
{
  "auth": {
    "token": "...",
    "refresh_token": "...",
    "expires_at": "..."
  },
  "theme": "default",
  "providers": {
    "github": { "token": "..." },
    "gitlab": { "token": "..." }
  },
  "workspaces": {
    "default": { "type": "local", "repos": [...] }
  }
}
```

## API Endpoints (To be determined)
- GitKraken API base URL needed
- OAuth endpoints
- Workspace endpoints
- Patch endpoints
- Provider endpoints

## Dependencies
- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration
- `github.com/go-resty/resty/v2` - HTTP client (optional)
- `golang.org/x/oauth2` - OAuth support
