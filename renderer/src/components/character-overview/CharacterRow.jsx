// src/components/character-overview/CharacterRow.jsx
//
// One row per character: status dot, portrait, name, account, SP, queue
// ETA, location, role, and a kebab. Stateless w.r.t. the page; the
// expanded detail panel is rendered as a sibling when isExpanded.
import PropTypes from 'prop-types';
import {
    OpenInNewOutlined,
    ExpandMoreOutlined,
} from '@mui/icons-material';
import { IconButton, Tooltip } from '@mui/material';
import StatusDot from '../ui/StatusDot.jsx';
import RowMenu from './RowMenu.jsx';
import ExpandedRow from './ExpandedRow.jsx';
import { formatSP, formatDuration, deriveQueueEta, deriveStatus } from './utils';

const CharacterRow = ({
    character,
    isLast,
    isExpanded,
    onToggleExpand,
    isFocused,
    onFocus,
    roles,
    skillConversions,
    onUpdateCharacter,
    onRemoveCharacter,
}) => {
    const c = character.Character || {};
    const id = c.CharacterID;
    const portrait = `https://images.evetech.net/characters/${id}/portrait?size=64`;
    const sp = c.CharacterSkillsResponse?.total_sp || 0;
    const eta = deriveQueueEta(character);
    const status = deriveStatus(character);
    const role = character.Role || '';

    const openZkill = (e) => {
        e.stopPropagation();
        const url = `https://zkillboard.com/character/${id}/`;
        if (window.electronAPI?.openExternal) window.electronAPI.openExternal(url);
        else window.open(url, '_blank', 'noopener,noreferrer');
    };

    return (
        <>
            <div
                role="row"
                tabIndex={0}
                aria-expanded={isExpanded}
                onClick={onToggleExpand}
                onFocus={onFocus}
                onKeyDown={(e) => {
                    if (e.key !== 'Enter' && e.key !== ' ') return;
                    // Don't toggle when the key originated from a nested
                    // interactive control (zKill button, kebab, etc.) — let
                    // that control handle the key itself.
                    if (e.target !== e.currentTarget) return;
                    e.preventDefault();
                    onToggleExpand();
                }}
                className={[
                    'grid grid-cols-[28px_minmax(180px,1.4fr)_1fr_84px_92px_minmax(140px,1fr)_minmax(120px,0.8fr)_28px]',
                    'gap-3 px-4 items-center',
                    'h-10 cursor-pointer select-none',
                    'transition-colors duration-fast ease-out-quart',
                    !isLast ? 'border-b border-rule-1' : '',
                    isFocused ? 'bg-surface-2 shadow-rail-accent' : 'hover:bg-surface-2',
                ].join(' ')}
            >
                <div className="flex items-center justify-center">
                    {isExpanded ? (
                        <ExpandMoreOutlined fontSize="small" sx={{ color: 'var(--accent)' }} />
                    ) : (
                        <StatusDot state={status} />
                    )}
                </div>
                <div className="flex items-center gap-2.5 min-w-0">
                    <img
                        src={portrait}
                        alt=""
                        aria-hidden
                        loading="lazy"
                        className="h-6 w-6 rounded-sm border border-rule-1 shrink-0"
                    />
                    <span className="text-body text-ink-1 truncate">{c.CharacterName}</span>
                    <Tooltip title="Open zKillboard">
                        <IconButton
                            size="small"
                            onClick={openZkill}
                            aria-label={`Open zKillboard for ${c.CharacterName}`}
                            className="!h-6 !w-6 !text-ink-3 hover:!bg-surface-3 hover:!text-ink-1"
                        >
                            <OpenInNewOutlined sx={{ fontSize: 14 }} />
                        </IconButton>
                    </Tooltip>
                </div>
                <span className="text-meta text-ink-3 truncate">{character.accountName}</span>
                <span className="text-body text-ink-1 font-mono tabular text-right" data-numeric="true">
                    {formatSP(sp)}
                </span>
                <span
                    className={`text-body font-mono tabular text-right ${
                        status === 'training' ? 'text-status-training' : 'text-ink-2'
                    }`}
                    data-numeric="true"
                    title={status === 'idle' ? 'No skill queue' : 'Time until queue completes'}
                >
                    {status === 'idle' ? '—' : (eta != null ? formatDuration(eta) : 'queued')}
                </span>
                <span className="text-meta text-ink-2 truncate">{c.LocationName || '—'}</span>
                <span className="text-meta text-ink-2 truncate">
                    {role || <span className="text-ink-3">—</span>}
                </span>
                <div className="flex items-center justify-center" onClick={(e) => e.stopPropagation()}>
                    <RowMenu character={character} onRemove={onRemoveCharacter} />
                </div>
            </div>

            {isExpanded ? (
                <ExpandedRow
                    character={character}
                    roles={roles}
                    skillConversions={skillConversions}
                    onUpdateCharacter={onUpdateCharacter}
                />
            ) : null}
        </>
    );
};

CharacterRow.propTypes = {
    character: PropTypes.object.isRequired,
    isLast: PropTypes.bool.isRequired,
    isExpanded: PropTypes.bool.isRequired,
    onToggleExpand: PropTypes.func.isRequired,
    isFocused: PropTypes.bool.isRequired,
    onFocus: PropTypes.func.isRequired,
    roles: PropTypes.array.isRequired,
    skillConversions: PropTypes.object.isRequired,
    onUpdateCharacter: PropTypes.func.isRequired,
    onRemoveCharacter: PropTypes.func.isRequired,
};

export default CharacterRow;
