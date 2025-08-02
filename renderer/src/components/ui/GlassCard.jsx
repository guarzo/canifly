import React from 'react';
import { motion } from 'framer-motion';
import { cn } from '../../utils/cn';

const GlassCard = React.forwardRef(({ 
  children, 
  className, 
  hover = true,
  onClick,
  ...props 
}, ref) => {
  const baseClasses = cn(
    "glass rounded-lg p-4",
    "transition-all duration-300",
    hover && "hover-lift hover:shadow-glow-lg hover:border-teal-500/50",
    onClick && "cursor-pointer",
    className
  );

  return (
    <motion.div
      ref={ref}
      className={baseClasses}
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.5 }}
      onClick={onClick}
      {...props}
    >
      {children}
    </motion.div>
  );
});

GlassCard.displayName = 'GlassCard';

export default GlassCard;