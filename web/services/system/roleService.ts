import { ApiResponse } from '../../types';
import { httpClient } from '../base/http';
import { SystemRoleListResponse, SystemRoleDetail, SystemRoleAddReq, EnableReq, Options } from '../../types';

// Service for system role management based on core swagger APIs
export const roleService = {
  /**
   * 角色分页列表
   * GET /v1/system/role/page
   */
  async fetchRolePage(params: {
    page: number;
    page_size: number;
    name?: string;
    status?: number;
  }): Promise<ApiResponse<SystemRoleListResponse>> {
    return httpClient.get<SystemRoleListResponse>('/v1/system/role/page', { params });
  },

  /**
   * 全部角色列表
   * GET /v1/system/role/list
   */
  async fetchRoleList(params?: { name?: string; status?: number }): Promise<ApiResponse<SystemRoleListResponse>> {
    return httpClient.get<SystemRoleListResponse>('/v1/system/role/list', { params });
  },

  /**
   * 角色标签
   * GET /v1/system/role/label
   */
  async fetchRoleLabels(): Promise<ApiResponse<Options[]>> {
    return httpClient.get<Options[]>('/v1/system/role/label');
  },

  /**
   * 获取角色详情
   * GET /v1/system/role/{id}
   */
  async getRoleDetail(id: string): Promise<ApiResponse<SystemRoleDetail>> {
    return httpClient.get<SystemRoleDetail>(`/v1/system/role/${id}`);
  },

  /**
   * 添加角色
   * POST /v1/system/role
   */
  async createRole(payload: SystemRoleAddReq): Promise<ApiResponse<void>> {
    return httpClient.post<void>('/v1/system/role', payload);
  },

  /**
   * 编辑角色
   * PUT /v1/system/role/{id}
   */
  async updateRole(id: string, payload: SystemRoleAddReq): Promise<ApiResponse<void>> {
    return httpClient.put<void>(`/v1/system/role/${id}`, payload);
  },

  /**
   * 删除角色
   * DELETE /v1/system/role/{id}
   */
  async deleteRole(id: string): Promise<ApiResponse<void>> {
    return httpClient.delete<void>(`/v1/system/role/${id}`);
  },

  /**
   * 启停角色
   * PUT /v1/system/role/enable/{id}
   */
  async toggleRoleStatus(id: string, payload: EnableReq): Promise<ApiResponse<void>> {
    return httpClient.put<void>(`/v1/system/role/enable/${id}`, payload);
  },
};


