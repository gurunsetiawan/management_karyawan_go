import axios, { AxiosError, AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { logout } from './auth';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8083';

class ApiService {
  private instance: AxiosInstance;

  constructor() {
    this.instance = axios.create({
      baseURL: API_URL,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.setupInterceptors();
  }

  private setupInterceptors() {
    // Request interceptor
    this.instance.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('token');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response interceptor
    this.instance.interceptors.response.use(
      (response: AxiosResponse) => response.data,
      (error: AxiosError) => {
        if (error.response?.status === 401) {
          logout();
          window.location.href = '/login';
        }
        return Promise.reject(error);
      }
    );
  }

  // HTTP Methods
  public async get<T = any>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await this.instance.get<T>(url, config);
    return response as unknown as T;
  }

  public async post<T = any>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    const response = await this.instance.post<T>(url, data, config);
    return response as unknown as T;
  }

  public async put<T = any>(
    url: string,
    data?: any,
    config?: AxiosRequestConfig
  ): Promise<T> {
    const response = await this.instance.put<T>(url, data, config);
    return response as unknown as T;
  }

  public async delete<T = any>(
    url: string,
    config?: AxiosRequestConfig
  ): Promise<T> {
    const response = await this.instance.delete<T>(url, config);
    return response as unknown as T;
  }

  // Auth methods
  public setAuthToken(token: string): void {
    this.instance.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    localStorage.setItem('token', token);
  }

  public clearAuthToken(): void {
    delete this.instance.defaults.headers.common['Authorization'];
    localStorage.removeItem('token');
  }

  // Authentication
  public async login(username: string, password: string): Promise<{ token: string; user: any }> {
    return this.post('/api/auth/login', { username, password });
  }

  // Employee methods
  public async getEmployees(page: number = 1, limit: number = 10, search: string = '', config?: AxiosRequestConfig): Promise<PaginatedEmployeeResponse> {
    const params = new URLSearchParams({
      page: page.toString(),
      limit: limit.toString(),
    });
    if (search) {
      params.append('search', search);
    }
    return this.get<PaginatedEmployeeResponse>(`/api/employees?${params.toString()}`, config);
  }

  public async getEmployee(id: number): Promise<Employee> {
    return this.get<Employee>(`/api/employees/${id}`);
  }

  public async createEmployee(employee: Omit<Employee, 'id'>): Promise<Employee> {
    return this.post<Employee>('/api/employees', employee);
  }

  public async updateEmployee(id: number, employee: Partial<Employee>): Promise<Employee> {
    return this.put<Employee>(`/api/employees/${id}`, employee);
  }

  public async deleteEmployee(id: number): Promise<void> {
    return this.delete(`/api/employees/${id}`);
  }


  private handleError = (error: AxiosError): ApiResponse => {
    if (error.response) {
      const errorData = error.response.data as any;
      return {
        error: (errorData?.message || errorData?.error || 'An error occurred') as string,
        status: error.response.status,
      };
    }
    return {
      error: error.message || 'Network error',
      status: 0,
    };
  }
}

// Export a singleton instance
const api = new ApiService();
export default api;

// Interfaces
export interface Employee {
  id: number;
  name: string;
  email: string;
  position: string;
  phone: string;
  alamat: string;
  role?: string;
  created_at: string;
  updated_at?: string;
}

export interface PaginationMeta {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface PaginatedEmployeeResponse {
  data: Employee[];
  meta: PaginationMeta;
}

export interface ApiResponse<T = any> {
  data?: T;
  error?: string;
  status: number;
}
