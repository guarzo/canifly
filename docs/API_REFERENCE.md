# CanIFly API Reference

This document provides a comprehensive reference for all API endpoints in the CanIFly application.

## Base URL

```
http://localhost:42423
```

## Authentication

All API endpoints (except `/auth/*`) require a valid session cookie. The session is established through EVE SSO authentication.

## Endpoints

### Authentication

#### Login with EVE SSO
```
GET /auth/login

Redirects to EVE SSO for authentication
```

#### EVE SSO Callback
```
GET /auth/callback

Handles the OAuth2 callback from EVE SSO
Query Parameters:
- code: Authorization code from EVE
- state: State parameter for CSRF protection
```

#### Logout
```
POST /auth/logout

Terminates the current session
```

### Account Management

#### List Accounts
```
GET /api/accounts

Returns a paginated list of all accounts

Query Parameters:
- page (int): Page number (default: 1)
- limit (int): Items per page (default: 10)

Response:
{
  "data": [
    {
      "id": "account-id",
      "name": "Account Name",
      "active": true,
      "characters": [...],
      "order": 1
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "pages": 3
  }
}
```

#### Get Account
```
GET /api/accounts/{id}

Returns a specific account by ID

Response:
{
  "id": "account-id",
  "name": "Account Name",
  "active": true,
  "characters": [...],
  "order": 1
}
```

#### Update Account
```
PATCH /api/accounts/{id}

Updates account properties

Request Body:
{
  "name": "New Account Name",
  "active": false,
  "order": 2
}

Response:
{
  "id": "account-id",
  "name": "New Account Name",
  "active": false,
  "characters": [...],
  "order": 2
}
```

#### Delete Account
```
DELETE /api/accounts/{id}

Deletes an account and all associated characters

Response:
{
  "message": "Account deleted successfully"
}
```

### Configuration

#### Get Configuration
```
GET /api/config

Returns the current application configuration

Response:
{
  "eveConfigPath": "/path/to/eve/settings",
  "autoUpdate": true,
  "refreshInterval": 3600,
  "skillPlansRepoURL": "https://github.com/..."
}
```

#### Update Configuration
```
PATCH /api/config

Updates application configuration

Request Body:
{
  "eveConfigPath": "/new/path/to/eve/settings",
  "autoUpdate": false
}

Response:
{
  "eveConfigPath": "/new/path/to/eve/settings",
  "autoUpdate": false,
  "refreshInterval": 3600,
  "skillPlansRepoURL": "https://github.com/..."
}
```

### Skill Plans

#### List Skill Plans
```
GET /api/skill-plans

Returns a paginated list of skill plans

Query Parameters:
- page (int): Page number (default: 1)
- limit (int): Items per page (default: 10)

Response:
{
  "data": [
    {
      "name": "Magic 14",
      "skills": {
        "3449": {
          "name": "CPU Management",
          "level": 5
        }
      },
      "icon": "skill_book_mastery"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 9,
    "pages": 1
  }
}
```

#### Get Skill Plan
```
GET /api/skill-plans/{name}

Returns a specific skill plan by name

Response:
{
  "name": "Bifrost",
  "skills": {
    "11569": {
      "name": "Command Destroyers",
      "level": 1
    }
  },
  "icon": "37480"
}
```

#### Create Skill Plan
```
POST /api/skill-plans

Creates a new skill plan

Request Body:
{
  "name": "Custom Plan",
  "content": "# icon: custom_icon\nCPU Management 5\nPower Grid Management 5"
}

Response:
{
  "name": "Custom Plan",
  "skills": {...},
  "icon": "custom_icon"
}
```

#### Delete Skill Plan
```
DELETE /api/skill-plans/{name}

Deletes a skill plan

Response:
{
  "message": "Skill plan deleted successfully"
}
```

#### Refresh Skill Plans
```
POST /api/skill-plans/refresh

Forces a refresh of skill plans from GitHub

Response:
{
  "status": "success",
  "message": "Skill plans refreshed"
}
```

### Fuzzworks Integration

