import PropTypes from 'prop-types';
import {
    Dialog,
    DialogContent,
    IconButton,
    Tooltip,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import OpenInNewIcon from '@mui/icons-material/OpenInNew';
import StatusDot from '../ui/StatusDot.jsx';
import { calculateDaysFromToday } from '../../utils/formatter.jsx';

const CharacterDetailModal = ({ open, onClose, character, skillConversions }) => {
    if (!character || !character.Character) return null;

    const c = character.Character;
    const charId = c.CharacterID;
    const charName = c.CharacterName;
    const portraitUrl = `https://images.evetech.net/characters/${charId}/portrait?size=128`;
    const totalSp = c.CharacterSkillsResponse?.total_sp || 0;
    const formattedSP = totalSp.toLocaleString();
    const zKillUrl = `https://zkillboard.com/character/${charId}/`;
    const queue = Array.isArray(c.SkillQueue) ? c.SkillQueue : [];
    const status = queue.length === 0 ? 'idle' : (character.MCT ? 'training' : 'queued');
    const statusLabel = status === 'training' ? `Training: ${character.Training || 'Unknown'}`
        : status === 'queued' ? 'Queue paused (MCT off)'
        : 'No skill queue';

    const openZkill = () => {
        if (window.electronAPI?.openExternal) window.electronAPI.openExternal(zKillUrl);
        else window.open(zKillUrl, '_blank', 'noopener,noreferrer');
    };

    return (
        <Dialog
            open={open}
            onClose={onClose}
            maxWidth="sm"
            fullWidth
            slotProps={{
                paper: {
                    sx: {
                        bgcolor: 'var(--surface-1)',
                        backgroundImage: 'none',
                        color: 'var(--ink-1)',
                        border: '1px solid var(--rule-1)',
                        borderRadius: '10px',
                        boxShadow: '0 6px 20px oklch(0 0 0 / 0.4)',
                    },
                },
            }}
        >
            <header className="flex items-center justify-between px-5 py-3 border-b border-rule-1">
                <h2 className="text-h3 text-ink-1 truncate">{charName}</h2>
                <IconButton
                    aria-label="Close"
                    onClick={onClose}
                    size="small"
                    sx={{ color: 'var(--ink-3)', '&:hover': { color: 'var(--ink-1)', bgcolor: 'var(--surface-2)' } }}
                >
                    <CloseIcon fontSize="small" />
                </IconButton>
            </header>

            <DialogContent sx={{ p: 0 }}>
                <div className="grid grid-cols-[88px_1fr] gap-5 px-5 py-5 border-b border-rule-1">
                    <img
                        src={portraitUrl}
                        alt=""
                        aria-hidden
                        className="h-22 w-22 rounded-md border border-rule-1"
                        style={{ height: 88, width: 88 }}
                    />
                    <div className="min-w-0">
                        <div className="flex items-center gap-2 mb-2">
                            <StatusDot state={status} label={statusLabel} />
                            <span className="text-meta text-ink-3">{statusLabel}</span>
                            <span className="ml-auto">
                                <Tooltip title="Open zKillboard">
                                    <IconButton
                                        size="small"
                                        onClick={openZkill}
                                        aria-label="Open zKillboard"
                                        sx={{ color: 'var(--ink-3)', '&:hover': { color: 'var(--ink-1)', bgcolor: 'var(--surface-2)' } }}
                                    >
                                        <OpenInNewIcon sx={{ fontSize: 16 }} />
                                    </IconButton>
                                </Tooltip>
                            </span>
                        </div>
                        <dl className="grid grid-cols-[max-content_1fr] gap-x-4 gap-y-1.5 text-meta">
                            {character.Role ? (
                                <>
                                    <dt className="text-ink-3">Role</dt>
                                    <dd className="text-ink-1">{character.Role}</dd>
                                </>
                            ) : null}
                            <dt className="text-ink-3">Location</dt>
                            <dd className="text-ink-1 truncate">{c.LocationName || 'Unknown'}</dd>
                            <dt className="text-ink-3">Total SP</dt>
                            <dd className="text-ink-1 font-mono tabular">{formattedSP}</dd>
                            {character.CorporationName ? (
                                <>
                                    <dt className="text-ink-3">Corporation</dt>
                                    <dd className="text-ink-1 truncate">{character.CorporationName}</dd>
                                </>
                            ) : null}
                            {character.AllianceName ? (
                                <>
                                    <dt className="text-ink-3">Alliance</dt>
                                    <dd className="text-ink-1 truncate">{character.AllianceName}</dd>
                                </>
                            ) : null}
                        </dl>
                    </div>
                </div>

                <div className="px-5 py-4">
                    <div className="flex items-baseline justify-between mb-2">
                        <h3 className="text-meta text-ink-3 uppercase tracking-wide">Skill queue</h3>
                        <span className="text-micro text-ink-3 tabular">
                            {queue.length} item{queue.length === 1 ? '' : 's'}
                        </span>
                    </div>
                    {queue.length === 0 ? (
                        <p className="text-meta text-ink-3">No skill queue.</p>
                    ) : (
                        <ol className="rounded-md border border-rule-1 max-h-72 overflow-y-auto">
                            {queue.map((item, i) => {
                                const skillName =
                                    skillConversions?.[String(item.skill_id)] || `Skill #${item.skill_id}`;
                                const finish = item.finish_date
                                    ? calculateDaysFromToday(item.finish_date)
                                    : '—';
                                return (
                                    <li
                                        key={`${item.skill_id}-${i}`}
                                        className="grid grid-cols-[1fr_44px_92px] gap-3 items-center px-3 h-9 border-b border-rule-1 last:border-b-0 bg-surface-1"
                                    >
                                        <span className="text-meta text-ink-1 truncate">{skillName}</span>
                                        <span className="text-meta text-ink-2 font-mono tabular text-right">L{item.finished_level}</span>
                                        <span className="text-meta text-ink-2 font-mono tabular text-right">{finish}</span>
                                    </li>
                                );
                            })}
                        </ol>
                    )}
                </div>
            </DialogContent>

            <footer className="flex justify-end gap-2 px-5 py-3 border-t border-rule-1">
                <button
                    type="button"
                    onClick={onClose}
                    className="h-8 px-3 rounded-md bg-accent text-accent-ink text-meta hover:bg-accent-strong transition-colors duration-fast ease-out-quart"
                >
                    Close
                </button>
            </footer>
        </Dialog>
    );
};

CharacterDetailModal.propTypes = {
    open: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired,
    character: PropTypes.object.isRequired,
    skillConversions: PropTypes.object.isRequired,
};

export default CharacterDetailModal;
