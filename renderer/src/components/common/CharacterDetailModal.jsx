import React from 'react';
import PropTypes from 'prop-types';
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogActions,
    IconButton,
    Tooltip,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import OpenInNewIcon from '@mui/icons-material/OpenInNew';
import { calculateDaysFromToday } from "../../utils/formatter.jsx";


const CharacterDetailModal = ({
                                  open,
                                  onClose,
                                  character,
                                  skillConversions,
                              }) => {
    if (!character || !character.Character) {
        return null;
    }

    const charId = character.Character.CharacterID;
    const charName = character.Character.CharacterName;
    const portraitUrl = `https://images.evetech.net/characters/${charId}/portrait`;
    const totalSp = character?.Character?.CharacterSkillsResponse?.total_sp;
    const formattedSP = totalSp ? totalSp.toLocaleString() : '0';
    const zKillUrl = `https://zkillboard.com/character/${charId}/`;

    const mctTooltip = character.MCT
        ? `Training: ${character.Training || 'Unknown'}`
        : 'Skill queue paused';

    // Use the entire skill queue
    const skillQueueItems = character.Character.SkillQueue || [];

    return (
        <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
            {/* Modal Title Bar */}
            <div className="flex justify-between items-center px-4 py-2 bg-gray-700">
                <DialogTitle className="text-teal-200 p-0">
                    {charName}
                </DialogTitle>
                <IconButton aria-label="close" onClick={onClose} sx={{ color: '#99f6e4' }}>
                    <CloseIcon />
                </IconButton>
            </div>

            <DialogContent className="bg-gray-800 text-teal-100">
                {/* Character Header Info */}
                <div className="flex flex-col sm:flex-row sm:space-x-4 mt-4">
                    <img
                        src={portraitUrl}
                        alt={`${charName} portrait`}
                        className="w-32 h-32 rounded-md shadow-md border border-teal-600"
                    />
                    <div className="mt-2 sm:mt-0 flex-1">
                        <div className="flex items-center space-x-2">
                            <Tooltip title="Open zKillboard">
                                <IconButton
                                    aria-label="Open zKillboard"
                                    size="small"
                                    onClick={() => {
                                        if (window.electronAPI && window.electronAPI.openExternal) {
                                            window.electronAPI.openExternal(zKillUrl);
                                        } else {
                                            window.open(zKillUrl, '_blank', 'noopener,noreferrer');
                                        }
                                    }}
                                    sx={{ color: '#99f6e4', '&:hover': { color: '#ffffff' } }}
                                >
                                    <OpenInNewIcon fontSize="inherit" />
                                </IconButton>
                            </Tooltip>

                            <Tooltip title={mctTooltip}>
                                <div
                                    data-testid="mct-indicator"
                                    className={`w-3 h-3 rounded-full ${character.MCT ? 'bg-green-400' : 'bg-gray-400'}`}
                                ></div>
                            </Tooltip>
                        </div>
                        <div className="flex flex-col mt-2 space-y-1">
                            {character.Role && (
                                <div className="text-sm">
                                    <span className="text-teal-400 font-medium">Role:</span> {character.Role}
                                </div>
                            )}
                            <div className="text-sm">
                                <span className="text-teal-400 font-medium">Location:</span> {character.Character.LocationName || 'Unknown'}
                            </div>
                            <div className="text-sm">
                                <span className="text-teal-400 font-medium">Total SP:</span> {formattedSP}
                            </div>
                            {character.CorporationName && (
                                <div className="text-sm">
                                    <span className="text-teal-400 font-medium">Corporation:</span>
                                    <span className="text-gray-300 italic"> {character.CorporationName} </span>
                                </div>
                            )}
                            {character.AllianceName && (
                                <div className="text-sm">
                                    <span className="text-teal-400 font-medium">Alliance:</span>
                                    <span className="text-gray-300 italic"> {character.AllianceName} </span>
                                </div>
                            )}
                        </div>
                    </div>
                </div>

                {/* Skill Queue Section */}
                <h3 className="text-teal-200 font-semibold text-md mt-6 mb-2">
                    Skill Queue
                </h3>
                {skillQueueItems.length === 0 ? (
                    <div className="text-gray-300 text-sm">No skill queue data available.</div>
                ) : (
                    <div className="overflow-x-auto max-h-64 overflow-y-auto border border-gray-700 rounded">
                        <table className="min-w-full text-sm">
                            <thead className="bg-gray-700 text-teal-300">
                            <tr>
                                <th className="px-2 py-1 text-left">Skill</th>
                                <th className="px-2 py-1 text-left">Level</th>
                                <th className="px-2 py-1 text-left">Completion</th>
                            </tr>
                            </thead>
                            <tbody>
                            {skillQueueItems.map((item, index) => {
                                const finishDate = item.finish_date ? calculateDaysFromToday(item.finish_date) : 'N/A';

                                // Convert skill_id to skill name
                                const skillName = skillConversions[item.skill_id] || `Skill #${item.skill_id}`;

                                return (
                                    <tr key={index} className="border-b border-gray-600 text-gray-200">
                                        <td className="px-2 py-1">{skillName}</td>
                                        <td className="px-2 py-1">{item.finished_level}</td>
                                        <td className="px-2 py-1">{finishDate}</td>
                                    </tr>
                                );
                            })}
                            </tbody>
                        </table>
                    </div>
                )}
            </DialogContent>

            <DialogActions className="bg-gray-800">
                <button
                    onClick={onClose}
                    className="px-3 py-1 bg-teal-600 hover:bg-teal-500 text-white rounded text-sm"
                >
                    Close
                </button>
            </DialogActions>
        </Dialog>
    );
};

CharacterDetailModal.propTypes = {
    open: PropTypes.bool.isRequired,
    onClose: PropTypes.func.isRequired,
    character: PropTypes.object.isRequired,
    skillConversions: PropTypes.object.isRequired,
};

export default CharacterDetailModal;
