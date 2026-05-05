import Subheader from './Subheader.jsx';

/**
 * PageShell — standard page wrapper. Owns the outer padding, the max-width
 * column, and the Subheader. Pages render their filter bar (if any) and main
 * body as children. Replaces the `px-6 pb-12 pt-8 max-w-[1280px] mx-auto` +
 * Subheader opening that was duplicated across SkillPlans, Mapping, Sync.
 */
const PageShell = ({ title, meta, actions, tip, children, className = '' }) => (
    <div className={`px-6 pb-12 pt-8 max-w-[1280px] mx-auto ${className}`}>
        <Subheader title={title} meta={meta} actions={actions} tip={tip} />
        {children}
    </div>
);

export default PageShell;
