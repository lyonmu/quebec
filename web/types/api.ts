// API Response Types
export interface ApiResponse<T> {
  code: number;
  data: T;
  message: string;
}

export interface CaptchaData {
  id: string;
  pictures: string; // base64 image
  length: number;
}

// 通用 Options 结构，用于标签等
export interface Options {
  label: string;
  value: string;
  children?: Options[];
}

// 后端的 YesOrNo 状态 [1: 启用, 2: 禁用]
export type YesOrNo = 1 | 2;

// 菜单类型枚举 [1: 目录, 2: 菜单, 3: 按钮]
export type MenuType = 1 | 2 | 3;

// 操作类型枚举
export type OperationType = 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 |
                            11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19 | 20;

// 操作类型描述映射
export const OperationTypeDescriptions: Record<number, string> = {
  1: '登录操作',
  2: '登出操作',
  3: '创建用户',
  4: '更新用户',
  5: '删除用户',
  6: '创建角色',
  7: '更新角色',
  8: '删除角色',
  9: '创建菜单',
  10: '更新菜单',
  11: '删除菜单',
  12: '踢出在线用户',
  13: '查询操作日志',
  14: '导出操作日志',
  15: '清理操作日志',
  16: '修改密码',
  17: '启用/禁用用户',
  18: '启用/禁用角色',
  19: '启用/禁用菜单',
  20: '角色绑定菜单',
};
