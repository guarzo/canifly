# Product Requirements Document (PRD) - CanIFly

## Document Information
- **Product Name**: CanIFly
- **Version**: 1.0
- **Date**: January 2025
- **Author**: Product Team
- **Status**: Current Implementation

## Executive Summary

CanIFly is a desktop application designed for EVE Online players to manage their characters, skill plans, and game settings efficiently. It provides a centralized interface for tracking character progression, determining ship eligibility, and synchronizing in-game settings across multiple accounts and profiles.

## Product Overview

### Vision
To be the essential companion tool for EVE Online players who manage multiple characters and accounts, providing seamless character management, skill planning, and settings synchronization.

### Mission
Simplify the complex task of managing multiple EVE Online characters by providing intuitive tools for skill tracking, ship eligibility verification, and game settings management.

### Target Users
- **Primary**: EVE Online players with multiple accounts and characters
- **Secondary**: New players planning their skill progression
- **Tertiary**: Corporation leaders managing member capabilities

## Business Context

### Problem Statement
EVE Online players face challenges in:
1. Tracking skill progress across multiple characters
2. Determining which ships characters can fly
3. Planning optimal skill training paths
4. Managing game settings across different accounts
5. Synchronizing configurations between game profiles

### Solution
CanIFly addresses these challenges by:
1. Integrating with EVE's ESI API for real-time character data
2. Providing visual skill plan management
3. Automating settings synchronization
4. Offering backup and restore capabilities
5. Creating a unified dashboard for all characters

### Key Features
- Real-time character data synchronization via EVE ESI API
- Comprehensive skill plan management and tracking
- Automated game settings synchronization across profiles
- Multi-account character organization and management
- Ship eligibility verification and skill gap analysis

## Core System Components

### Backend Architecture (Go)
- **RESTful API Server**: Handles all client requests on port 42423
- **Service Layer**: Business logic for accounts, characters, skill plans, and synchronization
- **Data Persistence**: File-based storage with JSON serialization
- **ESI Integration**: OAuth2-based EVE API client with automatic token refresh
- **WebSocket Server**: Real-time updates for character data and sync status
- **Session Management**: Secure session handling with middleware

### Frontend Architecture (React)
- **Component Library**: Reusable UI components with Material-UI and Tailwind CSS
- **State Management**: Zustand for global state, React hooks for local state
- **API Client**: Axios-based service layer for backend communication
- **Real-time Updates**: WebSocket integration for live data
- **Routing**: React Router for multi-page navigation
- **Animations**: Framer Motion for smooth UI transitions

### Desktop Integration (Electron)
- **Cross-platform Support**: Windows, macOS, and Linux distributions
- **File System Access**: Direct access to EVE game settings directories
- **Auto-updater**: Built-in update mechanism for new releases
- **Native Menus**: Platform-specific menu integration
- **Window Management**: Multi-window support with state persistence

## Functional Requirements

### 1. Authentication & Authorization

#### 1.1 EVE SSO Integration
- **Description**: Secure login via EVE Online's Single Sign-On
- **Implementation Details**:
  - OAuth2 authorization code flow with PKCE
  - Automatic token refresh before expiration
  - Secure token storage in local filesystem
  - Session persistence across application restarts
  - Multiple character authentication per account
  - Scope management for required ESI endpoints
- **ESI Scopes Required**:
  - publicData
  - esi-skills.read_skills.v1
  - esi-skills.read_skillqueue.v1
  - esi-clones.read_clones.v1
  - esi-location.read_location.v1
  - esi-ui.open_window.v1
  - esi-universe.read_structures.v1
- **Priority**: P0 (Critical)

#### 1.2 Account Management
- **Description**: Support for multiple EVE accounts
- **Features**:
  - Create new accounts with custom names
  - Add multiple characters to single account
  - Account visibility toggle (show/hide)
  - Account activation/deactivation
  - Account renaming capabilities
  - Account deletion with confirmation
  - MCT (Multiple Character Training) status display
  - Account omega/alpha status tracking
  - Account-level refresh operations
- **Data Persistence**:
  - JSON file storage in user data directory
  - Automatic backup before modifications
  - Migration support for data structure changes
- **Priority**: P0 (Critical)

### 2. Character Management

#### 2.1 Character Overview Dashboard
- **Description**: Central view of all characters with real-time data
- **Display Information**:
  - Character portrait (256x256) with refresh capability
  - Character name, corporation, and alliance
  - Total skill points and unallocated SP
  - Current location (station/structure name)
  - Active skill queue status
  - Account association and MCT status
  - Custom role assignment
  - Last update timestamp
