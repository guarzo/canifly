# Technical Debt Implementation Plan (Simplified)

## Overview

Simplified implementation plan for addressing technical debt in CanIFly. Since there are no current users, we can make breaking changes without migration concerns.

## Phase 1: Critical Fixes (2 days)

### 1.1 Remove Embedded .env and Implement First-Run Configuration

**Priority**: 🔴 Critical  
**Effort**: 1 day

#### Tasks:
1. **Delete embedded .env file**
   ```bash
   rm internal/embed/config/.env
   ```

2. **Create first-run configuration service**
   ```go
   // internal/services/config/first_run.go
   type FirstRunService struct {
       configService *ConfigurationService
       logger        interfaces.Logger
   }
   
   func (s *FirstRunService) NeedsConfiguration() bool {
       config, _ := s.configService.GetConfig()
       return config.EVEClientID == "" || config.EVEClientSecret == ""
   }
   
   func (s *FirstRunService) SaveEVECredentials(clientID, clientSecret string) error {
       return s.configService.UpdateConfig(map[string]interface{}{
           "EVEClientID": clientID,
           "EVEClientSecret": clientSecret,
           "EVECallbackURL": "http://localhost:42423/callback",
       })
   }
   ```

3. **Add first-run UI dialog**
   ```jsx
   // renderer/src/components/setup/FirstRunDialog.jsx
   import React, { useState } from 'react';
   import { Dialog, TextField, Button } from '@mui/material';
   
   const FirstRunDialog = ({ open, onComplete }) => {
       const [clientId, setClientId] = useState('');
       const [clientSecret, setClientSecret] = useState('');
       
       const handleSubmit = async () => {
           await apiService.saveEVECredentials(clientId, clientSecret);
           onComplete();
       };
       
       return (
           <Dialog open={open} maxWidth="sm" fullWidth>
               <DialogTitle>EVE Online Application Setup</DialogTitle>
               <DialogContent>
                   <TextField
                       label="Client ID"
                       value={clientId}
                       onChange={(e) => setClientId(e.target.value)}
                       fullWidth
                       margin="normal"
                   />
                   <TextField
                       label="Client Secret"
                       value={clientSecret}
                       onChange={(e) => setClientSecret(e.target.value)}
                       type="password"
                       fullWidth
                       margin="normal"
                   />
                   <Typography variant="caption" color="textSecondary">
                       Callback URL: http://localhost:42423/callback
                   </Typography>
               </DialogContent>
               <DialogActions>
                   <Button onClick={handleSubmit} variant="contained">
                       Save Configuration
                   </Button>
               </DialogActions>
           </Dialog>
       );
   };
   ```

4. **Update App.jsx to check configuration**
   ```jsx
   // In App.jsx
   const [needsConfig, setNeedsConfig] = useState(false);
   
   useEffect(() => {
       checkConfiguration().then(needs => setNeedsConfig(needs));
   }, []);
   
   if (needsConfig) {
       return <FirstRunDialog open={true} onComplete={() => setNeedsConfig(false)} />;
   }
   ```

### 1.2 Fix Port Configuration Everywhere

**Priority**: 🔴 Critical  
**Effort**: 0.5 days

#### Single source of truth - change all 8713 references to 42423:

1. **Backend**
   ```go
   // internal/server/env.go
   if port == "" {
       port = "42423"  // Changed from 8713
   }
   ```

2. **Frontend**
   ```javascript
   // renderer/src/Config.jsx
   const backEndURL = isDev ? '' : 'http://localhost:42423';
   
   // renderer/vite.config.js
   proxy: {
       '/api': {
           target: 'http://localhost:42423',  // Changed from 8713
       }
   }
   ```

3. **Main process**
   ```javascript
   // main.js:149
   http.get('http://localhost:42423/static/', (res) => {
   ```

4. **Tests**
   - Update all test files referencing 8713

### 1.3 Replace All console.log Statements

**Priority**: 🔴 Critical  
**Effort**: 0.5 days

#### Quick fix approach:

1. **Create simple logger**
   ```javascript
   // renderer/src/utils/logger.js
   export const logger = {
       debug: (...args) => {
           if (import.meta.env.DEV) console.log(...args);
       },
       info: (...args) => console.info(...args),
       warn: (...args) => console.warn(...args),
       error: (...args) => console.error(...args)
   };
   ```

