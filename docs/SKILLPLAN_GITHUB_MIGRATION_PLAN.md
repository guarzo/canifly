# CanIFly Skill Plan Migration: Direct GitHub Implementation

## Current State Analysis

### Embedded File Structure
- **Location**: `/internal/embed/static/plans/`
- **Files**: 9 predefined skill plans (Bifrost.txt, Flycatcher.txt, Jaguar.txt, etc.)
- **Loading**: Files are embedded at compile time using Go's `embed` package
- **Runtime**: Copied to writable directory on first run, then managed locally

### Key Components
1. **SkillStore** (`internal/persist/eve/skill_store.go`):
   - Manages skill plan loading, saving, and deletion
   - Currently copies embedded plans to writable directory

2. **Handlers** (`internal/handlers/skillplan.go`):
   - REST API endpoints for CRUD operations
   - Existing caching layer for performance

3. **Services**: 
   - EVEDataService interface defines skill plan operations
   - Already handles cases where no skill plans exist

## GitHub Repository Structure

```
github.com/[YourOrg]/canifly-skillplans/
├── README.md
├── plans/
│   ├── Bifrost.txt
│   ├── Flycatcher.txt
│   ├── Jaguar.txt
│   ├── Keres.txt
│   ├── Kiki.txt
│   ├── Leshak.txt
│   ├── Magic_14.txt
│   ├── Manticore.txt
│   └── Naga.txt
└── metadata.json          # Optional: plan descriptions and metadata
```

## Skill Plan Format Enhancement

### Current Format
```
CPU Management 1
CPU Management 2
Power Grid Management 1
...
```

### Enhanced Format (with optional icon)
```
# icon: skill_book_mastery
CPU Management 1
CPU Management 2
Power Grid Management 1
...
```

Or for ship-specific plans:
```
# icon: ship_frigate_interceptor
Interdiction Nullification 1
Interceptors 1
...
```

### Icon Format Specification
- Icon field is optional - maintains backward compatibility
- Specified as a comment line at the top: `# icon: <icon_identifier>`
- Falls back to current filename-based logic if not present
- Icon identifiers could be:
  - EVE item type IDs (e.g., `11395` for Interdictors skill book)
  - Predefined categories (e.g., `skill_book_mastery`, `ship_frigate`)
  - Ship type names (e.g., `bifrost`, `flycatcher`)

## Implementation Plan

### Phase 1: Remove Embedded Files

#### 1.1 Clean up embed package
```go
// Remove from internal/embed/embed.go:
- //go:embed static/plans/*
- All plan-related embed logic

// Delete directory:
- internal/embed/static/plans/
```

#### 1.2 Update SkillStore
```go
// Remove from skill_store.go:
- copyEmbeddedPlansToWritable()
- copyEmbeddedFile()
- loadDeletedEmbeddedPlans()
- saveDeletedEmbeddedPlans()
```

### Phase 2: Update Skill Plan Parser

#### 2.1 Enhance readSkillsFromFile to parse icon
```go
// Update skill_store.go
type ParsedSkillPlan struct {
    Skills map[string]model.Skill
    Icon   string // Optional icon identifier
}

func (s *SkillStore) readSkillsFromFile(filePath string) (*ParsedSkillPlan, error) {
    data, err := s.fs.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read skill plan file %s: %w", filePath, err)
    }
    
    result := &ParsedSkillPlan{
        Skills: make(map[string]model.Skill),
    }
    
    scanner := bufio.NewScanner(strings.NewReader(string(data)))
    lineNumber := 0
    
    for scanner.Scan() {
        lineNumber++
        line := scanner.Text()
        trimmed := strings.TrimSpace(line)
        
        // Check for icon directive
        if strings.HasPrefix(trimmed, "# icon:") {
            result.Icon = strings.TrimSpace(strings.TrimPrefix(trimmed, "# icon:"))
            continue
        }
        
        // Skip empty lines and other comments
        if trimmed == "" || strings.HasPrefix(trimmed, "#") {
            continue
        }
        
        // Parse skill line as before
        parts := strings.Fields(line)
        if len(parts) < 2 {
            return nil, fmt.Errorf("invalid format in file %s at line %d: %s", filePath, lineNumber, line)
        }
        
        // ... rest of skill parsing logic
    }
    
    return result, nil
}
```

#### 2.2 Update SkillPlan model
```go
// internal/model/skill.go
type SkillPlan struct {
    Name   string
    Skills map[string]Skill
    Icon   string `json:"icon,omitempty"` // Optional icon identifier
}
```

### Phase 3: Add GitHub Downloader

#### 3.1 Create Simple Downloader Service
```go
// internal/services/github_plans.go
type GitHubPlansService struct {
    repoURL    string  // "https://raw.githubusercontent.com/[org]/[repo]/main"
    httpClient *http.Client
    logger     interfaces.Logger
}

func (g *GitHubPlansService) DownloadPlans(destDir string) error {
    // Simple implementation:
    // 1. Fetch metadata.json or hardcoded list
    // 2. Download each .txt file to destDir
    // 3. Return error only if critical failure
}

func (g *GitHubPlansService) DownloadPlan(planName, destPath string) error {
    url := fmt.Sprintf("%s/plans/%s.txt", g.repoURL, planName)
    // Download and save to destPath
}
```

### Phase 4: Integrate with Startup

