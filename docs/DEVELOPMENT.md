# CanIFly Development Guide

This guide provides detailed information for developers working on the CanIFly project.

## Architecture Overview

CanIFly is built with:
- **Backend**: Go 1.23.2 with Gin framework
- **Frontend**: React 18 with Material-UI and Tailwind CSS
- **Desktop**: Electron 32.2.0
- **State Management**: Zustand
- **Build Tools**: Vite, npm

## Project Structure

```
canifly/
├── internal/           # Go backend code
│   ├── cmd/           # Application entry point
│   ├── handlers/      # HTTP request handlers
│   ├── services/      # Business logic
│   │   ├── eve/       # EVE data service
│   │   ├── fuzzworks/ # Fuzzworks integration
│   │   └── skillplans/# GitHub skill plan downloads
│   ├── persist/       # Data persistence layer
│   ├── model/         # Data models
│   ├── http/          # HTTP utilities
│   └── server/        # Server setup and routing
├── renderer/          # React frontend
│   └── src/
│       ├── api/       # API client services
│       ├── components/# UI components
│       ├── pages/     # Route pages
│       ├── hooks/     # Custom React hooks
│       └── stores/    # Zustand state stores
├── scripts/           # Build and utility scripts
├── docs/              # Documentation
├── static/            # Static assets
└── dist/              # Build output
```

## Development Setup

### 1. Prerequisites

- Go 1.23.2 or higher
- Node.js 20.0.0 or higher
- npm 10.7.0 or higher
- Git

### 2. Clone the Repository

```bash
git clone https://github.com/guarzo/canifly.git
cd canifly
```

### 3. Install Dependencies

```bash
# Install all dependencies
make install

# Or manually:
npm install
cd renderer && npm install
go mod download
```

### 4. Environment Configuration

Create a `.env` file in the root directory:

```env
EVE_CLIENT_ID=<your_client_id>
EVE_CLIENT_SECRET=<your_client_secret>
EVE_CALLBACK_URL=<your_callback_url>
SECRET_KEY=<your_generated_secret_key>
DEV_MODE=true
SKILLPLANS_REPO_URL=https://raw.githubusercontent.com/[org]/[repo]/main
```

Generate a secret key:
```bash
openssl rand -base64 32
```

### 5. EVE Developer Application

1. Go to [EVE Online Developers](https://developers.eveonline.com)
2. Create a new application
3. Set callback URL (e.g., `http://localhost:42423/auth/callback`)
4. Request required scopes:
   - publicData
   - esi-skills.read_skills.v1
   - esi-skills.read_skillqueue.v1
   - esi-clones.read_clones.v1
   - esi-location.read_location.v1
   - esi-ui.open_window.v1
   - esi-universe.read_structures.v1

## Development Workflow

### Running the Application

```bash
# Start all services (recommended)
make dev
# or
npm start

# Run individual services:
npm run dev:go      # Backend only
npm run dev:react   # Frontend only
npm run dev:electron # Electron (needs React running)
```

### Testing

```bash
# Run all tests
make test

# Run specific tests
npm run test:go    # Go tests
npm run test:react # React tests

# Run a single Go test
go test ./internal/services -run TestSpecificFunction
```

### Building

```bash
# Build everything
make build

# Build specific components
npm run build:go    # Backend for all platforms
npm run build:react # Frontend
```

## Key Features and Implementation

### 1. Fuzzworks Integration

The application downloads EVE static data from Fuzzworks:

- **Service**: `internal/services/fuzzworks/service.go`
- **Features**:
  - Automatic download on startup
  - ETag-based caching
  - Manual refresh endpoint
  - No embedded CSV files

### 2. GitHub Skill Plans

Skill plans are downloaded from a GitHub repository:

- **Service**: `internal/services/skillplans/github_downloader.go`
- **Features**:
  - 24-hour cache
  - Graceful fallback to local files
  - Support for icon identifiers in plans
  - Manual refresh endpoint

### 3. WebSocket Real-time Updates

- **Backend**: WebSocket server in handlers
- **Frontend**: WebSocket hook for real-time updates
- **Events**: Character updates, skill plan changes

### 4. Session Management

- **Backend**: Session middleware with secure cookies
- **Frontend**: Auth store with automatic token refresh
- **Security**: Encrypted token storage

## API Endpoints

### Core Endpoints

```
GET    /api/accounts              # List accounts
GET    /api/accounts/{id}          # Get account
PATCH  /api/accounts/{id}          # Update account
DELETE /api/accounts/{id}          # Delete account

GET    /api/config                 # Get configuration
PATCH  /api/config                 # Update configuration

GET    /api/skill-plans            # List skill plans
POST   /api/skill-plans            # Create skill plan
GET    /api/skill-plans/{name}     # Get skill plan
DELETE /api/skill-plans/{name}     # Delete skill plan
POST   /api/skill-plans/refresh    # Refresh from GitHub

GET    /api/fuzzworks/status       # Fuzzworks data status
POST   /api/fuzzworks/update       # Manual update

GET    /api/ws                     # WebSocket endpoint
```

## Code Style and Standards

### Go Code

- Use `golangci-lint` for linting
- Follow standard Go conventions
- Use structured logging with logrus
- Handle errors explicitly
- Write tests for new features

### React Code

- Follow ESLint configuration
- Use TypeScript for type safety
- Implement proper error boundaries
- Use React hooks appropriately
- Keep components focused and reusable

### Git Workflow

1. Create feature branches from `main`
2. Write descriptive commit messages
3. Include tests with new features
4. Update documentation as needed
5. Create pull requests for review

## Common Tasks

### Adding a New API Endpoint

1. Define handler in `internal/handlers/`
2. Add route in `internal/server/router.go`
3. Implement service logic in `internal/services/`
4. Add frontend API client in `renderer/src/api/`
5. Update TypeScript types
6. Add tests

### Adding a New Skill Plan

1. Create plan file in GitHub repository
2. Optional: Add icon identifier as first line
3. Trigger manual refresh or wait for cache expiry

### Debugging

- Backend logs to console with logrus
- Frontend uses console logging
- Check browser DevTools for network requests
- Use React Developer Tools for component inspection

## Troubleshooting

### Common Issues

1. **Port already in use**
   - Kill existing processes on ports 42423 (backend) or 3113 (frontend)

2. **ESI API errors**
   - Check token expiration
   - Verify scopes are correct
   - Check EVE server status

3. **Build failures**
   - Clean build artifacts: `make clean`
   - Reinstall dependencies
   - Check Go and Node versions

### Getting Help

- Check existing issues on GitHub
- Review logs for error messages
- Join EVE developer communities
- Create detailed bug reports

## Performance Considerations

- Use pagination for large datasets
- Implement proper caching strategies
- Minimize API calls with batch operations
- Use WebSocket for real-time updates instead of polling
- Lazy load components and data

## Security Best Practices

- Never commit credentials
- Use environment variables for secrets
- Validate all user input
- Implement proper CORS policies
- Keep dependencies updated
- Use HTTPS in production