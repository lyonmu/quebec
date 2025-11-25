// System Management Types
export interface Certificate {
  id: string;
  domain: string;
  issuer: string;
  expiryDate: string;
  isAutoRenew: boolean;
  status: 'VALID' | 'EXPIRING_SOON' | 'EXPIRED';
}

export interface OperationLog {
  id: string;
  user: string;
  action: string;
  target: string;
  ip: string;
  status: 'SUCCESS' | 'FAILURE';
  timestamp: string;
  details: string;
}

export interface Notification {
  id: string;
  title: string;
  message: string;
  timestamp: string;
  read: boolean;
  type: 'info' | 'warning' | 'error' | 'success';
}

export interface OnlineUser {
  id: string;
  access_ip: string;
  browser_engine_name: string;
  browser_engine_version: string;
  browser_name: string;
  browser_version: string;
  last_operation_time: number;
  nickname: string;
  operation_type: number;
  os: string;
  platform: string;
}

export interface OnlineUserListResponse {
  items: OnlineUser[];
  page: number;
  page_size: number;
  total: number;
}

export interface OnlineUserLabel {
  label: string;
  value: string;
  children?: OnlineUserLabel[];
}
