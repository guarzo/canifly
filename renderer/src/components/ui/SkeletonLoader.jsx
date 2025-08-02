import React from 'react';
import { motion } from 'framer-motion';

const SkeletonLoader = ({ variant = 'text', className = '' }) => {
    const baseClasses = 'bg-gradient-to-r from-gray-800 via-gray-700 to-gray-800 rounded animate-shimmer';
    
    const variants = {
        text: 'h-4 w-full',
        title: 'h-8 w-3/4',
        avatar: 'h-16 w-16 rounded-full',
        card: 'h-32 w-full rounded-lg',
        button: 'h-10 w-24 rounded',
    };
    
    return (
        <motion.div
            className={`${baseClasses} ${variants[variant]} ${className}`}
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.3 }}
        />
    );
};

// Account Card Skeleton
export const AccountCardSkeleton = () => (
    <div className="glass rounded-lg p-6 space-y-4">
        <div className="flex items-center space-x-4">
            <SkeletonLoader variant="avatar" />
            <div className="flex-1 space-y-2">
                <SkeletonLoader variant="title" />
                <SkeletonLoader variant="text" className="w-1/2" />
            </div>
        </div>
        <div className="grid grid-cols-3 gap-4">
            {[1, 2, 3].map(i => (
                <SkeletonLoader key={i} variant="card" />
            ))}
        </div>
    </div>
);

// Character Item Skeleton
export const CharacterItemSkeleton = () => (
    <div className="glass p-3 rounded-lg space-y-2">
        <div className="flex justify-between items-center">
            <div className="flex items-center space-x-3">
                <SkeletonLoader variant="text" className="w-32" />
                <SkeletonLoader variant="text" className="w-16" />
            </div>
            <SkeletonLoader className="w-3 h-3 rounded-full" />
        </div>
        <div className="flex items-center space-x-2">
            <SkeletonLoader variant="text" className="w-20" />
            <SkeletonLoader variant="text" className="w-24" />
        </div>
    </div>
);

// Skill Plan Table Skeleton
export const SkillPlanTableSkeleton = () => (
    <div className="space-y-2">
        <SkeletonLoader variant="title" className="mb-4" />
        {[1, 2, 3, 4, 5].map(i => (
            <div key={i} className="flex space-x-4">
                <SkeletonLoader variant="text" className="w-1/3" />
                <SkeletonLoader variant="text" className="w-1/4" />
                <SkeletonLoader variant="text" className="w-1/4" />
            </div>
        ))}
    </div>
);

export default SkeletonLoader;