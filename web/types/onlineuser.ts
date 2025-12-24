import { YesOrNo, OperationType } from './api';

// Re-export for backward compatibility
export type { YesOrNo };

// System Management Types
export interface Certificate {
  id: string;
  domain: string;
  issuer: string;
  expiryDate: string;
  isAutoRenew: boolean;
  status: 'VALID' | 'EXPIRING_SOON' | 'EXPIRED';
}

export interface Notification {
  username: string;
  access_ip: string;
  operation_type: OperationType;
  operation_time: number;
  os: string;
  platform: string;
  browser_name: string;
  browser_engine_name: string;
  browser_version?: string;
  browser_engine_version?: string;
  nickname?: string;
  read?: boolean;
}

// 在线用户响应 (对应 response.SystemOnlineUserResp)
export interface OnlineUser {
  id: string;
  access_ip: string;
  browser_engine_name: string;
  browser_engine_version: string;
  browser_name: string;
  browser_version: string;
  last_operation_time: number;
  username: string;
  nickname: string;
  operation_type: OperationType;
  os: string;
  platform: string;
}

// 在线用户列表响应 (对应 response.SystemOnlineUserListResp)
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

// ---- 操作日志类型 (from core swagger) ----

// 操作日志响应 (对应 response.SystemOperationLogResp)
export interface SystemOperationLog {
  id: string;
  access_ip: string;
  browser_engine_name: string;
  browser_engine_version: string;
  browser_name: string;
  browser_version: string;
  nickname: string;
  operation_time: number;
  operation_type: OperationType;
  os: string;
  platform: string;
  username: string;
}

// 操作日志列表响应 (对应 response.SystemOperationLogListResp)
export interface SystemOperationLogListResponse {
  items: SystemOperationLog[];
  page: number;
  page_size: number;
  total: number;
}

// 操作日志查询请求
export interface SystemOperationLogPageReq {
  id?: string;
  access_ip?: string;
  operation_type?: OperationType;
  user_id?: string;
  start_time?: number;
  end_time?: number;
  page: number;
  page_size: number;
}

// ---- System User & Role Types (from core swagger) ----

// 用户详情响应（对应 response.SystemUserResp）
export interface SystemUserDetail {
  id: string;
  username: string;
  nickname: string;
  email: string;
  role_id: string;
  role_name: string;
  status: YesOrNo;
  remark?: string;
  last_operation_time?: number;
  last_password_change?: number;
  operation_type?: OperationType;
}

// 用户列表响应（对应 response.SystemUserListResp）
export interface SystemUserListResponse {
  items: SystemUserDetail[];
  page: number;
  page_size: number;
  total: number;
}

// 角色详情响应（对应 response.SystemRoleResp）
export interface SystemRoleDetail {
  id: string;
  name: string;
  remark?: string;
  status: YesOrNo;
  users_count?: number;
}

// 角色列表响应（对应 response.SystemRoleListResp）
export interface SystemRoleListResponse {
  items: SystemRoleDetail[];
  page: number;
  page_size: number;
  total: number;
}

// 请求体 - 启停（对应 request.EnableReq）
export interface EnableReq {
  status: YesOrNo;
}

// 请求体 - 添加用户（对应 request.SystemUserAddReq）
export interface SystemUserAddReq {
  email: string;
  nickname: string;
  password: string;
  remark?: string;
  role_id: string;
  status?: YesOrNo;
  username: string;
}

// 请求体 - 更新用户（对应 request.SystemUserUpdateReq）
export interface SystemUserUpdateReq {
  email?: string;
  nickname?: string;
  remark?: string;
  role_id?: string;
  username?: string;
}

// 请求体 - 修改密码（对应 request.SystemUserEditPasswordReq）
export interface SystemUserEditPasswordReq {
  confirm_password: string;
  new_password: string;
  pre_password: string;
}

// 请求体 - 添加/编辑角色（对应 request.SystemRoleAddReq）
export interface SystemRoleAddReq {
  name: string;
  remark?: string;
}
