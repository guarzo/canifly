import React, { useState, useEffect } from 'react';
import { 
  Card, 
  CardContent, 
  Typography, 
  Switch, 
  FormControlLabel, 
  Button, 
  Box,
  CircularProgress,
  Alert,
  Chip
} from '@mui/material';
import { CloudDownload, Refresh } from '@mui/icons-material';
import { useAppData } from '../../hooks/useAppData';
import { useAsyncOperation } from '../../hooks/useAsyncOperation';
import apiService from '../../api/apiService';

const FuzzworksSettings = () => {
  const { config, updateConfig } = useAppData();
  const { execute, loading } = useAsyncOperation();
  const [status, setStatus] = useState(null);
  const [loadingStatus, setLoadingStatus] = useState(false);

  // Default auto-update to true if not set
  const autoUpdate = config?.AutoUpdateFuzzworks !== false;

  useEffect(() => {
    fetchStatus();
  }, []);

  const fetchStatus = async () => {
    setLoadingStatus(true);
    try {
      const response = await apiService.get('/api/fuzzworks/status');
      setStatus(response.data);
    } catch (error) {
      console.error('Failed to fetch Fuzzworks status:', error);
    } finally {
      setLoadingStatus(false);
    }
  };

  const handleAutoUpdateToggle = async () => {
    await execute(
      () => updateConfig({ AutoUpdateFuzzworks: !autoUpdate }),
      {
        successMessage: 'Auto-update setting saved',
        errorMessage: 'Failed to update setting'
      }
    );
  };

  const handleManualUpdate = async () => {
    await execute(
      async () => {
        const response = await apiService.post('/api/fuzzworks/update');
        await fetchStatus(); // Refresh status after update
        return response.data;
      },
      {
        successMessage: 'Fuzzworks data updated successfully',
        errorMessage: 'Failed to update Fuzzworks data'
      }
    );
  };

  const formatDate = (dateString) => {
    if (!dateString) return 'Never';
    const date = new Date(dateString);
    return date.toLocaleString();
  };

  const formatFileSize = (bytes) => {
    if (!bytes) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  return (
    <Card>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Fuzzworks Data Updates
        </Typography>
        
        <Typography variant="body2" color="textSecondary" paragraph>
          Fuzzworks provides up-to-date EVE Online static data including item types and solar systems.
        </Typography>

        <Box sx={{ mb: 3 }}>
          <FormControlLabel
            control={
              <Switch
                checked={autoUpdate}
                onChange={handleAutoUpdateToggle}
                disabled={loading}
              />
            }
            label="Auto-update data on startup"
          />
        </Box>

        {loadingStatus ? (
          <Box display="flex" justifyContent="center" p={2}>
            <CircularProgress />
          </Box>
        ) : status && (
          <Box sx={{ mb: 3 }}>
            {status.hasData ? (
              <>
                {status.invTypes && (
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="subtitle2" gutterBottom>
                      Item Types (invTypes.csv)
                    </Typography>
                    <Box sx={{ pl: 2 }}>
                      <Typography variant="body2">
                        Last Updated: {formatDate(status.invTypes.lastUpdated)}
                      </Typography>
                      <Typography variant="body2">
                        File Size: {formatFileSize(status.invTypes.fileSize)}
                      </Typography>
                      {status.invTypes.etag && (
                        <Typography variant="caption" color="textSecondary">
                          ETag: {status.invTypes.etag}
                        </Typography>
                      )}
                    </Box>
                  </Box>
                )}
                
                {status.solarSystems && (
                  <Box sx={{ mb: 2 }}>
                    <Typography variant="subtitle2" gutterBottom>
                      Solar Systems (mapSolarSystems.csv)
                    </Typography>
                    <Box sx={{ pl: 2 }}>
                      <Typography variant="body2">
                        Last Updated: {formatDate(status.solarSystems.lastUpdated)}
                      </Typography>
                      <Typography variant="body2">
                        File Size: {formatFileSize(status.solarSystems.fileSize)}
                      </Typography>
                      {status.solarSystems.etag && (
                        <Typography variant="caption" color="textSecondary">
                          ETag: {status.solarSystems.etag}
                        </Typography>
                      )}
                    </Box>
                  </Box>
                )}
              </>
            ) : (
              <Alert severity="info">
                No Fuzzworks data downloaded yet. Click "Update Now" to download the latest data.
              </Alert>
            )}
          </Box>
        )}

        <Button
          variant="contained"
          color="primary"
          onClick={handleManualUpdate}
          disabled={loading}
          startIcon={loading ? <CircularProgress size={20} /> : <CloudDownload />}
          fullWidth
        >
          {loading ? 'Updating...' : 'Update Now'}
        </Button>

        <Box sx={{ mt: 2 }}>
          <Typography variant="caption" color="textSecondary">
            Updates check for new data from fuzzwork.co.uk using ETag headers for efficiency.
            Only changed files are downloaded.
          </Typography>
        </Box>
      </CardContent>
    </Card>
  );
};

export default FuzzworksSettings;