- **Organization Features**:
  - Group by: Account, Role, or Location
  - Collapsible account groups
  - Drag-and-drop character reordering
  - Hide/show individual accounts
  - Search/filter by character name
  - Bulk selection for operations
- **Interactive Elements**:
  - Click character for detailed view
  - Right-click context menu
  - Refresh individual or all characters
  - Remove character with confirmation
  - Copy character info to clipboard
- **Priority**: P0 (Critical)

#### 2.2 Character Details
- **Description**: Detailed character information modal
- **Acceptance Criteria**:
  - Full skill list with levels
  - Current skill queue
  - Skill points distribution
  - Character attributes
  - Clone and implant information
- **Priority**: P1 (High)

#### 2.3 Role Assignment
- **Description**: Categorize characters by role
- **Acceptance Criteria**:
  - Assign custom roles to characters
  - Filter by role
  - Predefined role suggestions
- **Priority**: P2 (Medium)

### 3. Skill Plan Management

#### 3.1 Skill Plan Management System
- **Description**: Comprehensive skill plan creation and management
- **Pre-built Skill Plans**:
  - Magic 14 (Core fitting and capacitor skills)
  - Bifrost Command Destroyer specialization
  - Flycatcher Interdictor specialization
  - Jaguar Assault Frigate specialization
  - Keres Electronic Warfare specialization
  - Kikimora Destroyer specialization
  - Leshak Battleship specialization
  - Manticore Stealth Bomber specialization
  - Naga Attack Battlecruiser specialization
- **Custom Plan Features**:
  - Create new plans from scratch
  - Import from EVEMon XML format
  - Copy existing plans as templates
  - Edit plan name and description
  - Add/remove individual skills
  - Set target skill levels (I-V)
  - Reorder skill priorities
  - Plan validation against game rules
- **Plan Operations**:
  - Duplicate plans
  - Delete plans with confirmation
  - Export to shareable format
  - Plan comparison tool
- **Priority**: P0 (Critical)

#### 3.2 Skill Plan Analysis Engine
- **Description**: Real-time analysis of character skill plan compatibility
- **Analysis Views**:
  - **By Character View**: Shows all plans for selected character
    - Qualified plans (all skills met)
    - Partially qualified (some skills met)
    - Missing skills with levels needed
    - Estimated training time per plan
  - **By Plan View**: Shows all characters for selected plan
    - Character qualification status
    - Color-coded status indicators
    - Missing skill count per character
    - Quick qualification summary
- **Detailed Analysis Modal**:
  - Complete skill breakdown
  - Current level vs required level
  - Skill points needed
  - Training time with current attributes
  - Prerequisite skill chains
  - Optimal training order suggestion
- **Real-time Updates**:
  - Automatic recalculation on data refresh
  - WebSocket notifications for changes
  - Visual indicators for outdated data
- **Priority**: P0 (Critical)

#### 3.3 Skill Plan Operations
- **Description**: Manage existing skill plans
- **Acceptance Criteria**:
  - Copy skill plans
  - Delete skill plans
  - Export skill plans
  - Share skill plans
- **Priority**: P1 (High)

### 4. Settings Synchronization

#### 4.1 EVE Settings File Management
- **Description**: Intelligent mapping between character and user settings files
- **File Detection**:
  - Automatic discovery of EVE settings directory
  - Support for custom Tranquility directory paths
  - Detection of all character_*.dat files
  - Detection of all user_*.dat files
  - File modification time tracking
  - File size and hash comparison
- **Visual Mapping Interface**:
  - Drag-and-drop file association
  - Color coding by modification recency:
    - Green: Modified < 1 day ago
    - Yellow: Modified 1-7 days ago
    - Orange: Modified 7-30 days ago
    - Red: Modified > 30 days ago
  - Side-by-side file comparison
  - Association preview before saving
  - Bulk association tools
- **Association Management**:
  - Save character-to-user mappings
  - Clear individual associations
  - Reset all associations
  - Import/export mapping configurations
- **Priority**: P0 (Critical)

#### 4.2 Settings Synchronization Engine
- **Description**: Automated synchronization of EVE game settings
- **Sync Operations**:
  - Individual profile synchronization
  - Bulk sync all profiles
  - Selective file sync options
  - Bidirectional sync support
  - Conflict resolution options
- **Sync Features**:
  - Pre-sync backup creation
  - File integrity verification
  - Atomic operations (all-or-nothing)
  - Progress tracking with status bar
  - Detailed sync log generation
  - Rollback capability
