// src/components/common/HeaderToolbarActions.jsx
//
// Left-side toolbar cluster: hamburger menu, Add Character, Add Skill Plan,
// and the Fuzzworks status badge. Visible only when authenticated.
import PropTypes from 'prop-types';
import { Box, IconButton, Tooltip } from '@mui/material';
import {
    Menu as MenuIcon,
    PersonAddAlt1Outlined,
    AddchartOutlined,
} from '@mui/icons-material';
import HeaderStatusBadges from './HeaderStatusBadges.jsx';

export const headerIconSx = {
    color: 'var(--ink-2)',
    width: 32,
    height: 32,
    borderRadius: '6px',
    transition: 'background-color var(--motion-duration-fast, 120ms) var(--motion-ease, ease), color var(--motion-duration-fast, 120ms) var(--motion-ease, ease)',
    '&:hover': {
        backgroundColor: 'var(--surface-2)',
        color: 'var(--ink-1)',
        transform: 'none',
    },
};

const HeaderToolbarActions = ({ onOpenDrawer, onAddCharacter, onAddSkillPlan }) => (
    <>
        <Tooltip title="Menu">
            <IconButton
                edge="start"
                aria-label="Menu"
                onClick={onOpenDrawer}
                sx={{ ...headerIconSx, WebkitAppRegion: 'no-drag' }}
            >
                <MenuIcon fontSize="small" />
            </IconButton>
        </Tooltip>
        <Tooltip title="Add character">
            <IconButton
                onClick={onAddCharacter}
                aria-label="Add character"
                sx={{ ...headerIconSx, WebkitAppRegion: 'no-drag' }}
            >
                <PersonAddAlt1Outlined fontSize="small" />
            </IconButton>
        </Tooltip>
        <Tooltip title="Add skill plan">
            <IconButton
                onClick={onAddSkillPlan}
                aria-label="Add skill plan"
                sx={{ ...headerIconSx, WebkitAppRegion: 'no-drag' }}
            >
                <AddchartOutlined fontSize="small" />
            </IconButton>
        </Tooltip>
        <Box sx={{ ml: 0.75, WebkitAppRegion: 'no-drag' }}>
            <HeaderStatusBadges />
        </Box>
    </>
);

HeaderToolbarActions.propTypes = {
    onOpenDrawer: PropTypes.func.isRequired,
    onAddCharacter: PropTypes.func.isRequired,
    onAddSkillPlan: PropTypes.func.isRequired,
};

export default HeaderToolbarActions;
