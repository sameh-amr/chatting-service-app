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
    ? message.media_url?.split('/').pop() || 'file'
    : '';

  // Helper: check if content is a file placeholder
  const isFilePlaceholder = hasMedia && message.content && message.content.startsWith('Sent a file:');
  // Helper: check if file is an image
  const isImage = hasMedia && /\.(jpg|jpeg|png|gif|bmp|webp)$/i.test(filename);

  return (
    <div className={`${containerClasses} mb-4`}>
      <div className="max-w-[80%] md:max-w-[70%]">
        <div className={`${messageClasses} px-4 py-2 shadow-sm`}>
          {/* Show text only if not a file placeholder, or if no media */}
          {!isFilePlaceholder && message.content}

          {/* If it's an image attachment, show preview */}
          {hasMedia && isImage && (
            <div className="mt-2 flex flex-col items-start">
              <img
                src={message.media_url}
                alt={filename}
                className="max-h-48 max-w-full rounded mb-2 border"
                style={{ objectFit: 'contain' }}
              />
              <div
                className="flex items-center p-2 bg-black bg-opacity-5 rounded cursor-pointer hover:bg-opacity-10"
                onClick={handleDownload}
              >
                <Download size={16} className="mr-2 flex-shrink-0" />
                <span className="text-sm truncate">{filename}</span>
              </div>
            </div>
          )}

          {/* If it's a non-image file, show download link */}
          {hasMedia && !isImage && (
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
            {message.created_at && message.created_at !== '0001-01-01T00:00:00Z' ? (
              format(new Date(message.created_at), 'MMM d, yyyy h:mm a')
            ) : (
              ''
            )}
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