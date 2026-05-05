import PropTypes from 'prop-types';
import { Select, MenuItem } from '@mui/material';
import { ArrowForwardOutlined } from '@mui/icons-material';
import StatusDot from './StatusDot.jsx';

const selectSx = {
    height: 32,
    minHeight: 32,
    fontSize: 14,
    color: 'var(--ink-1)',
    backgroundColor: 'var(--surface-1)',
    borderRadius: '6px',
    '& .MuiOutlinedInput-notchedOutline': { borderColor: 'var(--rule-1)' },
    '&:hover .MuiOutlinedInput-notchedOutline': { borderColor: 'var(--rule-2)' },
    '&.Mui-focused .MuiOutlinedInput-notchedOutline': {
        borderColor: 'var(--accent)',
        borderWidth: '1px',
    },
    '& .MuiSelect-select': { paddingTop: 6, paddingBottom: 6 },
};

/**
 * PairSelect — a [StatusDot] [Select] → [Select] triad. Used for paired
 * file selection (Sync mode) and similar source→destination pickers. Each
 * option is `{ value, primary, secondary? }`; the closed Select displays
 * the option's `primary` label so users see a name, not a bare id.
 */
const PairSelect = ({
    state,
    leftValue,
    leftOptions,
    onLeftChange,
    leftLabel,
    leftPlaceholder = 'Select…',
    rightValue,
    rightOptions,
    onRightChange,
    rightLabel,
    rightPlaceholder = 'Select…',
    arrowActive,
    statusLabel,
}) => {
    const renderValueFor = (options, placeholder) => (val) => {
        if (!val) return <span className="text-ink-3">{placeholder}</span>;
        const opt = options.find((o) => o.value === val);
        if (!opt) return <span className="font-mono text-body text-ink-1 tabular">{val}</span>;
        return (
            <span className="font-mono text-body text-ink-1 truncate">
                {opt.primary}
                {opt.secondary ? (
                    <span className="ml-2 font-mono text-meta text-ink-3 tabular">{opt.secondary}</span>
                ) : null}
            </span>
        );
    };

    return (
        <div className="grid items-center gap-3 grid-cols-[20px_minmax(180px,1.4fr)_20px_minmax(180px,1.4fr)]">
            <div className="flex items-center justify-center">
                <StatusDot state={state} label={statusLabel} />
            </div>
            <Select
                size="small"
                value={leftValue}
                onChange={(e) => onLeftChange(e.target.value)}
                displayEmpty
                inputProps={{ 'aria-label': leftLabel }}
                renderValue={renderValueFor(leftOptions, leftPlaceholder)}
                sx={selectSx}
                fullWidth
            >
                <MenuItem value="">
                    <span className="text-ink-3">{leftPlaceholder}</span>
                </MenuItem>
                {leftOptions.map((o) => (
                    <MenuItem key={o.value} value={o.value}>
                        <span className="font-mono text-body text-ink-1">{o.primary}</span>
                        {o.secondary ? (
                            <span className="ml-2 font-mono text-meta text-ink-3 tabular">{o.secondary}</span>
                        ) : null}
                    </MenuItem>
                ))}
            </Select>
            <div className="flex items-center justify-center" aria-hidden>
                <ArrowForwardOutlined sx={{ fontSize: 16, color: arrowActive ? 'var(--accent)' : 'var(--ink-4)' }} />
            </div>
            <Select
                size="small"
                value={rightValue}
                onChange={(e) => onRightChange(e.target.value)}
                displayEmpty
                inputProps={{ 'aria-label': rightLabel }}
                renderValue={renderValueFor(rightOptions, rightPlaceholder)}
                sx={selectSx}
                fullWidth
            >
                <MenuItem value="">
                    <span className="text-ink-3">{rightPlaceholder}</span>
                </MenuItem>
                {rightOptions.map((o) => (
                    <MenuItem key={o.value} value={o.value}>
                        <span className="font-mono text-body text-ink-1">{o.primary}</span>
                        {o.secondary ? (
                            <span className="ml-2 font-mono text-meta text-ink-3 tabular">{o.secondary}</span>
                        ) : null}
                    </MenuItem>
                ))}
            </Select>
        </div>
    );
};

const optionShape = PropTypes.shape({
    value: PropTypes.string.isRequired,
    primary: PropTypes.string.isRequired,
    secondary: PropTypes.string,
});

PairSelect.propTypes = {
    state: PropTypes.oneOf(['ready', 'training', 'queued', 'idle', 'error']).isRequired,
    leftValue: PropTypes.string.isRequired,
    leftOptions: PropTypes.arrayOf(optionShape).isRequired,
    onLeftChange: PropTypes.func.isRequired,
    leftLabel: PropTypes.string.isRequired,
    leftPlaceholder: PropTypes.string,
    rightValue: PropTypes.string.isRequired,
    rightOptions: PropTypes.arrayOf(optionShape).isRequired,
    onRightChange: PropTypes.func.isRequired,
    rightLabel: PropTypes.string.isRequired,
    rightPlaceholder: PropTypes.string,
    arrowActive: PropTypes.bool,
    statusLabel: PropTypes.string,
};

export default PairSelect;
