import React from 'react';
import PropTypes from 'prop-types';
import AccountCard from './AccountCard';
import UnassignedCharacters from './UnassignedCharacters';

const Dashboard = ({
                       accounts,
                       unassignedCharacters,
                       onAssignCharacter,
                       onToggleAccountStatus,
                       onUpdateCharacter,
                       onUpdateAccountName,
                       onRemoveCharacter,
                       roles,
                   }) => {
    return (
        <div className="space-y-6 mt-16">
            {/* Accounts Section */}
            {accounts.map((account) => (
                <AccountCard
                    key={account.ID}
                    account={account}
                    onToggleAccountStatus={onToggleAccountStatus}
                    onUpdateAccountName={onUpdateAccountName}
                    onUpdateCharacter={onUpdateCharacter}
                    onRemoveCharacter={onRemoveCharacter}
                    roles={roles}
                />
            ))}

            {/* Unassigned Characters Section */}
            {unassignedCharacters.length > 0 && (
                <UnassignedCharacters
                    characters={unassignedCharacters}
                    accounts={accounts}
                    onAssignCharacter={onAssignCharacter}
                />
            )}
        </div>
    );
};

Dashboard.propTypes = {
    accounts: PropTypes.array.isRequired,
    unassignedCharacters: PropTypes.array.isRequired,
    onAssignCharacter: PropTypes.func.isRequired,
    onToggleAccountStatus: PropTypes.func.isRequired,
    onUpdateCharacter: PropTypes.func.isRequired,
    onUpdateAccountName: PropTypes.func.isRequired,
    onRemoveCharacter: PropTypes.func.isRequired,
};

export default Dashboard;
