import type { ApiResponse } from '../../types';
import type { SystemOperationLog, SystemOperationLogListResponse, SystemOperationLogPageReq } from '../../types/onlineuser';
import { httpClient } from '../base/http';

// Re-export for backward compatibility
export type { SystemOperationLog as OperationLogItem };
export type { SystemOperationLogListResponse as OperationLogListResponse };

export const operationLogService = {
  /**
   * 操作日志分页列表
   * GET /v1/system/operation-log/page
   */
  async fetchLogPage(params: Partial<SystemOperationLogPageReq>): Promise<ApiResponse<SystemOperationLogListResponse>> {
    return httpClient.get<SystemOperationLogListResponse>('/v1/system/operation-log/page', { params });
  },
};
