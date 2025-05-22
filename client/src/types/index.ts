export interface User {
  id: string;
  username: string;
  email: string;
}

export interface Message {
  id: string;
  sender_id: string;
  recipient_id: string;
  content: string;
  media_url?: string;
  is_broadcast: boolean;
  created_at: string;
  delivered: boolean;
  read: boolean;
}