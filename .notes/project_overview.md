# CanIFly Project Overview

## Introduction

**CanIFly** is a desktop application designed to help EVE Online players manage their characters, skill plans, and game settings. It integrates with EVE Online's ESI (EVE Swagger Interface) API to provide real-time character data while also offering local profile management capabilities. The application helps players determine which ships they can pilot, what skills they need to train, and simplifies the process of managing settings across multiple characters and accounts.

## Architecture

CanIFly is built as an Electron application with a Go backend and React frontend:

1. **Electron Shell**: Provides the desktop application wrapper
2. **Go Backend**: Handles API integration, data persistence, and business logic
3. **React Frontend**: Delivers a responsive and intuitive user interface

### System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     Electron Application                     │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐            ┌───────────────────────┐   │
│  │                 │            │                       │   │
│  │   Go Backend    │◄─────────►│    React Frontend     │   │
│  │                 │   REST    │                       │   │
│  └────────┬────────┘            └───────────────────────┘   │
│           │                                                 │
└───────────┼─────────────────────────────────────────────────┘
            │
            ▼
┌─────────────────────┐    ┌─────────────────────┐
│                     │    │                     │
│  EVE Online ESI    │    │  Local EVE Profiles │
│                     │    │                     │
└─────────────────────┘    └─────────────────────┘
```

## Key Components

### Backend (Go)

The Go backend is organized into several key packages:

1. **`cmd`**: Contains application entry points and initialization logic
2. **`handlers`**: REST API handlers that process HTTP requests from the frontend
3. **`http`**: HTTP client implementations, middleware, and session management
4. **`model`**: Data models and structures used throughout the application
5. **`persist`**: Data persistence layer for storing application state, user preferences, and cached data
6. **`services`**: Business logic implementation, including:
   - Account management
   - Character data processing
   - EVE Online API integration
   - Skill plan management
   - Profile synchronization
7. **`server`**: Server configuration, routing, and service initialization
8. **`errors`**: Custom error types and error handling utilities

### Frontend (React)

The React frontend is organized into a component-based architecture:

1. **`api`**: API client for communicating with the Go backend
2. **`components`**: Reusable UI components organized by feature:
   - Common components (buttons, modals, etc.)
   - Dashboard components
   - Landing page components
   - Mapping interface components
   - Skill plan components
   - Sync interface components
3. **`hooks`**: Custom React hooks for shared functionality
4. **`pages`**: Main application pages:
   - Character Overview
   - Landing
   - Mapping
   - Skill Plans
   - Sync
5. **`utils`**: Utility functions for data transformation, formatting, and logging

## Core Features

### 1. Character Management

The application allows users to:
- Add characters via EVE Online OAuth authentication
- View character details including skills, attributes, and training queue
- Group characters by account, role, or location
- Sort and filter characters for easy management

### 2. Skill Plans

Users can:
- Create, view, and delete skill plans
- See which characters qualify for specific plans
- Identify missing skills and training time required
- Copy skill plans for sharing or backup

### 3. Profile Mapping & Synchronization

The application provides tools to:
- Associate character files with local user files
- Sync in-game settings across multiple profiles
- Back up and restore EVE Online settings
- Automatically detect and suggest appropriate mappings

## Data Flow

1. **Authentication Flow**:
   - User initiates login through the frontend
   - Backend generates OAuth URL for EVE Online
   - User authenticates with EVE Online
   - EVE Online redirects back with authorization code
   - Backend exchanges code for access token
   - Character data is retrieved and stored

2. **Character Data Flow**:
   - Backend periodically refreshes character data from EVE Online API
   - Data is processed, normalized, and stored in the persistence layer
   - Frontend requests data through REST API endpoints
   - UI components render the data in a user-friendly format

3. **Settings Synchronization Flow**:
   - User selects character and user profiles to sync
   - Backend identifies relevant files based on mappings
   - Files are copied between directories while preserving structure
   - Success/failure status is reported back to the frontend

## Technical Details

### Authentication & Security

- OAuth 2.0 integration with EVE Online SSO
- Secure storage of access tokens
- Encryption of sensitive data using a secret key
- Session management for API access

### Data Persistence

- Local file-based storage for application state
- Caching of EVE Online API responses to reduce API calls
- User preferences and settings stored locally

### Performance Considerations

- Asynchronous API calls to prevent UI blocking
- Efficient data normalization for quick rendering
- Pagination and virtualization for handling large datasets

## Development Workflow

1. **Local Development**:
   - Run the application in development mode with hot reloading
   - Backend and frontend can be developed independently

2. **Building**:
   - Frontend is built with Vite
   - Go backend is compiled for the target platform
   - Electron packages everything into a distributable application

3. **Testing**:
   - Unit tests for both backend and frontend components
   - Integration tests for API endpoints
   - End-to-end testing for critical user flows

## Future Enhancements

Based on the TODO section in the README:
- Support for clone management
- Implant tracking and optimization
- Corporation-based sorting and filtering

## Conclusion

CanIFly provides EVE Online players with a comprehensive tool for managing characters, skills, and settings. By combining real-time data from the EVE Online API with local profile management, it offers a unique solution that streamlines many aspects of the EVE Online experience. 