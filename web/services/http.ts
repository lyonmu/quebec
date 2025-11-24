import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse, InternalAxiosRequestConfig } from 'axios';
import { ApiResponse } from '../types';

// Create axios instance
const http: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/core/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor
http.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    // Add authentication token if available
    const token = localStorage.getItem('x-quebec-token');
    if (token && config.headers) {
      config.headers['x-quebec-token'] = token;
    }
    return config;
  },
  (error) => {
    console.error('Request error:', error);
    return Promise.reject(error);
  }
);

// Response interceptor
http.interceptors.response.use(
  (response: AxiosResponse<ApiResponse<any>>) => {
    const { data } = response;
    
    // Check if response follows the standard API format
    if (data && typeof data.code !== 'undefined') {
      // Success code is 50000
      if (data.code === 50000) {
        return response;
      } else {
        // API returned error code
        const error = new Error(data.message || 'Request failed');
        console.error('API error:', data);
        return Promise.reject(error);
      }
    }
    
    // If response doesn't follow standard format, return as is
    return response;
  },
  (error) => {
    // Handle HTTP errors
    if (error.response) {
      const { status, data } = error.response;
      
      switch (status) {
        case 401:
          console.error('Unauthorized - redirecting to login');
          // Clear token and redirect to login
          localStorage.removeItem('x-quebec-token');
          window.location.href = '/';
          break;
        case 403:
          console.error('Forbidden - insufficient permissions');
          break;
        case 404:
          console.error('Resource not found');
          break;
        case 500:
          console.error('Server error');
          break;
        default:
          console.error(`HTTP error ${status}:`, data);
      }
    } else if (error.request) {
      console.error('Network error - no response received');
    } else {
      console.error('Request setup error:', error.message);
    }
    
    return Promise.reject(error);
  }
);

// Typed request methods
export const httpClient = {
  get: <T = any>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> => {
    return http.get<ApiResponse<T>>(url, config).then(res => res.data);
  },
  
  post: <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> => {
    return http.post<ApiResponse<T>>(url, data, config).then(res => res.data);
  },
  
  put: <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> => {
    return http.put<ApiResponse<T>>(url, data, config).then(res => res.data);
  },
  
  delete: <T = any>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> => {
    return http.delete<ApiResponse<T>>(url, config).then(res => res.data);
  },
  
  patch: <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> => {
    return http.patch<ApiResponse<T>>(url, data, config).then(res => res.data);
  },
};

export default http;
