
import React, { useState, useEffect, useCallback } from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { operationLogService, OperationLogItem } from '../../services/system/operationLogService';
import { Search, Filter, Download, CheckCircle, XCircle, Clock } from 'lucide-react';

const AuditLogs: React.FC = () => {
  const { t } = useLanguage();
  const [logs, setLogs] = useState<OperationLogItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [searchKeyword, setSearchKeyword] = useState('');
  const [operationType, setOperationType] = useState<number>(0);

  const fetchLogs = useCallback(async () => {
    setLoading(true);
    try {
      const res = await operationLogService.fetchLogPage({
        page,
        page_size: pageSize,
        operation_type: operationType || undefined,
      });
      if (res.code === 50000 && res.data) {
        setLogs(res.data.items || []);
        setTotal(res.data.total || 0);
      }
    } catch (e) {
      console.error('Failed to fetch logs', e);
    } finally {
      setLoading(false);
    }
  }, [page, pageSize, operationType]);

  useEffect(() => {
    fetchLogs();
  }, [fetchLogs]);

  // 操作类型映射
  const getOperationTypeText = (type: number): string => {
    const typeMap: Record<number, string> = {
      1: t('operation.types.login') || '登录',
      2: t('operation.types.logout') || '登出',
      3: t('operation.types.user_create') || '创建用户',
      4: t('operation.types.user_update') || '更新用户',
      5: t('operation.types.user_delete') || '删除用户',
      6: t('operation.types.role_create') || '创建角色',
      7: t('operation.types.role_update') || '更新角色',
      8: t('operation.types.role_delete') || '删除角色',
      9: t('operation.types.menu_create') || '创建菜单',
      10: t('operation.types.menu_update') || '更新菜单',
      11: t('operation.types.menu_delete') || '删除菜单',
    };
    return typeMap[type] || t('operation.types.unknown') || '未知操作';
  };

  // 格式化时间
  const formatTime = (timestamp: number): string => {
    const date = new Date(timestamp * 1000);
    return date.toLocaleString();
  };

  return (
    <div className="p-6 h-full flex flex-col">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
        <div>
           <h2 className="text-2xl font-bold text-slate-900 dark:text-white">{t('logs.title')}</h2>
           <p className="text-slate-500 dark:text-slate-400 text-sm mt-1">{t('logs.subtitle')}</p>
        </div>
        <div className="flex gap-3 w-full md:w-auto">
          <div className="relative flex-1 md:w-64">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
            <input
              type="text"
              placeholder={t('common.search')}
              value={searchKeyword}
              onChange={(e) => setSearchKeyword(e.target.value)}
              className="w-full bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-900 dark:text-white pl-10 pr-4 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm shadow-sm"
            />
          </div>
          <select
            value={operationType}
            onChange={(e) => setOperationType(Number(e.target.value))}
            className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="0">{t('operation.types.all') || '全部操作'}</option>
            <option value="1">{t('operation.types.login') || '登录'}</option>
            <option value="2">{t('operation.types.logout') || '登出'}</option>
            <option value="3">{t('operation.types.user_create') || '创建用户'}</option>
            <option value="4">{t('operation.types.user_update') || '更新用户'}</option>
            <option value="5">{t('operation.types.user_delete') || '删除用户'}</option>
            <option value="6">{t('operation.types.role_create') || '创建角色'}</option>
            <option value="7">{t('operation.types.role_update') || '更新角色'}</option>
            <option value="8">{t('operation.types.role_delete') || '删除角色'}</option>
            <option value="9">{t('operation.types.menu_create') || '创建菜单'}</option>
            <option value="10">{t('operation.types.menu_update') || '更新菜单'}</option>
            <option value="11">{t('operation.types.menu_delete') || '删除菜单'}</option>
          </select>
        </div>
      </div>

      <div className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl overflow-hidden flex-1 shadow-sm">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-slate-50 dark:bg-slate-900/50 border-b border-slate-200 dark:border-slate-700 text-slate-500 dark:text-slate-400 text-xs uppercase tracking-wider">
                <th className="p-4 font-medium">{t('logs.table.action')}</th>
                <th className="p-4 font-medium">{t('logs.table.user')}</th>
                <th className="p-4 font-medium">{t('logs.table.target')}</th>
                <th className="p-4 font-medium">{t('logs.table.ip')}</th>
                <th className="p-4 font-medium">{t('logs.table.time')}</th>
                <th className="p-4 font-medium">{t('logs.table.details')}</th>
              </tr>
            </thead>
            <tbody className="text-sm divide-y divide-slate-200 dark:divide-slate-700">
              {loading ? (
                <tr>
                  <td colSpan={6} className="p-6 text-center text-slate-500 dark:text-slate-400">
                    {t('common.loading')}
                  </td>
                </tr>
              ) : logs.length === 0 ? (
                <tr>
                  <td colSpan={6} className="p-6 text-center text-slate-500 dark:text-slate-400">
                    {t('common.empty') || 'No logs'}
                  </td>
                </tr>
              ) : (
                logs
                  .filter(log =>
                    !searchKeyword ||
                    log.username?.toLowerCase().includes(searchKeyword.toLowerCase()) ||
                    log.operation_action?.toLowerCase().includes(searchKeyword.toLowerCase()) ||
                    log.ip?.includes(searchKeyword)
                  )
                  .map(log => (
                    <tr key={log.id} className="hover:bg-slate-50 dark:hover:bg-slate-700/30 transition-colors">
                      <td className="p-4">
                        <div className="font-medium text-slate-900 dark:text-white">
                          {getOperationTypeText(log.operation_type)}
                        </div>
                        <div className="text-xs text-slate-500 font-mono">{log.id}</div>
                      </td>
                      <td className="p-4 text-slate-600 dark:text-slate-300">
                        <span className="bg-slate-100 dark:bg-slate-700 px-2 py-0.5 rounded text-xs">
                          {log.username || '-'}
                        </span>
                      </td>
                      <td className="p-4 text-slate-600 dark:text-slate-300 font-mono text-xs">
                        {log.target_type && log.target_id ? `${log.target_type}:${log.target_id}` : '-'}
                      </td>
                      <td className="p-4 text-slate-600 dark:text-slate-300 font-mono text-xs">
                        {log.ip || '-'}
                      </td>
                      <td className="p-4 text-slate-500 dark:text-slate-400 text-xs flex items-center gap-1">
                        <Clock size={12} />
                        {formatTime(log.created_at)}
                      </td>
                      <td className="p-4 text-slate-500 dark:text-slate-400 text-xs max-w-xs truncate" title={log.details}>
                        {log.details || '-'}
                      </td>
                    </tr>
                  ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Pagination */}
      <div className="mt-4 flex items-center justify-between text-sm text-slate-500 dark:text-slate-400">
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
            {page}
          </span>
          <button
            onClick={() => setPage(p => p + 1)}
            disabled={logs.length === 0 || logs.length < pageSize}
            className="px-3 py-1 border border-slate-200 dark:border-slate-700 rounded hover:bg-white dark:hover:bg-slate-800 disabled:opacity-50"
          >
            {t('common.next')}
          </button>
        </div>
      </div>
    </div>
  );
};

export default AuditLogs;
