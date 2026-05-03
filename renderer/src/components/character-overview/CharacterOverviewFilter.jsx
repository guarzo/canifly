// src/components/character-overview/CharacterOverviewFilter.jsx
//
// Filter input bar with the slash-key hint and match count.
import PropTypes from 'prop-types';
import { Search as SearchIcon } from '@mui/icons-material';
import Kbd from '../ui/Kbd.jsx';

const CharacterOverviewFilter = ({ filter, onFilterChange, onClear, matchCount, filterRef }) => (
    <div className="mb-4 flex items-center gap-3">
        <div className="flex-1 max-w-md flex items-center gap-2 rounded-md border border-rule-1 bg-surface-1 px-2.5 py-1.5 focus-within:border-accent">
            <SearchIcon fontSize="small" sx={{ color: 'var(--ink-3)' }} />
            <input
                ref={filterRef}
                value={filter}
                onChange={(e) => onFilterChange(e.target.value)}
                placeholder="Filter by name, location, role, or ship"
                aria-label="Filter characters"
                className="bg-transparent flex-1 outline-hidden text-body text-ink-1 placeholder:text-ink-3"
            />
            {filter ? (
                <button
                    type="button"
                    onClick={onClear}
                    className="text-ink-3 hover:text-ink-1 text-meta"
                    aria-label="Clear filter"
                >
                    clear
                </button>
            ) : (
                <Kbd>/</Kbd>
            )}
        </div>
        <span className="text-meta text-ink-3 tabular">
            {filter ? `${matchCount} match${matchCount === 1 ? '' : 'es'}` : null}
        </span>
    </div>
);

CharacterOverviewFilter.propTypes = {
    filter: PropTypes.string.isRequired,
    onFilterChange: PropTypes.func.isRequired,
    onClear: PropTypes.func.isRequired,
    matchCount: PropTypes.number.isRequired,
    filterRef: PropTypes.object,
};

export default CharacterOverviewFilter;
