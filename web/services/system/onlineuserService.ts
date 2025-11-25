import { ApiResponse, OnlineUserListResponse, OnlineUserLabel } from '../../types';
import { httpClient } from '../base/http';

export const onlineuserService = {
  /**
   * Get online users list
   * @param params - Query parameters including page, page_size, etc.
   * @returns Online users list response
   */
  async fetchOnlineUsers(params: {
    page: number;
    page_size: number;
    access_ip?: string;
    last_operation_time?: number;
    user_id?: string;
    start_time?: number;
    end_time?: number;
  }): Promise<ApiResponse<OnlineUserListResponse>> {
    return httpClient.get<OnlineUserListResponse>('/v1/system/onlineuser/list', { params });
  },

  /**
   * Get online user labels
   * @returns List of online user labels
   */
  async fetchOnlineUserLabels(): Promise<ApiResponse<OnlineUserLabel[]>> {
    return httpClient.get<OnlineUserLabel[]>('/v1/system/onlineuser/label');
  },

  /**
   * Clear (kick) an online user
   * @param id - User ID to clear
   * @returns Success response
   */
  async clearOnlineUser(id: string): Promise<ApiResponse<void>> {
    return httpClient.delete<void>(`/v1/system/onlineuser/clearance/${id}`);
  },
};
