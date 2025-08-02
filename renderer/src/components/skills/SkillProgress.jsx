import React from 'react';
import { motion } from 'framer-motion';

const SkillProgress = ({ skill, progress, isTraining }) => {
    const circumference = 2 * Math.PI * 45; // radius = 45
    const strokeDashoffset = circumference - (progress / 100) * circumference;
    
    return (
        <div className="relative w-32 h-32">
            <svg className="transform -rotate-90 w-full h-full">
                {/* Background circle */}
                <circle
                    cx="64"
                    cy="64"
                    r="45"
                    stroke="rgba(30, 41, 59, 0.5)"
                    strokeWidth="8"
                    fill="none"
                />
                
                {/* Progress circle */}
                <motion.circle
                    cx="64"
                    cy="64"
                    r="45"
                    stroke="url(#skillGradient)"
                    strokeWidth="8"
                    fill="none"
                    strokeLinecap="round"
                    initial={{ strokeDashoffset: circumference }}
                    animate={{ strokeDashoffset }}
                    transition={{ duration: 1, ease: "easeOut" }}
                    style={{
                        strokeDasharray: circumference,
                        filter: isTraining ? 'drop-shadow(0 0 10px rgba(20, 184, 166, 0.5))' : 'none'
                    }}
                />
                
                {/* Gradient definition */}
                <defs>
                    <linearGradient id="skillGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                        <stop offset="0%" stopColor="#14b8a6" />
                        <stop offset="100%" stopColor="#0891b2" />
                    </linearGradient>
                </defs>
            </svg>
            
            {/* Center content */}
            <div className="absolute inset-0 flex flex-col items-center justify-center">
                <motion.div
                    className="text-2xl font-bold font-mono"
                    initial={{ scale: 0 }}
                    animate={{ scale: 1 }}
                    transition={{ delay: 0.5, type: "spring" }}
                >
                    {progress}%
                </motion.div>
                <div className="text-xs text-gray-400">{skill.name}</div>
            </div>
            
            {/* Training pulse effect */}
            {isTraining && (
                <motion.div
                    className="absolute inset-0 rounded-full"
                    style={{
                        background: 'radial-gradient(circle, rgba(20,184,166,0.2) 0%, transparent 70%)',
                    }}
                    animate={{
                        scale: [1, 1.3, 1],
                        opacity: [0.5, 0, 0.5]
                    }}
                    transition={{
                        duration: 2,
                        repeat: Infinity,
                        ease: "easeInOut"
                    }}
                />
            )}
        </div>
    );
};

export default SkillProgress;