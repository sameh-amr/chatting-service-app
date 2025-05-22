import React from 'react';
import { format } from 'date-fns';
import { Check, CheckCheck, Download } from 'lucide-react';
import { Message as MessageType } from '../types';
import { useChat } from '../contexts/ChatContext';

interface MessageProps {
  message: MessageType;
  isCurrentUser: boolean;
}

const Message: React.FC<MessageProps> = ({ message, isCurrentUser }) => {
  console.log(message)
  
  const { downloadFile } = useChat();
  
  const handleDownload = () => {
    if (message.media_url) {
      const filename = message.media_url.split('/').pop() || 'file';
      downloadFile(filename);
    }
  };

  const containerClasses = isCurrentUser
    ? 'flex justify-end'
    : 'flex justify-start';

  const messageClasses = isCurrentUser
    ? 'bg-blue-600 text-white rounded-t-xl rounded-l-xl'
    : 'bg-white text-gray-800 border rounded-t-xl rounded-r-xl';
    
  const timeClasses = isCurrentUser
    ? 'text-right text-gray-400'
    : 'text-left text-gray-500';
  
  // Extract filename from media URL if it exists
  const hasMedia = !!message.media_url;
  const filename = hasMedia 
    ? message.media_url.split('/').pop() || 'file'
    : '';

  
  return (
    <div className={`${containerClasses} mb-4`}>
      <div className="max-w-[80%] md:max-w-[70%]">
        <div className={`${messageClasses} px-4 py-2 shadow-sm`}>
          {message.content}
          
          {hasMedia && (
            <div 
              className="mt-2 flex items-center p-2 bg-black bg-opacity-5 rounded cursor-pointer hover:bg-opacity-10"
              onClick={handleDownload}
            >
              <Download size={16} className="mr-2 flex-shrink-0" />
              <span className="text-sm truncate">{filename}</span>
            </div>
          )}
        </div>
        
        <div className={`${timeClasses} text-xs mt-1 flex items-center`}>
          <span>
            
          </span>
          
          {isCurrentUser && (
            <span className="ml-1">
              {message.read ? (
                <CheckCheck size={14} className="inline text-blue-500" />
              ) : message.delivered ? (
                <CheckCheck size={14} className="inline text-gray-400" />
              ) : (
                <Check size={14} className="inline text-gray-400" />
              )}
            </span>
          )}
        </div>
      </div>
    </div>
  );
};

export default Message;