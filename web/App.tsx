
import React, { useState, useEffect } from 'react';
import { HashRouter } from 'react-router-dom';
import Sidebar from './components/layout/Sidebar';
import Dashboard from './components/dashboard/Dashboard';
import ProxyManager from './components/proxy/ProxyManager';
import CertManager from './components/proxy/CertManager';
import NodeList from './components/proxy/NodeList';

import Login from './components/auth/Login';
import NotificationDropdown from './components/layout/NotificationDropdown';
import GlobalToast from './components/common/GlobalToast';
import UserManager from './components/system/UserManager';
import RoleManager from './components/system/RoleManager';
import MenuManager from './components/system/MenuManager';
import AuditLogs from './components/system/AuditLogs';
import OnlineUsers from './components/system/OnlineUsers';
import ErrorBoundary from './components/common/ErrorBoundary';
import { ViewState } from './types';
import { Languages, Sun, Moon, Users, Settings } from 'lucide-react';
import { useLanguage } from './contexts/LanguageContext';
import { useTheme } from './contexts/ThemeContext';

import { loginService } from './services/system/loginService';

// Inner App component to use the context
const EnvoyNexusApp: React.FC = () => {
  const [isAuthenticated, setIsAuthenticated] = useState(() => {
    return !!localStorage.getItem('x-quebec-token');
  });
  const [currentView, setCurrentView] = useState<ViewState>(() => {
    // 从 localStorage 恢复上次访问的页面
    const savedView = localStorage.getItem('quebec-current-view');
    if (savedView && Object.values(ViewState).includes(savedView as ViewState)) {
      return savedView as ViewState;
    }
    return ViewState.DASHBOARD;
  });
  const { t, language, setLanguage } = useLanguage();
  const { theme, toggleTheme } = useTheme();
  
  // 获取用户名
  const userInfo = useState(() => {
    return JSON.parse(localStorage.getItem('quebec-user-info') || '{}');
  });

  const toggleLanguage = () => {
    setLanguage(language === 'zh' ? 'en' : 'zh');
  };

  // 保存当前页面到 localStorage
  const handleViewChange = (view: ViewState) => {
    setCurrentView(view);
    localStorage.setItem('quebec-current-view', view);
  };

  // 获取用户名首字母（最多2个字符）
  // 优化：中文字符只取1个，英文字符取2个
  const getUserInitials = (name: string) => {
    if (!name) return 'AD';
    // 检查是否包含中文字符
    const hasChinese = /[\u4e00-\u9fa5]/.test(name);
    if (hasChinese) {
      // 中文名只取1个字符
      return name.substring(0, 1).toUpperCase();
    }
    // 英文名取2个字符
    return name.substring(0, 2).toUpperCase();
  };

  const handleLogout = async () => {
    try {
      const response = await loginService.logout();
      if (response.code === 50000) {
        localStorage.removeItem('x-quebec-token');
        localStorage.removeItem('quebec-username');
        localStorage.removeItem('quebec-role-name');
        localStorage.removeItem('quebec-nickname');
        localStorage.removeItem('quebec-current-view');
        setIsAuthenticated(false);
      } else {
        console.error('Logout failed:', response.message);
      }
    } catch (error) {
      console.error('Logout error:', error);
    }
  };

  const renderContent = () => {
    switch (currentView) {
      case ViewState.DASHBOARD:
        return <Dashboard />;
      case ViewState.PROXY_L4:
        return <ProxyManager layer="L4" />;
      case ViewState.PROXY_L7:
        return <ProxyManager layer="L7" />;
      case ViewState.NODES:
        return <NodeList />;
      case ViewState.CERTS:
        return <CertManager />;
      case ViewState.SYSTEM_USERS:
        return <UserManager />;
      case ViewState.SYSTEM_ROLES:
        return <RoleManager />;
      case ViewState.SYSTEM_MENUS:
        return <MenuManager />;
      case ViewState.SYSTEM_ONLINE_USERS:
        return <OnlineUsers />;
      case ViewState.SYSTEM_LOGS:
        return <AuditLogs />;
      default:
        return <Dashboard />;
    }
  };

  const getViewTitle = (view: ViewState) => {
      // Map viewstate to sidebar translation keys manually for cleaner titles
      const mapping: Record<string, string> = {
          [ViewState.DASHBOARD]: 'sidebar.dashboard',
          [ViewState.PROXY_L4]: 'sidebar.proxy_l4',
          [ViewState.PROXY_L7]: 'sidebar.proxy_l7',
          [ViewState.NODES]: 'sidebar.nodes',
          [ViewState.CERTS]: 'sidebar.certs',
          [ViewState.SYSTEM_USERS]: 'sidebar.system_users',
          [ViewState.SYSTEM_ONLINE_USERS]: 'sidebar.system_online_users',
          [ViewState.SYSTEM_ROLES]: 'sidebar.system_roles',
          [ViewState.SYSTEM_MENUS]: 'sidebar.system_menus',
          [ViewState.SYSTEM_LOGS]: 'sidebar.system_logs',
      }
      return t(mapping[view] as any) || view;
  }

  if (!isAuthenticated) {
    return (
      <>
        <div className="absolute top-6 right-6 z-50 flex gap-3">
            <button
                onClick={toggleTheme}
                className="flex items-center justify-center w-10 h-10 bg-white/80 dark:bg-slate-800/50 hover:bg-white dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300 rounded-lg border border-slate-200 dark:border-slate-700 transition-colors shadow-sm"
                title={theme === 'dark' ? 'Switch to Light Mode' : '切换到亮色模式'}
            >
                {theme === 'dark' ? <Sun size={20} /> : <Moon size={20} />}
            </button>
            <button 
                onClick={toggleLanguage}
                className="flex items-center gap-2 bg-white/80 dark:bg-slate-800/50 hover:bg-white dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300 px-3 py-2 rounded-lg border border-slate-200 dark:border-slate-700 transition-colors text-sm shadow-sm"
            >
                <Languages size={16} />
                {language === 'zh' ? 'English' : '中文'}
            </button>
        </div>
        <Login onLogin={() => {
          setIsAuthenticated(true);
          localStorage.setItem('quebec-user-info', JSON.stringify(userInfo));
        }} />
      </>
    );
  }

  return (
    <div className="flex h-screen bg-slate-50 dark:bg-slate-950 overflow-hidden font-sans text-slate-900 dark:text-slate-200 transition-colors duration-300">
      <Sidebar currentView={currentView} onChangeView={handleViewChange} onLogout={handleLogout} />
      
      <main className="flex-1 flex flex-col min-w-0">
        {/* Top Header */}
        <header className="h-16 border-b border-slate-200 dark:border-slate-800 bg-white/80 dark:bg-slate-950/80 backdrop-blur-md flex items-center justify-between px-6 sticky top-0 z-10 transition-colors duration-300">
            <div className="flex items-center gap-2 text-slate-500 dark:text-slate-400 text-sm">
              <span className="opacity-50">{t('header.platform')}</span>
              <span>/</span>
              <span className="text-slate-800 dark:text-white font-medium capitalize">
                {getViewTitle(currentView)}
              </span>
            </div>

            <div className="flex items-center gap-4">
              <button 
                  onClick={toggleTheme}
                  className="p-2 text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white transition-colors"
                  title={theme === 'dark' ? 'Switch to Light Mode' : '切换到亮色模式'}
              >
                  {theme === 'dark' ? <Sun size={20} /> : <Moon size={20} />}
              </button>
              <button 
                  onClick={toggleLanguage}
                  className="p-2 text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white transition-colors"
                  title={language === 'zh' ? 'Switch to English' : '切换到中文'}
              >
                  <Languages size={20} />
              </button>
              <div className="h-6 w-[1px] bg-slate-200 dark:bg-slate-800"></div>
              
              {/* Notification Dropdown */}
              <NotificationDropdown />

              <div className="h-8 w-[1px] bg-slate-200 dark:bg-slate-800"></div>
              <button className="flex items-center gap-3 pl-2">
                  <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-blue-500 to-purple-500 flex items-center justify-center text-xs font-bold text-white shadow-md">
                    {getUserInitials(localStorage.getItem('quebec-role-name') || '')}
                  </div>
                  <div className="hidden md:block text-left">
                    <p className="text-sm font-medium text-slate-800 dark:text-white leading-none">{localStorage.getItem('quebec-nickname') || ''}</p>
                    <p className="text-xs text-slate-500 dark:text-slate-500 mt-1">{localStorage.getItem('quebec-username') || ''}</p>
                  </div>
              </button>
            </div>
        </header>

        {/* Scrollable Content Area */}
        <div className="flex-1 overflow-y-auto scroll-smooth">
          {renderContent()}
        </div>
      </main>

      {/* AI Assistant Overlay */}
    </div>
  );
};

const App: React.FC = () => {
  return (
    <ErrorBoundary>
      <HashRouter>
        <EnvoyNexusApp />
        <GlobalToast />
      </HashRouter>
    </ErrorBoundary>
  );
};

export default App;
