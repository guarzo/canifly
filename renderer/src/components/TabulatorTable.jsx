import React, { useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import { TabulatorFull as Tabulator } from 'tabulator-tables';
import 'tabulator-tables/dist/css/tabulator_midnight.min.css';

// ForwardRef to properly handle ref passing to the div container
const TabulatorTable = React.forwardRef(({ id, data, columns, options }, ref) => {
    const tableRef = useRef(null); // Local ref for the table DOM element
    const tabulatorRef = useRef(null); // Local ref for the Tabulator instance

    useEffect(() => {
        if (tableRef.current && !tabulatorRef.current) {
            // Initialize Tabulator only once, save it to the ref
            tabulatorRef.current = new Tabulator(tableRef.current, {
                data,
                columns,
                layout: options.layout || 'fitColumns',
                ...options,
            });
        }
    }, [data, columns, options]);

    // Expose the Tabulator instance to the parent component via the forwarded ref
    useEffect(() => {
        if (ref && tabulatorRef.current) {
            ref.current = tabulatorRef.current;
        }
    }, [ref]);

    return (
        <div id={id} ref={tableRef} />
    );
});

TabulatorTable.propTypes = {
    id: PropTypes.string.isRequired,
    data: PropTypes.array.isRequired,
    columns: PropTypes.array.isRequired,
    options: PropTypes.object,
};

export default TabulatorTable;
