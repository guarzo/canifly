import React from 'react';
import PropTypes from 'prop-types';
import { Button, Tooltip } from '@mui/material';
import BackupIcon from '@mui/icons-material/Backup';
import FolderOpenIcon from '@mui/icons-material/FolderOpen';
import UndoIcon from '@mui/icons-material/Undo';

const SyncActionsBar = ({
                            handleBackup,
                            handleChooseSettingsDir,
                            handleResetToDefault,
                            isDefaultDir,
                            isLoading
                        }) => {
    return (
        <div className="mb-4 flex justify-center space-x-2">
            <Tooltip title="Backup Settings">
                <span>
                    <Button
                        aria-label="Backup Settings"
                        variant="contained"
                        color="primary"
                        onClick={handleBackup}
                        disabled={isLoading}
                        className="w-10 h-10 p-0 flex items-center justify-center"
                    >
                        <BackupIcon fontSize="small" />
                    </Button>
                </span>
            </Tooltip>
            <Tooltip title="Choose Settings Directory">
                <span>
                    <Button
                        aria-label="Choose Settings Directory"
                        variant="contained"
                        color="info"
                        onClick={handleChooseSettingsDir}
                        disabled={isLoading}
                        className="w-10 h-10 p-0 flex items-center justify-center"
                    >
                        <FolderOpenIcon fontSize="small" />
                    </Button>
                </span>
            </Tooltip>
            {!isDefaultDir && (
                <Tooltip title="Reset to Default Directory">
                    <span>
                        <Button
                            aria-label="Reset to Default Directory"
                            variant="contained"
                            color="warning"
                            onClick={handleResetToDefault}
                            disabled={isLoading}
                            className="w-10 h-10 p-0 flex items-center justify-center"
                        >
                            <UndoIcon fontSize="small" />
                        </Button>
                    </span>
                </Tooltip>
            )}
        </div>
    );
};

SyncActionsBar.propTypes = {
    handleBackup: PropTypes.func.isRequired,
    handleChooseSettingsDir: PropTypes.func.isRequired,
    handleResetToDefault: PropTypes.func.isRequired,
    isDefaultDir: PropTypes.bool.isRequired,
    isLoading: PropTypes.bool.isRequired,
};

export default SyncActionsBar;
