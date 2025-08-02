# Technical Debt Report - CanIFly Codebase

Generated: August 2, 2025

## Executive Summary

This report documents technical debt identified in the CanIFly codebase through comprehensive analysis. The codebase shows signs of being in transition with mixed patterns, legacy code, and incomplete migrations. While generally well-structured, several critical issues need immediate attention.

## Critical Issues

### 1. Port Configuration Mismatch
**Severity**: 🔴 Critical

The codebase has conflicting port configurations:
- Documentation (CLAUDE.md, PRD) states backend runs on port **42423**
- Implementation uses port **8713** throughout:
  - `internal/server/env.go:100`: Default port 8713
  - `renderer/src/Config.jsx:10`: Hardcoded `http://localhost:8713`
  - `renderer/vite.config.js:16`: Proxy target 8713
  - `main.js:149`: Health check on 8713
  - Test files reference 8713

**Impact**: Developer confusion, potential runtime failures, deployment issues

### 2. Console.log Pollution
**Severity**: 🔴 Critical

Found **80 console.log statements** across 19 files:
- Production code contains debug statements
- Sensitive data potentially exposed in logs
- No structured logging in frontend

**Examples**:
```javascript
// Routes.jsx:26-27
console.log('Routes - associations:', associations);
console.log('Routes - eveProfiles:', eveProfiles);

// main.js - Multiple instances
console.log('[CanIFly] Main process starting...');
console.log('[CanIFly] App is ready');
```

**Impact**: Performance overhead, security risks, unprofessional production logs

### 3. Embedded Environment Variables
**Severity**: 🔴 Critical

- `internal/embed/config/.env` contains embedded environment template
- Requires users to manually configure before building
- No first-run configuration UI

**Impact**: Poor user experience, security concerns, build complexity

## Major Issues

### 4. Incomplete API Migration
**Severity**: 🟡 High

`renderer/src/api/apiService.jsx` explicitly states:
> "This file contains both new RESTful API functions and legacy functions. The backend supports both patterns during the transition period."

- Mixed RESTful and legacy endpoints
- No deprecation timeline
- Both patterns actively used

**Impact**: Technical debt accumulation, maintenance overhead

### 5. Legacy npm Scripts
**Severity**: 🟡 Medium

`package.json:18` contains legacy scripts marked with:
```json
"// Legacy scripts (use make commands instead)"
```

Duplicate scripts maintained:
- `dev:react`, `dev:electron`, `dev:go`
- `test:go`, `test:react`

**Impact**: Confusion about which build system to use, maintenance overhead

### 6. Outdated React Patterns
**Severity**: 🟡 Medium

`ErrorBoundary.jsx` uses class component syntax:
- Only class component in otherwise modern codebase
- Uses deprecated `componentDidCatch`
- Should use functional component with error boundary hook

**Impact**: Inconsistent codebase, harder onboarding for new developers

## Minor Issues

### 7. Development Files in Repository
**Severity**: 🟢 Low

- `notes` file contains TODOs and .env template
- Should be in `.gitignore` or proper documentation

### 8. Naming Inconsistencies
**Severity**: 🟢 Low

- Test file typo: `SkillPlans.text.jsx` → `SkillPlans.test.jsx`
- Mixed API endpoint naming:
  - `/api/skill-plans` (kebab-case)
  - `/api/accounts` (lowercase)

### 9. Untracked TODOs
**Severity**: 🟢 Low

TODOs scattered in files without issue tracking:
- `README.md`: "clones, implants, corp sort"
- `notes`: Multiple items including "test all the things", "make video"

### 10. Deprecated Dependencies
**Severity**: 🟢 Low

Package-lock.json shows multiple deprecated packages:
- Q promise library (use native promises)
- Glob versions < v9
- Several packages marked "no longer supported"

## Technical Debt Metrics

| Category | Count | Files Affected | Severity |
|----------|-------|----------------|----------|
| Port misconfigurations | 7 | 5 files | Critical |
| Console.log statements | 80 | 19 files | Critical |
| Legacy API endpoints | ~10+ | Multiple | High |
| Legacy npm scripts | 6 | 1 file | Medium |
| Class components | 1 | 1 file | Low |
| TODO comments | 7 | 2 files | Low |
| Deprecated packages | ~5 | package-lock.json | Low |

## Code Smells Identified

1. **Mixed Error Handling Patterns**
   - Backend uses structured logging with logrus
   - Frontend uses console.log/error
   - Main process uses raw console.log

2. **Environment Configuration**
   - Hardcoded URLs in multiple places
   - Inconsistent DEV_MODE usage
   - No centralized configuration management

3. **Test Coverage Gaps**
   - Some test files exist but coverage unclear
   - No visible test strategy documentation
   - Missing integration tests

## Security Concerns

1. **Embedded Credentials Template**
   - .env file embedded in binary
   - No secure credential input mechanism

2. **Debug Information Leakage**
   - Console.log statements may expose sensitive data
   - No log level control in production

3. **Hardcoded Configuration**
   - URLs and ports hardcoded
   - No environment-based configuration

## Recommendations Summary

### Immediate Actions (Critical)
1. Standardize port configuration across entire codebase
2. Remove all console.log statements, implement proper logging
3. Implement first-run configuration UI for EVE credentials
4. Remove embedded .env file

### Short-term (1-2 weeks)
1. Complete API migration to RESTful patterns
2. Remove legacy npm scripts
3. Convert ErrorBoundary to functional component
4. Fix naming inconsistencies

### Long-term (1 month)
1. Implement comprehensive test strategy
2. Update all deprecated dependencies
3. Create proper configuration management system
4. Add pre-commit hooks for code quality

## Conclusion

While the CanIFly codebase demonstrates good overall architecture, it's clearly in a transition phase with several critical issues that need immediate attention. The primary concerns are around configuration management, logging practices, and incomplete migrations. Addressing these issues will significantly improve code quality, security, and maintainability.