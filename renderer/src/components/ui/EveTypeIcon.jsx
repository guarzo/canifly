import PropTypes from 'prop-types';

/**
 * EveTypeIcon — a small square thumbnail for an EVE type, looked up by name
 * via the conversions map. Falls back to a hairline-bordered placeholder when
 * the name is unknown. Decorative; rendered with aria-hidden.
 */
const EveTypeIcon = ({ name, conversions, size = 6, className = '' }) => {
    const typeID = conversions?.[name];
    const url = typeID ? `https://images.evetech.net/types/${typeID}/icon` : null;
    const sizeClass = `h-${size} w-${size}`;
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
    size: PropTypes.number,
    className: PropTypes.string,
};

export default EveTypeIcon;
