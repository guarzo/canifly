# CanIFly Build Guide

This document describes how to build and package the CanIFly application.

## Prerequisites

- **Go** 1.23.2 or higher
- **Node.js** 20.0.0 or higher
- **npm** 10.7.0 or higher
- **Make** (optional but recommended)

## Quick Start

The project now uses a Makefile for simplified build commands:

```bash
# Install all dependencies
make deps

# Run development mode
make dev

# Build everything
make build

# Run tests
make test

# Package for distribution
make package
```

## Development

### Running in Development Mode

```bash
# Start all services (Go backend, React frontend, Electron)
make dev

# Or run individually:
make dev-go      # Backend only
make dev-react   # Frontend only
make dev-electron # Electron only
```

### Environment Variables

Create a `.env` file in the root directory:

```env
EVE_CLIENT_ID=<your_client_id>
EVE_CLIENT_SECRET=<your_client_secret>
EVE_CALLBACK_URL=<your_callback_url>
SECRET_KEY=<your_generated_secret_key>
DEV_MODE=true  # For development
```

## Building

### Build All Components

```bash
make build
```

This will:
1. Build Go backend for all platforms (macOS, Windows, Linux)
2. Build React frontend

### Platform-Specific Builds

```bash
# Backend only
make build-go         # All platforms
make build-go-mac     # macOS only
make build-go-win     # Windows only
make build-go-linux   # Linux only

# Frontend only
make build-react
```

### Using the Build Script Directly

```bash
# Build everything
./scripts/build.sh all

# Build specific targets
./scripts/build.sh mac      # macOS backend
./scripts/build.sh win      # Windows backend
./scripts/build.sh linux    # Linux backend
./scripts/build.sh frontend # React frontend
./scripts/build.sh backend  # All backend platforms
```

## Testing

```bash
# Run all tests
make test

# Run specific tests
make test-go    # Go tests
make test-react # React tests
```

## Linting

```bash
# Run all linters
make lint

# Run specific linters
make lint-go    # Go linter (requires golangci-lint)
make lint-react # React linter
```

## Packaging for Distribution

### Package for All Platforms

```bash
make package
```

This will:
1. Clean previous builds
2. Build all components
3. Create distribution packages using electron-builder

### Platform-Specific Packages

```bash
make package-mac    # macOS DMG
make package-win    # Windows installer
make package-linux  # Linux AppImage and deb
```

### Output Locations

- **Built binaries**: `dist/` directory
  - `dist/mac/canifly-backend` - macOS binary
  - `dist/win/canifly-backend.exe` - Windows binary
  - `dist/linux/canifly-backend` - Linux binary
- **React build**: `renderer/dist/`
- **Distribution packages**: `release/` directory

## Version Management

```bash
# Create a release tag
make release
# Note: Version must be manually updated in 'version' file and package.json first
```

## Cleaning

```bash
# Remove all build artifacts
make clean
```

## Troubleshooting

### Check Environment

```bash
make check
```

This will display versions of all required tools.

### Common Issues

1. **Go build fails**: Ensure Go modules are downloaded
   ```bash
   go mod download
   ```

2. **React build fails**: Clear cache and reinstall
   ```bash
   cd renderer
   rm -rf node_modules
   npm install
   ```

3. **Electron packaging fails**: Ensure all platforms are built
   ```bash
   make clean
   make build
   make package
   ```

## Legacy npm Scripts

The following npm scripts are still available for backward compatibility:

```bash
npm start        # Start development
npm run build    # Build all
npm run test     # Run tests
npm run dist     # Package application
```

However, using the Makefile commands is recommended for consistency.

## CI/CD

For continuous integration, use:

```bash
# Full pipeline
make all

# Or step by step
make clean
make deps
make build
make test
make package
```