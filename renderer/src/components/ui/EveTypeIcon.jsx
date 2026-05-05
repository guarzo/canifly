import PropTypes from 'prop-types';

/**
 * EveTypeIcon — a small square thumbnail for an EVE type, looked up by name
 * via the conversions map. Falls back to a hairline-bordered placeholder when
 * the name is unknown. Decorative; rendered with aria-hidden.
 *
 * Size must be one of the keys in SIZE_CLASS so Tailwind picks up the
 * h-N/w-N utilities at build time. Adding a new size means adding it here.
 */
const SIZE_CLASS = {
    4: 'h-4 w-4',
    5: 'h-5 w-5',
    6: 'h-6 w-6',
    8: 'h-8 w-8',
};

const EveTypeIcon = ({ name, conversions, size = 6, className = '' }) => {
    const typeID = conversions?.[name];
    const url = typeID ? `https://images.evetech.net/types/${typeID}/icon` : null;
    const sizeClass = SIZE_CLASS[size] || SIZE_CLASS[6];
    if (url) {
        return (
            <img
                src={url}
                alt=""
                aria-hidden
                loading="lazy"
                className={`${sizeClass} rounded-sm border border-rule-1 shrink-0 ${className}`}
            />
        );
    }
    return (
        <span
            aria-hidden
            className={`${sizeClass} rounded-sm border border-rule-1 bg-surface-2 shrink-0 ${className}`}
        />
    );
};

EveTypeIcon.propTypes = {
    name: PropTypes.string.isRequired,
    conversions: PropTypes.object,
    size: PropTypes.oneOf(Object.keys(SIZE_CLASS).map(Number)),
    className: PropTypes.string,
};

export default EveTypeIcon;
