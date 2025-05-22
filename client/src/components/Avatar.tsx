import React from 'react';
import { User } from '../types';

interface AvatarProps {
  user: User;
  size?: 'sm' | 'md' | 'lg' | 'xl';
  showStatus?: boolean;
  isOnline?: boolean;
  className?: string;
}

const Avatar: React.FC<AvatarProps> = ({
  user,
  size = 'md',
  showStatus = false,
  isOnline = false,
  className = '',
}) => {
  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map(part => part[0])
      .join('')
      .toUpperCase()
      .substring(0, 2);
  };
  
  const sizeClasses = {
    sm: 'h-8 w-8 text-xs',
    md: 'h-10 w-10 text-sm',
    lg: 'h-12 w-12 text-base',
    xl: 'h-16 w-16 text-lg',
  };
  
  const statusSizeClasses = {
    sm: 'h-2 w-2',
    md: 'h-2.5 w-2.5',
    lg: 'h-3 w-3',
    xl: 'h-4 w-4',
  };
  
  const statusPositionClasses = {
    sm: '-right-0.5 -bottom-0.5',
    md: '-right-0.5 -bottom-0.5',
    lg: 'right-0 bottom-0',
    xl: 'right-0.5 bottom-0.5',
  };
  
  // Generate color based on username
  const getColor = (name: string) => {
    const colors = [
      'bg-blue-500',
      'bg-green-500',
      'bg-purple-500',
      'bg-yellow-500',
      'bg-pink-500',
      'bg-indigo-500',
      'bg-red-500',
      'bg-teal-500',
    ];
    
    const index = name.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0) % colors.length;
    return colors[index];
  };
  
  return (
    <div className={`relative inline-flex ${className}`}>
      <div
        className={`
          ${sizeClasses[size]}
          ${getColor(user.username)}
          rounded-full flex items-center justify-center text-white font-medium
        `}
      >
        {user.avatar ? (
          <img
            src={user.avatar}
            alt={user.username}
            className="h-full w-full object-cover rounded-full"
          />
        ) : (
          getInitials(user.username)
        )}
      </div>
      
      {showStatus && (
        <span
          className={`
            absolute ${statusPositionClasses[size]} ${statusSizeClasses[size]}
            rounded-full border-2 border-white
            ${isOnline ? 'bg-green-500' : 'bg-gray-400'}
          `}
        />
      )}
    </div>
  );
};

export default Avatar;