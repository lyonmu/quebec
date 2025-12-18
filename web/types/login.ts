// Authentication Types
export type UserRole = 'ADMIN' | 'EDITOR' | 'VIEWER';

export interface User {
  id: string;
  username: string;
  nickname: string;
  role: UserRole;
  status: 'ACTIVE' | 'INACTIVE';
  lastLogin: string;
  email?: string;
  // 上次密码修改时间（毫秒时间戳，可选）
  lastPasswordChange?: number;
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
  nickname: string;
  role_name: string;
}
