import { cloneElement, useRef, useState } from 'react';
import { Popover } from '@mui/material';

const MissingSkillsPopover = ({ planName, missingSkills, trigger }) => {
    const [open, setOpen] = useState(false);
    const anchorRef = useRef(null);
    const skills = Object.entries(missingSkills || {});

    const enhanced = cloneElement(trigger, {
        ref: anchorRef,
        onClick: (e) => {
            trigger.props.onClick?.(e);
            setOpen(true);
        },
        'aria-haspopup': 'dialog',
        'aria-expanded': open,
    });

    return (
        <>
            {enhanced}
            <Popover
                open={open}
                anchorEl={anchorRef.current}
                onClose={() => setOpen(false)}
                anchorOrigin={{ vertical: 'bottom', horizontal: 'left' }}
                slotProps={{
                    paper: {
                        className: 'rounded-md border border-rule-1 bg-surface-1',
                        style: { backgroundImage: 'none' },
                    },
                }}
            >
                <div className="min-w-[220px] max-w-[320px] p-3">
                    <div className="text-meta text-ink-3 uppercase tracking-wide mb-2">{planName}</div>
                    {skills.length === 0 ? (
                        <div className="text-meta text-ink-3 italic">No missing skills.</div>
                    ) : (
                        <ul role="list" className="flex flex-col gap-1">
                            {skills.map(([name, level]) => (
                                <li key={name} className="grid grid-cols-[1fr_auto] gap-3 items-center">
                                    <span className="text-body text-ink-1 truncate">{name}</span>
                                    <span className="text-meta font-mono tabular text-ink-3">{level}</span>
                                </li>
                            ))}
                        </ul>
                    )}
                </div>
            </Popover>
        </>
    );
};

export default MissingSkillsPopover;
