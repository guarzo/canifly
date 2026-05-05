// MapCharacterCard.jsx — redesigned per DESIGN.md.
// One dense draggable row representing an unassociated character file.

import PropTypes from 'prop-types';
import { formatDate } from '../../utils/formatter.jsx';
import MtimeSwatch from './MtimeSwatch.jsx';

const CharacterCard = ({ char, handleDragStart, mtimeToColor }) => {
    const swatch = mtimeToColor[char.mtime];

    return (
        <div
            draggable
            onDragStart={(e) => handleDragStart(e, char.charId)}
            role="listitem"
            aria-label={`Drag ${char.name} to associate`}
            className={[
                'grid grid-cols-[16px_minmax(0,1fr)_auto_auto] gap-3 items-center',
                'h-10 px-4',
                'cursor-grab active:cursor-grabbing select-none',
                'border-b border-rule-1 last:border-b-0',
                'bg-surface-1 hover:bg-surface-2',
                'transition-colors duration-fast ease-out-quart',
            ].join(' ')}
        >
            <MtimeSwatch color={swatch} title={`mtime ${char.mtime}`} />
            <span className="text-body text-ink-1 truncate" title={char.name}>
                {char.name}
            </span>
            <span className="text-meta text-ink-3 font-mono tabular">
                {char.charId}
            </span>
            <span className="text-meta text-ink-3 font-mono tabular">
                {formatDate(char.mtime)}
            </span>
        </div>
    );
};

CharacterCard.propTypes = {
    char: PropTypes.shape({
        name: PropTypes.string.isRequired,
        charId: PropTypes.string.isRequired,
        mtime: PropTypes.string.isRequired,
    }).isRequired,
    handleDragStart: PropTypes.func.isRequired,
    mtimeToColor: PropTypes.objectOf(PropTypes.string).isRequired,
};

export default CharacterCard;
