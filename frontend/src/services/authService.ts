import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8083/api';

// Create axios instance with default config
const api = axios.create({
  baseURL: API_URL,
  withCredentials: true, // This is important for sending cookies with requests
  headers: {
    'Content-Type': 'application/json',
  },
});

interface LoginResponse {
  token: string;
  user: {
    id: number;
    username: string;
    email: string;
    role: string;
    isActive: boolean;
    lastLogin: string;
    createdAt: string;
  };
}

export const login = async (username: string, password: string) => {
  const response = await api.post<LoginResponse>('/auth/login', {
    username,
    password,
  });

  if (response.data.token) {
    localStorage.setItem('token', response.data.token);
    localStorage.setItem('user', JSON.stringify(response.data.user));
  }

  return response.data;
};

export const logout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
};

export const getCurrentUser = () => {
  const userStr = localStorage.getItem('user');
  if (!userStr) return null;
  
  try {
    return JSON.parse(userStr);
  } catch (e) {
    return null;
  }
};

export const getAuthToken = () => {
  return localStorage.getItem('token');
};

export const isAuthenticated = () => {
  return !!getAuthToken();
};

export const isAdmin = () => {
  const user = getCurrentUser();
  return user && user.role === 'admin';
};

// Request interceptor to add auth token to requests
api.interceptors.request.use(
  (config) => {
    const token = getAuthToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle 401 Unauthorized
try {
  api.interceptors.response.use(
    (response) => response,
    (error) => {
      if (error.response?.status === 401) {
        // If we get a 401, clear the auth data and redirect to login
        logout();
        window.location.href = '/login';
      }
      return Promise.reject(error);
    }
  );
} catch (error) {
  console.error('Error setting up response interceptor:', error);
}

export const setupAxiosInterceptors = () => {
  // Interceptors are now set up with the api instance
};
