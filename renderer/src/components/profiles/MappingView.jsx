// renderer/src/components/profiles/MappingView.jsx
import { useCallback, useState } from 'react';
import PropTypes from 'prop-types';
import { toast } from 'react-toastify';

import MapAccountCard from './MapAccountCard.jsx';
import MapCharacterCard from './MapCharacterCard.jsx';
import { useMappingDerivations } from '../../hooks/useMappingDerivations.js';
import { useConfirmDialog } from '../../hooks/useConfirmDialog.jsx';
import { useAppData } from '../../hooks/useAppData';
import { associateCharacter, unassociateCharacter } from '../../api/accountsApi';

const sortRows = (list, sortOrder) => {
    const copy = [...list];
    copy.sort((a, b) => {
        switch (sortOrder) {
            case 'mtime-asc': return new Date(a.mtime) - new Date(b.mtime);
            case 'name-asc': return (a.name || '').localeCompare(b.name || '');
            case 'name-desc': return (b.name || '').localeCompare(a.name || '');
            case 'mtime-desc':
            default: return new Date(b.mtime) - new Date(a.mtime);
        }
    });
    return copy;
};

const matchesQuery = (q) => (str) => !q || (str || '').toLowerCase().includes(q);

const MappingView = ({ subDirs, associations: initialAssociations, filter, view, sortOrder }) => {
    const { refreshData } = useAppData();
    const [associations, setAssociations] = useState(initialAssociations);
    const [showConfirmDialog, confirmDialog] = useConfirmDialog();

    const { accounts, availableCharacters, mtimeToColor } =
        useMappingDerivations(subDirs, associations);

    const handleDragStart = (event, charId) => {
        event.dataTransfer.setData('text/plain', charId);
    };

    const handleDrop = useCallback(async (event, userId, userName) => {
        event.preventDefault();
        const charId = event.dataTransfer.getData('text/plain');
        if (associations.some((assoc) => assoc.charId === charId)) {
            toast.error('That character is already associated!');
            return;
        }
        const char = availableCharacters.find((c) => c.charId === charId);
        if (!char) {
            toast.error('Character not found.');
            return;
        }
        const ok = await showConfirmDialog({
            title: 'Confirm Association',
            message: `Associate "${char.name}" with account "${userName}"?`,
        });
        if (!ok.isConfirmed) return;
        const result = await associateCharacter(userId, charId, userName, char.name);
        if (result?.success) {
            toast.success(result.message);
            setAssociations((prev) => [...prev, { userId, charId, charName: char.name, mtime: char.mtime }]);
            if (refreshData) await refreshData();
        }
    }, [associations, availableCharacters, showConfirmDialog, refreshData]);

    const handleUnassociate = useCallback(async (userId, charId, charName, userName) => {
        const ok = await showConfirmDialog({
            title: 'Confirm Unassociation',
            message: `Unassociate "${charName}" from account "${userName}"?`,
        });
        if (!ok.isConfirmed) return;
        const result = await unassociateCharacter(userId, charId, userName, charName);
        if (result?.success) {
            toast.success(result.message);
            setAssociations((prev) =>
                prev.filter((a) => a.charId !== charId || a.userId !== userId),
            );
            if (refreshData) await refreshData();
        }
    }, [showConfirmDialog, refreshData]);

    const q = filter.trim().toLowerCase();
    const m = matchesQuery(q);

    const filteredAccounts = (() => {
        const list = accounts.filter((a) => {
            const assocs = associations.filter((x) => x.userId === a.userId);
            const hasMatchedFilter =
                view === 'all' ||
                (view === 'matched' && assocs.length > 0) ||
                (view === 'unmatched' && assocs.length === 0);
            if (!hasMatchedFilter) return false;
            if (!q) return true;
            const nameHit = m(a.name) || m(a.userId);
            const assocHit = assocs.some((x) => m(x.charName) || m(x.charId));
            return nameHit || assocHit;
        });
        return sortRows(list, sortOrder);
    })();

    const filteredCharacters = (() => {
        if (view === 'matched') return [];
        const list = availableCharacters.filter((c) => !q || m(c.name) || m(c.charId));
        return sortRows(list, sortOrder);
    })();

    if (accounts.length === 0) {
        return (
            <>
                <div className="mt-16 text-center">
                    <p className="text-body text-ink-2">No accounts found.</p>
                    <p className="mt-1 text-meta text-ink-3">EVE settings directories will appear here once detected.</p>
                </div>
                {confirmDialog}
            </>
        );
    }

    return (
        <>
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                <section aria-label="User files">
                    <div className="mb-2 flex items-baseline justify-between">
                        <h2 className="text-h3 text-ink-1 uppercase tracking-wide">User files</h2>
                        <span className="text-meta text-ink-3 tabular">
                            {filteredAccounts.length} of {accounts.length}
                        </span>
                    </div>
                    {filteredAccounts.length === 0 ? (
                        <div className="rounded-lg border border-rule-1 bg-surface-1 px-4 py-8 text-center text-meta text-ink-3">
                            No user files match the current filter.
                        </div>
                    ) : (
                        <div className="flex flex-col gap-3">
                            {filteredAccounts.map((mapping) => (
                                <MapAccountCard
                                    key={`${mapping.userId}-${mapping.mtime}`}
                                    mapping={mapping}
                                    associations={associations}
                                    handleUnassociate={handleUnassociate}
                                    handleDrop={handleDrop}
                                    mtimeToColor={mtimeToColor}
                                />
                            ))}
                        </div>
                    )}
                </section>

                {view !== 'matched' ? (
                    <section aria-label="Unassociated characters">
                        <div className="mb-2 flex items-baseline justify-between">
                            <h2 className="text-h3 text-ink-1 uppercase tracking-wide">Unassociated characters</h2>
                            <span className="text-meta text-ink-3 tabular">{filteredCharacters.length}</span>
                        </div>
                        {availableCharacters.length === 0 ? (
                            <div className="rounded-lg border border-rule-1 bg-surface-1 px-4 py-8 text-center text-meta text-ink-3">
                                All characters mapped.
                            </div>
                        ) : filteredCharacters.length === 0 ? (
                            <div className="rounded-lg border border-rule-1 bg-surface-1 px-4 py-8 text-center text-meta text-ink-3">
                                No characters match the current filter.
                            </div>
                        ) : (
                            <div role="list" data-testid="available-characters" className="rounded-lg border border-rule-1 bg-surface-1 overflow-hidden">
                                {filteredCharacters.map((char) => (
                                    <MapCharacterCard
                                        key={`${char.charId}-${char.mtime}`}
                                        char={char}
                                        handleDragStart={handleDragStart}
                                        mtimeToColor={mtimeToColor}
                                    />
                                ))}
                            </div>
                        )}
                    </section>
                ) : null}
            </div>
            {confirmDialog}
        </>
    );
};

MappingView.propTypes = {
    subDirs: PropTypes.array.isRequired,
    associations: PropTypes.array.isRequired,
    filter: PropTypes.string.isRequired,
    view: PropTypes.oneOf(['all', 'unmatched', 'matched']).isRequired,
    sortOrder: PropTypes.oneOf(['mtime-desc', 'mtime-asc', 'name-asc', 'name-desc']).isRequired,
};

export default MappingView;
