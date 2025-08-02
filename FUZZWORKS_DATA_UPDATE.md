# Fuzzworks Data Update Implementation

## Overview

The application now downloads EVE Online static data from Fuzzworks on startup instead of relying solely on embedded CSV files. This ensures the application always has the latest item types and solar system information.

## Data Sources

- **Item Types**: Downloaded from `https://www.fuzzwork.co.uk/dump/latest/invTypes.csv.bz2`
- **Solar Systems**: Downloaded from `https://www.fuzzwork.co.uk/dump/latest/mapSolarSystems.csv.bz2`

## Implementation Details

### Fuzzworks Service (`internal/services/fuzzworks/service.go`)

The service handles:
- Downloading compressed CSV files from Fuzzworks
- Decompressing bzip2 data
- Validating downloaded data
- Caching and change detection
- Automatic fallback to embedded data if download fails

### Integration Points

1. **Startup**: The service initializes during application startup in `GetServices()`
2. **Skill Store**: Checks for downloaded `invTypes.csv` before falling back to embedded data
3. **System Store**: Checks for downloaded `mapSolarSystems.csv` before falling back to embedded data

### Data Storage

Downloaded files are stored in:
- `<basePath>/config/fuzzworks/invTypes.csv`
- `<basePath>/config/fuzzworks/mapSolarSystems.csv`
- `<basePath>/config/fuzzworks/fuzzworks_metadata.json` (tracks download times and checksums)

### Features

1. **Automatic Updates**: Data is checked and updated on each startup
2. **Daily Refresh**: Files older than 24 hours are re-downloaded
3. **Backup Creation**: Previous versions are backed up before updating
4. **Validation**: Downloaded data is validated for completeness
5. **Graceful Fallback**: If download fails, embedded data is used

### Configuration

The service can be configured with:
- `forceUpdate`: Forces re-download regardless of age
- Custom timeout values for downloads
- Retry logic with exponential backoff

## Benefits

1. **Always Current**: Users get the latest EVE data without app updates
2. **Reduced Binary Size**: Eventually, embedded files can be removed
3. **Reliability**: Fallback ensures the app works even offline
4. **Performance**: Data is cached locally after download

## Future Improvements

1. Add configuration option to disable auto-updates
2. Implement ETag support for efficient change detection
3. Add manual update trigger in UI
4. Consider downloading additional data files from Fuzzworks