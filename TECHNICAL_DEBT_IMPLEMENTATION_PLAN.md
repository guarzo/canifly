# Technical Debt Implementation Plan

## Overview

This plan addresses all technical debt identified in the Technical Debt Report, organized by priority and estimated effort. Each phase builds upon the previous one to ensure smooth implementation.

## Phase 1: Critical Security & Configuration (2-3 days)

### 1.1 Remove Embedded .env and Implement First-Run Configuration

**Priority**: 🔴 Critical  
**Effort**: 1 day

#### Tasks:
1. **Remove embedded .env file**
   ```bash
   rm internal/embed/config/.env
   ```

2. **Create first-run configuration UI**
   ```go
   // internal/services/config/first_run.go
   type FirstRunService struct {
       configService *ConfigurationService
       logger        interfaces.Logger
   }
   
   func (s *FirstRunService) IsFirstRun() bool {
       // Check if EVE credentials exist
   }
   
   func (s *FirstRunService) SaveCredentials(clientID, clientSecret string) error {
       // Securely store credentials
   }
   ```

3. **Add frontend configuration dialog**
   ```jsx
   // renderer/src/components/setup/FirstRunDialog.jsx
   const FirstRunDialog = () => {
       // Modal with client ID/secret input fields
       // Callback URL displayed as read-only
   };
   ```

4. **Update main.js to check first-run**
   ```javascript
   // Check if configuration exists
   const isFirstRun = await checkFirstRun();
   if (isFirstRun) {
       // Show configuration window
   }
   ```

5. **Hardcode callback URL**
   ```go
   const DefaultCallbackURL = "http://localhost:42423/callback"
   ```

### 1.2 Standardize Port Configuration

**Priority**: 🔴 Critical  
**Effort**: 0.5 days

#### Tasks:
1. **Create constants file**
   ```go
   // internal/config/constants.go
   const (
       DefaultPort = "42423"
       DefaultHost = "localhost"
   )
   ```

2. **Update all port references**
   - `internal/server/env.go`: Use DefaultPort
   - `renderer/src/Config.jsx`: Update to 42423
   - `renderer/vite.config.js`: Update proxy
   - `main.js`: Update health check
   - All test files

3. **Add environment variable support**
   ```go
   port := os.Getenv("CANIFLY_PORT")
   if port == "" {
       port = DefaultPort
   }
   ```

### 1.3 Remove Console.log Statements

**Priority**: 🔴 Critical  
**Effort**: 1 day

#### Tasks:
1. **Create frontend logger utility**
   ```javascript
   // renderer/src/utils/logger.js
   const logger = {
       debug: (...args) => isDev && console.log(...args),
       info: (...args) => console.info(...args),
       warn: (...args) => console.warn(...args),
       error: (...args) => console.error(...args)
   };
   ```

2. **Replace all console.log statements**
   ```bash
   # Find and replace script
   find renderer/src -name "*.jsx" -o -name "*.js" | \
     xargs sed -i 's/console\.log/logger.debug/g'
   ```

3. **Add ESLint rule**
   ```javascript
   // .eslintrc.js
   rules: {
       'no-console': ['error', { allow: ['warn', 'error'] }]
   }
   ```

4. **Update main.js logging**
   - Use electron-log package
   - Configure log levels based on environment

## Phase 2: Complete API Migration (3-4 days)

### 2.1 Remove Legacy API Endpoints

**Priority**: 🟡 High  
**Effort**: 2 days

#### Tasks:
1. **Identify all legacy endpoints**
   ```bash
   grep -r "legacy" renderer/src/api/
   ```

2. **Update frontend to use RESTful endpoints only**
   - Remove legacy functions from apiService.jsx
   - Update all components using legacy APIs
   - Test each endpoint conversion

3. **Remove legacy handlers from backend**
   - Delete deprecated handler functions
   - Update router to remove old routes
   - Clean up unused service methods

4. **Update API documentation**
   ```markdown
   // docs/API.md
   ## Deprecated Endpoints (Removed in v2.0)
   - POST /api/legacy/... → Use REST endpoints
   ```

### 2.2 API Naming Consistency

**Priority**: 🟢 Low  
**Effort**: 0.5 days

#### Tasks:
1. **Standardize endpoint naming**
   - Use kebab-case for all endpoints
   - `/api/skillplans` → `/api/skill-plans`
   - Update frontend and backend

2. **Update OpenAPI spec if exists**

## Phase 3: Code Modernization (2-3 days)

### 3.1 Convert ErrorBoundary to Functional Component

**Priority**: 🟡 Medium  
**Effort**: 0.5 days

#### Tasks:
1. **Install react-error-boundary**
   ```bash
   cd renderer && npm install react-error-boundary
   ```

