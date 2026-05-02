// FuzzworksSettings — a settings section presented as a reading panel
// with a list-like row of file statuses underneath. No MUI Card, no glass.

import { useEffect, useState } from 'react';
import { Refresh as RefreshIcon } from '@mui/icons-material';
import Surface from '../ui/Surface.jsx';
import { useAppData } from '../../hooks/useAppData';
import { useAsyncOperation } from '../../hooks/useAsyncOperation';
import { getFuzzworksStatus, updateFuzzworks } from '../../api/apiService';
import { logger } from '../../utils/logger';

const formatDate = (iso) => {
    if (!iso) return 'Never';
    const d = new Date(iso);
    if (Number.isNaN(d.getTime())) return 'Never';
    return d.toLocaleString();
};

const formatSize = (bytes) => {
    if (!bytes) return '—';
    const k = 1024;
    const units = ['B', 'KB', 'MB', 'GB'];
    const i = Math.min(units.length - 1, Math.floor(Math.log(bytes) / Math.log(k)));
    return `${(bytes / Math.pow(k, i)).toFixed(i === 0 ? 0 : 1)} ${units[i]}`;
};

const FILES = [
    { key: 'invTypes', label: 'Item types', file: 'invTypes.csv' },
    { key: 'solarSystems', label: 'Solar systems', file: 'mapSolarSystems.csv' },
];

const FuzzworksSettings = () => {
    const { config, updateConfig } = useAppData();
    const { execute, loading } = useAsyncOperation();
    const [status, setStatus] = useState(null);
    const [statusLoading, setStatusLoading] = useState(false);

    const autoUpdate = config?.AutoUpdateFuzzworks !== false;

    const fetchStatus = async () => {
        setStatusLoading(true);
        try {
            const response = await getFuzzworksStatus();
            setStatus(response.data);
        } catch (err) {
            logger.error('Failed to fetch Fuzzworks status:', err);
        } finally {
            setStatusLoading(false);
        }
    };

    useEffect(() => {
        fetchStatus();
    }, []);

    const handleAutoToggle = async () => {
        await execute(
            () => updateConfig({ AutoUpdateFuzzworks: !autoUpdate }),
            { successMessage: 'Auto-update setting saved' },
        );
    };

    const handleManualUpdate = async () => {
        await execute(
            async () => {
                const response = await updateFuzzworks();
                await fetchStatus();
                return response.data;
            },
            { successMessage: 'Fuzzworks data updated' },
        );
    };

    return (
        <section aria-labelledby="fuzzworks-heading" className="space-y-4">
            <header>
                <h2 id="fuzzworks-heading" className="text-h3 text-ink-1">
                    Fuzzworks data
                </h2>
                <p className="mt-1 max-w-[65ch] text-body text-ink-3">
                    EVE static data — item types and solar systems — fetched from{' '}
                    <span className="font-mono text-ink-2">fuzzwork.co.uk</span>. ETag-aware:
                    only changed files are downloaded.
                </p>
            </header>

            {/* Auto-update row */}
            <Surface tier={1} bordered radius="lg" className="overflow-hidden">
                <label className="flex items-center justify-between gap-4 px-4 h-12 cursor-pointer">
                    <span className="text-body text-ink-1">Auto-update on startup</span>
                    <input
                        type="checkbox"
                        checked={autoUpdate}
                        onChange={handleAutoToggle}
                        disabled={loading}
                        className="h-4 w-4 accent-accent cursor-pointer"
                        aria-label="Auto-update Fuzzworks data on startup"
                    />
                </label>

                <div className="border-t border-rule-1">
                    {/* Column header */}
                    <div className="grid grid-cols-[minmax(180px,1.4fr)_minmax(160px,1fr)_80px] gap-3 px-4 h-8 items-center bg-surface-2 text-meta text-ink-3 border-b border-rule-1">
                        <span>File</span>
                        <span>Last updated</span>
                        <span className="text-right">Size</span>
                    </div>

                    {statusLoading ? (
                        <div className="px-4 py-3 text-meta text-ink-3">Checking…</div>
                    ) : status?.hasData ? (
                        FILES.map((f, idx) => {
                            const s = status[f.key];
                            return (
                                <div
                                    key={f.key}
                                    className={`grid grid-cols-[minmax(180px,1.4fr)_minmax(160px,1fr)_80px] gap-3 px-4 h-10 items-center ${
                                        idx < FILES.length - 1 ? 'border-b border-rule-1' : ''
                                    }`}
                                >
                                    <div className="min-w-0">
                                        <div className="text-body text-ink-1 truncate">{f.label}</div>
                                        <div className="text-micro text-ink-3 font-mono truncate">{f.file}</div>
                                    </div>
                                    <span className="text-meta text-ink-2 tabular truncate">
                                        {s ? formatDate(s.lastUpdated) : '—'}
                                    </span>
                                    <span className="text-meta text-ink-2 font-mono tabular text-right">
                                        {s ? formatSize(s.fileSize) : '—'}
                                    </span>
                                </div>
                            );
                        })
                    ) : (
                        <div className="px-4 py-3 text-meta text-ink-3">
                            No data downloaded yet.
                        </div>
                    )}
                </div>

                <div className="flex items-center justify-end gap-3 px-4 h-12 border-t border-rule-1">
                    <button
                        type="button"
                        onClick={handleManualUpdate}
                        disabled={loading}
                        className="inline-flex items-center gap-1.5 h-8 px-3 rounded-md border border-rule-1 bg-surface-1 text-meta text-ink-1 hover:bg-surface-2 disabled:opacity-50"
                    >
                        <RefreshIcon
                            sx={{
                                fontSize: 16,
                                animation: loading ? 'fz-spin 1s linear infinite' : undefined,
                                '@keyframes fz-spin': { to: { transform: 'rotate(360deg)' } },
                            }}
                        />
                        {loading ? 'Updating…' : 'Update now'}
                    </button>
                </div>
            </Surface>
        </section>
    );
};

export default FuzzworksSettings;
