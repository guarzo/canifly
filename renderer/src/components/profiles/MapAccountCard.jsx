// MapAccountCard.jsx — redesigned per DESIGN.md.
// A user-file row: dense header (swatch + name + mtime + count) with the
// associated characters listed beneath as inline rows. Drop target.

import { forwardRef } from 'react';
import PropTypes from 'prop-types';
import { DeleteOutlined as DeleteOutline } from '@mui/icons-material';
import { IconButton, Tooltip } from '@mui/material';
import { formatDate } from '../../utils/formatter.jsx';
import MtimeSwatch from './MtimeSwatch.jsx';

const AccountCard = forwardRef(function AccountCard(
    { mapping, associations, handleUnassociate, handleDrop, mtimeToColor },
    ref,
) {
    const userId = mapping.userId;
    const userName = mapping.name || `Account ${userId}`;
    const associatedChars = associations.filter((a) => a.userId === userId);
    const mtime = mapping.mtime;
    const swatch = mtimeToColor[mtime];

    return (
        <div
            onDragOver={(e) => e.preventDefault()}
            onDrop={(e) => handleDrop(e, userId, userName)}
            className="rounded-lg border border-rule-1 bg-surface-1 overflow-hidden"
        >
            <div ref={ref} className="grid grid-cols-[16px_minmax(0,1fr)_auto_60px] gap-3 items-center h-10 px-4 bg-surface-2 border-b border-rule-1">
                <MtimeSwatch color={swatch} title={`mtime ${mtime}`} />
                <span className="text-body text-ink-1 truncate" title={userName}>
                    {userName}
                </span>
                <span className="text-meta text-ink-3 font-mono tabular">
                    {mtime ? formatDate(mtime) : '—'}
                </span>
                <span className="text-meta text-ink-3 font-mono tabular text-right">
                    {associatedChars.length}/—
                </span>
            </div>

            {associatedChars.length === 0 ? (
                <div className="h-10 px-4 flex items-center text-meta text-ink-3 italic">
                    Drop a character here to associate.
                </div>
            ) : (
                <ul role="list">
                    {associatedChars.map((assoc, idx) => (
                        <li
                            key={`assoc-${assoc.charId}`}
                            className={[
                                'grid grid-cols-[16px_minmax(0,1fr)_auto_28px] gap-3 items-center',
                                'h-10 px-4',
                                idx !== associatedChars.length - 1 ? 'border-b border-rule-1' : '',
                            ].join(' ')}
                        >
                            <MtimeSwatch
                                color={mtimeToColor[assoc.mtime]}
                                title={assoc.mtime ? `mtime ${assoc.mtime}` : 'mtime unknown'}
                            />
                            <span className="text-body text-ink-1 truncate">{assoc.charName}</span>
                            <span className="text-meta text-ink-3 font-mono tabular">
                                {assoc.charId}
                            </span>
                            <div className="flex items-center justify-center">
                                <Tooltip title={`Unassociate ${assoc.charName}`}>
                                    <IconButton
                                        size="small"
                                        aria-label={`Unassociate ${assoc.charName}`}
                                        onClick={() =>
                                            handleUnassociate(userId, assoc.charId, assoc.charName, userName)
                                        }
                                        className="h-6! w-6! rounded-md! text-ink-3! hover:bg-surface-3! hover:text-status-error!"
                                    >
                                        <DeleteOutline sx={{ fontSize: 16 }} />
                                    </IconButton>
                                </Tooltip>
                            </div>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
});

AccountCard.propTypes = {
    mapping: PropTypes.shape({
        userId: PropTypes.string.isRequired,
        name: PropTypes.string,
        mtime: PropTypes.string,
    }).isRequired,
    associations: PropTypes.arrayOf(
        PropTypes.shape({
            userId: PropTypes.string.isRequired,
            charId: PropTypes.string.isRequired,
            charName: PropTypes.string.isRequired,
            mtime: PropTypes.string,
        }),
    ).isRequired,
    handleUnassociate: PropTypes.func.isRequired,
    handleDrop: PropTypes.func.isRequired,
    mtimeToColor: PropTypes.objectOf(PropTypes.string).isRequired,
};

export default AccountCard;
