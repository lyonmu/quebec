import { ApiResponse } from '../../types';
import { httpClient } from '../base/http';

// 操作日志响应类型
export interface OperationLogItem {
  id: string;
  user_id: string;
  username: string;
  operation_type: number;
  operation_action: string;
  target_id: string;
  target_type: string;
  ip: string;
  user_agent: string;
  details: string;
  created_at: number;
}

export interface OperationLogListResponse {
  items: OperationLogItem[];
  page: number;
  page_size: number;
  total: number;
}

export interface OperationLogPageReq {
  user_id?: string;
  operation_type?: number;
  target_type?: string;
  page?: number;
  page_size?: number;
}

// Service for operation log management based on core APIs
export const operationLogService = {
  /**
   * 操作日志分页列表
   * GET /v1/system/operation-log/page
   */
  async fetchLogPage(params: OperationLogPageReq): Promise<ApiResponse<OperationLogListResponse>> {
    return httpClient.get<OperationLogListResponse>('/v1/system/operation-log/page', { params });
  },
};
