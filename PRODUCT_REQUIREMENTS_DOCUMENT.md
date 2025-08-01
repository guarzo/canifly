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

### Success Metrics
- User retention rate > 70%
- Daily active users
- Average session duration > 15 minutes
- Character sync operations per user
- User satisfaction score > 4/5

## User Personas

### 1. The Fleet Commander
- **Name**: Alex "FC" Rodriguez
- **Characteristics**: 
  - Manages 10+ accounts
  - Needs to track pilot capabilities for fleet composition
  - Values efficiency and quick access to information
- **Pain Points**: 
  - Difficulty tracking who can fly what ships
  - Time-consuming manual skill checks
- **Goals**: 
  - Quickly assess fleet capabilities
  - Plan training for doctrine ships

### 2. The Industry Mogul
- **Name**: Sarah Chen
- **Characteristics**: 
  - Focuses on industrial and trading activities
  - Manages specialized alts for different tasks
  - Needs precise skill planning
- **Pain Points**: 
  - Complex skill requirements for industrial operations
  - Managing settings across mining/hauling alts
- **Goals**: 
  - Optimize training for industrial efficiency
  - Maintain consistent UI settings across alts

### 3. The New Pilot
- **Name**: Marcus Johnson
- **Characteristics**: 
  - Recently started playing EVE
  - Overwhelmed by skill complexity
  - Wants clear progression path
- **Pain Points**: 
  - Unclear on what to train next
  - Doesn't understand ship requirements
- **Goals**: 
  - Clear skill training roadmap
  - Understand ship prerequisites

## Functional Requirements

### 1. Authentication & Authorization

#### 1.1 EVE SSO Integration
- **Description**: Secure login via EVE Online's Single Sign-On
- **Acceptance Criteria**:
  - Users can authenticate using EVE SSO
  - OAuth2 flow implementation
  - Token refresh handling
  - Session management
- **Priority**: P0 (Critical)

#### 1.2 Account Management
- **Description**: Support for multiple EVE accounts
- **Acceptance Criteria**:
  - Add/remove accounts
  - Account naming and organization
  - Account status tracking (Alpha/Omega)
- **Priority**: P0 (Critical)

### 2. Character Management

#### 2.1 Character Overview Dashboard
- **Description**: Central view of all characters
- **Acceptance Criteria**:
  - Display character portraits and basic info
  - Show skill points and training status
  - Group by account/role/location
  - Sort and filter capabilities
  - Hide/show accounts
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

#### 3.1 Skill Plan Creation
- **Description**: Create and manage skill training plans
- **Acceptance Criteria**:
  - Import skill plans from various formats
  - Create custom skill plans
  - Name and describe plans
  - Set skill order and priorities
- **Priority**: P0 (Critical)

#### 3.2 Skill Plan Analysis
- **Description**: Analyze character eligibility for plans
- **Acceptance Criteria**:
  - Show which characters meet requirements
  - Display missing skills
  - Calculate training time
  - Show partial completion status
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

#### 4.1 Profile Mapping
- **Description**: Associate game files with characters
- **Acceptance Criteria**:
  - Detect EVE profile directories
  - Map user files to character files
  - Visual file association interface
  - Color-coded by modification date
- **Priority**: P0 (Critical)

#### 4.2 Settings Sync
- **Description**: Synchronize settings between profiles
- **Acceptance Criteria**:
  - Sync individual profiles
  - Sync all profiles at once
  - Preserve user selections
  - Show sync status
- **Priority**: P0 (Critical)

#### 4.3 Backup & Restore
- **Description**: Backup game settings
- **Acceptance Criteria**:
  - Choose backup location
  - Create timestamped backups
  - Restore from backup
  - Backup history
- **Priority**: P1 (High)

### 5. Data Management

#### 5.1 Auto-refresh
- **Description**: Keep character data current
- **Acceptance Criteria**:
  - Manual refresh button
  - Auto-refresh on app launch
  - Background refresh option
  - Refresh status indicators
- **Priority**: P1 (High)

#### 5.2 Offline Mode
- **Description**: Function with cached data
- **Acceptance Criteria**:
  - Cache character data locally
  - Show last update timestamp
  - Indicate stale data
  - Work without internet
- **Priority**: P2 (Medium)

### 6. User Interface

#### 6.1 Dark Theme
- **Description**: EVE-appropriate dark interface
- **Acceptance Criteria**:
  - Consistent dark color scheme
  - Proper contrast ratios
  - Themed components
- **Priority**: P0 (Critical)

#### 6.2 Responsive Design
- **Description**: Adapt to different screen sizes
- **Acceptance Criteria**:
  - Desktop optimization
  - Tablet support
  - Minimum 1280x720 resolution
- **Priority**: P1 (High)

#### 6.3 Help & Instructions
- **Description**: In-app guidance
- **Acceptance Criteria**:
  - Dismissible instruction panels
  - Contextual help
  - Tooltips
  - First-time user guide
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

## Future Enhancements

### Phase 2 (Q2 2025)
- Fleet composition tools
- Doctrine compliance checking
- Skill plan sharing via links
- Mobile companion app
- Discord bot integration

### Phase 3 (Q3 2025)
- Corporation management features
- Skill plan marketplace
- Advanced analytics
- API for third-party tools
- Cloud sync option

### Phase 4 (Q4 2025)
- AI-powered skill recommendations
- Market integration
- Fitting management
- Asset tracking
- Notification system

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
- [ ] Core authentication working
- [ ] Character management functional
- [ ] Basic skill plan features
- [ ] Settings sync operational
- [ ] Windows/Mac builds available

### GA (v1.1)
- [ ] All P0 requirements complete
- [ ] 90% P1 requirements complete
- [ ] Performance targets met
- [ ] Security audit passed
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