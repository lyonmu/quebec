
import React, { useState } from 'react';
import { LayoutDashboard, Network, Server, Settings, LogOut, ChevronDown, ChevronRight, Users, ShieldCheck, FileKey, Globe, Layers, ScrollText, Monitor } from 'lucide-react';
import { ViewState } from '../../types';
import { useLanguage } from '../../contexts/LanguageContext';

interface SidebarProps {
  currentView: ViewState;
  onChangeView: (view: ViewState) => void;
  onLogout: () => void;
}

type MenuItem = {
  id?: ViewState;
  label: string;
  icon: any;
  children?: { id: ViewState; label: string; icon: any }[];
};

const Sidebar: React.FC<SidebarProps> = ({ currentView, onChangeView, onLogout }) => {
  const { t, language } = useLanguage();
  const [expandedMenus, setExpandedMenus] = useState<string[]>(['proxy', 'system']);

  const toggleMenu = (menuKey: string) => {
    setExpandedMenus(prev => 
      prev.includes(menuKey) ? prev.filter(k => k !== menuKey) : [...prev, menuKey]
    );
  };

  const menuItems: MenuItem[] = [
    { id: ViewState.DASHBOARD, label: t('sidebar.dashboard'), icon: LayoutDashboard },
    { 
      label: t('sidebar.proxy'), 
      icon: Network,
      children: [
        { id: ViewState.NODES, label: t('sidebar.nodes'), icon: Server },
        { id: ViewState.PROXY_L4, label: t('sidebar.proxy_l4'), icon: Layers },
        { id: ViewState.PROXY_L7, label: t('sidebar.proxy_l7'), icon: Globe },
        { id: ViewState.CERTS, label: t('sidebar.certs'), icon: FileKey }
      ]
    },
    { 
      label: t('sidebar.system'), 
      icon: Settings,
      children: [
        { id: ViewState.SYSTEM_USERS, label: t('sidebar.system_users'), icon: Users },
        { id: ViewState.SYSTEM_ONLINE_USERS, label: t('sidebar.system_online_users'), icon: Monitor },
        { id: ViewState.SYSTEM_ROLES, label: t('sidebar.system_roles'), icon: ShieldCheck },
        { id: ViewState.SYSTEM_LOGS, label: t('sidebar.system_logs'), icon: ScrollText }
      ]
    },
  ];

  const isChildActive = (item: MenuItem) => {
    return item.children?.some(child => child.id === currentView);
  };

  return (
    <div className={`${language === 'en' ? 'w-72' : 'w-64'} h-screen bg-white dark:bg-slate-900 border-r border-slate-200 dark:border-slate-800 flex flex-col flex-shrink-0 transition-all duration-300`}>
      <div className="p-6 flex items-center gap-3">
        <img src="/quebec.png" alt="Quebec Logo" className="w-8 h-8 rounded-lg shadow-lg shadow-blue-900/20 flex-shrink-0" />
        <h1 className="text-xl font-bold text-slate-900 dark:text-white tracking-tight truncate">Quebec</h1>
      </div>

      <nav className="flex-1 px-4 space-y-1 overflow-y-auto overflow-x-hidden">
        {menuItems.map((item, index) => {
          const Icon = item.icon;
          const isActive = item.id === currentView;
          const hasChildren = !!item.children;
          const isExpanded = expandedMenus.includes(item.label) || isChildActive(item);
          const childActive = isChildActive(item);

          if (hasChildren) {
            return (
              <div key={index} className="space-y-1">
                <button
                  onClick={() => toggleMenu(item.label)}
                  className={`w-full flex items-center justify-between px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 whitespace-nowrap ${
                    childActive
                      ? 'text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-900/10'
                      : 'text-slate-500 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 hover:text-slate-900 dark:hover:text-slate-200'
                  }`}
                >
                  <div className="flex items-center gap-3 overflow-hidden">
                    <Icon size={20} className="flex-shrink-0" />
                    <span className="truncate">{item.label}</span>
                  </div>
                  {isExpanded ? <ChevronDown size={16} className="flex-shrink-0" /> : <ChevronRight size={16} className="flex-shrink-0" />}
                </button>
                
                {isExpanded && (
                  <div className="pl-4 space-y-1 animate-fade-in-down">
                    {item.children!.map((child) => {
                      const ChildIcon = child.icon;
                      const isChildSelected = child.id === currentView;
                      return (
                        <button
                          key={child.id}
                          onClick={() => onChangeView(child.id)}
                          className={`w-full flex items-center gap-3 px-4 py-2 text-sm font-medium rounded-lg transition-all duration-200 border-l-2 whitespace-nowrap ${
                            isChildSelected
                              ? 'border-blue-600 text-blue-600 dark:text-blue-400 bg-blue-50/50 dark:bg-blue-900/20'
                              : 'border-transparent text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-slate-200'
                          }`}
                        >
                           {ChildIcon && <ChildIcon size={16} className={`flex-shrink-0 ${isChildSelected ? 'text-blue-600 dark:text-blue-400' : 'text-slate-400'}`} />}
                           <span className="truncate">{child.label}</span>
                        </button>
                      );
                    })}
                  </div>
                )}
              </div>
            );
          }

          return (
            <button
              key={item.id}
              onClick={() => onChangeView(item.id!)}
              className={`w-full flex items-center gap-3 px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 whitespace-nowrap ${
                isActive 
                  ? 'bg-blue-600 text-white shadow-lg shadow-blue-900/20' 
                  : 'text-slate-500 dark:text-slate-400 hover:bg-slate-100 dark:hover:bg-slate-800 hover:text-slate-900 dark:hover:text-slate-200'
              }`}
            >
              <Icon size={20} className={`flex-shrink-0 ${isActive ? 'text-white' : 'text-slate-500 dark:text-slate-500'}`} />
              <span className="truncate">{item.label}</span>
            </button>
          );
        })}
      </nav>

      <div className="p-4 border-t border-slate-200 dark:border-slate-800">
        <div className="bg-slate-50 dark:bg-slate-800 rounded-xl p-4 mb-4 border border-slate-100 dark:border-slate-700">
           <div className="flex items-center gap-3 mb-2">
             <div className="w-2 h-2 rounded-full bg-emerald-500 flex-shrink-0"></div>
             <span className="text-xs font-medium text-slate-600 dark:text-slate-400 whitespace-nowrap truncate">{t('sidebar.controlPlane')}: {t('sidebar.healthy')}</span>
           </div>
           <div className="h-1 w-full bg-slate-200 dark:bg-slate-700 rounded-full overflow-hidden">
             <div className="h-full bg-emerald-500 w-[98%]"></div>
           </div>
        </div>
        
        <button 
          onClick={onLogout}
          className="w-full flex items-center gap-3 px-4 py-2 text-sm font-medium text-slate-500 dark:text-slate-400 hover:text-rose-600 dark:hover:text-rose-400 hover:bg-rose-50 dark:hover:bg-rose-500/10 rounded-lg transition-colors whitespace-nowrap"
        >
          <LogOut size={18} className="flex-shrink-0" />
          <span>{t('sidebar.signOut')}</span>
        </button>
      </div>
    </div>
  );
};

export default Sidebar;
