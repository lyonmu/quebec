import { YesOrNo, MenuType } from './api';

// Re-export for backward compatibility
export type { MenuType };

// Menu Types
export interface SystemMenu {
  id: string;
  name: string;
  menu_type: MenuType;
  api_path: string;
  api_path_method: string;
  order: number;
  parent_id: string;
  parent_name: string;
  component: string;
  status: YesOrNo;
  remark: string;
}

export interface SystemMenuPageReq {
  name?: string;
  menu_type?: MenuType;
  status?: YesOrNo;
  parent_id?: string;
  page: number;
  page_size: number;
}

export interface SystemMenuListReq {
  name?: string;
  menu_type?: MenuType;
  status?: YesOrNo;
}

export interface SystemMenuAddReq {
  name: string;
  menu_type: MenuType;
  api_path?: string;
  api_path_method?: string;
  order?: number;
  parent_id?: string;
  component?: string;
  status?: YesOrNo;
  remark?: string;
}

export interface SystemMenuUpdateReq {
  name?: string;
  menu_type?: MenuType;
  api_path?: string;
  api_path_method?: string;
  order?: number;
  parent_id?: string;
  component?: string;
  status?: YesOrNo;
  remark?: string;
}

export interface SystemMenuListResponse {
  total: number;
  items: SystemMenu[];
  page: number;
  page_size: number;
}

export interface SystemMenuTreeItem {
  id: string;
  name: string;
  menu_type: MenuType;
  api_path: string;
  api_path_method: string;
  order: number;
  component: string;
  status: YesOrNo;
  children: SystemMenuTreeItem[];
}

export interface MenuOptions {
  label: string;
  value: string;
}

// Role Menu Binding Types
export interface SystemRoleMenu {
  menu_id: string;
  menu_name: string;
}

export interface SystemRoleMenuBindReq {
  menu_ids: string[];
}
