import PropTypes from 'prop-types';

/**
 * MtimeSwatch — a small filled square that color-codes a file by its
 * (rounded-to-the-minute) mtime bucket. Two files written in the same minute
 * share a color; that's the whole signal users rely on when pairing a char
 * file with the user file it was saved alongside.
 *
 * Not a side-stripe. Not a hover affordance. Just a token rendered inline at
 * the start of a row.
 */
const MtimeSwatch = ({ color, title }) => (
    <span
        role="img"
        aria-label={title || 'mtime group'}
        title={title}
        className="inline-block h-2.5 w-2.5 rounded-sm border border-rule-1"
        style={{ backgroundColor: color || 'var(--surface-3)' }}
    />
);

MtimeSwatch.propTypes = {
    color: PropTypes.string,
    title: PropTypes.string,
};

export default MtimeSwatch;
