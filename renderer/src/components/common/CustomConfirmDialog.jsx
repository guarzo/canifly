// CustomConfirmDialog.jsx
import React from 'react';
import PropTypes from 'prop-types';
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    Button,
    Typography,
    useTheme,
    Divider
} from '@mui/material';
import CheckIcon from '@mui/icons-material/Check';
import CloseIcon from '@mui/icons-material/Close';

const CustomConfirmDialog = ({
                                 open,
                                 title,
                                 message,
                                 onConfirm,
                                 onCancel
                             }) => {
    const theme = useTheme();

    return (
        <Dialog
            open={open}
            onClose={onCancel}
            PaperProps={{
                sx: {
                    borderRadius: 2,
                    boxShadow: theme.shadows[5],
                    backgroundColor: theme.palette.background.paper
                }
            }}
        >
            <DialogTitle
                sx={{
                    backgroundColor: theme.palette.primary.main,
                    color: theme.palette.primary.contrastText,
                    padding: theme.spacing(2),
                    fontWeight: 600,
                }}
            >
                {title}
            </DialogTitle>

            <Divider />

            <DialogContent sx={{ padding: theme.spacing(3) }}>
                <Typography variant="body1">{message}</Typography>
            </DialogContent>

            <Divider />

            <DialogActions sx={{ padding: theme.spacing(2) }}>
                <Button
                    onClick={onCancel}
                    color="inherit"
                    startIcon={<CloseIcon />}
                    variant="outlined"
                    sx={{ mr: 1 }}
                >
                    Cancel
                </Button>
                <Button
                    onClick={onConfirm}
                    color="primary"
                    startIcon={<CheckIcon />}
                    variant="contained"
                >
                    Confirm
                </Button>
            </DialogActions>
        </Dialog>
    );
};

CustomConfirmDialog.propTypes = {
    open: PropTypes.bool.isRequired,
    title: PropTypes.string.isRequired,
    message: PropTypes.string.isRequired,
    onConfirm: PropTypes.func.isRequired,
    onCancel: PropTypes.func.isRequired,
};

export default CustomConfirmDialog;
