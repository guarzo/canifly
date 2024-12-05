import React, { useState, useEffect } from 'react';
import { toast, ToastContainer } from 'react-toastify';
import Swal from 'sweetalert2';
import withReactContent from 'sweetalert2-react-content';
import 'react-toastify/dist/ReactToastify.css';
import 'sweetalert2/dist/sweetalert2.min.css';
import {
    Button,
    Select,
    MenuItem,
    FormControl,
    InputLabel,
    CircularProgress,
    Grid,
    Typography,
    Card,
    CardContent,
    Box,
    Divider,
    Tooltip,
    Snackbar,
    Alert,
    Collapse,
} from '@mui/material';
import LoadingButton from '@mui/lab/LoadingButton';
import BackupIcon from '@mui/icons-material/Backup';
import DeleteIcon from '@mui/icons-material/Delete';
import SyncIcon from '@mui/icons-material/Sync';
import SyncAllIcon from '@mui/icons-material/SyncAlt';

const MySwal = withReactContent(Swal);

const Sync = () => {
    const [settingsData, setSettingsData] = useState([]);
    const [associations, setAssociations] = useState([]);
    const [currentSettingsDir, setCurrentSettingsDir] = useState('');
    const [isLoading, setIsLoading] = useState(false);
    const [selections, setSelections] = useState({});
    const [successMessage, setSuccessMessage] = useState('');

    useEffect(() => {
        loadDataAndRender();
    }, []);

    const loadDataAndRender = async () => {
        try {
            setIsLoading(true);
            const data = await window.electronAPI.loadSettings();
            if (data && data.settingsData) {
                setSettingsData(data.settingsData);
                setAssociations(data.associations);
                setCurrentSettingsDir(data.currentSettingsDir);

                const initialSelections = {};
                data.settingsData.forEach(subDir => {
                    initialSelections[subDir.subDir] = { charId: '', userId: '' };
                });
                setSelections(initialSelections);
            } else {
                toast.error('Failed to load settings: No data received.');
            }
        } catch (error) {
            console.error('Error loading settings:', error);
            toast.error('Failed to load settings: ' + error.message);
        } finally {
            setIsLoading(false);
        }
    };

    const handleSelectionChange = (subDir, field, value) => {
        setSelections(prev => ({
            ...prev,
            [subDir]: { ...prev[subDir], [field]: value },
        }));
    };

    const handleSync = async (subDir) => {
        setIsLoading(true);
        setTimeout(() => {
            setIsLoading(false);
            setSuccessMessage(`Successfully synced ${subDir}!`);
        }, 2000); // Simulating async operation
    };

    const handleSyncAll = async (subDir) => {
        setIsLoading(true);
        setTimeout(() => {
            setIsLoading(false);
            setSuccessMessage(`Successfully synced all profiles for ${subDir}!`);
        }, 2000); // Simulating async operation
    };

    return (
        <div className="pt-16">
            <ToastContainer />
            <Box sx={{ padding: '16px' }}>
                {/* Snackbar for Success Messages */}
                <Snackbar
                    open={!!successMessage}
                    autoHideDuration={3000}
                    onClose={() => setSuccessMessage('')}
                    anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
                >
                    <Alert onClose={() => setSuccessMessage('')} severity="success" variant="filled">
                        {successMessage}
                    </Alert>
                </Snackbar>

                {/* Subdirectories */}
                <Grid container spacing={3}>
                    {settingsData.map(subDir => (
                        <Grid item xs={12} sm={6} md={4} key={subDir.subDir}>
                            <Collapse in={!isLoading}>
                                <Card
                                    sx={{
                                        boxShadow: 3,
                                        transition: 'transform 0.2s, box-shadow 0.2s',
                                        '&:hover': { transform: 'scale(1.02)', boxShadow: 6 },
                                    }}
                                >
                                    <CardContent>
                                        {/* Profile Header */}
                                        <Typography
                                            variant="h6"
                                            gutterBottom
                                            sx={{ textAlign: 'center', fontWeight: 600 }}
                                        >
                                            {subDir.subDir.replace('settings_', '')}
                                        </Typography>
                                        <Divider sx={{ marginY: 2 }} />

                                        {/* Input Group */}
                                        <Box
                                            sx={{
                                                padding: 2,
                                                border: '1px solid rgba(0, 0, 0, 0.12)',
                                                borderRadius: 2,
                                                bgcolor: 'rgba(0, 0, 0, 0.03)',
                                                marginBottom: 2,
                                            }}
                                        >
                                            <FormControl fullWidth margin="normal">
                                                <InputLabel id={`char-select-label-${subDir.subDir}`}>
                                                    Select Character
                                                </InputLabel>
                                                <Select
                                                    labelId={`char-select-label-${subDir.subDir}`}
                                                    value={selections[subDir.subDir]?.charId || ''}
                                                    onChange={(e) =>
                                                        handleSelectionChange(
                                                            subDir.subDir,
                                                            'charId',
                                                            e.target.value
                                                        )
                                                    }
                                                >
                                                    <MenuItem value="">
                                                        <em>-- Select Character --</em>
                                                    </MenuItem>
                                                    {subDir.availableCharFiles.map(char => (
                                                        <MenuItem key={char.charId} value={char.charId}>
                                                            {char.name}
                                                        </MenuItem>
                                                    ))}
                                                </Select>
                                            </FormControl>
                                            <FormControl fullWidth margin="normal">
                                                <InputLabel id={`user-select-label-${subDir.subDir}`}>
                                                    Select User
                                                </InputLabel>
                                                <Select
                                                    labelId={`user-select-label-${subDir.subDir}`}
                                                    value={selections[subDir.subDir]?.userId || ''}
                                                    onChange={(e) =>
                                                        handleSelectionChange(
                                                            subDir.subDir,
                                                            'userId',
                                                            e.target.value
                                                        )
                                                    }
                                                >
                                                    <MenuItem value="">
                                                        <em>-- Select User --</em>
                                                    </MenuItem>
                                                    {subDir.availableUserFiles.map(user => (
                                                        <MenuItem key={user.userId} value={user.userId}>
                                                            {user.userId}
                                                        </MenuItem>
                                                    ))}
                                                </Select>
                                            </FormControl>
                                        </Box>

                                        {/* Buttons */}
                                        <Box display="flex" justifyContent="space-between" gap={2}>
                                            <Tooltip title="Sync">
                                                <LoadingButton
                                                    variant="contained"
                                                    color="primary"
                                                    loading={isLoading}
                                                    onClick={() => handleSync(subDir.subDir)}
                                                    fullWidth
                                                >
                                                    <SyncIcon />
                                                </LoadingButton>
                                            </Tooltip>
                                            <Tooltip title="Sync All">
                                                <LoadingButton
                                                    variant="contained"
                                                    color="secondary"
                                                    loading={isLoading}
                                                    onClick={() => handleSyncAll(subDir.subDir)}
                                                    fullWidth
                                                >
                                                    <SyncAllIcon />
                                                </LoadingButton>
                                            </Tooltip>
                                        </Box>
                                    </CardContent>
                                </Card>
                            </Collapse>
                        </Grid>
                    ))}
                </Grid>
            </Box>
        </div>
    );
};

export default Sync;
