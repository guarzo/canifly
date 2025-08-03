# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

CanIFly is an Electron desktop application for EVE Online players that helps manage characters, check flyable ships, manage skill plans, and sync in-game settings. It consists of:
- Go backend (v1.23.2) serving REST APIs
- React frontend (v18) with Material-UI and Tailwind CSS
- Electron wrapper (v32.2.0) for desktop distribution

## Key Commands

### Development
```bash
# Start all components (Go backend, React dev server, Electron)
npm start

# Run individual components
npm run dev:go      # Go backend only
npm run dev:react   # React frontend only
npm run dev:electron # Electron (requires React server running)
```

### Testing
```bash
# Run all tests
npm test

# Run specific tests
npm run test:go    # Go tests: go test ./...
npm run test:react # React tests with Vitest

# Run a single Go test
go test ./internal/services -run TestSpecificFunction
```

### Building & Packaging
```bash
# Build both backend and frontend
npm run build

# Build Go backend for all platforms (Mac, Windows, Linux)
npm run go:build

# Create Electron distribution packages
npm run dist

# Full packaging (build + distribute)
npm run package:app
```

### Version Management
```bash
# Create a release tag
npm run release
# Note: Version is manually updated in 'version' file and package.json
```

## Architecture

### Backend Structure (Go)
The Go backend follows internal package pattern with clean architecture:
- `internal/cmd/` - Application entry point, initializes server
- `internal/handlers/` - HTTP handlers for API endpoints
- `internal/services/` - Business logic layer
- `internal/persist/` - Data persistence (stores for accounts, config, EVE data)
- `internal/model/` - Data models and structures
- `internal/http/` - HTTP client, middleware, session management
- `internal/embed/` - (Removed - skill plans now downloaded from GitHub)

### Frontend Structure (React)
- `renderer/src/api/` - API client services
- `renderer/src/components/` - UI components organized by feature
- `renderer/src/pages/` - Route page components
- `renderer/src/hooks/` - Custom React hooks
- `renderer/src/utils/` - Utility functions

### Key API Patterns
Backend serves REST APIs at `http://localhost:42423` with routes:
- `/api/accounts/*` - Account management
- `/api/config/*` - Configuration endpoints
- `/api/esi/*` - EVE ESI proxy endpoints
- `/api/skill-plans/*` - Skill plan management
- `/api/fuzzworks/*` - Fuzzworks data updates
- `/api/sync/*` - Settings synchronization

Frontend uses axios-based API services in `renderer/src/api/`.

## Development Setup

1. **Environment Variables**: Create `.env` file in root:
```
EVE_CLIENT_ID=<your_client_id>
EVE_CLIENT_SECRET=<your_client_secret>
EVE_CALLBACK_URL=<your_callback_url>
SECRET_KEY=<your_generated_secret_key>
```

2. **EVE Developer App**: Register at https://developers.eveonline.com
   - Set callback URL to match your environment
   - Request required scopes for character/skill/asset access

3. **Dependencies**: 
   - Go 1.23.2+
   - npm 10.7.0+
   - Run `npm install` in root directory

## Code Quality Standards

- **Go**: Use `golangci-lint` for linting
- **React**: ESLint configuration in `renderer/.eslintrc.cjs`
- **Testing**: Write tests for new features, maintain existing test coverage
- **Error Handling**: Use structured errors with proper logging
- **Logging**: Use logrus with appropriate log levels

## Common Development Tasks

### Adding New API Endpoint
1. Define handler in `internal/handlers/`
2. Add route in `internal/server/server.go`
3. Implement service logic in `internal/services/`
4. Add frontend API client in `renderer/src/api/`

### Working with EVE Data
- Fuzzworks integration for EVE static data downloads
- GitHub integration for skill plan downloads
- EVE ESI client integrated in `internal/services/eve/`
- Character data models in `internal/model/`

### Electron Updates
- Main process: `main.js`
- Preload script: `preload.js`
- Auto-updater configuration in main process