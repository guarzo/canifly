// Dashboard.jsx
import PropTypes from 'prop-types';
import AccountCard from './AccountCard';

const Dashboard = ({
                       accounts,
                       onToggleAccountStatus,
                       onUpdateCharacter,
                       onUpdateAccountName,
                       onRemoveCharacter,
                        onRemoveAccount,
                       roles,
                   }) => {
    return (
        <div className="space-y-6 mt-16">
            {/* Accounts Section */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {accounts.map((account) => (
                    <AccountCard
                        key={account.ID}
                        account={account}
                        onToggleAccountStatus={onToggleAccountStatus}
                        onUpdateAccountName={onUpdateAccountName}
                        onUpdateCharacter={onUpdateCharacter}
                        onRemoveCharacter={onRemoveCharacter}
                        onRemoveAccount={onRemoveAccount}
                        roles={roles}
                    />
                ))}
            </div>
        </div>
    );
};

Dashboard.propTypes = {
    accounts: PropTypes.array.isRequired,
    onToggleAccountStatus: PropTypes.func.isRequired,
    onUpdateCharacter: PropTypes.func.isRequired,
    onUpdateAccountName: PropTypes.func.isRequired,
    onRemoveCharacter: PropTypes.func.isRequired,
    onRemoveAccount: PropTypes.func.isRequired,
    roles: PropTypes.array.isRequired,
};

export default Dashboard;
