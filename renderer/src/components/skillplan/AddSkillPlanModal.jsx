// AddSkillPlanModal.jsx
import React, { useState } from 'react';
import PropTypes from 'prop-types';
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    TextField,
    Button,
} from '@mui/material';

const AddSkillPlanModal = ({ onClose, onSave }) => {
    const [planName, setPlanName] = useState('');
    const [planContents, setPlanContents] = useState('');

    const handleSubmit = (e) => {
        e.preventDefault();
        onSave(planName, planContents);
    };

    return (
        <Dialog open onClose={onClose} maxWidth="sm" fullWidth>
            <DialogTitle className="bg-gray-800 text-teal-200">Add Skill Plan</DialogTitle>
            <form onSubmit={handleSubmit}>
                <DialogContent className="bg-gray-800">
                    <TextField
                        label="Skill Plan Name"
                        value={planName}
                        onChange={(e) => setPlanName(e.target.value)}
                        fullWidth
                        margin="normal"
                        required
                        InputProps={{
                            className: 'text-teal-200',
                        }}
                        InputLabelProps={{
                            className: 'text-gray-300',
                        }}
                        className="bg-gray-700"
                    />
                    <TextField
                        label="Skill Plan Contents"
                        value={planContents}
                        onChange={(e) => setPlanContents(e.target.value)}
                        fullWidth
                        margin="normal"
                        multiline
                        rows={5}
                        required
                        InputProps={{
                            className: 'text-teal-200',
                        }}
                        InputLabelProps={{
                            className: 'text-gray-300',
                        }}
                        className="bg-gray-700"
                    />
                </DialogContent>
                <DialogActions className="bg-gray-800">
                    <Button onClick={onClose} className="text-gray-100">
                        Cancel
                    </Button>
                    <Button type="submit" variant="contained" className="bg-teal-500 text-gray-100 hover:bg-teal-400">
                        Save
                    </Button>
                </DialogActions>
            </form>
        </Dialog>
    );
};

AddSkillPlanModal.propTypes = {
    onClose: PropTypes.func.isRequired,
    onSave: PropTypes.func.isRequired,
};

export default AddSkillPlanModal;
