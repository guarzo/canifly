//
// The sticky-scroll container layout used by the character × plan matrix.
// Pulled out of PlanMatrix so the readiness layer is the only thing that
// changes per-cell and so this shell can host alternate dense matrices in
// the future.

const FROZEN_W = 220 + 72 + 92; // nameW + spW + summaryW from the original
const COL_W = 36;

const MatrixShell = ({ ariaLabel, frozenHeader, columnHeaders, rows, frozenWidth = FROZEN_W, colWidth = COL_W }) => (
    <div
        role="table"
        aria-label={ariaLabel}
        className="rounded-lg border border-rule-1 bg-surface-1 overflow-auto max-h-[calc(100vh-220px)]"
    >
        <div role="row" className="sticky top-0 z-20 flex items-end bg-surface-2 border-b border-rule-1">
            <div
                role="columnheader"
                className="sticky left-0 z-30 bg-surface-2 border-r border-rule-1 flex items-end gap-3 px-4 py-2"
                style={{ width: frozenWidth }}
            >
                {frozenHeader}
            </div>
            {columnHeaders.map((ch, i) => (
                <div
                    key={i}
                    role="columnheader"
                    className="shrink-0 flex flex-col items-center justify-end gap-1.5 px-1 pb-2 pt-3 border-r border-rule-1 last:border-r-0"
                    style={{ width: colWidth, minHeight: 132 }}
                >
                    {ch}
                </div>
            ))}
        </div>

        {rows.map((row, rIdx) => {
            const isLast = rIdx === rows.length - 1;
            return (
                <div
                    role="row"
                    key={row.key}
                    className={`flex items-stretch ${isLast ? '' : 'border-b border-rule-1'} hover:bg-surface-2`}
                >
                    <div
                        role="rowheader"
                        className="sticky left-0 z-10 bg-surface-1 border-r border-rule-1 flex items-center gap-3 px-4 h-10"
                        style={{ width: frozenWidth }}
                    >
                        {row.frozenCell}
                    </div>
                    {row.cells.map((cell, i) => (
                        <div
                            key={i}
                            role="cell"
                            className="shrink-0 flex items-center justify-center border-r border-rule-1 last:border-r-0 h-10"
                            style={{ width: colWidth }}
                        >
                            {cell}
                        </div>
                    ))}
                </div>
            );
        })}
    </div>
);

export default MatrixShell;
