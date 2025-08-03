# CanIFly Architecture

## Overview

CanIFly is an Electron desktop application for EVE Online players that helps manage characters, check flyable ships, manage skill plans, and sync in-game settings. The architecture follows a clean, service-oriented approach with clear separation of concerns.

## Technology Stack

- **Backend**: Go 1.23.2 with RESTful APIs
- **Frontend**: React 18 with Material-UI and Tailwind CSS
- **Desktop**: Electron 32.2.0
- **State Management**: Zustand
- **Real-time Updates**: WebSockets
- **Data Storage**: File-based persistence with atomic writes

## Architecture Principles

1. **Service-Oriented Architecture**: Each domain (accounts, configuration, EVE data) is managed by dedicated services
2. **RESTful APIs**: Clean, resource-based endpoints following REST conventions
3. **Separation of Concerns**: Clear boundaries between data access, business logic, and presentation
4. **Real-time Updates**: WebSocket integration for live data synchronization
5. **Efficient Caching**: HTTP caching with ETag support for optimal performance

## Backend Architecture

### Service Layer (`internal/services/`)

The backend is organized around domain-specific services:

- **AccountManagementService**: Manages EVE accounts and characters
- **ConfigurationService**: Handles application configuration and user preferences
- **EVEDataService**: Manages EVE Online static data (skills, ships, systems)
- **FuzzworksService**: Downloads and updates EVE static data from Fuzzworks
- **SyncService**: Synchronizes in-game settings
- **HTTPCacheService**: Provides server-side caching with ETag support
- **StorageService**: Unified file-based storage with atomic operations

### API Layer (`internal/handlers/`)

RESTful handlers provide clean HTTP endpoints:

```
GET    /api/accounts              # List all accounts (paginated)
GET    /api/accounts/{id}          # Get specific account
PATCH  /api/accounts/{id}          # Update account
DELETE /api/accounts/{id}          # Delete account

GET    /api/config                 # Get configuration
PATCH  /api/config                 # Update configuration

GET    /api/skill-plans            # List skill plans (paginated)
POST   /api/skill-plans            # Create skill plan
GET    /api/skill-plans/{name}     # Get specific skill plan
PUT    /api/skill-plans/{name}     # Update skill plan
DELETE /api/skill-plans/{name}     # Delete skill plan

GET    /api/fuzzworks/status       # Get Fuzzworks data status
POST   /api/fuzzworks/update       # Trigger manual update

GET    /api/ws                     # WebSocket endpoint for real-time updates
```

### Data Persistence (`internal/persist/`)

File-based storage with:
- Atomic writes to prevent data corruption
- File locking for exclusive access
- Automatic backup creation
- JSON format for human readability

## Frontend Architecture

### State Management

Zustand stores provide centralized state management:

- **authStore**: Authentication state and session management
- **appDataStore**: Account and configuration data
- **eveDataStore**: EVE static data (skill plans, profiles, conversions)

### Custom Hooks

React hooks encapsulate business logic:

- **useAuth**: Authentication operations
- **useAppData**: Account and config data access
- **useEveData**: EVE static data access
- **useWebSocket**: Real-time updates via WebSocket
- **useAsyncOperation**: Consistent async operation handling

### Component Structure

```
renderer/src/
├── components/           # Reusable UI components
│   ├── common/          # Shared components
│   ├── account/         # Account management
│   ├── character/       # Character views
│   ├── settings/        # Settings components
│   └── skillplan/       # Skill plan management
├── pages/               # Route-based page components
├── hooks/               # Custom React hooks
├── stores/              # Zustand state stores
├── api/                 # API client services
└── utils/               # Utility functions
```

## Key Features

### 1. Real-time Updates
- WebSocket connection for live data synchronization
- Automatic reconnection with exponential backoff
- Event-based updates (account changes, skill plan updates)

### 2. Efficient Data Loading
- Server-side HTTP caching with ETag support
- Client-side pagination for large datasets
- Lazy loading of data components

### 3. Fuzzworks Integration
- Automatic updates of EVE static data on startup
- ETag-based change detection to minimize downloads
- Manual update trigger via UI
- Configurable auto-update behavior
- No embedded CSV files - all data downloaded from Fuzzworks
- Startup validation ensures data availability
- Binary size reduced by 10MB+ through removal of embedded data

### 4. GitHub Skill Plans Integration
- Skill plans downloaded from GitHub repository
- 24-hour cache to prevent frequent re-downloads
- Manual refresh endpoint available
- Graceful fallback to local files if download fails
- Support for custom icon identifiers in skill plan format
- Community-contributed skill plans via pull requests

### 5. Security
- Session-based authentication
- Secure token storage
- CORS protection
- Input validation and sanitization

## Data Flow

1. **User Action** → React Component → Custom Hook → API Service
2. **API Service** → HTTP Request → Go Handler → Service Layer
3. **Service Layer** → Storage/External API → Response
4. **Response** → Handler → HTTP Response → API Service
5. **API Service** → State Update → React Re-render
6. **WebSocket Events** → Broadcast → Connected Clients → State Update

## Performance Optimizations

1. **HTTP Caching**: ETag-based caching reduces redundant data transfers
2. **Pagination**: Large datasets are paginated for better performance
3. **WebSocket**: Real-time updates eliminate the need for polling
4. **Atomic Writes**: File operations are atomic to prevent corruption
5. **Lazy Loading**: Data is loaded on-demand to improve startup time
6. **Binary Size**: Reduced by 10MB+ through removal of embedded CSV files
7. **UI Enhancements**: 
   - Particle effects for visual appeal
   - Smooth page transitions
   - Skeleton loaders for better perceived performance
   - Enhanced character card animations

## Future Considerations

1. **Database Migration**: Consider migrating from file-based to database storage for better scalability
2. **GraphQL**: Evaluate GraphQL for more flexible data fetching
3. **Service Mesh**: Consider microservices architecture for larger scale
4. **Kubernetes**: Container orchestration for cloud deployment
5. **Monitoring**: Add comprehensive logging and metrics collection