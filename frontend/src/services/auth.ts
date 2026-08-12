import api from './api';

// Auth related types and services
export interface User {
  id: number;
  username: string;
  role: string;
  email: string;
  created_at: string;
  updated_at?: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

const AUTH_TOKEN_KEY = 'token';
const USER_KEY = 'user';

/**
 * Logs in a user with the provided credentials
 */
export const login = async (username: string, password: string): Promise<{ token: string; user: User }> => {
  try {
    const response = await api.post<AuthResponse>('/api/auth/login', { username, password });
    const { token, user } = response;
    
    // Store token and user in localStorage
    localStorage.setItem(AUTH_TOKEN_KEY, token);
    localStorage.setItem(USER_KEY, JSON.stringify(user));
    
    // Set default auth header for future requests
    api.setAuthToken(token);
    
    return { token, user };
  } catch (error) {
    const errorMessage = (error as any).response?.data?.message || 'Login failed. Please try again.';
    throw new Error(errorMessage);
  }
};

/**
 * Logs out the current user
 */
export const logout = (): void => {
  localStorage.removeItem(AUTH_TOKEN_KEY);
  localStorage.removeItem(USER_KEY);
  api.clearAuthToken();
  // Redirect to login page
  window.location.href = '/login';
};

/**
 * Gets the current authenticated user
 */
export const getCurrentUser = (): User | null => {
  if (typeof window === 'undefined') return null;
  
  const userStr = localStorage.getItem(USER_KEY);
  if (!userStr) return null;
  
  try {
    return JSON.parse(userStr);
  } catch (e) {
    console.error('Error parsing user data:', e);
    return null;
  }
};

/**
 * Gets the current authentication token
 */
export const getAuthToken = (): string | null => {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem(AUTH_TOKEN_KEY);
};

/**
 * Checks if a user is authenticated
 */
export const isAuthenticated = (): boolean => {
  return !!getAuthToken();
};

/**
 * Checks if the current user has a specific role
 */
export const hasRole = (role: string): boolean => {
  const user = getCurrentUser();
  return user?.role === role;
};

/**
 * Checks if the current user has any of the specified roles
 */
export const hasAnyRole = (roles: string[]): boolean => {
  const user = getCurrentUser();
  return roles.includes(user?.role || '');
};
