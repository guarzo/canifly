// src/components/common/HeaderStatusBadges.jsx
//
// Status chips that live in the Header. Currently renders the Fuzzworks
// update state; new chips can be added here without changing the page.
// Styled to match the calm AppBar — outlined chip, small text.
import { Chip, Tooltip } from '@mui/material';
import { useFuzzworksStatus } from '../../hooks/useFuzzworksStatus';

const FUZZWORKS_VARIANTS = {
    idle:     { label: 'Fuzzworks',          tooltip: 'Awaiting status from server',     color: 'var(--ink-3)' },
    updating: { label: 'Fuzzworks · updating', tooltip: 'Fuzzworks data is updating',     color: 'var(--accent)' },
    ready:    { label: 'Fuzzworks · ready',    tooltip: 'Fuzzworks data is up to date',   color: 'var(--ink-2)' },
    error:    { label: 'Fuzzworks · error',    tooltip: 'Fuzzworks update failed',        color: 'var(--status-error)' },
};

const HeaderStatusBadges = () => {
    const { state, error } = useFuzzworksStatus();
    const variant = FUZZWORKS_VARIANTS[state] || FUZZWORKS_VARIANTS.idle;
    const tooltip = state === 'error' && error ? `Fuzzworks: ${error}` : variant.tooltip;
    return (
        <Tooltip title={tooltip}>
            <Chip
                size="small"
                label={variant.label}
                variant="outlined"
                style={{ WebkitAppRegion: 'no-drag' }}
                sx={{
                    height: 22,
                    fontSize: 11,
                    color: variant.color,
                    borderColor: 'var(--rule-1)',
                    backgroundColor: 'transparent',
                    '.MuiChip-label': { px: 1 },
                }}
            />
        </Tooltip>
    );
};

export default HeaderStatusBadges;
