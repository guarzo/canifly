// src/hooks/useCharacterOverview.js
//
// Owns the Character Overview page's state, persisted UI toggles,
// derivations (visibility, filter, grouping, summary), action handlers,
// and the global keyboard-shortcut handler. The page is a pure renderer
// over what this hook returns.
import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { useAppData } from './useAppData';
import { useAsyncOperation } from './useAsyncOperation';
import { updateCharacter, deleteCharacter, refreshCharacter } from '../api/accountsApi';
import { logger } from '../utils/logger';
import { LS, readLS, writeLS } from '../components/character-overview/utils';

export function useCharacterOverview({ roles }) {
    const { accounts = [], updateAccount, deleteAccount, fetchAccounts } = useAppData();
    const { execute } = useAsyncOperation();

    // Persisted UI state.
    const [view, setView] = useState(() => readLS(LS.view, 'account'));
    const [sortOrder, setSortOrder] = useState(() => readLS(LS.sort, 'asc'));
    const [showHidden, setShowHidden] = useState(() => readLS(LS.showHidden, false));
    const [filter, setFilter] = useState(() => readLS(LS.filter, ''));
    const [isRefreshing, setIsRefreshing] = useState(false);
    const [expanded, setExpanded] = useState(() => new Set());
    const [focusedRowId, setFocusedRowId] = useState(null);

    useEffect(() => writeLS(LS.view, view), [view]);
    useEffect(() => writeLS(LS.sort, sortOrder), [sortOrder]);
    useEffect(() => writeLS(LS.showHidden, showHidden), [showHidden]);
    useEffect(() => writeLS(LS.filter, filter), [filter]);

    const filterRef = useRef(null);
    const refreshRef = useRef(null);

    const visibleAccounts = useMemo(() => {
        if (!accounts) return [];
        return showHidden ? accounts : accounts.filter((a) => a.Visible !== false);
    }, [accounts, showHidden]);

    const allCharacters = useMemo(() => {
        const out = [];
        for (const acc of visibleAccounts) {
            const accountName = acc.Name || 'Unknown Account';
            for (const ch of acc.Characters || []) {
                out.push({
                    ...ch,
                    accountId: acc.ID,
                    accountName,
                    accountVisible: acc.Visible !== false,
                });
            }
        }
        return out;
    }, [visibleAccounts]);

    const filteredCharacters = useMemo(() => {
        const q = filter.trim().toLowerCase();
        if (!q) return allCharacters;
        return allCharacters.filter((ch) => {
            const c = ch.Character || {};
            const name = (c.CharacterName || '').toLowerCase();
            const loc = (c.LocationName || '').toLowerCase();
            const role = (ch.Role || '').toLowerCase();
            const ship = (c.ShipTypeName || c.CurrentShip || '').toLowerCase();
            return name.includes(q) || loc.includes(q) || role.includes(q) || ship.includes(q);
        });
    }, [allCharacters, filter]);

    const groups = useMemo(() => {
        const map = new Map();
        // Group-by-account uses a stable key so two accounts that happen to
        // share a display name don't collapse into one bucket. The display
        // name (`accountName`) is derived from the key when needed.
        const accountKey = (ch) => (
            ch.accountId != null ? `${ch.accountId}::${ch.accountName || 'Unknown Account'}`
                                  : (ch.accountName || 'Unknown Account')
        );
        if (view === 'account') {
            for (const ch of filteredCharacters) {
                const key = accountKey(ch);
                if (!map.has(key)) map.set(key, []);
                map.get(key).push(ch);
            }
        } else if (view === 'role') {
            for (const r of roles) map.set(r, []);
            map.set('Unassigned', []);
            for (const ch of filteredCharacters) {
                const key = ch.Role || 'Unassigned';
                if (!map.has(key)) map.set(key, []);
                map.get(key).push(ch);
            }
        } else {
            for (const ch of filteredCharacters) {
                const key = ch.Character?.LocationName || 'Unknown Location';
                if (!map.has(key)) map.set(key, []);
                map.get(key).push(ch);
            }
        }
        let keys = [...map.keys()];
        if (view === 'role') keys = keys.filter((k) => (map.get(k) || []).length > 0);
        // Sort by the human-facing label, not the storage key. In account
        // view the key is `${accountId}::${accountName}`, so sorting on the
        // raw key would order by id; sort on the display label instead.
        const labelFor = (key) => {
            if (view !== 'account') return key;
            const chars = map.get(key) || [];
            return chars[0]?.accountName || key;
        };
        keys.sort((a, b) => {
            const la = labelFor(a);
            const lb = labelFor(b);
            return sortOrder === 'asc' ? la.localeCompare(lb) : lb.localeCompare(la);
        });
        return keys.map((k) => ({
            key: k,
            characters: [...(map.get(k) || [])].sort((a, b) =>
                (a.Character?.CharacterName || '').localeCompare(b.Character?.CharacterName || ''),
            ),
        }));
    }, [filteredCharacters, view, sortOrder, roles]);

    const summary = useMemo(() => {
        const charCount = allCharacters.length;
        const accountCount = visibleAccounts.length;
        const totalSp = allCharacters.reduce(
            (sum, ch) => sum + (ch.Character?.CharacterSkillsResponse?.total_sp || 0), 0,
        );
        return { charCount, accountCount, totalSp };
    }, [allCharacters, visibleAccounts]);

    const toggleExpanded = useCallback((id) => {
        setExpanded((prev) => {
            const next = new Set(prev);
            if (next.has(id)) next.delete(id); else next.add(id);
            return next;
        });
    }, []);

    const handleRefreshAll = useCallback(async () => {
        setIsRefreshing(true);
        try {
            for (const acc of accounts || []) {
                for (const ch of acc.Characters || []) {
                    try {
                        await refreshCharacter(ch.Character.CharacterID);
                    } catch (err) {
                        logger.error(`Refresh failed for ${ch.Character.CharacterName}:`, err);
                    }
                }
            }
            await fetchAccounts();
        } finally {
            setIsRefreshing(false);
        }
    }, [accounts, fetchAccounts]);

    const handleUpdateCharacter = useCallback(async (characterId, updates) => {
        await execute(
            () => updateCharacter(characterId, updates),
            { successMessage: 'Character updated' },
        );
    }, [execute]);

    const handleRemoveCharacter = useCallback(async (characterId) => {
        await execute(
            () => deleteCharacter(characterId),
            { successMessage: 'Character removed' },
        );
    }, [execute]);

    const handleUpdateAccount = useCallback((accountId, updates) =>
        execute(() => updateAccount(accountId, updates), { successMessage: 'Account updated' }),
    [execute, updateAccount]);

    const handleRemoveAccount = useCallback((accountId) =>
        execute(() => deleteAccount(accountId), { successMessage: 'Account removed' }),
    [execute, deleteAccount]);

    // ── Keyboard shortcuts ────────────────────────────────────────────────
    const flatRowIds = useMemo(
        () => groups.flatMap((g) => g.characters.map((c) => c.Character.CharacterID)),
        [groups],
    );

    useEffect(() => {
        const onKey = (e) => {
            const target = e.target;
            const tag = (target?.tagName || '').toLowerCase();
            const isTyping = tag === 'input' || tag === 'textarea' || target?.isContentEditable;
            // Slash always focuses filter, even from inputs.
            if (e.key === '/' && !isTyping) {
                e.preventDefault();
                filterRef.current?.focus();
                return;
            }
            if (isTyping) return;
            // Ignore shortcuts while focus is on a button, inside an open
            // menu/dialog, or on a link — those controls handle their own
            // keys. Rows (role="row", tabIndex=0) and the document body are
            // still allowed so j/k navigation keeps working.
            const isInNestedControl = typeof target?.closest === 'function' && (
                target.closest('button') ||
                target.closest('[role="menu"]') ||
                target.closest('[role="menuitem"]') ||
                target.closest('[role="dialog"]') ||
                target.closest('a[href]')
            );
            if (isInNestedControl) return;
            if (e.key === 'j' || e.key === 'ArrowDown') {
                e.preventDefault();
                setFocusedRowId((cur) => {
                    const idx = flatRowIds.indexOf(cur);
                    return flatRowIds[Math.min(flatRowIds.length - 1, idx < 0 ? 0 : idx + 1)] ?? cur;
                });
            } else if (e.key === 'k' || e.key === 'ArrowUp') {
                e.preventDefault();
                setFocusedRowId((cur) => {
                    const idx = flatRowIds.indexOf(cur);
                    return flatRowIds[Math.max(0, idx <= 0 ? 0 : idx - 1)] ?? cur;
                });
            } else if (e.key === 'Enter' && focusedRowId != null) {
                e.preventDefault();
                toggleExpanded(focusedRowId);
            } else if (e.key === 'Escape') {
                // Escape clears the filter first (matching the empty-state
                // "Clear filter (Esc)" button), then collapses rows on a
                // subsequent press (matching the "Esc collapse" legend).
                if (filter) setFilter('');
                else setExpanded(new Set());
            } else if (e.key === 'r') {
                e.preventDefault();
                refreshRef.current?.click();
            } else if (e.key === 'g') {
                // chord: g a / g r / g l
                const handler = (e2) => {
                    // Suppress the second key so it doesn't bubble to the
                    // page-level shortcut handler (e.g., "g r" must not
                    // also trigger the "r = refresh" shortcut).
                    if (['a', 'r', 'l'].includes(e2.key)) {
                        e2.preventDefault();
                        e2.stopImmediatePropagation();
                    }
                    if (e2.key === 'a') setView('account');
                    else if (e2.key === 'r') setView('role');
                    else if (e2.key === 'l') setView('location');
                    window.removeEventListener('keydown', handler, true);
                };
                window.addEventListener('keydown', handler, true);
                setTimeout(() => window.removeEventListener('keydown', handler, true), 1200);
            }
        };
        window.addEventListener('keydown', onKey);
        return () => window.removeEventListener('keydown', onKey);
    }, [flatRowIds, focusedRowId, toggleExpanded, filter]);

    return {
        // refs
        filterRef, refreshRef,
        // raw
        accounts,
        // ui state
        view, setView,
        sortOrder, setSortOrder,
        showHidden, setShowHidden,
        filter, setFilter,
        isRefreshing,
        expanded,
        focusedRowId, setFocusedRowId,
        // derived
        filteredCharacters,
        groups,
        summary,
        // handlers
        toggleExpanded,
        handleRefreshAll,
        handleUpdateCharacter,
        handleRemoveCharacter,
        handleUpdateAccount,
        handleRemoveAccount,
    };
}
