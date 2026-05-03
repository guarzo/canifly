// src/components/character-overview/GroupBlock.jsx
//
// One section per group key (account name / role / location). Renders a
// sticky header (group name, char count, group SP) and the per-character
// rows underneath. Account-view groups also expose an account kebab menu.
import { useState } from 'react';
import PropTypes from 'prop-types';
import { IconButton, Tooltip, Menu, MenuItem } from '@mui/material';
import { MoreVert as MoreVertIcon } from '@mui/icons-material';
import CharacterRow from './CharacterRow.jsx';
import { formatSP } from './utils';

function AccountMenu({ account, onUpdate, onRemove }) {
    const [anchor, setAnchor] = useState(null);
    const open = Boolean(anchor);
    const close = () => setAnchor(null);
    return (
        <>
            <Tooltip title="Account actions">
                <IconButton
                    aria-label="Account actions"
                    size="small"
                    onClick={(e) => setAnchor(e.currentTarget)}
                    className="h-7! w-7! rounded-md! text-ink-3! hover:bg-surface-3! hover:text-ink-1!"
                >
                    <MoreVertIcon fontSize="small" />
                </IconButton>
            </Tooltip>
            <Menu anchorEl={anchor} open={open} onClose={close}>
                <MenuItem onClick={() => { onUpdate(account.ID, { Visible: !(account.Visible !== false) }); close(); }}>
                    {account.Visible === false ? 'Show account' : 'Hide account'}
                </MenuItem>
                <MenuItem onClick={() => { onUpdate(account.ID, { Status: account.Status === 'Alpha' ? 'Omega' : 'Alpha' }); close(); }}>
                    Toggle Alpha / Omega ({account.Status === 'Alpha' ? 'α' : 'Ω'})
                </MenuItem>
                <MenuItem
                    onClick={() => { onRemove(account.ID); close(); }}
                    sx={{ color: 'var(--status-error)' }}
                >
                    Remove account
                </MenuItem>
            </Menu>
        </>
    );
}
AccountMenu.propTypes = {
    account: PropTypes.object.isRequired,
    onUpdate: PropTypes.func.isRequired,
    onRemove: PropTypes.func.isRequired,
};

const GroupBlock = ({
    groupKey,
    view,
    characters,
    roles,
    skillConversions,
    expanded,
    onToggleExpand,
    focusedRowId,
    setFocusedRowId,
    onUpdateCharacter,
    onRemoveCharacter,
    onUpdateAccount,
    onRemoveAccount,
    allAccounts,
}) => {
    const groupSp = characters.reduce(
        (sum, ch) => sum + (ch.Character?.CharacterSkillsResponse?.total_sp || 0), 0,
    );
    // The hook tags each character with accountId, so look up by id rather
    // than by name (two accounts can share a display name).
    const account = view === 'account' && characters.length > 0
        ? allAccounts.find((a) => a.ID === characters[0].accountId)
        : null;
    const groupLabel = account?.Name || characters[0]?.accountName || groupKey;

    return (
        <section role="rowgroup">
            <header
                role="row"
                className="sticky top-0 z-10 flex items-center justify-between gap-3 px-4 py-2 border-b border-rule-1 bg-surface-2"
            >
                <div className="flex items-center gap-2 min-w-0">
                    <h2 className="text-h3 text-ink-1 truncate uppercase tracking-wide">{view === 'account' ? groupLabel : groupKey}</h2>
                    {view === 'account' && account && account.Visible === false ? (
                        <span className="px-1.5 py-0.5 rounded-sm border border-rule-1 text-micro text-ink-3 uppercase tracking-wide">
                            Hidden
                        </span>
                    ) : null}
                </div>
                <div className="flex items-center gap-3 shrink-0">
                    <span className="text-meta text-ink-3 tabular">
                        {characters.length} {characters.length === 1 ? 'char' : 'chars'} · {formatSP(groupSp)} SP
                    </span>
                    {view === 'account' && account ? (
                        <AccountMenu
                            account={account}
                            onUpdate={onUpdateAccount}
                            onRemove={onRemoveAccount}
                        />
                    ) : null}
                </div>
            </header>

            {characters.map((ch, idx) => (
                <CharacterRow
                    key={ch.Character.CharacterID}
                    character={ch}
                    isLast={idx === characters.length - 1}
                    isExpanded={expanded.has(ch.Character.CharacterID)}
                    onToggleExpand={() => onToggleExpand(ch.Character.CharacterID)}
                    isFocused={focusedRowId === ch.Character.CharacterID}
                    onFocus={() => setFocusedRowId(ch.Character.CharacterID)}
                    roles={roles}
                    skillConversions={skillConversions}
                    onUpdateCharacter={onUpdateCharacter}
                    onRemoveCharacter={onRemoveCharacter}
                />
            ))}
        </section>
    );
};

GroupBlock.propTypes = {
    groupKey: PropTypes.string.isRequired,
    view: PropTypes.string.isRequired,
    characters: PropTypes.array.isRequired,
    roles: PropTypes.array.isRequired,
    skillConversions: PropTypes.object.isRequired,
    expanded: PropTypes.instanceOf(Set).isRequired,
    onToggleExpand: PropTypes.func.isRequired,
    focusedRowId: PropTypes.number,
    setFocusedRowId: PropTypes.func.isRequired,
    onUpdateCharacter: PropTypes.func.isRequired,
    onRemoveCharacter: PropTypes.func.isRequired,
    onUpdateAccount: PropTypes.func.isRequired,
    onRemoveAccount: PropTypes.func.isRequired,
    allAccounts: PropTypes.array.isRequired,
};

export default GroupBlock;