- **File Handling**:
  - Preserve file permissions
  - Maintain file timestamps
  - Handle locked files gracefully
  - Skip corrupted files with warning
  - Compression for backup storage
- **User Feedback**:
  - Real-time sync progress
  - Success/failure notifications
  - Detailed error messages
  - Sync history tracking
- **Priority**: P0 (Critical)

#### 4.3 Backup & Restore System
- **Description**: Comprehensive backup solution for EVE settings
- **Backup Features**:
  - Custom backup location selection
  - Automatic timestamped naming
  - Incremental backup support
  - Compression options (ZIP format)
  - Selective file backup
  - Scheduled backup options
  - Pre-sync automatic backups
- **Backup Management**:
  - Backup history browser
  - Backup size information
  - Backup integrity verification
  - Old backup pruning
  - Backup metadata storage
- **Restore Capabilities**:
  - Full restore from any backup
  - Selective file restoration
  - Preview before restore
  - Restore to different location
  - Merge with existing settings
- **Safety Features**:
  - Pre-restore backup creation
  - Restore validation
  - Rollback on failure
  - Detailed restore logs
- **Priority**: P1 (High)

### 5. Data Management & Integration

#### 5.1 Real-time Data Synchronization
- **Description**: Comprehensive data management with ESI integration
- **Refresh Mechanisms**:
  - Manual refresh per character
  - Bulk refresh all characters
  - Auto-refresh on application start
  - Periodic background refresh (configurable)
  - Smart refresh based on data age
  - Refresh queue management
- **ESI Data Retrieved**:
  - Character public information
  - Skill list and levels
  - Skill queue details
  - Clone and implant data
  - Current location
  - Corporation/alliance info
  - Character attributes
- **WebSocket Features**:
  - Real-time update notifications
  - Multi-client synchronization
  - Connection status monitoring
  - Automatic reconnection
  - Event-based updates
- **Priority**: P1 (High)

#### 5.2 Fuzzworks Market Data Integration
- **Description**: Integration with Fuzzworks for EVE market data
- **Features**:
  - Automatic invTypes.csv updates
  - Solar system data synchronization
  - ETag-based update checking
  - Manual update triggers
  - Update status tracking
  - File size monitoring
  - Last update timestamps
  - Error retry logic
- **Data Usage**:
  - Item type information
  - System name resolution
  - Market data enrichment
- **Priority**: P2 (Medium)

### 6. User Interface & Experience

