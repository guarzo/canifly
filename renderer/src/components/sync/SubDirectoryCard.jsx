import React from 'react';
import PropTypes from 'prop-types';
import { Card, Typography, Divider, FormControl, InputLabel, Select, MenuItem, Grid, Tooltip } from '@mui/material';
import LoadingButton from '@mui/lab/LoadingButton';
import SyncIcon from '@mui/icons-material/Sync';
import SyncAllIcon from '@mui/icons-material/SyncAlt';

const SubDirectoryCard = ({
                              subDir,
                              selections,
                              handleSelectionChange,
                              handleSync,
                              handleSyncAll,
                              isLoading
                          }) => {
    const sortedUserFiles = [...subDir.availableUserFiles].sort((a, b) => a.name.localeCompare(b.name));
    const sortedCharFiles = [...subDir.availableCharFiles].sort((a, b) => a.name.localeCompare(b.name));

    const displaySubDir = subDir.profile.replace('settings_', '');
    const selectedCharId = selections[subDir.profile]?.charId || '';
    const selectedUserId = selections[subDir.profile]?.userId || '';

    return (
        <Card className="bg-gray-800 text-teal-200 p-4 rounded-md shadow-md flex flex-col justify-between h-full transform transition-transform duration-200 ease-in-out hover:scale-105 hover:shadow-lg">
            <div>
                <Typography
                    variant="h6"
                    gutterBottom
                    sx={{ textAlign: 'center', fontWeight: 700 }}
                >
                    {displaySubDir}
                </Typography>
                <Divider sx={{ my: 2, borderColor: 'rgba(255,255,255,0.2)' }} />
                <div className="p-2 rounded-md border border-gray-600 bg-gray-700 mb-2">
                    <FormControl fullWidth margin="normal">
                        <InputLabel id={`char-select-label-${subDir.profile}`} sx={{ color: '#99f6e4' }}>
                            Select Character
                        </InputLabel>
                        <Select
                            labelId={`char-select-label-${subDir.profile}`}
                            id={`char-select-${subDir.profile}`}
                            value={selectedCharId}
                            label="Select Character"
                            onChange={(e) => handleSelectionChange(subDir.profile, 'charId', e.target.value)}
                            sx={{
                                borderRadius: 1,
                                color: '#fff',
                                '& .MuiOutlinedInput-notchedOutline': {
                                    borderColor: 'rgba(255,255,255,0.2)',
                                },
                                '&:hover .MuiOutlinedInput-notchedOutline': {
                                    borderColor: '#ffffff',
                                },
                                '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                                    borderColor: '#ffffff',
                                },
                            }}
                        >
                            <MenuItem value="">
                                <em>-- Select Character --</em>
                            </MenuItem>
                            {sortedCharFiles.map(char => (
                                <MenuItem key={char.charId} value={char.charId}>
                                    {char.name}
                                </MenuItem>
                            ))}
                        </Select>
                    </FormControl>
                    <FormControl fullWidth margin="normal">
                        <InputLabel id={`user-select-label-${subDir.profile}`} sx={{ color: '#99f6e4' }}>
                            Select User
                        </InputLabel>
                        <Select
                            labelId={`user-select-label-${subDir.profile}`}
                            id={`user-select-${subDir.profile}`}
                            value={selectedUserId}
                            label="Select User"
                            onChange={(e) => handleSelectionChange(subDir.profile, 'userId', e.target.value)}
                            sx={{
                                borderRadius: 1,
                                color: '#fff',
                                '& .MuiOutlinedInput-notchedOutline': {
                                    borderColor: 'rgba(255,255,255,0.2)',
                                },
                                '&:hover .MuiOutlinedInput-notchedOutline': {
                                    borderColor: '#ffffff',
                                },
                                '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
                                    borderColor: '#ffffff',
                                },
                            }}
                        >
                            <MenuItem value="">
                                <em>-- Select User --</em>
                            </MenuItem>
                            {sortedUserFiles.map(user => (
                                <MenuItem key={user.userId} value={user.userId}>
                                    {user.userId}
                                </MenuItem>
                            ))}
                        </Select>
                    </FormControl>
                </div>
            </div>
            <div>
                <Grid container spacing={2}>
                    <Grid item xs={6}>
                        <Tooltip title="Sync this specific profile">
                            <span style={{ width: '100%', display: 'block' }}>
                                <LoadingButton
                                    aria-label="Sync this specific profile"
                                    variant="contained"
                                    color="primary"
                                    onClick={() => handleSync(subDir.profile)}
                                    loading={isLoading}
                                    fullWidth
                                    disabled={isLoading || !selectedCharId || !selectedUserId}
                                    className="p-0"
                                >
                                    <SyncIcon fontSize="small" />
                                </LoadingButton>
                            </span>
                        </Tooltip>
                    </Grid>
                    <Grid item xs={6}>
                        <Tooltip title="Sync all profiles based on this selection">
                            <span style={{ width: '100%', display: 'block' }}>
                                <LoadingButton
                                    aria-label="Sync all profiles based on this selection"
                                    variant="contained"
                                    color="secondary"
                                    onClick={() => handleSyncAll(subDir.profile)}
                                    loading={isLoading}
                                    fullWidth
                                    disabled={isLoading || !selectedCharId || !selectedUserId}
                                    className="p-0"
                                >
                                    <SyncAllIcon fontSize="small" />
                                </LoadingButton>
                            </span>
                        </Tooltip>
                    </Grid>
                </Grid>
            </div>
        </Card>
    );
};

SubDirectoryCard.propTypes = {
    subDir: PropTypes.shape({
        profile: PropTypes.string.isRequired,
        availableCharFiles: PropTypes.array.isRequired,
        availableUserFiles: PropTypes.array.isRequired,
    }).isRequired,
    selections: PropTypes.object.isRequired,
    handleSelectionChange: PropTypes.func.isRequired,
    handleSync: PropTypes.func.isRequired,
    handleSyncAll: PropTypes.func.isRequired,
    isLoading: PropTypes.bool.isRequired,
};

export default SubDirectoryCard;
