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

export type MenuType = 1 | 2 | 3; // 1: 目录, 2: 菜单, 3: 按钮

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
