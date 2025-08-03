// useConfirmDialog.js
import { useState } from 'react';
import CustomConfirmDialog from '../components/common/CustomConfirmDialog.jsx';

export const useConfirmDialog = () => {
    const [confirmDialogOptions, setConfirmDialogOptions] = useState({
        open: false,
        title: '',
        message: '',
        resolve: null
    });

    const showConfirmDialog = ({ title, message }) => {
        return new Promise((resolve) => {
            setConfirmDialogOptions({
                open: true,
                title,
                message,
                resolve
            });
        });
    };

    const handleConfirm = () => {
        confirmDialogOptions.resolve({ isConfirmed: true });
        setConfirmDialogOptions((prev) => ({ ...prev, open: false }));
    };

    const handleCancel = () => {
        confirmDialogOptions.resolve({ isConfirmed: false });
        setConfirmDialogOptions((prev) => ({ ...prev, open: false }));
    };

    const confirmDialog = (
        <CustomConfirmDialog
            open={confirmDialogOptions.open}
            title={confirmDialogOptions.title}
            message={confirmDialogOptions.message}
            onConfirm={handleConfirm}
            onCancel={handleCancel}
        />
    );

    return [showConfirmDialog, confirmDialog];
};
