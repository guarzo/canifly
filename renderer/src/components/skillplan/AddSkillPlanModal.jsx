// AddSkillPlanModal — design-system rewrite. Plain dialog: surface tiers,
// rule hairlines, accent button. No glass, gradient, or motion.

import { useState } from 'react';
import PropTypes from 'prop-types';
import { Dialog } from '@mui/material';
import { useAsyncOperation } from '../../hooks/useAsyncOperation';
import { createSkillPlan } from '../../api/skillPlansApi';

const AddSkillPlanModal = ({ onClose }) => {
    const [planName, setPlanName] = useState('');
    const [planContents, setPlanContents] = useState('');
    const { execute } = useAsyncOperation();

    const handleSubmit = async (e) => {
        e.preventDefault();
        await execute(
            () => createSkillPlan(planName.trim(), planContents.trim()),
            { successMessage: 'Skill plan created', onSuccess: onClose },
        );
    };

    return (
        <Dialog
            open
            onClose={onClose}
            maxWidth="sm"
            fullWidth
            slotProps={{
                paper: {
                    className: 'bg-surface-1 border border-rule-1 rounded-lg',
                    sx: { backgroundImage: 'none' },
                },
            }}
        >
            <form onSubmit={handleSubmit}>
                <header className="px-5 py-4 border-b border-rule-1">
                    <h2 className="text-h3 text-ink-1">Add skill plan</h2>
                    <p className="mt-1 text-meta text-ink-3">
                        Paste skills as <span className="font-mono">SkillName Level</span>, one per line.
                    </p>
                </header>

                <div className="px-5 py-4 space-y-4 bg-surface-1">
                    <label className="block">
                        <span className="text-meta text-ink-3">Plan name</span>
                        <input
                            value={planName}
                            onChange={(e) => setPlanName(e.target.value)}
                            placeholder="e.g. Battleship Mastery"
                            required
                            autoFocus
                            className="mt-1 w-full h-9 px-2.5 rounded-md border border-rule-1 bg-surface-0 text-body text-ink-1 placeholder:text-ink-3 focus:outline-hidden focus:border-accent"
                        />
                    </label>
                    <label className="block">
                        <span className="text-meta text-ink-3">Skill plan contents</span>
                        <textarea
                            value={planContents}
                            onChange={(e) => setPlanContents(e.target.value)}
                            placeholder={'Gunnery 5\nLarge Hybrid Turret 4'}
                            required
                            rows={8}
                            className="mt-1 w-full px-2.5 py-2 rounded-md border border-rule-1 bg-surface-0 font-mono text-body text-ink-1 placeholder:text-ink-3 focus:outline-hidden focus:border-accent"
                        />
                    </label>
                </div>

                <footer className="px-5 py-3 border-t border-rule-1 flex items-center justify-end gap-2">
                    <button
                        type="button"
                        onClick={onClose}
                        className="h-8 px-3 rounded-md border border-rule-1 text-meta text-ink-2 hover:bg-surface-2 hover:text-ink-1"
                    >
                        Cancel
                    </button>
                    <button
                        type="submit"
                        className="h-8 px-3 rounded-md bg-accent text-accent-ink text-meta font-medium hover:bg-accent-strong"
                    >
                        Save
                    </button>
                </footer>
            </form>
        </Dialog>
    );
};

AddSkillPlanModal.propTypes = {
    onClose: PropTypes.func.isRequired,
};

export default AddSkillPlanModal;
