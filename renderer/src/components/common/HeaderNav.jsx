// src/components/common/HeaderNav.jsx
//
// Drawer-based primary navigation for the Header. Stateless — drawer
// open/close lives in the parent.
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
import {
    Drawer,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
    ListItemButton,
    Box,
} from '@mui/material';
import {
    SpaceDashboardOutlined as CharacterOverviewIcon,
    ListAltOutlined as SkillPlansIcon,
    ManageAccountsOutlined as ProfilesIcon,
    SettingsOutlined as SettingsIcon,
} from '@mui/icons-material';
import { styled } from '@mui/material/styles';

const StyledDrawer = styled(Drawer)(() => ({
    '& .MuiPaper-root': {
        background: 'var(--surface-1)',
        backgroundImage: 'none',
        backdropFilter: 'none',
        borderRight: '1px solid var(--rule-1)',
        boxShadow: 'none',
        color: 'var(--ink-1)',
        // 240px on narrow screens, 260px from sm-up. Fixed width is intentional —
        // the drawer is overlay-mode, not push, so it never crowds content.
        width: 240,
        '@media (min-width: 640px)': { width: 260 },
    },
}));

const NAV_LINKS = [
    { text: 'Overview',    icon: <CharacterOverviewIcon fontSize="small" />, path: '/' },
    { text: 'Skill Plans', icon: <SkillPlansIcon fontSize="small" />,        path: '/skill-plans' },
    { text: 'Profiles',    icon: <ProfilesIcon fontSize="small" />,          path: '/profiles' },
    { text: 'Settings',    icon: <SettingsIcon fontSize="small" />,          path: '/settings' },
];

const navItemSx = (selected) => ({
    borderRadius: '6px',
    px: 1.5,
    py: 1,
    transition: 'background-color var(--motion-duration-fast, 120ms) var(--motion-ease, ease)',
    color: 'var(--ink-2)',
    '& .MuiListItemIcon-root': { color: 'var(--ink-3)', minWidth: 28 },
    '&:hover': {
        backgroundColor: 'var(--surface-2)',
        color: 'var(--ink-1)',
        transform: 'none',
        '& .MuiListItemIcon-root': { color: 'var(--ink-1)', transform: 'none' },
    },
    '&.Mui-selected': {
        backgroundColor: 'var(--accent-soft)',
        color: 'var(--ink-1)',
        '& .MuiListItemIcon-root': { color: 'var(--accent)' },
        '&:hover': { backgroundColor: 'var(--accent-soft)' },
    },
    ...(selected ? {} : {}),
});

const HeaderNav = ({ open, onClose, currentPath }) => {
    // Drawer's MUI onClose handles backdrop/Escape. The presentation div
    // mirrors clicks for closing on link selection — but we must NOT close on
    // Tab/Shift+Tab so keyboard users can move focus between nav items.
    const handleKeyDown = (event) => {
        if (event.key === 'Escape') onClose(event);
    };
    return (
        <StyledDrawer anchor="left" open={open} onClose={onClose} disableScrollLock>
            <div
                role="presentation"
                onClick={onClose}
                onKeyDown={handleKeyDown}
                style={{ flex: 1, display: 'flex', flexDirection: 'column' }}
            >
            <Box sx={{ px: 2, pt: 2, pb: 1.5 }}>
                <span
                    style={{
                        fontSize: 11,
                        fontWeight: 500,
                        letterSpacing: '0.08em',
                        textTransform: 'uppercase',
                        color: 'var(--ink-3)',
                    }}
                >
                    CanIFly
                </span>
            </Box>
            <Box sx={{ height: '1px', bgcolor: 'var(--rule-1)', mx: 2 }} />
            <List sx={{ flex: 1, py: 1, px: 1 }}>
                {NAV_LINKS.map((item) => {
                    const selected = currentPath === item.path;
                    return (
                        <ListItem key={item.text} disablePadding sx={{ mb: 0.25 }}>
                            <ListItemButton
                                component={Link}
                                to={item.path}
                                selected={selected}
                                sx={navItemSx(selected)}
                            >
                                <ListItemIcon>{item.icon}</ListItemIcon>
                                <ListItemText
                                    primary={item.text}
                                    primaryTypographyProps={{
                                        sx: {
                                            fontSize: 14,
                                            fontWeight: selected ? 500 : 400,
                                            color: 'inherit',
                                        },
                                    }}
                                />
                            </ListItemButton>
                        </ListItem>
                    );
                })}
            </List>
        </div>
    </StyledDrawer>
    );
};

HeaderNav.propTypes = {
    open: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired,
    currentPath: PropTypes.string.isRequired,
};

export default HeaderNav;
