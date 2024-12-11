// EditableAccountName.jsx
import React, { useState } from 'react';
import PropTypes from 'prop-types';

const EditableAccountName = ({ accountID, accountName, onNameUpdate }) => {
    const [name, setName] = useState(accountName);
    const [isEditing, setIsEditing] = useState(false);

    const handleNameChange = (e) => {
        setName(e.target.value);
    };

    const handleBlur = () => {
        if (name !== accountName) {
            onNameUpdate(accountID, name);
        }
        setIsEditing(false);
    };

    const handleKeyDown = (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            handleBlur();
        }
    };

    return isEditing ? (
        <input
            type="text"
            value={name}
            onChange={handleNameChange}
            onBlur={handleBlur}
            onKeyDown={handleKeyDown}
            autoFocus
        />
    ) : (
        <span onClick={() => setIsEditing(true)}>{accountName}</span>
    );
};

EditableAccountName.propTypes = {
    accountID: PropTypes.number.isRequired,
    accountName: PropTypes.string.isRequired,
    onNameUpdate: PropTypes.func.isRequired,
};

export default EditableAccountName;
