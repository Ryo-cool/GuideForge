import React from 'react';

interface ButtonProps {
  children: React.ReactNode;
  variant?: 'primary' | 'secondary' | 'default';
  size?: 'small' | 'medium' | 'large';
  onClick?: () => void;
  type?: 'button' | 'submit' | 'reset';
  disabled?: boolean;
  className?: string;
  fullWidth?: boolean;
}

export const Button: React.FC<ButtonProps> = ({
  children,
  variant = 'default',
  size = 'medium',
  onClick,
  type = 'button',
  disabled = false,
  className = '',
  fullWidth = false,
}) => {
  const getVariantClasses = () => {
    switch (variant) {
      case 'primary':
        return 'bg-blue-600 hover:bg-blue-700 text-white shadow-md';
      case 'secondary':
        return 'bg-gray-200 hover:bg-gray-300 text-gray-800 shadow-md';
      default:
        return 'bg-white hover:bg-gray-100 text-gray-800 border border-gray-300 shadow-sm';
    }
  };

  const getSizeClasses = () => {
    switch (size) {
      case 'small':
        return 'text-sm py-1.5 px-3';
      case 'large':
        return 'text-lg py-3 px-7';
      default:
        return 'text-base py-2 px-5';
    }
  };

  const baseClasses = getVariantClasses();
  const sizeClasses = getSizeClasses();
  const widthClass = fullWidth ? 'w-full' : '';
  
  const classes = `rounded-md font-medium transition-colors duration-200 ${baseClasses} ${sizeClasses} ${widthClass} ${className}`.trim();

  return (
    <button
      type={type}
      className={classes}
      onClick={onClick}
      disabled={disabled}
    >
      {children}
    </button>
  );
}; 