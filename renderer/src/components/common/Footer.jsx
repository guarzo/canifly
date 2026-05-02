// Footer — minimal. A hairline rule, version, and a couple of links.
// No glass, no glow, no animation. text-micro sizing per DESIGN.md.

import { version as pkgVersion } from '../../../../package.json';

const openExternal = (url) => (e) => {
    e.preventDefault();
    if (window.electronAPI?.openExternal) window.electronAPI.openExternal(url);
    else window.open(url, '_blank', 'noopener,noreferrer');
};

const Footer = () => (
    <footer
        className="mt-auto border-t border-rule-1 bg-surface-0"
        style={{ WebkitAppRegion: 'drag' }}
    >
        <div
            className="mx-auto flex max-w-[1280px] items-center justify-between gap-4 px-6 py-2 text-micro text-ink-3 tabular"
            style={{ WebkitAppRegion: 'no-drag' }}
        >
            <span className="font-mono">CanIFly v{pkgVersion}</span>
            <nav className="flex items-center gap-4">
                <a
                    href="https://github.com/guarzo/canifly"
                    onClick={openExternal('https://github.com/guarzo/canifly')}
                    className="hover:text-ink-1"
                >
                    GitHub
                </a>
                <a
                    href="https://developers.eveonline.com"
                    onClick={openExternal('https://developers.eveonline.com')}
                    className="hover:text-ink-1"
                >
                    EVE Developers
                </a>
            </nav>
        </div>
    </footer>
);

export default Footer;
