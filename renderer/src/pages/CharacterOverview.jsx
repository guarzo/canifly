// src/pages/CharacterOverview.jsx
//
// Page-level composition for the table-first character overview. State,
// derivations, action handlers, and the keyboard-shortcut handler all
// live in `useCharacterOverview`. Visual sections live in
// renderer/src/components/character-overview/*.
import PropTypes from 'prop-types';
import Subheader from '../components/ui/Subheader.jsx';
import Kbd from '../components/ui/Kbd.jsx';
import { overviewInstructions } from '../utils/instructions.jsx';
import { useCharacterOverview } from '../hooks/useCharacterOverview';
import CharacterOverviewToolbar from '../components/character-overview/CharacterOverviewToolbar.jsx';
import CharacterOverviewFilter from '../components/character-overview/CharacterOverviewFilter.jsx';
import GroupBlock from '../components/character-overview/GroupBlock.jsx';
import { formatSP } from '../components/character-overview/utils';

const CharacterOverview = ({ roles = [], skillConversions = {} }) => {
    const {
        filterRef, refreshRef,
        accounts,
        view, setView,
        sortOrder, setSortOrder,
        showHidden, setShowHidden,
        filter, setFilter,
        isRefreshing,
        expanded,
        focusedRowId, setFocusedRowId,
        filteredCharacters,
        groups,
        summary,
        toggleExpanded,
        handleRefreshAll,
        handleUpdateCharacter,
        handleRemoveCharacter,
        handleUpdateAccount,
        handleRemoveAccount,
    } = useCharacterOverview({ roles });

    const meta = `${summary.charCount} characters across ${summary.accountCount} ${
        summary.accountCount === 1 ? 'account' : 'accounts'
    } · ${formatSP(summary.totalSp)} SP total`;

    return (
        <div className="px-6 pb-12 pt-8 max-w-[1280px] mx-auto">
            <Subheader
                title="Character Overview"
                meta={meta}
                actions={
                    <CharacterOverviewToolbar
                        view={view}
                        onViewChange={setView}
                        sortOrder={sortOrder}
                        onToggleSortOrder={() => setSortOrder((s) => s === 'asc' ? 'desc' : 'asc')}
                        showHidden={showHidden}
                        onToggleShowHidden={() => setShowHidden((v) => !v)}
                        isRefreshing={isRefreshing}
                        onRefreshAll={handleRefreshAll}
                        refreshRef={refreshRef}
                    />
                }
                tip={overviewInstructions}
            />

            <CharacterOverviewFilter
                filter={filter}
                onFilterChange={setFilter}
                onClear={() => setFilter('')}
                matchCount={filteredCharacters.length}
                filterRef={filterRef}
            />

            {accounts.length === 0 ? (
                <div className="mt-16 text-center">
                    <p className="text-body text-ink-2">No accounts found.</p>
                    <p className="mt-1 text-meta text-ink-3">Add an EVE account from the header to get started.</p>
                </div>
            ) : filteredCharacters.length === 0 ? (
                <div className="mt-16 text-center">
                    <p className="text-body text-ink-2">
                        No characters match <span className="font-mono text-ink-1">{filter}</span>.
                    </p>
                    <button
                        type="button"
                        onClick={() => setFilter('')}
                        className="mt-2 text-meta text-accent hover:text-accent-strong"
                    >
                        Clear filter (Esc)
                    </button>
                </div>
            ) : (
                <div role="table" aria-label="Characters" className="rounded-lg border border-rule-1 bg-surface-1 overflow-hidden">
                    <div
                        role="row"
                        className="hidden sm:grid grid-cols-[28px_minmax(180px,1.4fr)_1fr_84px_92px_minmax(140px,1fr)_minmax(120px,0.8fr)_28px] gap-3 px-4 py-2 text-meta text-ink-3 border-b border-rule-1 bg-surface-2"
                    >
                        <div role="columnheader" aria-label="Status" />
                        <div role="columnheader">Character</div>
                        <div role="columnheader">Account</div>
                        <div role="columnheader" className="text-right">SP</div>
                        <div role="columnheader" className="text-right">Queue</div>
                        <div role="columnheader">Location</div>
                        <div role="columnheader">Role</div>
                        <div role="columnheader" aria-label="Actions" />
                    </div>
                    {groups.map((group) => (
                        <GroupBlock
                            key={group.key}
                            groupKey={group.key}
                            view={view}
                            characters={group.characters}
                            roles={roles}
                            skillConversions={skillConversions}
                            expanded={expanded}
                            onToggleExpand={toggleExpanded}
                            focusedRowId={focusedRowId}
                            setFocusedRowId={setFocusedRowId}
                            onUpdateCharacter={handleUpdateCharacter}
                            onRemoveCharacter={handleRemoveCharacter}
                            onUpdateAccount={handleUpdateAccount}
                            onRemoveAccount={handleRemoveAccount}
                            allAccounts={accounts}
                        />
                    ))}
                </div>
            )}

            <div className="mt-4 flex flex-wrap items-center gap-x-4 gap-y-1 text-micro text-ink-3 font-mono">
                <span><Kbd>/</Kbd> filter</span>
                <span><Kbd>j</Kbd>/<Kbd>k</Kbd> move</span>
                <span><Kbd>Enter</Kbd> expand</span>
                <span><Kbd>Esc</Kbd> collapse</span>
                <span><Kbd>g</Kbd> then <Kbd>a</Kbd>/<Kbd>r</Kbd>/<Kbd>l</Kbd> group</span>
                <span><Kbd>r</Kbd> refresh</span>
            </div>
        </div>
    );
};

CharacterOverview.propTypes = {
    roles: PropTypes.array,
    skillConversions: PropTypes.object,
};

export default CharacterOverview;
