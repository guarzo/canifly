import React, { useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import { TabulatorFull as Tabulator } from 'tabulator-tables';
import 'tabulator-tables/dist/css/tabulator_midnight.min.css';

const TabulatorTable = ({ id, data, columns, options }) => {
    const tableRef = useRef(null);
    const tabulatorRef = useRef(null);

    useEffect(() => {
        if (tableRef.current && !tabulatorRef.current) {
            // Initialize Tabulator only once
            tabulatorRef.current = new Tabulator(tableRef.current, {
                data,
                columns,
                layout: options.layout || 'fitColumns',
                ...options,
            });
        } else if (tabulatorRef.current) {
            // Update the data without destroying/reinitializing the table
            tabulatorRef.current.replaceData(data);
        }

        return () => {
            if (tabulatorRef.current) {
                tabulatorRef.current.destroy();
                tabulatorRef.current = null;
            }
        };
    }, [data, columns, options]);

    return <div id={id} ref={tableRef}></div>;
};

TabulatorTable.propTypes = {
    id: PropTypes.string.isRequired,
    data: PropTypes.array.isRequired,
    columns: PropTypes.array.isRequired,
    options: PropTypes.object,
};

export default TabulatorTable;
