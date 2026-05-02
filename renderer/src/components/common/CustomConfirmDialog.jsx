// CustomConfirmDialog — design-system rewrite. Matches AddSkillPlanModal:
// surface-1 paper, rule-1 hairlines, radius-lg, accent confirm button.
import PropTypes from 'prop-types';
import { useId } from 'react';
import { Dialog } from '@mui/material';

const CustomConfirmDialog = ({
    open,
    title,
    message,
    onConfirm,
    onCancel,
    confirmLabel = 'Confirm',
    cancelLabel = 'Cancel',
}) => {
    const titleId = useId();
    const descId = useId();
    return (
        <Dialog
            open={open}
            onClose={onCancel}
            maxWidth="xs"
            fullWidth
            aria-labelledby={titleId}
            aria-describedby={descId}
            PaperProps={{
                className: 'bg-surface-1 border border-rule-1 rounded-lg',
                sx: { backgroundImage: 'none' },
            }}
        >
            <header className="px-5 py-4 border-b border-rule-1">
                <h2 id={titleId} className="text-h3 text-ink-1">{title}</h2>
            </header>

            <div className="px-5 py-4 bg-surface-1">
                <p id={descId} className="text-body text-ink-2">{message}</p>
            </div>

            <footer className="px-5 py-3 border-t border-rule-1 flex items-center justify-end gap-2">
                <button
                    type="button"
                    onClick={onCancel}
                    className="h-8 px-3 rounded-md border border-rule-1 text-meta text-ink-2 hover:bg-surface-2 hover:text-ink-1"
                >
                    {cancelLabel}
                </button>
                <button
                    type="button"
                    onClick={onConfirm}
                    className="h-8 px-3 rounded-md bg-accent text-accent-ink text-meta font-medium hover:bg-accent-strong"
                >
                    {confirmLabel}
                </button>
            </footer>
        </Dialog>
    );
};

CustomConfirmDialog.propTypes = {
    open: PropTypes.bool.isRequired,
    title: PropTypes.string.isRequired,
    message: PropTypes.string.isRequired,
    onConfirm: PropTypes.func.isRequired,
    onCancel: PropTypes.func.isRequired,
    confirmLabel: PropTypes.string,
    cancelLabel: PropTypes.string,
};

export default CustomConfirmDialog;
