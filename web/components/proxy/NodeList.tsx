
import React from 'react';
import { Server, ShieldCheck, ShieldAlert, Cpu, MapPin } from 'lucide-react';
import { mockNodes } from '../../services/base/mockData';
import { useLanguage } from '../../contexts/LanguageContext';

const NodeList: React.FC = () => {
  const { t } = useLanguage();

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold text-slate-900 dark:text-white">{t('nodes.title')}</h2>
        <button className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors shadow-md shadow-blue-900/10 dark:shadow-blue-900/20">
          {t('nodes.deployNew')}
        </button>
      </div>

      <div className="grid grid-cols-1 gap-4">
        {mockNodes.map((node) => (
          <div key={node.id} className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl p-4 flex flex-col md:flex-row items-start md:items-center justify-between hover:bg-slate-50 dark:hover:bg-slate-750 transition-colors shadow-sm">
            <div className="flex items-center gap-4 mb-4 md:mb-0">
              <div className={`p-3 rounded-lg ${node.status === 'HEALTHY' ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-500' : 'bg-amber-500/10 text-amber-600 dark:text-amber-500'}`}>
                <Server size={24} />
              </div>
              <div>
                <h4 className="text-lg font-semibold text-slate-900 dark:text-white flex items-center gap-2">
                  {node.id}
                  <span className="text-xs font-normal text-slate-500 dark:text-slate-400 px-2 py-0.5 bg-slate-100 dark:bg-slate-700 rounded-full">{node.ip}</span>
                </h4>
                <p className="text-slate-500 dark:text-slate-400 text-sm flex items-center gap-4 mt-1">
                  <span className="flex items-center gap-1"><MapPin size={12} /> {node.region}</span>
                  <span className="flex items-center gap-1"><Cpu size={12} /> {node.version}</span>
                  <span>{t('nodes.uptime')}: {node.uptime}</span>
                </p>
              </div>
            </div>

            <div className="flex items-center gap-6 w-full md:w-auto justify-between md:justify-end">
               <div className="text-right">
                  <p className="text-slate-500 dark:text-slate-400 text-xs">{t('nodes.connections')}</p>
                  <p className="text-slate-900 dark:text-white font-mono font-medium">{node.connections.toLocaleString()}</p>
               </div>
               
               <div className={`px-3 py-1 rounded-full text-xs font-bold flex items-center gap-1 border ${
                 node.status === 'HEALTHY' 
                 ? 'bg-emerald-50 dark:bg-emerald-950 border-emerald-200 dark:border-emerald-900 text-emerald-700 dark:text-emerald-400' 
                 : 'bg-amber-50 dark:bg-amber-950 border-amber-200 dark:border-amber-900 text-amber-700 dark:text-amber-400'
               }`}>
                 {node.status === 'HEALTHY' ? <ShieldCheck size={14} /> : <ShieldAlert size={14} />}
                 {t(`nodes.status.${node.status.toLowerCase()}`)}
               </div>
               
               <button className="text-slate-400 hover:text-slate-600 dark:hover:text-white p-2 hover:bg-slate-100 dark:hover:bg-slate-700 rounded-lg transition-colors">
                 {t('nodes.manage')}
               </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default NodeList;
