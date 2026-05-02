// src/components/character-overview/utils.js
//
// Pure helpers shared by the page, hook, and subcomponents.
// localStorage keys, formatters, and queue/status derivations.

export const LS = {
    view: 'co.view',
    sort: 'co.sort',
    showHidden: 'co.showHidden',
    filter: 'co.filter',
};

export const readLS = (key, fallback) => {
    try {
        const raw = localStorage.getItem(key);
        return raw === null ? fallback : JSON.parse(raw);
    } catch {
        return fallback;
    }
};

export const writeLS = (key, val) => {
    try { localStorage.setItem(key, JSON.stringify(val)); } catch { /* swallow */ }
};

export const formatSP = (sp) => {
    if (!sp || sp <= 0) return '—';
    if (sp >= 1_000_000_000) return `${(sp / 1_000_000_000).toFixed(2)}B`;
    if (sp >= 1_000_000) return `${(sp / 1_000_000).toFixed(1)}M`;
    if (sp >= 1_000) return `${(sp / 1_000).toFixed(0)}K`;
    return String(sp);
};

export const formatDuration = (ms) => {
    if (ms == null || ms <= 0) return '—';
    const totalMin = Math.floor(ms / 60_000);
    const days = Math.floor(totalMin / (60 * 24));
    const hours = Math.floor((totalMin % (60 * 24)) / 60);
    const mins = totalMin % 60;
    if (days > 0) return `${days}d ${String(hours).padStart(2, '0')}h`;
    if (hours > 0) return `${hours}h ${String(mins).padStart(2, '0')}m`;
    return `${mins}m`;
};

// Derive queue ETA from the skill queue's last finish_date.
export const deriveQueueEta = (character) => {
    const queue = character?.Character?.SkillQueue;
    if (!Array.isArray(queue) || queue.length === 0) return null;
    const last = queue[queue.length - 1];
    if (!last?.finish_date) return null;
    const ms = new Date(last.finish_date).getTime() - Date.now();
    return ms > 0 ? ms : null;
};

export const deriveStatus = (character) => {
    const queue = character?.Character?.SkillQueue;
    const hasQueue = Array.isArray(queue) && queue.length > 0;
    if (!hasQueue) return 'idle';
    if (character.MCT) return 'training';
    return 'queued';
};
