import React from 'react';
import { motion } from 'framer-motion';
import { cn } from '../../utils/cn';

const FuturisticButton = React.forwardRef(({ 
  children, 
  className,
  variant = 'primary',
  size = 'md',
  disabled = false,
  loading = false,
  icon,
  iconPosition = 'left',
  onClick,
  ...props 
}, ref) => {
  const variants = {
    primary: 'bg-gradient-to-r from-teal-500 to-blue-600 hover:from-teal-400 hover:to-blue-500 text-white',
    secondary: 'bg-gradient-to-r from-gray-700 to-gray-800 hover:from-gray-600 hover:to-gray-700 text-gray-300',
    danger: 'bg-gradient-to-r from-red-500 to-orange-600 hover:from-red-400 hover:to-orange-500 text-white',
    success: 'bg-gradient-to-r from-green-500 to-teal-600 hover:from-green-400 hover:to-teal-500 text-white',
    ghost: 'bg-transparent hover:bg-gray-800/50 text-teal-400 border border-teal-400/50 hover:border-teal-400',
  };

  const sizes = {
    xs: 'px-2 py-1 text-xs',
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2 text-base',
    lg: 'px-6 py-3 text-lg',
    xl: 'px-8 py-4 text-xl',
  };

  const baseClasses = cn(
    "relative overflow-hidden rounded-lg font-medium",
    "transition-all duration-300 transform",
    "focus:outline-none focus:ring-2 focus:ring-teal-500 focus:ring-offset-2 focus:ring-offset-gray-900",
    !disabled && "hover:scale-105 active:scale-95",
    disabled && "opacity-50 cursor-not-allowed",
    variants[variant],
    sizes[size],
    className
  );

  const rippleClasses = cn(
    "absolute inset-0 rounded-lg",
    "before:content-[''] before:absolute before:top-1/2 before:left-1/2",
    "before:w-0 before:h-0 before:rounded-full",
    "before:bg-white/20 before:-translate-x-1/2 before:-translate-y-1/2",
    "before:transition-all before:duration-600",
    "hover:before:w-[300px] hover:before:h-[300px]"
  );

  return (
    <motion.button
      ref={ref}
      className={baseClasses}
      disabled={disabled || loading}
      onClick={onClick}
      whileTap={{ scale: disabled ? 1 : 0.95 }}
      {...props}
    >
      <span className={rippleClasses} />
      <span className="relative flex items-center justify-center gap-2">
        {loading ? (
          <svg className="animate-spin h-5 w-5" viewBox="0 0 24 24">
            <circle 
              className="opacity-25" 
              cx="12" 
              cy="12" 
              r="10" 
              stroke="currentColor" 
              strokeWidth="4"
              fill="none"
            />
            <path 
              className="opacity-75" 
              fill="currentColor" 
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
        ) : (
          <>
            {icon && iconPosition === 'left' && <span className="flex-shrink-0">{icon}</span>}
            {children}
            {icon && iconPosition === 'right' && <span className="flex-shrink-0">{icon}</span>}
          </>
        )}
      </span>
    </motion.button>
  );
});

FuturisticButton.displayName = 'FuturisticButton';

export default FuturisticButton;