#### 6.1 Modern Dark Theme Design
- **Description**: EVE Online-inspired dark interface with modern aesthetics
- **Design System**:
  - Dark background (#0a0a0a) with teal accents (#14b8a6)
  - Glass-morphism effects for panels
  - Consistent spacing using 8px grid
  - High contrast for readability
  - Custom scrollbar styling
  - Smooth transitions (200ms default)
- **Component Library**:
  - Material-UI base components
  - Custom styled variants
  - Tailwind utility classes
  - Framer Motion animations
  - Loading skeletons
  - Toast notifications
- **Priority**: P0 (Critical)

#### 6.2 Interactive Elements
- **Description**: Rich interactive features for enhanced UX
- **Interactions**:
  - Hover effects on all clickable elements
  - Smooth accordion animations
  - Drag-and-drop file associations
  - Context menus on right-click
  - Keyboard shortcuts support
  - Multi-select with checkboxes
  - Progressive disclosure patterns
- **Feedback Systems**:
  - Loading spinners and progress bars
  - Success/error toast messages
  - Confirmation dialogs for destructive actions
  - Inline validation messages
  - Status badges and indicators
  - Tooltips on hover
- **Priority**: P1 (High)

#### 6.3 Navigation & Information Architecture
- **Description**: Intuitive navigation and content organization
- **Navigation Features**:
  - Persistent header with breadcrumbs
  - Tab-based navigation for views
  - Floating action buttons
  - Contextual action bars
  - Search and filter controls
  - Pagination for large datasets
- **Help System**:
  - Collapsible instruction panels
  - Step-by-step setup guides
  - Contextual help tooltips
  - Keyboard shortcut reference
  - External documentation links
  - Video tutorial integration
- **Priority**: P2 (Medium)

## Non-Functional Requirements

### Performance
- Application launch time < 3 seconds
- Character data load < 2 seconds
- Sync operations < 5 seconds
- Memory usage < 500MB
- Smooth 60fps UI animations

### Security
- Encrypted token storage
- Secure API communication (HTTPS)
- No plaintext password storage
- Session timeout after inactivity
- Secure file operations

### Reliability
- 99.9% uptime (excluding EVE API downtime)
- Graceful error handling
- Data corruption prevention
- Automatic recovery from failures

### Usability
- Intuitive navigation
- Consistent UI patterns
- Keyboard shortcuts
- Accessibility support
- Minimal learning curve

### Compatibility
- Windows 10/11 support
- macOS 10.15+ support
- Linux (Ubuntu 20.04+) support
- EVE Online client compatibility

## Technical Architecture

### Technology Stack
- **Frontend**: React 18, Material-UI, Tailwind CSS
- **Backend**: Go 1.23, Gin framework
- **Desktop**: Electron 32
- **State Management**: Zustand
- **Build Tools**: Vite, npm/yarn
- **Testing**: Vitest, Go testing

### API Integration
- EVE ESI v1 endpoints
- OAuth2 authentication flow
- Rate limiting compliance
- Error retry logic

### Data Storage
- Local JSON storage
- Encrypted credentials
- Settings persistence
- Cache management

## User Flows

### 1. First-Time Setup
```
1. User launches application
2. Clicks "Login with EVE SSO"
3. Enters account name
4. Redirected to EVE login
5. Authorizes application
6. Returns to app authenticated
7. Views character dashboard
```

### 2. Adding a Character
```
1. User clicks "Add Character" (+) button
2. Enters/selects account name
3. Redirected to EVE SSO
4. Selects character to add
5. Character appears in dashboard
6. User assigns role (optional)
```

### 3. Checking Skill Plans
```
1. User navigates to Skill Plans
2. Views character/plan matrix
3. Clicks on character row
4. Sees missing skills modal
5. Reviews training time
6. Exports or copies plan
```

### 4. Syncing Settings
```
1. User navigates to Sync page
2. Selects profile to sync
3. Chooses character and user files
4. Clicks sync button
5. Confirms operation
6. Views success message
```

## Additional Features & Capabilities

### Session Management
- Persistent login across application restarts
- Secure session storage with encryption
- Automatic session refresh
- Multi-session support for different characters
- Session timeout configuration
- Force logout capabilities

### Error Handling & Recovery
- Comprehensive error catching and logging
- User-friendly error messages
- Automatic retry for failed API calls
- Graceful degradation on service failures
- Detailed error logs for debugging
- Crash recovery with data preservation

### Performance Optimizations
- Lazy loading of character data
- Request caching with TTL
- Debounced search inputs
- Virtual scrolling for large lists
- Image optimization and caching
- Background task queuing

### Developer Tools
- Structured logging with logrus
- Debug mode with verbose output
- API request/response logging
- Performance profiling endpoints
- Health check endpoints
- Metrics collection

## Dependencies

### External Services
- EVE Online ESI API
- EVE SSO authentication
- Operating system file APIs

### Third-Party Libraries
- Electron framework
- React ecosystem
- Go standard library
- OAuth2 libraries

## Risks & Mitigation

### Technical Risks
| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| ESI API changes | High | Medium | Version detection, graceful degradation |
| Electron security | High | Low | Regular updates, security audits |
| Performance issues | Medium | Medium | Profiling, optimization, lazy loading |

### Business Risks
| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| Low adoption | High | Medium | Marketing, community engagement |
| Competition | Medium | High | Unique features, better UX |
| EVE player decline | High | Low | Diversification options |

## Release Criteria

### MVP (v1.0)
- [x] Core authentication working
- [x] Character management functional
- [x] Basic skill plan features
- [x] Settings sync operational
- [x] Windows/Mac builds available

### GA (v1.1)
- [x] All P0 requirements complete
- [x] 90% P1 requirements complete
- [x] Performance targets met
- [ ] User documentation complete

## Appendices

### A. Glossary
- **ESI**: EVE Swagger Interface - EVE Online's API
- **SSO**: Single Sign-On
- **SP**: Skill Points
- **MCT**: Multiple Character Training
- **Alpha/Omega**: Free/Subscribed account status

### B. References
- [EVE Online Developers](https://developers.eveonline.com)
- [ESI Documentation](https://esi.evetech.net)
- [Electron Documentation](https://electronjs.org)
- [React Documentation](https://react.dev)

### C. Mockups
- See `/design` directory for UI mockups
- Figma link: [Design System](#)
- User flow diagrams: [Miro Board](#)

---

**Document Version History**
- v1.0 - Initial PRD (January 2025)
- v0.9 - Draft review (December 2024)
- v0.1 - Initial outline (November 2024)