// src/components/common/HeaderToolbarActions.jsx
//
// Left-side toolbar cluster: hamburger menu, Add Character, Add Skill Plan,
// and the Fuzzworks status badge. Visible only when authenticated.
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
                style={{ WebkitAppRegion: 'no-drag' }}
                sx={headerIconSx}
            >
                <MenuIcon fontSize="small" />
            </IconButton>
        </Tooltip>
        <Tooltip title="Add character">
            <IconButton
                onClick={onAddCharacter}
                aria-label="Add character"
                style={{ WebkitAppRegion: 'no-drag' }}
                sx={headerIconSx}
            >
                <PersonAddAlt1Outlined fontSize="small" />
            </IconButton>
        </Tooltip>
        <Tooltip title="Add skill plan">
            <IconButton
                onClick={onAddSkillPlan}
                aria-label="Add skill plan"
                style={{ WebkitAppRegion: 'no-drag' }}
                sx={headerIconSx}
            >
                <AddchartOutlined fontSize="small" />
            </IconButton>
        </Tooltip>
        <Box style={{ WebkitAppRegion: 'no-drag' }} sx={{ ml: 0.75 }}>
            <HeaderStatusBadges />
        </Box>
    </>
);

export default HeaderToolbarActions;
