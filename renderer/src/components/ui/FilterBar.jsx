import { Search as SearchIcon } from '@mui/icons-material';
import Kbd from './Kbd.jsx';

/**
 * FilterBar — the standard `/`-shortcut filter row used on dense list/table
 * pages. Caller owns the value, the setter, and the ref (so it can wire up
 * its own keyboard shortcut handler). The component renders the input, the
 * search icon, the clear button when non-empty, the `/` Kbd hint when empty,
 * and an optional right-aligned count slot.
 */
const FilterBar = ({
    inputRef,
    value,
    onChange,
    placeholder,
    ariaLabel,
    count,
    className = '',
}) => (
    <div className={`mb-4 flex items-center gap-3 ${className}`}>
        <div className="flex-1 max-w-md flex items-center gap-2 rounded-md border border-rule-1 bg-surface-1 px-2.5 py-1.5 focus-within:border-accent">
            <SearchIcon fontSize="small" sx={{ color: 'var(--ink-3)' }} />
            <input
                ref={inputRef}
                value={value}
                onChange={(e) => onChange(e.target.value)}
                placeholder={placeholder}
                aria-label={ariaLabel || placeholder}
                className="bg-transparent flex-1 outline-hidden text-body text-ink-1 placeholder:text-ink-3"
            />
            {value ? (
                <button
                    type="button"
                    onClick={() => onChange('')}
                    className="text-ink-3 hover:text-ink-1 text-meta"
                    aria-label="Clear filter"
                >
                    clear
                </button>
            ) : (
                <Kbd>/</Kbd>
            )}
        </div>
        {count != null ? (
            <span className="text-meta text-ink-3 tabular">{count}</span>
        ) : null}
    </div>
);

export default FilterBar;