2. **Bulk replace**
   ```bash
   # Replace console.log with logger.debug
   find . -name "*.js" -o -name "*.jsx" | \
     grep -v node_modules | \
     xargs sed -i 's/console\.log(/logger.debug(/g'
   
   # Add import where needed
   ```

3. **Add ESLint rule**
   ```javascript
   // eslint.config.mjs
   rules: {
       'no-console': 'error'
   }
   ```

## Phase 2: Complete API Migration (2 days)

### 2.1 Remove All Legacy Code

**Priority**: 🟡 High  
**Effort**: 2 days

Since we have no users, we can simply delete legacy code:

1. **Frontend - Clean up apiService.jsx**
   - Remove all functions marked as "legacy"
   - Remove the comment about transition period
   - Keep only RESTful functions

2. **Backend - Remove old handlers**
   - Delete any handler not following REST patterns
   - Remove corresponding routes
   - Clean up unused service methods

3. **Remove legacy npm scripts**
   ```json
   // package.json - remove these lines entirely:
   "// Legacy scripts (use make commands instead)"
   "dev:react": "...",
   "dev:electron": "...",
   "dev:go": "...",
   "test:go": "...",
   "test:react": "..."
   ```

## Phase 3: Quick Modernization (1 day)

### 3.1 Fix React Components

**Priority**: 🟡 Medium  
**Effort**: 0.5 days

1. **Replace ErrorBoundary class component**
   ```bash
   cd renderer && npm install react-error-boundary
   ```

   ```jsx
   // renderer/src/components/common/ErrorBoundary.jsx
   import { ErrorBoundary as ReactErrorBoundary } from 'react-error-boundary';
   
   export default function ErrorBoundary({ children }) {
       return (
           <ReactErrorBoundary
               fallback={<div>Something went wrong</div>}
               onError={(error) => logger.error('Error boundary:', error)}
           >
               {children}
           </ReactErrorBoundary>
       );
   }
   ```

### 3.2 Quick Fixes

**Priority**: 🟢 Low  
**Effort**: 0.5 days

1. **Batch fixes**
   ```bash
   # Fix test file name
   mv renderer/src/pages/SkillPlans.text.jsx \
      renderer/src/pages/SkillPlans.test.jsx
   
   # Remove notes file
   rm notes
   
   # Remove embedded .env
   rm internal/embed/config/.env
   ```

2. **Standardize API naming**
   - Quick find/replace to use kebab-case consistently
   - Update both frontend and backend

## Phase 4: Dependency Updates (1 day)

### 4.1 Update Everything

**Priority**: 🟢 Low  
**Effort**: 1 day

Since no users, we can aggressively update:

```bash
# Update npm dependencies
npm update
npm audit fix --force

# Update Go dependencies
go get -u ./...
go mod tidy

# Update Electron
npm install electron@latest
```

## Simplified Schedule

| Day | Tasks | Outcome |
|-----|-------|---------|
| Day 1 | Remove .env, add first-run UI, fix ports | Secure credential handling, consistent ports |
| Day 2 | Remove all console.logs, add logger | Clean production logs |
| Day 3-4 | Delete all legacy API code | Single, clean API pattern |
| Day 5 | Modernize components, quick fixes | Modern React patterns |
| Day 6 | Update all dependencies | Current, secure dependencies |

## Total Effort: 6 developer days

## Success Checklist

- [ ] No embedded .env file
- [ ] First-run configuration UI works
- [ ] All ports are 42423
- [ ] Zero console.log statements
- [ ] No legacy API code remains
- [ ] No legacy npm scripts
- [ ] No class components
- [ ] All dependencies current
- [ ] Test file named correctly
- [ ] Notes file removed

## Benefits of No Users

1. **No backwards compatibility needed**
2. **Can make breaking changes freely**
3. **No migration documentation required**
4. **Can delete instead of deprecate**
5. **Aggressive dependency updates possible**

This simplified plan cuts the implementation time from 10-15 days to just 6 days by removing all user migration concerns and taking a more aggressive approach to cleaning up technical debt.