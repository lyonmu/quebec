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
