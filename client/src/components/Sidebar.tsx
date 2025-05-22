import React, { useState, useEffect } from 'react';
import { Search, X } from 'lucide-react';
import Avatar from './Avatar';
import { User } from '../types';

interface SidebarProps {
  onlineUsers: User[]|null;
  selectedUser: User | null;
  setSelectedUser: (user: User | null) => void;
  showMobile: boolean;
  toggleMobile: () => void;
  currentUser: User | null;
}

const Sidebar: React.FC<SidebarProps> = ({
  onlineUsers,
  selectedUser,
  setSelectedUser,
  showMobile,
  toggleMobile,
  currentUser
}) => {
  const [searchTerm, setSearchTerm] = useState('');
  const [filteredUsers, setFilteredUsers] = useState<User[]|null>([]);

  // Filter users based on search term
  useEffect(() => {
    console.log(onlineUsers)
    const filtered = onlineUsers?.filter(user => 
      user.id !== currentUser?.id && // Filter out current user
      user.username.toLowerCase().includes(searchTerm.toLowerCase())
    )||null;
    setFilteredUsers(filtered);
  }, [onlineUsers, searchTerm, currentUser]);

  return (
    <>
      {/* Mobile Overlay */}
      {showMobile && (
        <div 
          className="md:hidden fixed inset-0 bg-black bg-opacity-50 z-40"
          onClick={toggleMobile}
        />
      )}
      
      {/* Sidebar */}
      <aside 
        className={`
          w-80 bg-white border-r border-gray-200 flex flex-col
          md:relative fixed inset-y-0 left-0 z-50 transition-transform duration-300 ease-in-out
          ${showMobile ? 'translate-x-0' : '-translate-x-full md:translate-x-0'}
        `}
      >
        {/* Header */}
        <div className="p-4 border-b border-gray-200 flex items-center justify-between">
          <h2 className="font-bold text-xl text-gray-800">Messages</h2>
          
          {/* Close button for mobile */}
          <button 
            onClick={toggleMobile}
            className="md:hidden text-gray-500 hover:text-gray-700"
          >
            <X size={20} />
          </button>
        </div>
        
        {/* Search */}
        <div className="p-4">
          <div className="relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" size={18} />
            <input
              type="text"
              placeholder="Search users..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-200"
            />
          </div>
        </div>
        
        {/* Users list */}
        <div className="flex-1 overflow-y-auto">
          {filteredUsers&&filteredUsers.length > 0 ? (
            <ul className="divide-y divide-gray-200">
              {filteredUsers.map(user => (
                <li key={user.id}>
                  <button
                    onClick={() => {
                      setSelectedUser(user);
                      if (showMobile) toggleMobile();
                    }}
                    className={`w-full p-3 flex items-center hover:bg-gray-50 transition-colors ${
                      selectedUser?.id === user.id ? 'bg-blue-50' : ''
                    }`}
                  >
                    <Avatar 
                      user={user} 
                      showStatus 
                      isOnline={onlineUsers?.some(u => u.id === user.id)} 
                    />
                    <div className="ml-3 text-left">
                      <p className="font-medium text-gray-900">{user.username}</p>
                      <p className="text-sm text-gray-500">
                        {onlineUsers?.some(u => u.id === user.id) ? 'Online' : 'Offline'}
                      </p>
                    </div>
                  </button>
                </li>
              ))}
            </ul>
          ) : (
            <div className="p-4 text-center text-gray-500">
              {searchTerm ? 'No users found' : 'No users available'}
            </div>
          )}
        </div>
        
        {/* Current user */}
        {currentUser && (
          <div className="p-4 border-t border-gray-200">
            <div className="flex items-center">
              <Avatar user={currentUser} />
              <div className="ml-3">
                <p className="font-medium text-gray-900">{currentUser.username}</p>
                <p className="text-xs text-green-500">Online</p>
              </div>
            </div>
          </div>
        )}
      </aside>
    </>
  );
};

export default Sidebar;