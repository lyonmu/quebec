// 操作类型映射常量
export const OPERATION_TYPE_MAP: Record<number, string> = {
  1: 'online_users.operation_types.login',
  // 可以根据实际需求添加更多操作类型
  // 2: 'online_users.operation_types.logout',
  // 3: 'online_users.operation_types.update',
};

// 默认操作类型
export const DEFAULT_OPERATION_TYPE = 'online_users.operation_types.other';

export const translations = {
  zh: {
    common: {
      search: "搜索...",
      loading: "加载中...",
      total: "共",
      items: "条",
      prev: "上一页",
      next: "下一页",
      actions: "操作",
      startTime: "开始时间",
      endTime: "结束时间",
      delete: "删除",
      select: "请选择",
      edit: "编辑",
      view: "详情",
      error: "请求出错"
    },
    dashboard: {
      title: "仪表盘",
      totalReq: "总请求数/秒",
      avgLatency: "平均延迟 (P99)",
      errorRate: "错误率 (5xx)",
      activeInstances: "活跃实例",
      activeConns: "活跃连接",
      trafficVolume: "流量趋势 (RPS)",
      latency: "延迟趋势 (ms)",
      spikeDetected: "检测到峰值",
      vsLastHour: "对比上小时",
      requests: "请求",
    },
    nodes: {
      title: "网关实例",
      deployNew: "部署新实例",
      connections: "连接数",
      uptime: "运行时间",
      manage: "管理",
      status: {
        healthy: "健康",
        degraded: "降级",
        draining: "排空"
      }
    },
    proxy: {
      title: "网关管理",
      titleL4: "L4 代理规则",
      titleL7: "L7 代理规则",
      subtitleL4: "管理 TCP/UDP 层面的流量转发。",
      subtitleL7: "管理 HTTP/GRPC 应用层路由。",
      searchPlaceholder: "搜索路由...",
      addRoute: "添加路由",
      table: {
        name: "路由名称",
        prefix: "前缀 / 端口",
        cluster: "上游集群",
        protocol: "协议",
        timeout: "超时",
        status: "状态",
        actions: "操作"
      },
      status: {
        active: "活跃",
        inactive: "非活跃"
      }
    },
    certs: {
      title: "证书管理",
      subtitle: "管理 TLS/SSL 上下文和密钥。",
      import: "导入证书",
      issuer: "颁发者",
      expires: "过期时间",
      autoRenew: "自动续期",
      renewNow: "立即续期",
      status: {
        valid: "有效",
        expiring_soon: "即将过期",
        expired: "已过期"
      }
    },
    users: {
      title: "用户管理",
      subtitle: "管理系统访问权限与用户账号。",
      addUser: "新增用户",
      table: {
        username: "用户名",
        role: "角色",
        status: "状态",
        lastLogin: "最后登录",
        lastPasswordChange: "上次密码修改时间",
        actions: "操作"
      },
      status: {
        active: "正常",
        inactive: "停用"
      },
      roles: {
        admin: "管理员",
        editor: "编辑员",
        viewer: "观察员"
      },
      modal: {
        addTitle: "新增用户",
        editPassTitle: "修改密码",
        editUserTitle: "编辑用户",
        detailTitle: "用户详情",
        username: "用户名",
        nickname: "昵称",
        password: "密码",
        newPassword: "新密码",
        confirmPassword: "确认密码",
        prePassword: "当前密码",
        role: "角色",
        email: "邮箱",
        cancel: "取消",
        save: "保存",
        update: "更新",
        passMismatch: "两次输入的密码不一致",
        roleRequired: "请选择角色"
      },
      actions: {
        disable: "停用",
        enable: "启用",
        changePass: "重置密码",
        delete: "删除"
      }
    },
    online_users: {
      title: "在线用户",
      subtitle: "查看当前在线用户及其会话信息。",
      table: {
        username: "用户名称",
        ip: "访问 IP",
        lastTime: "最后操作时间",
        type: "操作类型",
        os: "操作系统",
        platform: "操作平台",
        browser: "浏览器名称",
        version: "浏览器版本",
        engine: "引擎名称",
        engineVersion: "引擎版本"
      },
      confirm_clear: "确定要清除此在线用户吗？该会话将被强制下线。",
      operation_types: {
        login: "登录",
        other: "其他"
      }
    },
    roles: {
      title: "角色管理",
      subtitle: "定义系统角色及其操作权限。",
      addRole: "新增角色",
      confirmDelete: "确定要删除该角色吗？",
      modal: {
        addTitle: "新增角色",
        editTitle: "编辑角色",
        name: "角色名称",
        remark: "备注",
        cancel: "取消",
        save: "保存"
      },
      table: {
        name: "角色名称",
        users: "关联用户数",
        permissions: "权限范围",
        actions: "操作"
      },
      permissions: {
        all: "完全控制",
        read_write: "读写权限",
        read_only: "只读权限"
      }
    },
    logs: {
      title: "操作日志",
      subtitle: "审计用户操作行为与系统变更记录。",
      table: {
        action: "操作行为",
        user: "执行用户",
        ip: "来源 IP",
        target: "目标对象",
        status: "结果",
        time: "操作时间",
        details: "详情"
      },
      status: {
        success: "成功",
        failure: "失败"
      }
    },
    sidebar: {
      dashboard: "仪表盘",
      proxy: "网关管理",
      proxy_l4: "L4 代理",
      proxy_l7: "L7 代理",
      nodes: "实例列表",
      certs: "证书管理",
      system: "系统管理",
      system_users: "用户管理",
      system_online_users: "在线用户",
      system_roles: "角色管理",
      system_logs: "操作日志",
      controlPlane: "控制平面",
      healthy: "健康",
      signOut: "退出登录"
    },
    header: {
      platform: "平台",
      adminUser: "管理员",
      superAdmin: "超级管理员"
    },
    ai: {
      button: "AI 配置生成",
      title: "Quebec AI 架构师",
      placeholder: "描述您需要的路由、集群或监听器...",
      hint: '例如："创建一个指向 backend-svc 集群的 HTTP 路由，路径为 /api/v1，超时时间 5s"',
      inputPlaceholder: "让 AI 生成配置...",
      copy: "复制 YAML",
      copied: "已复制",
      error: "生成配置出错。",
      genButton: "发送"
    },
    login: {
      username: "用户名",
      password: "密码",
      enterCaptcha: "请输入验证码",
      signIn: "登录",
      signingIn: "登录中...",
      continueToQuebec: "可视化 Envoy 控制平面"
    },
    settings: {
      title: "系统管理",
      users: "用户管理",
      manage: "系统设置",
      placeholder: "设置模块占位符"
    },
    notifications: {
      title: "消息通知",
      empty: "暂无消息",
      noMore: "没有更多消息了",
      loading: "加载中..."
    }
  },
  en: {
    common: {
      search: "Search...",
      loading: "Loading...",
      total: "Total",
      items: "items",
      prev: "Previous",
      next: "Next",
      actions: "Actions",
      startTime: "Start Time",
      endTime: "End Time",
      delete: "Delete",
      select: "Select",
      edit: "Edit",
      view: "View",
      error: "Request error"
    },
    dashboard: {
      title: "Dashboard",
      totalReq: "Total Requests/Sec",
      avgLatency: "Avg Latency (P99)",
      errorRate: "Error Rate (5xx)",
      activeInstances: "Active Instances",
      activeConns: "active conns",
      trafficVolume: "Traffic Volume (RPS)",
      latency: "Latency (ms)",
      spikeDetected: "spike detected",
      vsLastHour: "vs last hour",
      requests: "requests",
    },
    nodes: {
      title: "Gateway Instances",
      deployNew: "Deploy New Instance",
      connections: "Connections",
      uptime: "Up",
      manage: "Manage",
       status: {
        healthy: "HEALTHY",
        degraded: "DEGRADED",
        draining: "DRAINING"
      }
    },
    proxy: {
      title: "Gateway Management",
      titleL4: "L4 Proxy Rules",
      titleL7: "L7 Proxy Rules",
      subtitleL4: "Manage TCP/UDP traffic forwarding.",
      subtitleL7: "Manage HTTP/GRPC application routing.",
      searchPlaceholder: "Search routes...",
      addRoute: "Add Route",
      table: {
        name: "Route Name",
        prefix: "Prefix / Port",
        cluster: "Upstream Cluster",
        protocol: "Protocol",
        timeout: "Timeout",
        status: "Status",
        actions: "Actions"
      },
      status: {
        active: "ACTIVE",
        inactive: "INACTIVE"
      }
    },
    certs: {
      title: "Certificates",
      subtitle: "Manage TLS/SSL contexts and secrets.",
      import: "Import Cert",
      issuer: "Issuer",
      expires: "Expires",
      autoRenew: "Auto-Renew",
      renewNow: "Renew Now",
       status: {
        valid: "VALID",
        expiring_soon: "EXPIRING SOON",
        expired: "EXPIRED"
      }
    },
    users: {
      title: "User Management",
      subtitle: "Manage system access and user accounts.",
      addUser: "Add User",
      table: {
        username: "Username",
        role: "Role",
        status: "Status",
        lastLogin: "Last Login",
        lastPasswordChange: "Last Password Change",
        actions: "Actions"
      },
      status: {
        active: "Active",
        inactive: "Inactive"
      },
      roles: {
        admin: "Admin",
        editor: "Editor",
        viewer: "Viewer"
      },
      modal: {
        addTitle: "Add New User",
        editPassTitle: "Change Password",
        editUserTitle: "Edit User",
        detailTitle: "User Details",
        username: "Username",
        nickname: "Nickname",
        password: "Password",
        newPassword: "New Password",
        confirmPassword: "Confirm Password",
        prePassword: "Current Password",
        role: "Role",
        email: "Email",
        cancel: "Cancel",
        save: "Save",
        update: "Update",
        passMismatch: "Passwords do not match",
        roleRequired: "Please select a role"
      },
      actions: {
        disable: "Disable",
        enable: "Enable",
        changePass: "Reset Password",
        delete: "Delete"
      }
    },
    online_users: {
      title: "Online Users",
      subtitle: "View currently online users and session info.",
      table: {
        username: "Username",
        ip: "Access IP",
        lastTime: "Last Operation Time",
        type: "Operation Type",
        os: "OS",
        platform: "Platform",
        browser: "Browser Name",
        version: "Browser Version",
        engine: "Engine Name",
        engineVersion: "Engine Version"
      },
      confirm_clear: "Are you sure you want to clear this online user? The session will be terminated.",
      operation_types: {
        login: "LOGIN",
        other: "OTHER"
      }
    },
    roles: {
      title: "Role Management",
      subtitle: "Define system roles and their permissions.",
      addRole: "Add Role",
      confirmDelete: "Are you sure you want to delete this role?",
      modal: {
        addTitle: "Add Role",
        editTitle: "Edit Role",
        name: "Role Name",
        remark: "Remark",
        cancel: "Cancel",
        save: "Save"
      },
      table: {
        name: "Role Name",
        users: "Assigned Users",
        permissions: "Permissions Scope",
        actions: "Actions"
      },
      permissions: {
        all: "Full Control",
        read_write: "Read & Write",
        read_only: "Read Only"
      }
    },
    logs: {
      title: "Operation Logs",
      subtitle: "Audit user actions and system changes.",
      table: {
        action: "Action",
        user: "User",
        ip: "Source IP",
        target: "Target",
        status: "Status",
        time: "Timestamp",
        details: "Details"
      },
      status: {
        success: "Success",
        failure: "Failure"
      }
    },
    sidebar: {
      dashboard: "Dashboard",
      proxy: "Gateway Management",
      proxy_l4: "L4 Proxy",
      proxy_l7: "L7 Proxy",
      nodes: "Instances",
      certs: "Certificates",
      system: "System Management",
      system_users: "User Management",
      system_online_users: "Online Users",
      system_roles: "Role Management",
      system_logs: "Operation Logs",
      controlPlane: "Control Plane",
      healthy: "Healthy",
      signOut: "Sign Out"
    },
    header: {
      platform: "Platform",
      adminUser: "Admin User",
      superAdmin: "Super Admin"
    },
    ai: {
      button: "AI Config Gen",
      title: "Quebec AI Architect",
      placeholder: "Describe the route, cluster, or listener you need.",
      hint: 'e.g., "Create an HTTP route for /api/v1 with a 5s timeout pointing to cluster backend-svc"',
      inputPlaceholder: "Ask AI to generate config...",
      copy: "Copy YAML",
      copied: "Copied",
      error: "Error generating configuration.",
      genButton: "Send"
    },
    login: {
      username: "Username",
      password: "Password",
      enterCaptcha: "Enter captcha",
      signIn: "Sign in",
      signingIn: "Signing in...",
      continueToQuebec: "Visualizing the Envoy control plane"
    },
    settings: {
      title: "System Management",
      users: "User Management",
      manage: "System Settings",
      placeholder: "Settings module placeholder"
    },
    notifications: {
      title: "Notifications",
      empty: "No notifications",
      noMore: "No more messages",
      loading: "Loading..."
    }
  }
};
