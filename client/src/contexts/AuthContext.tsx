import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useNavigate } from 'react-router-dom';
import toast from 'react-hot-toast';
import { User } from '../types';
import api from '../services/api';
import { disconnectWebSocket } from '../services/websocket';

interface AuthContextType {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  userId: string | null;
  login: (email: string, password: string) => Promise<void>;
  signup: (username: string, email: string, password: string) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(sessionStorage.getItem('token'));
  const [isLoading, setIsLoading] = useState(false);
  const navigate = useNavigate();

  // Persist token in sessionStorage when it changes, no navigation here
  useEffect(() => {
    if (token) {
      sessionStorage.setItem('token', token);
    } else {
      sessionStorage.removeItem('token');
    }
  }, [token]);

  // On mount, if token exists, verify and fetch user profile
  useEffect(() => {
    const loadUser = async () => {
      if (!token) return;
      try {
        setIsLoading(true);
        // Assuming your API has /auth/me or similar endpoint to get user info
        const userResp = await api.get('/auth/me', {
          headers: { Authorization: `Bearer ${token}` },
        });
        setUser(userResp.data);
      } catch (error) {
        // Token invalid or expired, clear it
        setToken(null);
        setUser(null);
      } finally {
        setIsLoading(false);
      }
    };

    loadUser();
  }, [token]);

  const login = async (email: string, password: string) => {
    try {
      setIsLoading(true);
      const response = await api.post('/auth/login', { email, password });
            console.log(response);

      const { token: newToken, user: userData } = response.data; // Expecting user info here
      setToken(newToken);
      setUser(userData);
      toast.success('Login successful!');
      navigate('/chat');
    } catch (error) {
      console.error('Login error:', error);
      toast.error('Login failed. Please check your credentials.');
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const signup = async (username: string, email: string, password: string) => {
    try {
      setIsLoading(true);
      const response = await api.post('/auth/signup', { username, email, password });
      const { token: newToken, user: userData } = response.data; // Expecting user info here
      setToken(newToken);
      setUser(userData);
      toast.success('Account created successfully!');
      navigate('/chat');
    } catch (error) {
      console.error('Signup error:', error);
      toast.error('Signup failed. Please try again.');
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = () => {
    setToken(null);
    setUser(null);
    sessionStorage.removeItem('token');
    disconnectWebSocket();
    navigate('/login');
    toast.success('Logged out successfully!');
  };

  const value = {
    user,
    token,
    isAuthenticated: !!token && !!user,
    isLoading,
    userId: user ? user.id : null,
    login,
    signup,
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
