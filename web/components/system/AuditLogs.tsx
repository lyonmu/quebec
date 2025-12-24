
import React, { useState, useEffect, useCallback } from 'react';
import { useLanguage } from '../../contexts/LanguageContext';
import { operationLogService, OperationLogItem } from '../../services/system/operationLogService';
import { userService } from '../../services/system/userService';
import { OperationTypeDescriptions } from '../../types/api';
import { Options } from '../../types';
import { Search, Clock, Monitor, Globe, RotateCcw } from 'lucide-react';

const AuditLogs: React.FC = () => {
  const { t } = useLanguage();
  const [logs, setLogs] = useState<OperationLogItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [searchKeyword, setSearchKeyword] = useState('');
  const [operationType, setOperationType] = useState<number>(0);
  const [selectedUserId, setSelectedUserId] = useState('');
  const [startTime, setStartTime] = useState('');
  const [endTime, setEndTime] = useState('');
  const [userLabels, setUserLabels] = useState<Options[]>([]);

  const fetchLogs = useCallback(async () => {
    setLoading(true);
    try {
      // 转换日期为Unix时间戳（秒）
      const startTimeTimestamp = startTime ? Math.floor(new Date(startTime).getTime() / 1000) : undefined;
      const endTimeTimestamp = endTime ? Math.floor(new Date(endTime).getTime() / 1000) : undefined;

      const res = await operationLogService.fetchLogPage({
        page,
        page_size: pageSize,
        operation_type: operationType || undefined,
        user_id: selectedUserId || undefined,
        start_time: startTimeTimestamp,
        end_time: endTimeTimestamp,
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
  }, [page, pageSize, operationType, searchKeyword, selectedUserId, startTime, endTime]);

  const fetchUserLabels = async () => {
    try {
      const res = await userService.fetchUserLabels();
      if (res.code === 50000 && res.data) {
        setUserLabels(res.data);
      }
    } catch (e) {
      console.error('Failed to fetch user labels', e);
    }
  };

  useEffect(() => {
    fetchLogs();
  }, [fetchLogs]);

  useEffect(() => {
    fetchUserLabels();
  }, []);

  // 操作类型映射
  const getOperationTypeText = (type: number): string => {
    return OperationTypeDescriptions[type] || t('operation.types.unknown') || '未知操作';
  };

  // 格式化时间
  const formatTime = (timestamp: number): string => {
    const date = new Date(timestamp * 1000);
    return date.toLocaleString();
  };

  // 过滤日志 - 如果搜索关键词不是IP地址，则前端过滤用户名和昵称
  const filteredLogs = logs.filter(log => {
    if (!searchKeyword) return true;

    // 如果是IP地址，后端已过滤，前端不需要再过滤
    if (/^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$/.test(searchKeyword)) {
      return true;
    }

    // 否则前端过滤用户名和昵称
    const keyword = searchKeyword.toLowerCase();
    return (
      log.username?.toLowerCase().includes(keyword) ||
      log.nickname?.toLowerCase().includes(keyword)
    );
  });

  // 重置过滤条件
  const resetFilters = () => {
    setSearchKeyword('');
    setOperationType(0);
    setSelectedUserId('');
    setStartTime('');
    setEndTime('');
    setPage(1);
  };

  return (
    <div className="p-6 max-w-[1600px] mx-auto space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
           <h2 className="text-2xl font-bold text-slate-900 dark:text-white">{t('logs.title')}</h2>
           <p className="text-slate-500 dark:text-slate-400 text-sm mt-1">{t('logs.subtitle')}</p>
        </div>
        <div className="flex flex-wrap gap-3 items-center justify-end">
          <select
            value={selectedUserId}
            onChange={(e) => setSelectedUserId(e.target.value)}
            className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 min-w-[150px]"
          >
            <option value="">{t('common.select') || '请选择用户'}</option>
            {userLabels.map((label) => (
              <option key={label.value} value={label.value}>{label.label}</option>
            ))}
          </select>
          <select
            value={operationType}
            onChange={(e) => setOperationType(Number(e.target.value))}
            className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="0">{t('operation.types.all') || '全部操作'}</option>
            <option value="1">{OperationTypeDescriptions[1]}</option>
            <option value="2">{OperationTypeDescriptions[2]}</option>
            <option value="3">{OperationTypeDescriptions[3]}</option>
            <option value="4">{OperationTypeDescriptions[4]}</option>
            <option value="5">{OperationTypeDescriptions[5]}</option>
            <option value="6">{OperationTypeDescriptions[6]}</option>
            <option value="7">{OperationTypeDescriptions[7]}</option>
            <option value="8">{OperationTypeDescriptions[8]}</option>
            <option value="9">{OperationTypeDescriptions[9]}</option>
            <option value="10">{OperationTypeDescriptions[10]}</option>
            <option value="11">{OperationTypeDescriptions[11]}</option>
            <option value="12">{OperationTypeDescriptions[12]}</option>
            <option value="16">{OperationTypeDescriptions[16]}</option>
            <option value="17">{OperationTypeDescriptions[17]}</option>
            <option value="18">{OperationTypeDescriptions[18]}</option>
            <option value="19">{OperationTypeDescriptions[19]}</option>
            <option value="20">{OperationTypeDescriptions[20]}</option>
          </select>
          <div className="flex items-center gap-2">
            <label className="text-sm text-slate-600 dark:text-slate-300 whitespace-nowrap">{t('common.startTime')}:</label>
            <input
              type="date"
              value={startTime}
              onChange={(e) => setStartTime(e.target.value)}
              className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-900 dark:text-white px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div className="flex items-center gap-2">
            <label className="text-sm text-slate-600 dark:text-slate-300 whitespace-nowrap">{t('common.endTime')}:</label>
            <input
              type="date"
              value={endTime}
              onChange={(e) => setEndTime(e.target.value)}
              className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-900 dark:text-white px-3 py-2 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <button 
              onClick={resetFilters}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors text-sm font-medium flex items-center gap-2"
            >
              <RotateCcw size={16} />
              {t('logs.reset')}
            </button>
        </div>
      </div>

      <div className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl overflow-hidden flex-1 shadow-sm">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-slate-50 dark:bg-slate-900/50 border-b border-slate-200 dark:border-slate-700 text-slate-500 dark:text-slate-400 text-xs uppercase tracking-wider">
                <th className="p-4 font-medium">{t('logs.table.user')}</th>
                <th className="p-4 font-medium">{t('logs.table.action')}</th>
                <th className="p-4 font-medium">{t('logs.table.ip')}</th>
                <th className="p-4 font-medium">{t('logs.table.operationType') || '操作类型'}</th>
                <th className="p-4 font-medium">{t('logs.table.time')}</th>
                <th className="p-4 font-medium">{t('logs.table.os') || '操作系统'}</th>
                <th className="p-4 font-medium">{t('logs.table.platform') || '操作平台'}</th>
                <th className="p-4 font-medium">{t('logs.table.browser') || '浏览器名称'}</th>
                <th className="p-4 font-medium">{t('logs.table.engine') || '引擎名称'}</th>
              </tr>
            </thead>
            <tbody className="text-sm divide-y divide-slate-200 dark:divide-slate-700">
              {loading ? (
                <tr>
                  <td colSpan={9} className="p-6 text-center text-slate-500 dark:text-slate-400">
                    {t('common.loading')}
                  </td>
                </tr>
              ) : filteredLogs.length === 0 ? (
                <tr>
                  <td colSpan={9} className="p-6 text-center text-slate-500 dark:text-slate-400">
                    {t('common.empty') || 'No logs'}
                  </td>
                </tr>
              ) : (
                filteredLogs.map(log => (
                  <tr key={log.id} className="hover:bg-slate-50 dark:hover:bg-slate-700/30 transition-colors">
                    
                    <td className="p-4 text-slate-600 dark:text-slate-300">
                      <div className="font-medium">{log.nickname || log.username || '-'}</div>
                      {log.nickname && log.username && (
                        <div className="text-xs text-slate-500">@{log.username}</div>
                      )}
                    </td>
                    <td className="p-4">
                      <span className="px-2.5 py-1 rounded-full text-xs font-medium bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 border border-slate-200 dark:border-slate-700 shadow-sm">
                        {getOperationTypeText(log.operation_type)}
                      </span>
                    </td>
                    <td className="p-4 text-slate-600 dark:text-slate-300">
                      <div className="flex items-center gap-1">
                        <Globe size={14} className="text-slate-400" />
                        <span className="font-mono text-xs">{log.access_ip || '-'}</span>
                      </div>
                    </td>
                    <td className="p-4 text-slate-500 dark:text-slate-400 text-xs">
                      <div className="flex items-center gap-1">
                        <Clock size={12} />
                        {formatTime(log.operation_time)}
                      </div>
                    </td>
                    <td className="p-4 text-slate-500 dark:text-slate-400 text-xs">
                      {log.os || '-'}
                    </td>
                    <td className="p-4 text-slate-500 dark:text-slate-400 text-xs">
                      {log.platform || '-'}
                    </td>
                    <td className="p-4 text-slate-600 dark:text-slate-300">
                      <div className="flex items-center gap-1">
                        <Monitor size={14} className="text-slate-400" />
                        <span className="text-xs">
                          {log.browser_name ? `${log.browser_name} ${log.browser_version || ''}` : '-'}
                        </span>
                      </div>
                    </td>
                    <td className="p-4 text-slate-500 dark:text-slate-400 text-xs">
                      <div className="flex items-center gap-1">
                        <span>{log.browser_engine_name || '-'}</span>
                        {log.browser_engine_version && (
                          <span className="text-slate-400"> {log.browser_engine_version}</span>
                        )}
                      </div>
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
