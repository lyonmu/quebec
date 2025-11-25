
import { EnvoyNode, Route, Certificate, Protocol, DashboardMetrics, Notification, User, OperationLog } from '../../types';

export const generateMetrics = (): DashboardMetrics => {
  const now = new Date();
  const rps = [];
  const latency = [];
  const errorRate = [];

  for (let i = 20; i >= 0; i--) {
    const time = new Date(now.getTime() - i * 60000).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    rps.push({ time, value: Math.floor(2000 + Math.random() * 1500) });
    latency.push({ time, value: Math.floor(20 + Math.random() * 40) });
    errorRate.push({ time, value: Number((Math.random() * 2).toFixed(2)) });
  }

  return { rps, latency, errorRate };
};

export const mockNodes: EnvoyNode[] = [
  { id: 'envoy-gw-01', ip: '10.0.1.15', version: 'v1.29.1', status: 'HEALTHY', uptime: '14d 2h', region: 'us-east-1', connections: 12450 },
  { id: 'envoy-gw-02', ip: '10.0.1.16', version: 'v1.29.1', status: 'HEALTHY', uptime: '14d 2h', region: 'us-east-1', connections: 11890 },
  { id: 'envoy-gw-03', ip: '10.0.2.40', version: 'v1.28.0', status: 'DEGRADED', uptime: '2d 5h', region: 'us-west-2', connections: 4500 },
  { id: 'envoy-gw-04', ip: '10.0.2.41', version: 'v1.29.1', status: 'HEALTHY', uptime: '5h 12m', region: 'us-west-2', connections: 8900 },
];

export const mockRoutes: Route[] = [
  { id: 'rt-001', name: 'payment-service-v1', prefix: '/api/v1/payments', targetCluster: 'payment-prod', protocol: Protocol.HTTP, timeout: '5s', status: 'ACTIVE' },
  { id: 'rt-002', name: 'user-service-grpc', prefix: '/proto.UserService', targetCluster: 'user-prod', protocol: Protocol.GRPC, timeout: '2s', status: 'ACTIVE' },
  { id: 'rt-003', name: 'legacy-monolith', prefix: '/legacy', targetCluster: 'monolith-vms', protocol: Protocol.HTTP, timeout: '30s', status: 'ACTIVE', trafficSplit: 100 },
  { id: 'rt-004', name: 'canary-search', prefix: '/api/v2/search', targetCluster: 'search-canary', protocol: Protocol.HTTPS, timeout: '1s', status: 'INACTIVE' },
  { id: 'rt-005', name: 'mysql-read-replica', prefix: '0.0.0.0:3306', targetCluster: 'db-cluster-ro', protocol: Protocol.TCP, timeout: '10s', status: 'ACTIVE' },
];

export const mockCerts: Certificate[] = [
  { id: 'crt-123', domain: '*.api.company.com', issuer: 'Let\'s Encrypt', expiryDate: '2024-12-01', isAutoRenew: true, status: 'VALID' },
  { id: 'crt-456', domain: 'admin.internal', issuer: 'Internal CA', expiryDate: '2023-11-15', isAutoRenew: false, status: 'EXPIRING_SOON' },
  { id: 'crt-789', domain: 'legacy.gateway.net', issuer: 'DigiCert', expiryDate: '2023-09-01', isAutoRenew: false, status: 'EXPIRED' },
];

export const mockUsers: User[] = [
  { id: 'u-001', username: 'admin', role: 'ADMIN', status: 'ACTIVE', lastLogin: '2024-03-10 09:42:00', email: 'admin@company.com' },
  { id: 'u-002', username: 'developer_lead', role: 'EDITOR', status: 'ACTIVE', lastLogin: '2024-03-09 14:20:00', email: 'dev@company.com' },
  { id: 'u-003', username: 'auditor', role: 'VIEWER', status: 'ACTIVE', lastLogin: '2024-03-01 10:00:00', email: 'audit@company.com' },
  { id: 'u-004', username: 'temp_contractor', role: 'VIEWER', status: 'INACTIVE', lastLogin: '2023-12-20 16:55:00' },
];

export const mockLogs: OperationLog[] = [
  { id: 'log-001', user: 'admin', action: 'UPDATE_ROUTE', target: 'payment-service-v1', ip: '10.20.1.5', status: 'SUCCESS', timestamp: '2024-03-10 14:32:15', details: 'Updated timeout configuration' },
  { id: 'log-002', user: 'developer_lead', action: 'DEPLOY_NODE', target: 'envoy-gw-05', ip: '10.20.1.12', status: 'FAILURE', timestamp: '2024-03-10 11:20:00', details: 'Connection timeout during bootstrap' },
  { id: 'log-003', user: 'admin', action: 'DELETE_USER', target: 'temp_contractor', ip: '10.20.1.5', status: 'SUCCESS', timestamp: '2024-03-09 16:45:10', details: 'Removed inactive user' },
  { id: 'log-004', user: 'auditor', action: 'EXPORT_REPORT', target: 'Monthly Audit', ip: '10.20.2.8', status: 'SUCCESS', timestamp: '2024-03-09 09:15:00', details: 'Downloaded PDF report' },
  { id: 'log-005', user: 'admin', action: 'RENEW_CERT', target: '*.api.company.com', ip: '10.20.1.5', status: 'SUCCESS', timestamp: '2024-03-08 10:00:22', details: 'Manual renewal triggered' },
  { id: 'log-006', user: 'developer_lead', action: 'ADD_ROUTE', target: 'user-service-grpc', ip: '10.20.1.12', status: 'SUCCESS', timestamp: '2024-03-08 08:45:00', details: 'Added new gRPC route' },
  { id: 'log-007', user: 'admin', action: 'LOGIN', target: 'System', ip: '10.20.1.5', status: 'SUCCESS', timestamp: '2024-03-08 08:00:00', details: 'User logged in successfully' },
];

export const fetchNotifications = (page: number, pageSize: number = 10): Promise<{ data: Notification[], hasMore: boolean }> => {
  return new Promise((resolve) => {
    setTimeout(() => {
      const totalMessages = 45;
      const start = (page - 1) * pageSize;
      const end = start + pageSize;
      
      const data: Notification[] = [];
      if (start < totalMessages) {
        for (let i = start; i < Math.min(end, totalMessages); i++) {
          const typeRandom = Math.random();
          let type: 'info' | 'warning' | 'error' | 'success' = 'info';
          if (typeRandom > 0.9) type = 'error';
          else if (typeRandom > 0.75) type = 'warning';
          else if (typeRandom > 0.6) type = 'success';

          data.push({
            id: `notif-${i}`,
            title: type === 'error' ? 'Connection Error' : type === 'warning' ? 'High Latency Warning' : type === 'success' ? 'Deployment Successful' : 'System Info',
            message: `System notification message #${i + 1}. This contains details about the event that occurred in the cluster.`,
            timestamp: new Date(Date.now() - i * 3600000 * 2).toISOString(), // Spread out over time
            read: i > 5,
            type
          });
        }
      }

      resolve({
        data,
        hasMore: end < totalMessages
      });
    }, 800); // Simulate network latency
  });
};