#### Get Fuzzworks Status
```
GET /api/fuzzworks/status

Returns the status of Fuzzworks data

Response:
{
  "invTypes": {
    "exists": true,
    "size": 15234567,
    "lastModified": "2024-01-15T10:30:00Z",
    "etag": "\"abc123\""
  },
  "mapSolarSystems": {
    "exists": true,
    "size": 1234567,
    "lastModified": "2024-01-15T10:30:00Z",
    "etag": "\"def456\""
  }
}
```

#### Update Fuzzworks Data
```
POST /api/fuzzworks/update

Triggers a manual update of Fuzzworks data

Response:
{
  "status": "success",
  "message": "Fuzzworks data updated",
  "details": {
    "invTypes": "updated",
    "mapSolarSystems": "unchanged"
  }
}
```

### ESI Proxy

#### Character Skills
```
GET /api/esi/characters/{characterID}/skills

Proxies to EVE ESI for character skills

Response:
{
  "skills": [
    {
      "skill_id": 3449,
      "skillpoints_in_skill": 256000,
      "trained_skill_level": 5,
      "active_skill_level": 5
    }
  ],
  "total_sp": 50000000,
  "unallocated_sp": 1000000
}
```

#### Character Skill Queue
```
GET /api/esi/characters/{characterID}/skillqueue

Proxies to EVE ESI for skill queue

Response:
[
  {
    "skill_id": 3450,
    "finished_level": 4,
    "start_date": "2024-01-15T10:00:00Z",
    "finish_date": "2024-01-16T15:30:00Z",
    "training_start_sp": 180000,
    "level_start_sp": 180000,
    "level_end_sp": 512000
  }
]
```

### Settings Synchronization

#### List EVE Profiles
```
GET /api/sync/profiles

Lists all EVE profile files

Response:
{
  "characters": [
    {
      "name": "character_12345678_Tranquility_Default.dat",
      "path": "/path/to/file",
      "modified": "2024-01-15T10:00:00Z",
      "size": 123456
    }
  ],
  "users": [
    {
      "name": "user_12345678_Tranquility_Default.dat",
      "path": "/path/to/file",
      "modified": "2024-01-15T10:00:00Z",
      "size": 234567
    }
  ]
}
```

#### Sync Settings
```
POST /api/sync/profile

Synchronizes settings between character and user files

Request Body:
{
  "characterFile": "character_12345678_Tranquility_Default.dat",
  "userFile": "user_12345678_Tranquility_Default.dat",
  "direction": "character_to_user"
}

Response:
{
  "status": "success",
  "filesProcessed": 5,
  "filesCopied": 3,
  "message": "Settings synchronized successfully"
}
```

#### Backup Settings
```
POST /api/sync/backup

Creates a backup of EVE settings

Request Body:
{
  "targetDir": "/path/to/eve/settings",
  "backupDir": "/path/to/backup"
}

Response:
{
  "status": "success",
  "filesBackedUp": 25,
  "backupPath": "/path/to/backup/backup_20240115_100000"
}
```

### WebSocket

#### Real-time Updates
```
GET /api/ws

WebSocket endpoint for real-time updates

Events:
- character_update: Character data has been updated
- skill_plan_update: Skill plans have changed
- config_update: Configuration has been modified
- sync_status: Sync operation status update

Example Message:
{
  "type": "character_update",
  "data": {
    "characterID": 12345678,
    "changes": ["skills", "location"]
  }
}
```

## Error Responses

All endpoints use consistent error responses:

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {}
}
```

Common HTTP status codes:
- 200: Success
- 400: Bad Request
- 401: Unauthorized
- 404: Not Found
- 500: Internal Server Error

## Rate Limiting

The API implements rate limiting to prevent abuse:
- 100 requests per minute per session
- WebSocket connections limited to 1 per session
- ESI proxy endpoints follow EVE's rate limits

## Caching

Several endpoints implement caching:
- Skill plans: 5 minutes
- Fuzzworks data: 24 hours
- Character data: 5 minutes
- Configuration: No cache

Cache can be bypassed with `Cache-Control: no-cache` header.