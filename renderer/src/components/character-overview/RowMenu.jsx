// src/components/character-overview/RowMenu.jsx
//
// Per-row kebab menu — currently just "Remove character".
import { useState } from 'react';
import PropTypes from 'prop-types';
import { IconButton, Menu, MenuItem } from '@mui/material';
import { DeleteOutline, MoreVert as MoreVertIcon } from '@mui/icons-material';

const RowMenu = ({ character, onRemove }) => {
    const [anchor, setAnchor] = useState(null);
    const open = Boolean(anchor);
    const close = () => setAnchor(null);
    return (
        <>
            <IconButton
                size="small"
                aria-label="Character actions"
                onClick={(e) => { e.stopPropagation(); setAnchor(e.currentTarget); }}
                className="!h-7 !w-7 !rounded-md !text-ink-3 hover:!bg-surface-3 hover:!text-ink-1"
            >
                <MoreVertIcon fontSize="small" />
            </IconButton>
            <Menu anchorEl={anchor} open={open} onClose={close}>
                <MenuItem
                    onClick={() => { onRemove(character.Character.CharacterID); close(); }}
                    sx={{ color: 'var(--status-error)' }}
                >
                    <DeleteOutline fontSize="small" sx={{ mr: 1 }} />
                    Remove character
                </MenuItem>
            </Menu>
        </>
    );
};

RowMenu.propTypes = {
    character: PropTypes.object.isRequired,
    onRemove: PropTypes.func.isRequired,
};

export default RowMenu;
