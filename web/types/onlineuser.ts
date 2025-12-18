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
  username: string;
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

// ---- System User & Role Types (from core swagger) ----

// 通用 Options 结构，用于标签等
export interface Options {
  label: string;
  value: string;
  children?: Options[];
}

// 后端的 YesOrNo 状态 [1: 启用, 2: 禁用]
export type YesOrNo = 1 | 2;

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