2. **Create new ErrorBoundary**
   ```jsx
   // renderer/src/components/common/ErrorBoundary.jsx
   import { ErrorBoundary } from 'react-error-boundary';
   
   function ErrorFallback({error, resetErrorBoundary}) {
       return (
           <div role="alert">
               <p>Something went wrong:</p>
               <pre>{error.message}</pre>
               <button onClick={resetErrorBoundary}>Try again</button>
           </div>
       );
   }
   
   export default function AppErrorBoundary({children}) {
       return (
           <ErrorBoundary
               FallbackComponent={ErrorFallback}
               onError={(error, info) => {
                   logger.error('Error boundary:', error, info);
               }}
           >
               {children}
           </ErrorBoundary>
       );
   }
   ```

### 3.2 Remove Legacy npm Scripts

**Priority**: 🟡 Medium  
**Effort**: 0.5 days

#### Tasks:
1. **Remove from package.json**
   ```json
   // Remove these scripts:
   "dev:react", "dev:electron", "dev:go",
   "test:go", "test:react"
   ```

2. **Update documentation**
   - README.md: Remove references to npm scripts
   - BUILD.md: Update to use make commands only

3. **Add deprecation notice if needed**
   ```javascript
   // If scripts must stay temporarily
   "dev:react": "echo 'DEPRECATED: Use make dev-react' && exit 1"
   ```

## Phase 4: Cleanup & Maintenance (1-2 days)

### 4.1 Fix Minor Issues

**Priority**: 🟢 Low  
**Effort**: 0.5 days

#### Tasks:
1. **Fix test file name**
   ```bash
   mv renderer/src/pages/SkillPlans.text.jsx \
      renderer/src/pages/SkillPlans.test.jsx
   ```

2. **Remove notes file**
   ```bash
   rm notes
   ```

3. **Create proper .env.example**
   ```bash
   # .env.example
   CANIFLY_PORT=42423
   CANIFLY_HOST=localhost
   ```

### 4.2 Update Dependencies

**Priority**: 🟢 Low  
**Effort**: 1 day

#### Tasks:
1. **Audit and fix dependencies**
   ```bash
   npm audit fix
   cd renderer && npm audit fix
   ```

2. **Update deprecated packages**
   - Replace Q with native promises
   - Update glob to v9+
   - Review and update other deprecated packages

### 4.3 Create GitHub Issues from TODOs

**Priority**: 🟢 Low  
**Effort**: 0.5 days

#### Tasks:
1. **Create issues for each TODO**
   - "clones, implants, corp sort"
   - "test all the things"
   - "test gh packaging"
   - "fix readme"
   - "make video"
   - "improve instructions and allow hiding"

2. **Add issue templates**
   ```yaml
   # .github/ISSUE_TEMPLATE/technical-debt.yml
   name: Technical Debt
   about: Track technical debt items
   ```

## Phase 5: Quality Assurance (2-3 days)

### 5.1 Add Pre-commit Hooks

**Priority**: 🟡 Medium  
**Effort**: 1 day

#### Tasks:
1. **Install husky and lint-staged**
   ```bash
   npm install --save-dev husky lint-staged
   ```

2. **Configure pre-commit hooks**
   ```json
   // package.json
   "lint-staged": {
       "*.{js,jsx}": ["eslint --fix", "git add"],
       "*.go": ["gofmt -w", "golangci-lint run"]
   }
   ```

3. **Add commit message linting**

### 5.2 Comprehensive Testing

**Priority**: 🟡 Medium  
**Effort**: 2 days

#### Tasks:
1. **Test all changes**
   - Port configuration changes
   - API migration completeness
   - First-run configuration flow
   - Error boundary functionality

2. **Add integration tests**
   - Test first-run flow
   - Test API endpoints
   - Test configuration management

## Implementation Schedule

| Week | Phase | Tasks |
|------|-------|-------|
| Week 1 | Phase 1 | Critical security & configuration fixes |
| Week 1-2 | Phase 2 | Complete API migration |
| Week 2 | Phase 3 | Code modernization |
| Week 3 | Phase 4 | Cleanup & maintenance |
| Week 3-4 | Phase 5 | Quality assurance & testing |

## Success Metrics

- [ ] Zero console.log statements in production code
- [ ] All port references use 42423
- [ ] No embedded .env file
- [ ] First-run configuration UI functional
- [ ] All legacy API endpoints removed
- [ ] Zero class components
- [ ] All dependencies up-to-date
- [ ] Pre-commit hooks preventing new tech debt

## Risk Mitigation

1. **Breaking Changes**
   - Test each change thoroughly
   - Implement changes incrementally
   - Keep backups of working state

2. **User Impact**
   - Document configuration changes
   - Provide migration guide
   - Test upgrade path from existing installations

3. **Development Workflow**
   - Communicate changes to team
   - Update development documentation
   - Ensure CI/CD pipelines updated

## Conclusion

This implementation plan provides a structured approach to eliminating technical debt while maintaining system stability. The phased approach ensures critical issues are addressed first, with each phase building upon the previous work. Total estimated effort: 10-15 developer days over 3-4 weeks.