// Settings — reading layout. Configuring a tool, not enrolling in a service.
// Max-width 65ch for prose, generous whitespace, list-like rows where dense.

import { useEffect, useState } from 'react';
import Subheader from '../components/ui/Subheader.jsx';
import FuzzworksSettings from '../components/settings/FuzzworksSettings.jsx';

const LS_SECTION = 'st.section';

const readLS = (key, fallback) => {
    try {
        const raw = localStorage.getItem(key);
        return raw === null ? fallback : JSON.parse(raw);
    } catch {
        return fallback;
    }
};
const writeLS = (key, val) => {
    try { localStorage.setItem(key, JSON.stringify(val)); } catch { /* swallow */ }
};

const SECTIONS = [
    { id: 'data', label: 'Data sources' },
];

const Settings = () => {
    const [section, setSection] = useState(() => readLS(LS_SECTION, 'data'));
    useEffect(() => writeLS(LS_SECTION, section), [section]);

    return (
        <div className="px-6 pb-12 pt-8 mx-auto max-w-[1024px]">
            <Subheader
                title="Settings"
                meta="Application preferences and data sources"
            />

            <div className="grid grid-cols-1 md:grid-cols-[180px_minmax(0,1fr)] gap-8">
                {/* Section nav. Single section today; structured for growth. */}
                <nav aria-label="Settings sections" className="md:sticky md:top-4 md:self-start">
                    <ul className="flex md:flex-col gap-1 text-meta">
                        {SECTIONS.map((s) => {
                            const active = s.id === section;
                            return (
                                <li key={s.id}>
                                    <button
                                        type="button"
                                        onClick={() => setSection(s.id)}
                                        aria-current={active ? 'page' : undefined}
                                        className={[
                                            'w-full text-left px-3 h-8 rounded-md transition-colors duration-fast ease-out-quart',
                                            active
                                                ? 'bg-surface-2 text-ink-1 shadow-rail-accent'
                                                : 'text-ink-3 hover:bg-surface-2 hover:text-ink-1',
                                        ].join(' ')}
                                    >
                                        {s.label}
                                    </button>
                                </li>
                            );
                        })}
                    </ul>
                </nav>

                <div className="min-w-0 space-y-10">
                    {section === 'data' ? <FuzzworksSettings /> : null}
                </div>
            </div>
        </div>
    );
};

export default Settings;
