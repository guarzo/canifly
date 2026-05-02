// src/components/common/Header.jsx
//
// Page-level composition for the application Header. State (drawer/modal
// open) lives here; data, side-effects, and the OAuth orchestration live
// in `useHeaderData`. Each visual section is its own subcomponent.
import { useState } from 'react';
import PropTypes from 'prop-types';
import { useLocation } from 'react-router-dom';
import { AppBar, Toolbar, Box } from '@mui/material';
import { styled } from '@mui/material/styles';
import AccountPromptModal from './AccountPromptModal.jsx';
import HeaderNav from './HeaderNav.jsx';
import HeaderUserMenu from './HeaderUserMenu.jsx';
import HeaderToolbarActions from './HeaderToolbarActions.jsx';
import { useHeaderData } from '../../hooks/useHeaderData';

const StyledAppBar = styled(AppBar)(() => ({
    background: 'var(--surface-1)',
    backgroundImage: 'none',
    backdropFilter: 'none',
    borderBottom: '1px solid var(--rule-1)',
    boxShadow: 'none',
    color: 'var(--ink-1)',
}));

const Header = ({ openSkillPlanModal, existingAccounts }) => {
    const location = useLocation();
    const [drawerOpen, setDrawerOpen] = useState(false);
    const [modalOpen, setModalOpen] = useState(false);
    const {
        isAuthenticated,
        logout,
        isRefreshing,
        handleRefreshClick,
        handleCloseWindow,
        handleAddCharacterSubmit,
    } = useHeaderData();

    const onSubmitAddCharacter = async (account) => {
        await handleAddCharacterSubmit(account);
        setModalOpen(false);
    };

    return (
        <>
            <StyledAppBar position="fixed">
                <Toolbar
                    variant="dense"
                    sx={{ WebkitAppRegion: 'drag', minHeight: 44, gap: 0.5, px: 1.5 }}
                >
                    {isAuthenticated && (
                        <HeaderToolbarActions
                            onOpenDrawer={() => setDrawerOpen(true)}
                            onAddCharacter={() => setModalOpen(true)}
                            onAddSkillPlan={openSkillPlanModal}
                        />
                    )}

                    <Box sx={{ flex: 1, display: 'flex', justifyContent: 'center', minWidth: 0 }}>
                        <span
                            style={{
                                fontFamily: 'Inter, system-ui, sans-serif',
                                fontWeight: 600,
                                fontSize: 13,
                                letterSpacing: '0.06em',
                                textTransform: 'uppercase',
                                color: 'var(--ink-2)',
                            }}
                        >
                            CanIFly
                        </span>
                    </Box>

                    <HeaderUserMenu
                        isAuthenticated={isAuthenticated}
                        isRefreshing={isRefreshing}
                        onRefresh={handleRefreshClick}
                        onLogout={logout}
                        onClose={handleCloseWindow}
                    />
                </Toolbar>
            </StyledAppBar>

            <HeaderNav
                open={drawerOpen}
                onClose={() => setDrawerOpen(false)}
                currentPath={location.pathname}
            />

            <AccountPromptModal
                isOpen={modalOpen}
                onClose={() => setModalOpen(false)}
                onSubmit={onSubmitAddCharacter}
                existingAccounts={existingAccounts}
                title="Add Character - Enter Account Name"
            />
        </>
    );
};

Header.propTypes = {
    openSkillPlanModal: PropTypes.func.isRequired,
    existingAccounts: PropTypes.array,
};

export default Header;
