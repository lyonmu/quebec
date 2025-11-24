// Dashboard and Metrics Types
export interface MetricPoint {
  time: string;
  value: number;
}

export interface DashboardMetrics {
  rps: MetricPoint[];
  latency: MetricPoint[];
  errorRate: MetricPoint[];
}

export enum ViewState {
  DASHBOARD = 'DASHBOARD',
  NODES = 'NODES',
  PROXY_L4 = 'PROXY_L4',
  PROXY_L7 = 'PROXY_L7',
  CERTS = 'CERTS',
  SYSTEM_USERS = 'SYSTEM_USERS',
  SYSTEM_ONLINE_USERS = 'SYSTEM_ONLINE_USERS',
  SYSTEM_ROLES = 'SYSTEM_ROLES',
  SYSTEM_LOGS = 'SYSTEM_LOGS'
}
