<div align="center">
  <img width="120" height="120" alt="Quebec Logo" src="./public/quebec.png" />
  <h1>Quebec</h1>
  <p><strong>高性能网关控制平面仪表板（前端 Web）</strong></p>
</div>

---

本目录（`web/`）是 **Quebec Dashboard** 的前端项目，基于 React + Vite 构建，用于管理 Envoy Gateway 控制平面能力（节点/代理/系统管理等）。

## 功能概览

- **Dashboard**：图表化展示关键指标（流量/延迟/错误率/活跃实例等）。
- **Proxy 管理**：L4/L7 配置与节点列表查看。
- **System 管理**：用户/角色、证书、审计日志、在线用户等。
- **基础体验**：中英文切换、深色模式、响应式布局。

## 技术栈

- **React**: 19.2.x
- **Vite**: 7.3.x
- **TypeScript**: 5.9.x
- **Tailwind CSS**: 4.1.x
- **React Router**: 7.10.x
- **Axios**: 1.13.x
- **Recharts**: 3.6.x
- **Lucide React**: 0.561.x

> 版本以 `web/package.json` 为准。

## 快速开始

### 前置要求

- **Bun**：建议使用最新稳定版
- **后端服务**：需要可访问的 Quebec Core API（默认本地 `http://127.0.0.1:59024`）

### 安装与启动

1. **进入前端目录**

```bash
cd web
```

2. **安装依赖**

```bash
bun install
```

3. **（可选）配置环境变量**

在 `web/` 下创建 `.env.local`：

```env
# API Base（默认值为 /core/api）
VITE_API_BASE_URL=/core/api
```

说明：

- 默认使用相对路径 `/core/api`，配合 Vite 代理转发到后端（推荐开发方式）。
- 也可以配置为绝对地址（例如 `http://127.0.0.1:59024/core/api`），此时可不依赖代理（可能涉及 CORS）。

4. **启动开发服务器**

```bash
bun run dev
```

默认端口：`3000`，并绑定 `0.0.0.0`（局域网可访问）。

> 代理：开发环境会把 `/core/api` 转发到 `vite.config.ts` 里配置的后端目标（默认 `http://127.0.0.1:59024`）。  
> 如需修改后端地址，编辑 `web/vite.config.ts` 的 `targetDict.local`。

## 常用命令（Scripts）

- **开发**：`bun run dev`
- **构建**：`bun run build`
- **预览构建产物**：`bun run preview`

## 关键配置说明

### API Base 与代理

- **API Base**：通过 `VITE_API_BASE_URL` 注入（默认 `/core/api`），代码中由 `web/services/base/http.ts` 使用。
- **本地代理**：`web/vite.config.ts` 将 `/core/api` 代理到 `http://127.0.0.1:59024`（可自行改为测试/开发环境）。

### 构建输出目录

`bun run build` 的产物不会输出到 `web/dist`，而是输出到：

- `cmd/core/internal/web/dist`

这样后端（Core）可以直接把前端静态资源一并打包/部署（具体以 Core 的静态资源加载方式为准）。

## 项目结构

```
web/
├── components/                 # 页面/模块组件
│   ├── auth/                   # 登录
│   ├── dashboard/              # 仪表板
│   ├── layout/                 # 布局组件（Sidebar/通知等）
│   ├── proxy/                  # 代理与节点
│   └── system/                 # 系统管理（用户/角色/证书/日志/在线用户）
├── contexts/                   # React Context（主题/语言）
├── services/                   # API 请求封装与业务服务
│   ├── base/                   # http client / mock / i18n
│   └── system/                 # system 相关服务
├── types/                      # TypeScript 类型
├── public/                     # 静态资源
├── App.tsx                     # App 根组件
├── index.tsx                   # 入口
├── vite.config.ts              # Vite/代理/构建配置
└── package.json                # 依赖与脚本
```
