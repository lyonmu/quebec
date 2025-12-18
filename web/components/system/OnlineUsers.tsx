import React, { useState, useEffect, useCallback } from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { Search, Monitor, Smartphone, Trash2, RefreshCw, Eye } from 'lucide-react';
import { onlineuserService } from '../../services/system/onlineuserService';
import { OnlineUser, OnlineUserLabel } from '../../types';
import DateTimePicker from '../common/DateTimePicker';
import { OPERATION_TYPE_MAP, DEFAULT_OPERATION_TYPE } from '../../services/base/translations';
import ConfirmDialog from '../common/ConfirmDialog';

const OnlineUsers: React.FC = () => {
  const { t, language } = useLanguage();
  const [users, setUsers] = useState<OnlineUser[]>([]);
  const [loading, setLoading] = useState(true);
  const [ipSearch, setIpSearch] = useState('');
  const [selectedUserId, setSelectedUserId] = useState('');
  const [startTime, setStartTime] = useState('');
  const [endTime, setEndTime] = useState('');
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [labels, setLabels] = useState<OnlineUserLabel[]>([]);
  const [confirmOpen, setConfirmOpen] = useState(false);
  const [confirmUser, setConfirmUser] = useState<OnlineUser | null>(null);
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [detailOpen, setDetailOpen] = useState(false);
  const [detailUser, setDetailUser] = useState<OnlineUser | null>(null);

  const fetchUsers = useCallback(async () => {
    setLoading(true);
    try {
      const response = await onlineuserService.fetchOnlineUsers({
        page,
        page_size: pageSize,
        access_ip: ipSearch || undefined,
        user_id: selectedUserId || undefined,
        start_time: startTime ? Math.floor(new Date(startTime).getTime() / 1000) : undefined,
        end_time: endTime ? Math.floor(new Date(endTime).getTime() / 1000) : undefined,
      });
      if (response.code === 50000 && response.data) {
        setUsers(response.data.items || []);
        setTotal(response.data.total);
      }
    } catch (error) {
      console.error('Failed to fetch online users:', error);
    } finally {
      setLoading(false);
    }
  }, [page, pageSize, ipSearch, selectedUserId, startTime, endTime]);

  const fetchLabels = async () => {
    try {
      const response = await onlineuserService.fetchOnlineUserLabels();
      if (response.code === 50000 && response.data) {
        setLabels(response.data);
      }
    } catch (error) {
      console.error('Failed to fetch labels:', error);
    }
  };

  useEffect(() => {
    fetchLabels();
  }, []);

  useEffect(() => {
    fetchUsers();
  }, [page, pageSize]);

  const handleSearch = () => {
    if (page === 1) {
      fetchUsers();
    } else {
      setPage(1);
    }
  };

  const handleRequestClearUser = (user: OnlineUser) => {
    setConfirmUser(user);
    setConfirmOpen(true);
  };

  const handleConfirmClearUser = async () => {
    if (!confirmUser) return;
    setConfirmLoading(true);
    try {
      const response = await onlineuserService.clearOnlineUser(confirmUser.id);
      if (response.code === 50000) {
        fetchUsers(); // Refresh list
      }
    } catch (error) {
      console.error('Failed to clear user:', error);
    } finally {
      setConfirmLoading(false);
      setConfirmOpen(false);
      setConfirmUser(null);
    }
  };

  const getPlatformIcon = (platform: string) => {
    if (platform === 'Mobile') return <Smartphone size={16} className="text-slate-400" />;
    return <Monitor size={16} className="text-slate-400" />;
  };

  const formatTime = (timestamp: number) => {
    return new Date(timestamp * 1000).toLocaleString();
  };

  const getOperationType = (operationType: number) => {
    const translationKey = OPERATION_TYPE_MAP[operationType] || DEFAULT_OPERATION_TYPE;
    return t(translationKey);
  };

  const totalPages = Math.ceil(total / pageSize);

  return (
    <div className="p-6 max-w-[1600px] mx-auto space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-slate-900 dark:text-white">{t('online_users.title')}</h1>
          <p className="text-slate-500 dark:text-slate-400 mt-1">{t('online_users.subtitle')}</p>
        </div>
        
        <div className="flex flex-col xl:flex-row gap-4 items-end xl:items-center">
          <div className="flex flex-wrap items-center gap-3">
            {/* User ID Select */}
            <div className="relative">
              <select
                value={selectedUserId}
                onChange={(e) => setSelectedUserId(e.target.value)}
                className="appearance-none pl-4 pr-10 py-2 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm min-w-[150px]"
              >
                <option value="">All Users</option>
                {labels.map((label) => (
                  label.children ? (
                    <optgroup key={label.value} label={label.label}>
                      {label.children.map(child => (
                        <option key={child.value} value={child.value}>{child.label}</option>
                      ))}
                    </optgroup>
                  ) : (
                    <option key={label.value} value={label.value}>{label.label}</option>
                  )
                ))}
              </select>
              <div className="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400">
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7"></path></svg>
              </div>
            </div>

            {/* IP Input */}
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={16} />
              <input 
                type="text" 
                placeholder="Access IP"
                value={ipSearch}
                onChange={(e) => setIpSearch(e.target.value)}
                className="pl-9 pr-4 py-2 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 w-40 text-sm"
              />
            </div>

            {/* Date Range */}
            <div className="flex items-center gap-2">
              <DateTimePicker 
                value={startTime}
                onChange={setStartTime}
                language={language}
                placeholder={t('common.startTime') || 'Start Time'}
              />
              <span className="text-slate-400">-</span>
              <DateTimePicker 
                value={endTime}
                onChange={setEndTime}
                language={language}
                placeholder={t('common.endTime') || 'End Time'}
              />
            </div>

            {/* Actions */}
            <button 
              onClick={handleSearch}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors text-sm font-medium flex items-center gap-2"
            >
              <Search size={16} />
              {t('common.search')}
            </button>
            
            {/* <button 
              onClick={fetchUsers}
              className="p-2 text-slate-500 hover:text-blue-600 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded-lg transition-colors"
              title={t('common.refresh')}
            >
              <RefreshCw size={20} />
            </button> */}
          </div>
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
                <th className="px-6 py-4 text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider whitespace-nowrap text-right">{t('common.actions')}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-200 dark:divide-slate-800">
              {loading ? (
                <tr>
                  <td colSpan={9} className="px-6 py-12 text-center text-slate-500 dark:text-slate-400">
                    {t('common.loading')}
                  </td>
                </tr>
              ) : users.length === 0 ? (
                <tr>
                  <td colSpan={9} className="px-6 py-12 text-center text-slate-500 dark:text-slate-400">
                    No online users found
                  </td>
                </tr>
              ) : (
                users.map((user) => (
                  <tr key={user.id} className="hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors">
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-3">
                        <div className="w-8 h-8 rounded-full bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center text-blue-600 dark:text-blue-400 font-medium text-xs">
                          {user.username ? user.username.substring(0, 2).toUpperCase() : 'UN'}
                        </div>
                        <span className="font-medium text-slate-900 dark:text-white">{user.nickname || 'Unknown'}</span>
                      </div>
                    </td>
                    <td className="px-6 py-4 text-sm text-slate-600 dark:text-slate-300 font-mono">
                      {user.access_ip}
                    </td>
                    <td className="px-6 py-4 text-sm text-slate-600 dark:text-slate-300 whitespace-nowrap">
                      {formatTime(user.last_operation_time)}
                    </td>
                    <td className="px-6 py-4">
                      <span className="px-2.5 py-1 rounded-full text-xs font-medium bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 border border-slate-200 dark:border-slate-700">
                        {getOperationType(user.operation_type)}
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
                        <span>{user.browser_name}</span>
                        <span className="text-xs text-slate-400">{user.browser_version}</span>
                      </div>
                    </td>
                    <td className="px-6 py-4 text-sm text-slate-600 dark:text-slate-300">
                      <div className="flex flex-col">
                        <span>{user.browser_engine_name}</span>
                        <span className="text-xs text-slate-400">{user.browser_engine_version}</span>
                      </div>
                    </td>
                    <td className="px-6 py-4 text-sm text-right">
                      <div className="flex justify-end gap-2">
                        <button
                          onClick={() => {
                            setDetailUser(user);
                            setDetailOpen(true);
                          }}
                          className="text-slate-500 hover:text-blue-600 hover:bg-blue-50 dark:text-slate-400 dark:hover:text-blue-400 dark:hover:bg-blue-900/20 p-2 rounded-lg transition-colors"
                          title={t('common.view') || 'View'}
                        >
                          <Eye size={16} />
                        </button>
                        <button 
                          onClick={() => handleRequestClearUser(user)}
                          className="text-red-500 hover:text-red-700 hover:bg-red-50 dark:hover:bg-red-900/20 p-2 rounded-lg transition-colors"
                          title={t('common.delete')}
                        >
                          <Trash2 size={16} />
                        </button>
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
            {t('common.total')} {total} {t('common.items')}
          </div>
          <div className="flex gap-2">
            <button 
              onClick={() => setPage(p => Math.max(1, p - 1))}
              disabled={page === 1}
              className="px-3 py-1 border border-slate-200 dark:border-slate-700 rounded hover:bg-white dark:hover:bg-slate-800 disabled:opacity-50"
            >
              {t('common.prev')}
            </button>
            <span className="px-3 py-1">
              {page} / {Math.max(1, totalPages)}
            </span>
            <button 
              onClick={() => setPage(p => Math.min(totalPages, p + 1))}
              disabled={page >= totalPages}
              className="px-3 py-1 border border-slate-200 dark:border-slate-700 rounded hover:bg-white dark:hover:bg-slate-800 disabled:opacity-50"
            >
              {t('common.next')}
            </button>
          </div>
        </div>
      </div>

      <ConfirmDialog
        open={confirmOpen}
        title={t('online_users.title')}
        description={
          confirmUser
            ? `${confirmUser.nickname || confirmUser.username || confirmUser.id} - ${t('online_users.confirm_clear')}`
            : t('online_users.confirm_clear')
        }
        confirmText={t('common.delete')}
        cancelText={t('users.modal.cancel')}
        loading={confirmLoading}
        onCancel={() => {
          if (confirmLoading) return;
          setConfirmOpen(false);
          setConfirmUser(null);
        }}
        onConfirm={handleConfirmClearUser}
      />

      {/* Detail Modal */}
      {detailOpen && detailUser && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4">
          <div className="bg-white dark:bg-slate-900 rounded-2xl shadow-2xl w-full max-w-xl border border-slate-200 dark:border-slate-800 overflow-hidden">
            <div className="px-6 py-4 border-b border-slate-200 dark:border-slate-800 flex items-center justify-between bg-slate-50 dark:bg-slate-900/80">
              <div>
                <h3 className="text-sm font-semibold text-slate-900 dark:text-slate-50">
                  {detailUser.nickname || detailUser.username || 'User'}
                </h3>
                <p className="text-xs text-slate-500 dark:text-slate-400 mt-1">
                  {detailUser.username} · {detailUser.access_ip}
                </p>
              </div>
              <button
                onClick={() => setDetailOpen(false)}
                className="text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 transition-colors text-xs"
              >
                {t('users.modal.cancel')}
              </button>
            </div>
            <div className="px-6 py-4 grid grid-cols-2 gap-3 text-xs text-slate-700 dark:text-slate-200">
              <div>
                <div className="text-slate-400 mb-1">{t('online_users.table.ip')}</div>
                <div className="font-mono">{detailUser.access_ip}</div>
              </div>
              <div>
                <div className="text-slate-400 mb-1">{t('online_users.table.lastTime')}</div>
                <div>{formatTime(detailUser.last_operation_time)}</div>
              </div>
              <div>
                <div className="text-slate-400 mb-1">{t('online_users.table.os')}</div>
                <div>{detailUser.os}</div>
              </div>
              <div>
                <div className="text-slate-400 mb-1">{t('online_users.table.platform')}</div>
                <div>{detailUser.platform}</div>
              </div>
              <div>
                <div className="text-slate-400 mb-1">{t('online_users.table.browser')}</div>
                <div>{detailUser.browser_name} {detailUser.browser_version}</div>
              </div>
              <div>
                <div className="text-slate-400 mb-1">{t('online_users.table.engine')}</div>
                <div>{detailUser.browser_engine_name} {detailUser.browser_engine_version}</div>
              </div>
              <div className="col-span-2">
                <div className="text-slate-400 mb-1">{t('online_users.table.type')}</div>
                <div>{getOperationType(detailUser.operation_type)}</div>
              </div>
            </div>
            <div className="px-6 py-3 border-t border-slate-200 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/80 flex justify-end">
              <button
                onClick={() => setDetailOpen(false)}
                className="px-4 py-1.5 text-xs font-medium rounded-full border border-slate-300 dark:border-slate-600 text-slate-700 dark:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
              >
                {t('users.modal.cancel')}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default OnlineUsers;
