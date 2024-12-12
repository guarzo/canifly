import { useState } from 'react';
import PropTypes from 'prop-types';
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    TextField,
    Button,
    Typography
} from '@mui/material';

const AddSkillPlanModal = ({ onClose, onSave }) => {
    const [planName, setPlanName] = useState('');
    const [planContents, setPlanContents] = useState('');

    const handleSubmit = (e) => {
        e.preventDefault();
        onSave(planName.trim(), planContents.trim());
    };

    return (
        <Dialog open onClose={onClose} maxWidth="sm" fullWidth>
            <form onSubmit={handleSubmit}>
                <DialogTitle
                    className="bg-gradient-to-r from-gray-900 to-gray-800 text-teal-200 border-b border-gray-700"
                    sx={{ paddingY: '0.75rem' }}
                >
                    <Typography variant="h6" sx={{ fontWeight: 'bold' }}>
                        Add Skill Plan
                    </Typography>
                </DialogTitle>
                <DialogContent className="bg-gray-800">
                    <div className="space-y-4">
                        <TextField
                            label="Skill Plan Name"
                            placeholder="e.g. Battleship Mastery"
                            value={planName}
                            onChange={(e) => setPlanName(e.target.value)}
                            fullWidth
                            required
                            variant="filled"
                            InputProps={{
                                className: 'text-teal-200 bg-gray-700 rounded-md px-3 py-2',
                            }}
                            InputLabelProps={{
                                className: 'text-gray-300',
                            }}
                        />
                        <TextField
                            label="Skill Plan Contents"
                            placeholder="Enter skills and levels, e.g. 'Gunnery 5\nLarge Hybrid Turret 4'"
                            value={planContents}
                            onChange={(e) => setPlanContents(e.target.value)}
                            fullWidth
                            required
                            multiline
                            rows={5}
                            variant="filled"
                            InputProps={{
                                className: 'text-teal-200 bg-gray-700 rounded-md px-3 py-2',
                            }}
                            InputLabelProps={{
                                className: 'text-gray-300',
                            }}
                        />
                    </div>
                </DialogContent>
                <DialogActions className="bg-gray-800 border-t border-gray-700 flex items-center justify-end space-x-2 py-2 px-3">
                    <Button
                        onClick={onClose}
                        className="text-gray-200 hover:text-white normal-case"
                        sx={{ textTransform: 'none' }}
                    >
                        Cancel
                    </Button>
                    <Button
                        type="submit"
                        variant="contained"
                        sx={{
                            textTransform: 'none',
                            backgroundColor: '#14b8a6',
                            '&:hover': {
                                backgroundColor: '#0d9488',
                            },
                            color: '#ffffff',
                            fontWeight: 'bold'
                        }}
                    >
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
