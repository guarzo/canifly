# CanIFly Project Review & Simplification Plan

## Executive Summary

This document outlines a comprehensive review of the CanIFly codebase with specific recommendations for simplification and improved file path handling to prevent startup crashes.

## 🚨 Critical Issues: File Path Handling

### Problems Identified

1. **No directory creation at startup**
   - The base config directory (`~/.config/canifly`) is never created
   - Causes crashes on first run when services attempt to write files
   - Location: `internal/server/env.go:58-62`

2. **Fatal error handling**
   - Uses `logger.Fatal()` on secret key generation failure
   - Terminates application instead of graceful degradation
   - Location: `internal/server/env.go:73`

3. **Missing validation**
   - No checks for directory existence before file operations
   - No fallback paths when primary directories fail

### Recommended Fixes

```go
// In internal/server/env.go after line 62
cfg.BasePath = filepath.Join(configDir, "canifly")
if err := os.MkdirAll(cfg.BasePath, 0755); err != nil {
    // Fallback to temp directory
    cfg.BasePath = filepath.Join(os.TempDir(), "canifly")
    logger.Warnf("Failed to create config dir, using temp: %v", err)
    os.MkdirAll(cfg.BasePath, 0755)
}

// Replace logger.Fatal with error return
if err := generateAndSaveSecretKey(cfg.BasePath); err != nil {
    return cfg, fmt.Errorf("failed to generate secret key: %w", err)
}
```

## 📦 Major Simplification Opportunities

### 1. Service Layer Consolidation

Current architecture has 15+ services with excessive abstraction layers:
- **Interfaces** → **Services** → **Repositories** → **Stores**

#### Proposed Consolidation

Reduce to 5-6 core services:

| New Service | Combines | Responsibility |
|------------|----------|----------------|
| `AccountManagementService` | Account, Character, Association | All account-related operations |
| `ConfigurationService` | Config, AppState, Dashboard | Application configuration and state |
| `EVEDataService` | ESI, Cache, EveData | EVE Online data and API interactions |
| `StorageService` | All file stores | Unified file persistence layer |
| `SyncService` | Sync operations | Settings synchronization |

### 2. Remove Unnecessary Abstractions

- Eliminate repository interfaces that only wrap stores
- Direct service-to-store communication
- Reduce from 4-layer to 2-3 layer architecture

### 3. Frontend State Management

Current issues:
- Prop drilling through multiple component levels
- Duplicate data fetching logic (`getAppData` vs `getAppDataNoCache`)
- Complex refresh logic with multiple flags

Solution:
- Implement Zustand for lightweight state management (3KB)
- Create unified authentication state
- Consolidate data fetching patterns

## 🔄 Code Duplication Remediation

### 1. Backend Utilities

Create shared error handling utilities:

```go
// internal/handlers/utils.go
package handlers

import (
    "fmt"
    "net/http"
)

func HandleServiceError(w http.ResponseWriter, logger Logger, err error, action string) {
    logger.Errorf("Failed to %s: %v", action, err)
    respondError(w, fmt.Sprintf("Failed to %s", action), http.StatusInternalServerError)
}

func DecodeAndValidate[T any](r *http.Request, w http.ResponseWriter) (*T, bool) {
    var req T
    if err := decodeJSONBody(r, &req); err != nil {
        respondError(w, "Invalid request body", http.StatusBadRequest)
        return nil, false
    }
    return &req, true
}
```

### 2. Frontend Custom Hooks

Create reusable async operation hook:

```jsx
// renderer/src/hooks/useAsyncOperation.js
import { useState, useCallback } from 'react';
import { toast } from 'react-toastify';

export function useAsyncOperation() {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState(null);
    
    const execute = useCallback(async (operation, successMessage) => {
        setIsLoading(true);
        setError(null);
        try {
            const result = await operation();
            if (result?.success) {
                toast.success(result.message || successMessage);
            }
            return result;
        } catch (err) {
            setError(err);
            console.error('Operation failed:', err);
            throw err;
        } finally {
            setIsLoading(false);
        }
    }, []);
    
    return { execute, isLoading, error };
}
```

