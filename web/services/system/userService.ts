import { ApiResponse, User } from '../../types';
import { httpClient } from '../base/http';
import {
  SystemUserListResponse,
  SystemUserDetail,
  SystemUserAddReq,
  SystemUserUpdateReq,
  SystemUserEditPasswordReq,
  EnableReq,
  Options,
} from '../../types';

// Service for system user management based on core swagger APIs
export const userService = {
  /**
   * 用户分页列表
   * GET /v1/system/user/page
   */
  async fetchUserPage(params: {
    page: number;
    page_size: number;
    email?: string;
    nickname?: string;
    username?: string;
    role_id?: string;
    status?: number;
  }): Promise<ApiResponse<SystemUserListResponse>> {
    return httpClient.get<SystemUserListResponse>('/v1/system/user/page', { params });
  },

  /**
   * 全部用户列表
   * GET /v1/system/user/list
   */
  async fetchUserList(params?: {
    email?: string;
    nickname?: string;
    username?: string;
    role_id?: string;
    status?: number;
  }): Promise<ApiResponse<SystemUserListResponse>> {
    return httpClient.get<SystemUserListResponse>('/v1/system/user/list', { params });
  },

  /**
   * 用户标签
   * GET /v1/system/user/label
   */
  async fetchUserLabels(): Promise<ApiResponse<Options[]>> {
    return httpClient.get<Options[]>('/v1/system/user/label');
  },

  /**
   * 获取用户详情
   * GET /v1/system/user/{id}
   */
  async getUserDetail(id: string): Promise<ApiResponse<SystemUserDetail>> {
    return httpClient.get<SystemUserDetail>(`/v1/system/user/${id}`);
  },

  /**
   * 添加用户
   * POST /v1/system/user
   */
  async createUser(payload: SystemUserAddReq): Promise<ApiResponse<void>> {
    return httpClient.post<void>('/v1/system/user', payload);
  },

  /**
   * 编辑用户
   * PUT /v1/system/user/{id}
   */
  async updateUser(id: string, payload: SystemUserUpdateReq): Promise<ApiResponse<void>> {
    return httpClient.put<void>(`/v1/system/user/${id}`, payload);
  },

  /**
   * 删除用户
   * DELETE /v1/system/user/{id}
   */
  async deleteUser(id: string): Promise<ApiResponse<void>> {
    return httpClient.delete<void>(`/v1/system/user/${id}`);
  },

  /**
   * 启停用户
   * PUT /v1/system/user/enable/{id}
   */
  async toggleUserStatus(id: string, payload: EnableReq): Promise<ApiResponse<void>> {
    return httpClient.put<void>(`/v1/system/user/enable/${id}`, payload);
  },

  /**
   * 修改用户密码
   * PUT /v1/system/user/password/{id}
   */
  async updateUserPassword(id: string, payload: SystemUserEditPasswordReq): Promise<ApiResponse<void>> {
    return httpClient.put<void>(`/v1/system/user/password/${id}`, payload);
  },
};


