// Authentication Types
export type UserRole = 'ADMIN' | 'EDITOR' | 'VIEWER';

export interface User {
  id: string;
  username: string;
  role: UserRole;
  status: 'ACTIVE' | 'INACTIVE';
  lastLogin: string;
  email?: string;
}

// Login Request/Response Types (matching backend API)
export interface LoginRequest {
  username: string;
  password: string;
  captcha: string;
  captcha_id: string;
}

export interface LoginResponse {
  token: string;
  username: string;
}
