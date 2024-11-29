import React, { useEffect, useRef, forwardRef, useImperativeHandle, memo } from "react";
import { TabulatorFull as Tabulator } from "tabulator-tables";
import "tabulator-tables/dist/css/tabulator.min.css";
import PropTypes from "prop-types";

const TabulatorTable = forwardRef(({ data, columns, options = {}, onCellClick = null }, ref) => {
    const tableRef = useRef(null);
    const tableInstance = useRef(null);
    const pendingData = useRef(null);

    useEffect(() => {
        if (!tableRef.current) {
            console.error("Table reference not found.");
            return;
        }

        console.log("Initializing Tabulator instance...");
        tableInstance.current = new Tabulator(tableRef.current, {
            data: [],
            columns,
            layout: "fitColumns",
            autoResize: true,
            ...options,
        });

        tableInstance.current.on("tableBuilt", () => {
            console.log("Tabulator table built.");
            if (pendingData.current) {
                console.log("Processing pending data update.");
                tableInstance.current.replaceData(pendingData.current).catch((err) =>
                    console.error("Error replacing data after initialization:", err)
                );
                pendingData.current = null;
            }
        });

        if (onCellClick) {
            tableInstance.current.on("cellClick", (e, cell) => onCellClick(e, cell));
        }

        return () => {
            if (tableInstance.current) {
                console.log("Destroying Tabulator instance...");
                tableInstance.current.destroy();
                tableInstance.current = null;
            }
        };
    }, [columns, options]);

    useEffect(() => {
        if (tableInstance.current?.initialized) {
            console.log("Updating Tabulator data...");
            tableInstance.current.replaceData(data).catch((err) =>
                console.error("Error replacing data:", err)
            );
        } else {
            console.log("Table not ready. Queuing data update.");
            pendingData.current = data;
        }
    }, [data]);

    useEffect(() => {
        const handleUpdate = (event) => {
            const { planName } = event.detail;
            console.log("Delete event received for:", planName);
            const rowDeleted = tableInstance.current?.deleteRow(planName);
            if (!rowDeleted) {
                console.error(`Delete Error - No matching row found: ${planName}`);
            }
        };
        document.addEventListener("updateSkillPlanTable", handleUpdate);
        return () => {
            document.removeEventListener("updateSkillPlanTable", handleUpdate);
        };
    }, []);

    useImperativeHandle(ref, () => ({
        deleteRow: (rowId) => {
            if (tableInstance.current) {
                console.log(`Deleting row with ID: ${rowId}`);
                const success = tableInstance.current.deleteRow(rowId);
                if (!success) console.error(`Delete Error - No matching row found: ${rowId}`);
            }
        },
    }));

    return <div ref={tableRef} style={{ width: "100%" }} />;
});

TabulatorTable.propTypes = {
    data: PropTypes.array.isRequired,
    columns: PropTypes.array.isRequired,
    options: PropTypes.object,
    onCellClick: PropTypes.func,
};

export default memo(TabulatorTable);
