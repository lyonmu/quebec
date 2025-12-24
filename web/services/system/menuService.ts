import { ApiResponse } from '../../types';
import { httpClient } from '../base/http';
import {
  SystemMenu,
  SystemMenuListResponse,
  SystemMenuPageReq,
  SystemMenuListReq,
  SystemMenuAddReq,
  SystemMenuUpdateReq,
  SystemMenuTreeItem,
  MenuOptions,
} from '../../types';

// Re-export roleService for convenience
export { roleService } from './roleService';

// Service for system menu management based on core swagger APIs
export const menuService = {
  /**
   * 菜单分页列表
   * GET /v1/system/menu/page
   */
  async fetchMenuPage(params: SystemMenuPageReq): Promise<ApiResponse<SystemMenuListResponse>> {
    return httpClient.get<SystemMenuListResponse>('/v1/system/menu/page', { params });
  },

  /**
   * 全部菜单列表
   * GET /v1/system/menu/list
   */
  async fetchMenuList(params?: SystemMenuListReq): Promise<ApiResponse<SystemMenu[]>> {
    return httpClient.get<SystemMenu[]>('/v1/system/menu/list', { params });
  },

  /**
   * 菜单树形列表
   * GET /v1/system/menu/tree
   */
  async fetchMenuTree(): Promise<ApiResponse<SystemMenuTreeItem[]>> {
    return httpClient.get<SystemMenuTreeItem[]>('/v1/system/menu/tree');
  },

  /**
   * 菜单标签
   * GET /v1/system/menu/label
   */
  async fetchMenuLabels(): Promise<ApiResponse<MenuOptions[]>> {
    return httpClient.get<MenuOptions[]>('/v1/system/menu/label');
  },

  /**
   * 获取菜单详情
   * GET /v1/system/menu/{id}
   */
  async getMenuDetail(id: string): Promise<ApiResponse<SystemMenu>> {
    return httpClient.get<SystemMenu>(`/v1/system/menu/${id}`);
  },

  /**
   * 添加菜单
   * POST /v1/system/menu
   */
  async createMenu(payload: SystemMenuAddReq): Promise<ApiResponse<void>> {
    return httpClient.post<void>('/v1/system/menu', payload);
  },

  /**
   * 编辑菜单
   * PUT /v1/system/menu/{id}
   */
  async updateMenu(id: string, payload: SystemMenuUpdateReq): Promise<ApiResponse<void>> {
    return httpClient.put<void>(`/v1/system/menu/${id}`, payload);
  },

  /**
   * 删除菜单
   * DELETE /v1/system/menu/{id}
   */
  async deleteMenu(id: string): Promise<ApiResponse<void>> {
    return httpClient.delete<void>(`/v1/system/menu/${id}`);
  },

  /**
   * 启停菜单
   * PUT /v1/system/menu/enable/{id}
   */
  async toggleMenuStatus(id: string, payload: { status: number }): Promise<ApiResponse<void>> {
    return httpClient.put<void>(`/v1/system/menu/enable/${id}`, payload);
  },
};

/**
 * Service for role-menu binding management
 */
export const roleMenuService = {
  /**
   * 获取角色菜单列表
   * GET /v1/system/role/{id}/menus
   */
  async getRoleMenus(roleId: string): Promise<ApiResponse<Array<{ menu_id: string; menu_name: string }>>> {
    return httpClient.get<Array<{ menu_id: string; menu_name: string }>>(`/v1/system/role/${roleId}/menus`);
  },

  /**
   * 绑定角色菜单
   * PUT /v1/system/role/{id}/menus
   */
  async bindRoleMenus(roleId: string, menuIds: string[]): Promise<ApiResponse<void>> {
    return httpClient.put<void>(`/v1/system/role/${roleId}/menus`, { menu_ids: menuIds });
  },

  /**
   * 添加菜单到角色
   * POST /v1/system/role/{role_id}/menu/{menu_id}
   */
  async addMenuToRole(roleId: string, menuId: string): Promise<ApiResponse<void>> {
    return httpClient.post<void>(`/v1/system/role/${roleId}/menu/${menuId}`);
  },

  /**
   * 从角色移除菜单
   * DELETE /v1/system/role/{role_id}/menu/{menu_id}
   */
  async removeMenuFromRole(roleId: string, menuId: string): Promise<ApiResponse<void>> {
    return httpClient.delete<void>(`/v1/system/role/${roleId}/menu/${menuId}`);
  },
};
