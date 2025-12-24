
import { EnvoyNode, Route, Certificate, Protocol, DashboardMetrics, Notification } from '../../types';
import { OperationType } from '../../types/api';

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



export const fetchNotifications = (page: number, pageSize: number = 10): Promise<{ data: Notification[], hasMore: boolean }> => {
  return new Promise((resolve) => {
    setTimeout(() => {
      const totalMessages = 45;
      const start = (page - 1) * pageSize;
      const end = start + pageSize;

      const data: Notification[] = [];
      if (start < totalMessages) {
        for (let i = start; i < Math.min(end, totalMessages); i++) {
          // 随机选择操作类型
          const operationTypes = [1, 2, 3, 4, 5, 6, 7, 16, 17, 18];
          const randomOpType = operationTypes[Math.floor(Math.random() * operationTypes.length)];

          // 随机用户名生成
          const usernames = ['admin', 'operator', 'viewer', 'editor', 'user123', 'manager'];
          const randomUsername = usernames[Math.floor(Math.random() * usernames.length)];

          // 随机操作系统
          const osList = ['Windows 10', 'macOS 14', 'Ubuntu 22.04', 'CentOS 7', 'Android 13', 'iOS 17'];
          const randomOS = osList[Math.floor(Math.random() * osList.length)];

          // 随机浏览器
          const browserList = ['Chrome', 'Firefox', 'Safari', 'Edge', 'Opera'];
          const randomBrowser = browserList[Math.floor(Math.random() * browserList.length)];

          // 随机引擎
          const engineList = ['Blink', 'Gecko', 'WebKit', 'Chromium'];
          const randomEngine = engineList[Math.floor(Math.random() * engineList.length)];

          // 随机 IP 生成
          const randomIP = `${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}.${Math.floor(Math.random() * 255)}`;

          data.push({
            username: randomUsername,
            nickname: `${randomUsername}_nickname`,
            access_ip: randomIP,
            operation_type: randomOpType as OperationType,
            operation_time: Math.floor(Date.now() / 1000) - (i * 3600),
            os: randomOS,
            platform: randomOS.includes('Android') || randomOS.includes('iOS') ? 'Mobile' : 'Desktop',
            browser_name: randomBrowser,
            browser_version: `${Math.floor(Math.random() * 100) + 90}.0.${Math.floor(Math.random() * 5000)}`,
            browser_engine_name: randomEngine,
            browser_engine_version: `${Math.floor(Math.random() * 100) + 90}.0.${Math.floor(Math.random() * 5000)}`,
            read: i > 5
          });
        }
      }

      resolve({
        data,
        hasMore: end < totalMessages
      });
    }, 800);
  });
};
