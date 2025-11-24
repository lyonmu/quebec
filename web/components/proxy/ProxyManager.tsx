
import React, { useState } from 'react';
import { Search, Plus, Filter, ArrowRight, Lock, Globe, FileJson, Database, Activity } from 'lucide-react';
import { mockRoutes } from '../../services/mockData';
import { Protocol } from '../../types';
import { useLanguage } from '../../contexts/LanguageContext';

interface ProxyManagerProps {
    layer: 'L4' | 'L7';
}

const ProxyManager: React.FC<ProxyManagerProps> = ({ layer }) => {
  const { t } = useLanguage();
  const [searchTerm, setSearchTerm] = useState('');

  const filteredRoutes = mockRoutes.filter(route => {
    const isL4 = layer === 'L4' && route.protocol === Protocol.TCP;
    const isL7 = layer === 'L7' && (route.protocol === Protocol.HTTP || route.protocol === Protocol.HTTPS || route.protocol === Protocol.GRPC);
    
    if (!isL4 && !isL7) return false;

    return route.name.toLowerCase().includes(searchTerm.toLowerCase()) || 
           route.prefix.toLowerCase().includes(searchTerm.toLowerCase());
  });

  return (
    <div className="p-6 h-full flex flex-col">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4">
        <div>
           <h2 className="text-2xl font-bold text-slate-900 dark:text-white">{layer === 'L4' ? t('proxy.titleL4') : t('proxy.titleL7')}</h2>
           <p className="text-slate-500 dark:text-slate-400 text-sm mt-1">{layer === 'L4' ? t('proxy.subtitleL4') : t('proxy.subtitleL7')}</p>
        </div>
        <div className="flex gap-3 w-full md:w-auto">
          <div className="relative flex-1 md:w-64">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" size={18} />
            <input 
              type="text" 
              placeholder={t('proxy.searchPlaceholder')} 
              className="w-full bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 text-slate-900 dark:text-white pl-10 pr-4 py-2 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm shadow-sm"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
          </div>
          <button className="bg-white dark:bg-slate-800 hover:bg-slate-50 dark:hover:bg-slate-700 border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 px-3 py-2 rounded-lg transition-colors shadow-sm">
            <Filter size={18} />
          </button>
          <button className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-medium flex items-center gap-2 transition-colors shadow-md shadow-blue-900/10 dark:shadow-blue-900/20">
            <Plus size={18} /> {t('proxy.addRoute')}
          </button>
        </div>
      </div>

      <div className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl overflow-hidden flex-1 shadow-sm">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-slate-50 dark:bg-slate-900/50 border-b border-slate-200 dark:border-slate-700 text-slate-500 dark:text-slate-400 text-xs uppercase tracking-wider">
                <th className="p-4 font-medium">{t('proxy.table.name')}</th>
                <th className="p-4 font-medium">{t('proxy.table.prefix')}</th>
                <th className="p-4 font-medium">{t('proxy.table.cluster')}</th>
                <th className="p-4 font-medium">{t('proxy.table.protocol')}</th>
                <th className="p-4 font-medium">{t('proxy.table.timeout')}</th>
                <th className="p-4 font-medium">{t('proxy.table.status')}</th>
                <th className="p-4 font-medium text-right">{t('proxy.table.actions')}</th>
              </tr>
            </thead>
            <tbody className="text-sm divide-y divide-slate-200 dark:divide-slate-700">
              {filteredRoutes.length > 0 ? (
                filteredRoutes.map(route => (
                    <tr key={route.id} className="hover:bg-slate-50 dark:hover:bg-slate-700/30 transition-colors group">
                    <td className="p-4">
                        <div className="font-medium text-slate-900 dark:text-white">{route.name}</div>
                        <div className="text-xs text-slate-500 font-mono">{route.id}</div>
                    </td>
                    <td className="p-4">
                        <code className="bg-slate-100 dark:bg-slate-900 px-2 py-1 rounded text-blue-600 dark:text-blue-400 font-mono text-xs border border-slate-200 dark:border-slate-700">
                        {route.prefix}
                        </code>
                    </td>
                    <td className="p-4">
                        <div className="flex items-center gap-2 text-slate-600 dark:text-slate-300">
                        <ArrowRight size={14} className="text-slate-400 dark:text-slate-500" />
                        {route.targetCluster}
                        </div>
                    </td>
                    <td className="p-4">
                        <span className="flex items-center gap-1.5 text-slate-600 dark:text-slate-300">
                        {route.protocol === Protocol.HTTPS || route.protocol === Protocol.GRPC ? 
                            <Lock size={14} className="text-emerald-500 dark:text-emerald-400" /> : 
                            route.protocol === Protocol.TCP ?
                            <Database size={14} className="text-indigo-500 dark:text-indigo-400" /> :
                            <Globe size={14} className="text-slate-400" />
                        }
                        {route.protocol}
                        </span>
                    </td>
                    <td className="p-4 text-slate-500 dark:text-slate-400">{route.timeout}</td>
                    <td className="p-4">
                        <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${
                        route.status === 'ACTIVE' ? 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400' : 'bg-slate-100 dark:bg-slate-700 text-slate-500 dark:text-slate-400'
                        }`}>
                        <span className={`w-1.5 h-1.5 rounded-full mr-1.5 ${
                            route.status === 'ACTIVE' ? 'bg-emerald-500 dark:bg-emerald-400' : 'bg-slate-400'
                        }`}></span>
                        {t(`proxy.status.${route.status.toLowerCase()}`)}
                        </span>
                    </td>
                    <td className="p-4 text-right">
                        <div className="flex justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                        <button title="View YAML" className="p-1.5 hover:bg-slate-100 dark:hover:bg-slate-600 rounded text-slate-400 hover:text-slate-700 dark:hover:text-white">
                            <FileJson size={16} />
                        </button>
                        <button className="text-blue-600 dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300 text-xs font-medium hover:underline px-2 py-1">
                            Edit
                        </button>
                        </div>
                    </td>
                    </tr>
                ))
              ) : (
                <tr>
                    <td colSpan={7} className="p-8 text-center text-slate-500 dark:text-slate-400">
                        <Activity size={32} className="mx-auto mb-2 opacity-20" />
                        <p>No {layer} routes found.</p>
                    </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
};

export default ProxyManager;