
import React from 'react';
import { mockLogs } from '../../services/base/mockData';
import { useLanguage } from '../../contexts/LanguageContext';
import { Search, Filter, Download, CheckCircle, XCircle, Clock } from 'lucide-react';

const AuditLogs: React.FC = () => {
  const { t } = useLanguage();

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
              className="w-full bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-900 dark:text-white pl-10 pr-4 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm shadow-sm"
            />
          </div>
          <button className="bg-white dark:bg-slate-800 hover:bg-slate-50 dark:hover:bg-slate-700 border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 px-3 py-2 rounded-lg transition-colors shadow-sm">
            <Filter size={18} />
          </button>
          <button className="bg-white dark:bg-slate-800 hover:bg-slate-50 dark:hover:bg-slate-700 border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 px-3 py-2 rounded-lg transition-colors shadow-sm flex items-center gap-2 text-sm font-medium">
            <Download size={18} /> Export
          </button>
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
                <th className="p-4 font-medium">{t('logs.table.status')}</th>
                <th className="p-4 font-medium">{t('logs.table.time')}</th>
                <th className="p-4 font-medium">{t('logs.table.details')}</th>
              </tr>
            </thead>
            <tbody className="text-sm divide-y divide-slate-200 dark:divide-slate-700">
              {mockLogs.map(log => (
                <tr key={log.id} className="hover:bg-slate-50 dark:hover:bg-slate-700/30 transition-colors">
                  <td className="p-4">
                    <div className="font-medium text-slate-900 dark:text-white">{log.action}</div>
                    <div className="text-xs text-slate-500 font-mono">{log.id}</div>
                  </td>
                  <td className="p-4 text-slate-600 dark:text-slate-300">
                    <span className="bg-slate-100 dark:bg-slate-700 px-2 py-0.5 rounded text-xs">
                      {log.user}
                    </span>
                  </td>
                   <td className="p-4 text-slate-600 dark:text-slate-300 font-mono text-xs">
                    {log.target}
                  </td>
                  <td className="p-4 text-slate-600 dark:text-slate-300 font-mono text-xs">
                    {log.ip}
                  </td>
                  <td className="p-4">
                    <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded text-xs font-medium ${
                        log.status === 'SUCCESS' ? 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400' : 'bg-rose-50 dark:bg-rose-500/10 text-rose-600 dark:text-rose-400'
                        }`}>
                        {log.status === 'SUCCESS' ? <CheckCircle size={12} /> : <XCircle size={12} />}
                        {t(`logs.status.${log.status.toLowerCase()}`)}
                      </span>
                  </td>
                  <td className="p-4 text-slate-500 dark:text-slate-400 text-xs flex items-center gap-1">
                    <Clock size={12} />
                    {log.timestamp}
                  </td>
                  <td className="p-4 text-slate-500 dark:text-slate-400 text-xs max-w-xs truncate" title={log.details}>
                    {log.details}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default AuditLogs;
