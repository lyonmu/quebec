<div align="center">
  <img width="120" height="120" alt="Quebec Logo" src="./public/quebec.png" />
  <h1>Quebec</h1>
  <p><strong>高性能网关控制平面仪表板</strong></p>
</div>

---

**Quebec** 是一个现代化的高性能控制平面仪表板，专为管理 Envoy Gateway 而设计。它提供实时指标监控、直观的 L4/L7 代理配置，以及 AI 辅助的 XDS 配置生成，所有功能都封装在一个优雅、响应式的界面中。

## ✨ 核心功能

- **实时仪表板**：通过动态图表可视化流量、延迟、错误率和活跃实例。
- **代理管理**：
  - **L4 (TCP)**：管理 TCP 路由和监听器。
  - **L7 (HTTP/GRPC)**：轻松配置 HTTP/HTTPS/GRPC 路由。
- **系统管理**：
  - **用户与角色管理**：基于 RBAC 的细粒度访问控制。
  - **证书管理**：监控和续期 SSL/TLS 证书。
  - **审计日志**：跟踪所有系统操作，确保安全性和合规性。
  - **在线用户**：实时查看和管理在线用户。
- **AI 助手**：内置 AI 助手，可生成 Envoy 配置并回答运维问题。
- **多语言支持**：支持中英文无缝切换。
- **深色模式**：完整支持深色模式，适合各种环境下的舒适查看。

## 🛠️ 技术栈

### 前端

- **框架**：[React](https://react.dev/) (v19.2.0)
- **构建工具**：[Vite](https://vitejs.dev/) (v6.2.0)
- **样式**：[Tailwind CSS](https://tailwindcss.com/) (v4.1.17)
- **图标**：[Lucide React](https://lucide.dev/) (v0.554.0)
- **图表**：[Recharts](https://recharts.org/) (v3.4.1)
- **路由**：[React Router](https://reactrouter.com/) (v7.9.6)
- **HTTP 客户端**：[Axios](https://axios-http.com/) (v1.13.2)
- **AI 集成**：[Google Generative AI](https://ai.google.dev/) (v1.30.0)
- **语言**：TypeScript (v5.8.2)

### 后端

- **语言**：Go (v1.24.6)
- **Web 框架**：[Gin](https://gin-gonic.com/) (v1.11.0)
- **ORM**：[Ent](https://entgo.io/) (v0.14.5)
- **控制平面**：[Envoy Go Control Plane](https://github.com/envoyproxy/go-control-plane) (v0.14.0)
- **数据库**：MySQL
- **缓存**：Redis (v9.17.0)
- **日志**：[Zap](https://github.com/uber-go/zap) (v1.27.0)
- **监控**：[Prometheus](https://prometheus.io/) (v1.23.2)

## 🚀 快速开始

### 前置要求

- **Node.js** (v18 或更高版本)
- **npm**、**yarn** 或 **bun**
- **Go** (v1.24.6 或更高版本，用于后端开发)
- **MySQL** (用于数据存储)
- **Redis** (用于缓存和会话管理)

### 安装步骤

1. **克隆仓库**

   ```bash
   git clone https://github.com/lyonmu/quebec.git
   cd quebec/web
   ```

2. **安装依赖**

   ```bash
   npm install
   # 或
   yarn install
   # 或
   bun install
   ```

3. **配置环境变量**

   在 `web` 目录下创建 `.env.local` 文件，添加以下配置（AI 功能为可选）：

   ```env

   # API 基础 URL（可选，默认为 /core/api）
   VITE_API_BASE_URL=/core/api
   ```

4. **启动开发服务器**

   ```bash
   npm run dev
   ```

   应用将在 `http://localhost:3000` 启动。

   > **注意**：前端开发服务器默认代理 `/core/api` 到后端服务 `http://127.0.0.1:59024`。如需修改代理目标，请编辑 `vite.config.ts` 中的 `targetDict.local`。

5. **构建生产版本**

   ```bash
   npm run build
   ```

   构建产物将输出到 `dist/` 目录。

6. **预览生产构建**

   ```bash
   npm run preview
   ```

## 📂 项目结构

```
web/
├── components/          # React 组件
│   ├── auth/           # 认证组件（登录）
│   ├── dashboard/      # 仪表板组件和图表
│   ├── layout/         # 布局组件（侧边栏、导航栏、通知）
│   ├── proxy/          # 代理管理（L4/L7、节点列表）
│   └── system/         # 系统管理（用户、角色、证书、日志、在线用户）
├── contexts/            # React Contexts（主题、语言）
├── services/            # API 服务和模拟数据
│   ├── authService.ts  # 认证服务
│   ├── http.ts         # HTTP 客户端配置
│   ├── mockData.ts     # 模拟数据
│   └── translations.ts # 国际化翻译
├── types/               # TypeScript 类型定义
│   ├── api.ts          # API 类型
│   ├── auth.ts         # 认证类型
│   ├── dashboard.ts    # 仪表板类型
│   ├── proxy.ts        # 代理类型
│   └── system.ts       # 系统管理类型
├── public/              # 静态资源
│   └── quebec.png      # 项目 Logo
├── index.html           # HTML 入口文件
├── index.tsx            # React 入口文件
├── App.tsx              # 主应用组件
├── vite.config.ts       # Vite 配置
├── tsconfig.json        # TypeScript 配置
└── package.json         # 项目依赖配置
```

## 🔧 开发配置

### Vite 代理配置

开发环境下的 API 代理配置位于 `vite.config.ts`：

```typescript
const targetDict = {
  local: "http://127.0.0.1:59024", // 本地后端服务地址
  // 可以添加其他环境配置
};
```

### 端口配置

- **前端开发服务器**：`3000`（可在 `vite.config.ts` 中修改）
- **后端 API 服务**：`59024`（默认，需与后端配置一致）

## 📝 开发说明

### 添加新功能

1. 在 `components/` 目录下创建对应的组件
2. 在 `types/` 目录下定义相关的 TypeScript 类型
3. 在 `services/` 目录下添加 API 服务调用
4. 在路由配置中添加新页面（如需要）

### 国际化

- 翻译文件位于 `services/translations.ts`
- 使用 `LanguageContext` 进行语言切换
- 支持中文和英文

### 主题切换

- 使用 `ThemeContext` 进行主题管理
- 支持浅色和深色模式
- 基于 Tailwind CSS 的深色模式实现

## 📄 许可证

本项目采用 MIT 许可证。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

<div align="center">
  <p>Made with ❤️ by the Quebec Team</p>
</div>
