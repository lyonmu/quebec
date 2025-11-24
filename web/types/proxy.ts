// Proxy and Envoy Types
export enum Protocol {
  HTTP = 'HTTP',
  HTTPS = 'HTTPS',
  TCP = 'TCP',
  GRPC = 'GRPC'
}

export interface EnvoyNode {
  id: string;
  ip: string;
  version: string;
  status: 'HEALTHY' | 'DEGRADED' | 'DRAINING';
  uptime: string;
  region: string;
  connections: number;
}

export interface Route {
  id: string;
  name: string;
  prefix: string;
  targetCluster: string;
  protocol: Protocol;
  timeout: string;
  status: 'ACTIVE' | 'INACTIVE';
  trafficSplit?: number;
}
