import React, { useState, useEffect } from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { Search, Monitor, Smartphone, Globe } from 'lucide-react';

interface OnlineUser {
  id: string;
  username: string;
  ip: string;
  lastOperationTime: string;
  operationType: string;
  os: string;
  platform: string;
  browserName: string;
  browserVersion: string;
  engineName: string;
  engineVersion: string;
}

const OnlineUsers: React.FC = () => {
  const { t } = useLanguage();
  const [users, setUsers] = useState<OnlineUser[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchTerm, setSearchTerm] = useState('');

  // Mock data generation
  useEffect(() => {
    // Simulate API call
    const timer = setTimeout(() => {
      const mockUsers: OnlineUser[] = [
        {
          id: '1',
          username: 'admin',
          ip: '192.168.1.100',
          lastOperationTime: '2023-11-24 13:30:45',
          operationType: 'QUERY',
          os: 'Windows 10',
          platform: 'Desktop',
          browserName: 'Chrome',
          browserVersion: '119.0.0.0',
          engineName: 'Blink',
          engineVersion: '119.0.0.0'
        },
        {
          id: '2',
          username: 'editor',
          ip: '192.168.1.105',
          lastOperationTime: '2023-11-24 13:28:12',
          operationType: 'UPDATE',
          os: 'macOS 14.1',
          platform: 'Desktop',
          browserName: 'Safari',
          browserVersion: '17.1',
          engineName: 'WebKit',
          engineVersion: '605.1.15'
        },
        {
          id: '3',
          username: 'viewer',
          ip: '10.0.0.55',
          lastOperationTime: '2023-11-24 13:15:30',
          operationType: 'LOGIN',
          os: 'iOS 17.1',
          platform: 'Mobile',
          browserName: 'Mobile Safari',
          browserVersion: '17.1',
          engineName: 'WebKit',
          engineVersion: '605.1.15'
        },
        {
          id: '4',
          username: 'test_user',
          ip: '172.16.0.20',
          lastOperationTime: '2023-11-24 13:05:10',
          operationType: 'QUERY',
          os: 'Linux x86_64',
          platform: 'Desktop',
          browserName: 'Firefox',
          browserVersion: '120.0',
          engineName: 'Gecko',
          engineVersion: '120.0'
        }
      ];
      setUsers(mockUsers);
      setLoading(false);
    }, 800);

    return () => clearTimeout(timer);
  }, []);

  const filteredUsers = users.filter(user => 
    user.username.toLowerCase().includes(searchTerm.toLowerCase()) ||
    user.ip.includes(searchTerm)
  );

  const getPlatformIcon = (platform: string) => {
    if (platform === 'Mobile') return <Smartphone size={16} className="text-slate-400" />;
    return <Monitor size={16} className="text-slate-400" />;
  };

  return (
    <div className="p-6 max-w-[1600px] mx-auto space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 dark:text-white">{t('online_users.title')}</h1>
          <p className="text-slate-500 dark:text-slate-400 mt-1">{t('online_users.subtitle')}</p>
        </div>
        
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={20} />
          <input 
            type="text" 
            placeholder={t('common.search')}
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="pl-10 pr-4 py-2 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 w-full md:w-64 transition-all"
          />
        </div>
      </div>

      <div className="bg-white dark:bg-slate-900 rounded-xl border border-slate-200 dark:border-slate-800 shadow-sm overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-slate-50 dark:bg-slate-800/50 border-b border-slate-200 dark:border-slate-800">
                <th className="px-6 py-4 text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider whitespace-nowrap">{t('online_users.table.username')}</th>
                <th className="px-6 py-4 text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider whitespace-nowrap">{t('online_users.table.ip')}</th>
                <th className="px-6 py-4 text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider whitespace-nowrap">{t('online_users.table.lastTime')}</th>
                <th className="px-6 py-4 text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider whitespace-nowrap">{t('online_users.table.type')}</th>
                <th className="px-6 py-4 text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider whitespace-nowrap">{t('online_users.table.os')}</th>
                <th className="px-6 py-4 text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider whitespace-nowrap">{t('online_users.table.platform')}</th>
                <th className="px-6 py-4 text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider whitespace-nowrap">{t('online_users.table.browser')}</th>
                <th className="px-6 py-4 text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider whitespace-nowrap">{t('online_users.table.engine')}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-200 dark:divide-slate-800">
              {loading ? (
                <tr>
                  <td colSpan={8} className="px-6 py-12 text-center text-slate-500 dark:text-slate-400">
                    {t('common.loading')}
                  </td>
                </tr>
              ) : filteredUsers.length === 0 ? (
                <tr>
                  <td colSpan={8} className="px-6 py-12 text-center text-slate-500 dark:text-slate-400">
                    No online users found
                  </td>
                </tr>
              ) : (
                filteredUsers.map((user) => (
                  <tr key={user.id} className="hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 rounded-full bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center text-blue-600 dark:text-blue-400 font-medium text-xs">
                          {user.username.substring(0, 2).toUpperCase()}
                        </div>
                        <span className="font-medium text-slate-900 dark:text-white">{user.username}</span>
                      </div>
                    </td>
                    <td className="px-6 py-4 text-sm text-slate-600 dark:text-slate-300 font-mono">
                      {user.ip}
                    </td>
                    <td className="px-6 py-4 text-sm text-slate-600 dark:text-slate-300 whitespace-nowrap">
                      {user.lastOperationTime}
                    </td>
                    <td className="px-6 py-4">
                      <span className="px-2.5 py-1 rounded-full text-xs font-medium bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 border border-slate-200 dark:border-slate-700">
                        {user.operationType}
                      </span>
                    </td>
                    <td className="px-6 py-4 text-sm text-slate-600 dark:text-slate-300">
                      {user.os}
                    </td>
                    <td className="px-6 py-4 text-sm text-slate-600 dark:text-slate-300">
                      <div className="flex items-center gap-2">
                        {getPlatformIcon(user.platform)}
                        <span>{user.platform}</span>
                      </div>
                    </td>
                    <td className="px-6 py-4 text-sm text-slate-600 dark:text-slate-300">
                      <div className="flex flex-col">
                        <span>{user.browserName}</span>
                        <span className="text-xs text-slate-400">{user.browserVersion}</span>
                      </div>
                    </td>
                    <td className="px-6 py-4 text-sm text-slate-600 dark:text-slate-300">
                      <div className="flex flex-col">
                        <span>{user.engineName}</span>
                        <span className="text-xs text-slate-400">{user.engineVersion}</span>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
        
        <div className="px-6 py-4 border-t border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-800/50 flex items-center justify-between text-sm text-slate-500 dark:text-slate-400">
          <div>
            {t('common.total')} {filteredUsers.length} {t('common.items')}
          </div>
          <div className="flex gap-2">
            <button className="px-3 py-1 border border-slate-200 dark:border-slate-700 rounded hover:bg-white dark:hover:bg-slate-800 disabled:opacity-50" disabled>{t('common.prev')}</button>
            <button className="px-3 py-1 border border-slate-200 dark:border-slate-700 rounded hover:bg-white dark:hover:bg-slate-800 disabled:opacity-50" disabled>{t('common.next')}</button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default OnlineUsers;
