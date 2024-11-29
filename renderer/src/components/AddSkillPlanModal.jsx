import React, { useState } from 'react';
import PropTypes from 'prop-types';

const AddSkillPlanModal = ({ onClose, onSave }) => {
    const [planName, setPlanName] = useState('');
    const [planContents, setPlanContents] = useState('');

    const handleSubmit = (e) => {
        e.preventDefault();
        onSave(planName, planContents);
    };

    return (
        <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50">
            <div className="bg-gradient-to-b from-gray-800 to-gray-700 text-gray-100 p-6 rounded-lg shadow-lg w-1/3">
                <h2 className="text-xl font-bold mb-4 text-teal-200">Add Skill Plan</h2>
                <form onSubmit={handleSubmit}>
                    <div className="mb-4">
                        <label className="block text-gray-300 mb-2">Skill Plan Name</label>
                        <input
                            type="text"
                            value={planName}
                            onChange={(e) => setPlanName(e.target.value)}
                            className="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring focus:border-teal-400 bg-gray-700 text-teal-200"
                            required
                        />
                    </div>
                    <div className="mb-4">
                        <label className="block text-gray-300 mb-2">Skill Plan Contents</label>
                        <textarea
                            value={planContents}
                            onChange={(e) => setPlanContents(e.target.value)}
                            className="w-full px-3 py-2 border rounded-md focus:outline-none focus:ring focus:border-teal-400 bg-gray-700 text-teal-200"
                            rows="5"
                            required
                        ></textarea>
                    </div>
                    <div className="flex justify-end space-x-4">
                        <button
                            type="button"
                            onClick={onClose}
                            className="px-4 py-2 bg-gray-600 text-gray-100 rounded hover:bg-gray-700 focus:outline-none"
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            className="px-4 py-2 bg-teal-500 text-gray-100 rounded hover:bg-teal-400 focus:outline-none"
                        >
                            Save
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

AddSkillPlanModal.propTypes = {
    onClose: PropTypes.func.isRequired,
    onSave: PropTypes.func.isRequired,
};

export default AddSkillPlanModal;
