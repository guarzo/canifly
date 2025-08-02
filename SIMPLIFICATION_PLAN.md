# CanIFly Build System Simplification Plan

## Current Issues

### 1. Version Management Complexity
- **Problem**: Version is maintained in 3 places:
  - `version` file (0.0.52)
  - `internal/cmd/start.go` (hardcoded "0.0.40" - out of sync!)
  - `package.json` (managed by scripts)
- **Impact**: Version mismatch causes confusion and potential bugs

### 2. GitHub Workflows Issues
- **Still creating .env files**: Both workflows create `internal/embed/config/.env` despite removal plan
- **Redundant test job**: `release.yml` duplicates the test job from `test.yml`
- **Go version**: Using non-existent Go 1.24.4 (should be 1.23.x)
- **Node version**: Using 24.2.0 (very specific, could use 22.x)

### 3. Build Process Complexity
- **Multiple entry points**: Makefile → npm scripts → bash scripts
- **Redundant npm scripts**: Most just call make commands
- **Platform-specific complexity**: Could be simplified

### 4. Bump Script Issues
- **Hardcoded branch**: Always pushes to 'main'
- **No dry-run option**: Immediately commits and pushes
- **Requires jq**: Additional dependency for simple version increment

## Simplification Recommendations

### 1. Unified Version Management

**Option A: Use Go's embed for version (Recommended)**
```go
//go:embed version.txt
var versionFile string
var Version = strings.TrimSpace(versionFile)
```

**Option B: Use build-time injection**
```bash
go build -ldflags "-X main.Version=$(cat version)"
```

### 2. Simplified GitHub Workflows

**Updated test.yml:**
```yaml
name: Test
on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.23'
      - uses: actions/setup-node@v4
        with:
          node-version: '22'
      - run: npm ci
      - run: npm test
```

**Updated release.yml:**
```yaml
name: Release
on:
  push:
    tags:
      - 'v*'

jobs:
  release:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      # ... build steps ...
```

### 3. Simplified Makefile

Remove redundant targets and focus on core functionality:

```makefile
.PHONY: dev test build clean

# Development
dev:
	npm start

# Testing
test:
	go test ./...
	cd renderer && npm test

# Building
build:
	bash scripts/build.sh all

# Cleaning
clean:
	rm -rf dist release renderer/dist
```

### 4. Simplified Version Bumping

**New bump.sh:**
```bash
#!/bin/bash
set -e

# Read current version
current=$(cat version)
echo "Current version: $current"

# Increment patch version
new=$(echo $current | awk -F. '{$NF = $NF + 1;} 1' OFS=.)
echo "New version: $new"

# Update files
echo $new > version
sed -i "s/\"version\": \".*\"/\"version\": \"$new\"/" package.json

# Commit and tag
git add version package.json
git commit -m "Release v$new"
git tag "v$new"

echo "Done! Run 'git push && git push --tags' to release"
```

### 5. Remove Redundancies

**package.json scripts:**
```json
{
  "scripts": {
    "start": "concurrently \"go run .\" \"cd renderer && npm run dev\" \"wait-on http://localhost:3113 && electron .\"",
    "test": "go test ./... && cd renderer && npm test",
    "build": "bash scripts/build.sh all",
    "dist": "electron-builder",
    "bump": "bash scripts/bump.sh"
  }
}
```

## Implementation Priority

1. **High Priority**:
   - Fix version synchronization issue
   - Remove .env file creation from workflows
   - Fix Go version in workflows

2. **Medium Priority**:
   - Simplify Makefile
   - Update bump script
   - Consolidate npm scripts

3. **Low Priority**:
   - Further workflow optimizations
   - Additional build script improvements

## Benefits

1. **Reduced Complexity**: Fewer files to maintain, clearer build flow
2. **Better Reliability**: Single source of truth for versions
3. **Easier Development**: Simpler commands, fewer dependencies
4. **Faster CI/CD**: Removed redundant steps
5. **No More .env Files**: Aligns with first-run configuration approach