#### 4.1 Modify SkillStore LoadSkillPlans
```go
func (s *SkillStore) LoadSkillPlans() error {
    s.logger.Infof("load skill plans")
    
    writableDir := filepath.Join(s.basePath, plansDir)
    if err := s.fs.MkdirAll(writableDir, os.ModePerm); err != nil {
        return fmt.Errorf("failed to ensure plans directory: %w", err)
    }
    
    // Try to download latest plans from GitHub
    if s.githubService != nil {
        if err := s.githubService.DownloadPlans(writableDir); err != nil {
            s.logger.Warnf("Failed to download plans from GitHub: %v", err)
            // Continue with whatever local files exist
        }
    }
    
    // Load whatever plans exist locally (downloaded or user-created)
    plans, err := s.loadSkillPlans(writableDir)
    if err != nil {
        return fmt.Errorf("failed to load skill plans: %w", err)
    }
    
    s.mut.Lock()
    s.skillPlans = plans
    s.mut.Unlock()
    
    s.logger.Debugf("Loaded %d skill plans", len(plans))
    return nil
}
```

### Phase 5: Simple Update Mechanism

#### 5.1 Add Refresh Endpoint
```go
// Add to skillplan.go handlers
func (h *SkillPlanHandler) RefreshPlans() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Re-download plans from GitHub
        if err := h.eveDataService.RefreshRemotePlans(); err != nil {
            respondError(w, "Failed to refresh plans", http.StatusInternalServerError)
            return
        }
        
        // Clear cache
        InvalidateCache(h.cache, "skillplans:")
        
        respondJSON(w, map[string]string{"status": "success"})
    }
}
```

#### 5.2 Add Route
```go
// POST /api/skill-plans/refresh
```

### Phase 6: Configuration

#### 6.1 Environment Variables
```bash
# Add to .env
SKILLPLANS_REPO_URL=https://raw.githubusercontent.com/[org]/[repo]/main
```

#### 6.2 Config Model
```go
// Add to existing config
type Config struct {
    // ... existing fields
    SkillPlansRepoURL string `json:"skillPlansRepoURL"`
}
```

## Implementation Details

### Error Handling
```go
// Simple approach - log errors but continue
func (g *GitHubPlansService) DownloadPlans(destDir string) error {
    planNames := []string{
        "Bifrost", "Flycatcher", "Jaguar", "Keres", 
        "Kiki", "Leshak", "Magic_14", "Manticore", "Naga",
    }
    
    for _, name := range planNames {
        destPath := filepath.Join(destDir, name+".txt")
        
        // Skip if file already exists and is recent (< 24 hours old)
        if info, err := os.Stat(destPath); err == nil {
            if time.Since(info.ModTime()) < 24*time.Hour {
                continue
            }
        }
        
        if err := g.DownloadPlan(name, destPath); err != nil {
            g.logger.Warnf("Failed to download %s: %v", name, err)
            // Continue with other plans
        }
    }
    return nil
}
```

### HTTP Client Configuration
```go
func NewGitHubPlansService(repoURL string, logger interfaces.Logger) *GitHubPlansService {
    return &GitHubPlansService{
        repoURL: repoURL,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        logger: logger,
    }
}
```

## API Changes

### New Endpoints
```
POST /api/skill-plans/refresh   # Force re-download from GitHub
```

### Existing Endpoints Remain Unchanged
- All existing CRUD operations work with local files
- No changes to response formats
- Cache behavior remains the same

## Benefits

1. **Zero Infrastructure Cost**: GitHub hosts the files
2. **Version Control**: Full history in git
3. **Community Contributions**: PRs for new plans
4. **Reduced Binary Size**: ~100KB smaller builds
5. **Simple Implementation**: Minimal code changes

## Considerations

### Network Failures
- App continues working with existing local files
- Failed downloads are logged but don't break functionality
- Manual refresh endpoint available when network returns

### GitHub Rate Limits
- Raw file access has generous limits (5000/hour)
- Each file is small (~1-2KB)
- 24-hour cache prevents frequent re-downloads

### First-Run Experience
- App starts normally even if download fails
- Users can create their own plans immediately
- Plans download in background when available

## Development Timeline

- **Day 1**: Remove embedded files and cleanup code
- **Day 2**: Update skill plan parser to support icon field
- **Day 3**: Implement GitHubPlansService
- **Day 4**: Integrate with SkillStore startup
- **Day 5**: Add refresh endpoint and testing
- **Day 6**: Update documentation and deploy

### Icon Field Examples

When creating skill plans in the GitHub repository, you can add icon identifiers:

**Magic_14.txt:**
```
# icon: skill_book_mastery
CPU Management 1
CPU Management 2
...
```

**Bifrost.txt:**
```
# icon: 37480
Command Destroyers 1
...
```

**Jaguar.txt:**
```
# icon: ship_frigate_assault
Assault Frigates 1
...
```

## Migration Steps

1. **Create GitHub Repository**
   - Upload current embedded plans
   - Add icon directives to each plan file
   - Set up public repository

2. **Update Application**
   - Remove embedded files
   - Update parser for icon support
   - Add GitHub downloader
   - Deploy new version

3. **User Impact**
   - Seamless transition
   - Existing local plans preserved
   - New plans auto-download on startup
   - Icons display when available, fallback to filename logic