### 3. API Service Simplification

Create helper functions for common request patterns:

```jsx
// renderer/src/api/helpers.js
export function postRequest(endpoint, data, options = {}) {
    return apiRequest(endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(data)
    }, options);
}

export function getRequest(endpoint, options = {}) {
    return apiRequest(endpoint, {
        method: 'GET',
        credentials: 'include'
    }, options);
}
```

## 🛠️ Implementation Plan

### Phase 1: Fix Critical Issues (1-2 days)

1. **Fix file path handling**
   - Add directory creation at startup
   - Implement fallback paths
   - Add validation before file operations

2. **Remove fatal error calls**
   - Replace `logger.Fatal()` with error returns
   - Implement graceful degradation

3. **Add startup validation**
   - Check all required paths
   - Verify permissions
   - Log warnings for non-critical failures

### Phase 2: Consolidate Services (3-5 days)

1. **Create new service structure**
   - Define new consolidated service interfaces
   - Implement service consolidation
   - Maintain backward compatibility during transition

2. **Migrate functionality**
   - Move logic from old services to new
   - Update dependency injection
   - Remove deprecated services

3. **Update handlers**
   - Point to new services
   - Simplify handler logic
   - Remove redundant code

### Phase 3: Frontend Simplification (2-3 days)

1. **Implement state management**
   - Add Zustand
   - Create global stores for auth, config, and data
   - Remove prop drilling

2. **Create custom hooks**
   - `useAsyncOperation` for API calls
   - `useAuth` for authentication state
   - `useAppData` for centralized data access

3. **Simplify components**
   - Remove unnecessary state
   - Use new hooks
   - Consolidate duplicate logic

### Phase 4: API & Build Simplification (2 days)

1. **RESTful API conversion**
   - From: `/api/toggle-account-status`
   - To: `PATCH /api/accounts/:id`
   - Consolidate 20+ endpoints to ~10

2. **Build script simplification**
   - Create Makefile for common tasks
   - Consolidate platform builds
   - Simplify npm scripts

## 📊 Metrics & Benefits

### Expected Improvements

| Metric | Current | Target | Improvement |
|--------|---------|--------|-------------|
| Service Count | 15+ | 5-6 | 67% reduction |
| API Endpoints | 20+ | ~10 | 50% reduction |
| Code Duplication | High | Low | ~40% reduction |
| Startup Reliability | Crashes on first run | Robust | 100% improvement |

### Code Complexity Reduction

- **Before**: 4-layer architecture with 15+ services
- **After**: 2-3 layer architecture with 5-6 services
- **Result**: 40% reduction in code complexity

## 🎯 Quick Wins

1. **Fix startup crashes** (30 minutes)
   - Add directory creation logic
   - Implement fallback paths

2. **Create error utilities** (1 hour)
   - Reduce 50+ duplicate error handling blocks
   - Improve consistency

3. **Consolidate API endpoints** (2 hours)
   - Convert to RESTful patterns
   - Reduce endpoint count by 50%

4. **Create unified FileStore** (3 hours)
   - Replace 5 separate file stores
   - Centralize file operations

## 📋 Dependencies

### Current Stack (Keep As-Is)

**Backend (Go)**
- gorilla/mux, sessions - Routing and session management
- logrus - Structured logging
- oauth2 - EVE SSO integration
- go-cache - In-memory caching

**Frontend (React)**
- Material-UI - Component library
- React Router - Navigation
- Axios - HTTP client
- Tailwind CSS - Styling

### Recommended Additions

- **zustand** (3KB) - Lightweight state management
- **zod** (optional) - Schema validation for both frontend/backend

## 🚀 Getting Started

1. **Immediate Actions**
   - Fix file path handling to prevent crashes
   - Create error handling utilities

2. **Short Term** (1 week)
   - Consolidate services
   - Implement state management

3. **Medium Term** (2 weeks)
   - Complete API simplification
   - Refactor build process

## 📝 Notes

- All changes maintain backward compatibility
- No functionality is removed, only simplified
- Focus on developer experience and maintainability
- Prioritize stability over feature additions during refactoring