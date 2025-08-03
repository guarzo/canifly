import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Button,
  Typography,
  Alert,
  Link,
  Box,
  CircularProgress
} from '@mui/material';
import { saveEVECredentials } from '../../api/apiService';

const FirstRunDialog = ({ open, onComplete }) => {
  const [clientId, setClientId] = useState('');
  const [clientSecret, setClientSecret] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async () => {
    if (!clientId || !clientSecret) {
      setError('Both Client ID and Client Secret are required');
      return;
    }

    setLoading(true);
    setError('');

    try {
      await saveEVECredentials(clientId, clientSecret);
      onComplete();
    } catch (err) {
      setError(err.message || 'Failed to save EVE credentials');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} maxWidth="sm" fullWidth disableEscapeKeyDown>
      <DialogTitle>EVE Online Application Setup</DialogTitle>
      <DialogContent>
        <Box sx={{ mt: 2 }}>
          <Typography variant="body1" paragraph>
            To use CanIFly, you need to register an EVE Online application to access character data.
          </Typography>
          <Typography variant="body2" paragraph>
            1. Go to{' '}
            <Link
              href="https://developers.eveonline.com"
              target="_blank"
              rel="noopener noreferrer"
            >
              EVE Online Developers
            </Link>{' '}
            and create a new application
          </Typography>
          <Typography variant="body2" paragraph>
            2. Set the callback URL to: <code>http://localhost:42423/callback</code>
          </Typography>
          <Typography variant="body2" paragraph>
            3. Select the required scopes for character and skill access
          </Typography>
          <Typography variant="body2" paragraph>
            4. Enter your Client ID and Client Secret below
          </Typography>

          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

          <TextField
            label="Client ID"
            value={clientId}
            onChange={(e) => setClientId(e.target.value)}
            fullWidth
            margin="normal"
            required
            disabled={loading}
          />
          <TextField
            label="Client Secret"
            value={clientSecret}
            onChange={(e) => setClientSecret(e.target.value)}
            type="password"
            fullWidth
            margin="normal"
            required
            disabled={loading}
          />
          <Typography variant="caption" color="textSecondary" sx={{ mt: 1, display: 'block' }}>
            Callback URL: http://localhost:42423/callback
          </Typography>
        </Box>
      </DialogContent>
      <DialogActions>
        <Button
          onClick={handleSubmit}
          variant="contained"
          disabled={loading || !clientId || !clientSecret}
        >
          {loading ? <CircularProgress size={24} /> : 'Save Configuration'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default FirstRunDialog;