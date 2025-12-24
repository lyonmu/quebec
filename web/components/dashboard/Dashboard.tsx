
import React, { useMemo } from 'react';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, LineChart, Line } from 'recharts';
import { Activity, Clock, AlertTriangle, Globe } from 'lucide-react';
import { generateMetrics, mockNodes } from '../../services/base/mockData';
import { useLanguage } from '../../contexts/LanguageContext';
import { useTheme } from '../../contexts/ThemeContext';

const Dashboard: React.FC = () => {
  const { t } = useLanguage();
  const { theme } = useTheme();
  const data = useMemo(() => generateMetrics(), []);

  const totalConnections = mockNodes.reduce((acc, node) => acc + node.connections, 0);
  const healthyNodes = mockNodes.filter(n => n.status === 'HEALTHY').length;

  // Chart styles based on theme
  const chartColors = {
     grid: theme === 'dark' ? '#334155' : '#e2e8f0',
     axis: theme === 'dark' ? '#94a3b8' : '#64748b',
     tooltipBg: theme === 'dark' ? '#1e293b' : '#ffffff',
     tooltipBorder: theme === 'dark' ? '#334155' : '#e2e8f0',
     tooltipText: theme === 'dark' ? '#f1f5f9' : '#0f172a',
  };

  return (
    <div className="p-6 max-w-[1600px] mx-auto space-y-6 animate-fade-in">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {/* Stat Cards */}
        <div className="bg-white dark:bg-slate-800 p-4 rounded-xl border border-slate-200 dark:border-slate-700 shadow-sm transition-colors">
          <div className="flex justify-between items-start">
            <div>
              <p className="text-slate-500 dark:text-slate-400 text-sm">{t('dashboard.totalReq')}</p>
              <h3 className="text-2xl font-bold text-slate-800 dark:text-white mt-1">4,235</h3>
              <span className="text-emerald-500 dark:text-emerald-400 text-xs font-medium flex items-center mt-1">
                <Activity size={12} className="mr-1" /> +12.5%
              </span>
            </div>
            <div className="p-2 bg-blue-50 dark:bg-blue-500/10 rounded-lg">
              <Globe className="text-blue-600 dark:text-blue-500" size={24} />
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-slate-800 p-4 rounded-xl border border-slate-200 dark:border-slate-700 shadow-sm transition-colors">
          <div className="flex justify-between items-start">
            <div>
              <p className="text-slate-500 dark:text-slate-400 text-sm">{t('dashboard.avgLatency')}</p>
              <h3 className="text-2xl font-bold text-slate-800 dark:text-white mt-1">42ms</h3>
              <span className="text-emerald-500 dark:text-emerald-400 text-xs font-medium flex items-center mt-1">
                -3ms {t('dashboard.vsLastHour')}
              </span>
            </div>
            <div className="p-2 bg-purple-50 dark:bg-purple-500/10 rounded-lg">
              <Clock className="text-purple-600 dark:text-purple-500" size={24} />
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-slate-800 p-4 rounded-xl border border-slate-200 dark:border-slate-700 shadow-sm transition-colors">
          <div className="flex justify-between items-start">
            <div>
              <p className="text-slate-500 dark:text-slate-400 text-sm">{t('dashboard.errorRate')}</p>
              <h3 className="text-2xl font-bold text-slate-800 dark:text-white mt-1">0.04%</h3>
              <span className="text-rose-500 dark:text-rose-400 text-xs font-medium flex items-center mt-1">
                +0.01% {t('dashboard.spikeDetected')}
              </span>
            </div>
            <div className="p-2 bg-rose-50 dark:bg-rose-500/10 rounded-lg">
              <AlertTriangle className="text-rose-600 dark:text-rose-500" size={24} />
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-slate-800 p-4 rounded-xl border border-slate-200 dark:border-slate-700 shadow-sm transition-colors">
          <div className="flex justify-between items-start">
            <div>
              <p className="text-slate-500 dark:text-slate-400 text-sm">{t('dashboard.activeInstances')}</p>
              <h3 className="text-2xl font-bold text-slate-800 dark:text-white mt-1">{healthyNodes} / {mockNodes.length}</h3>
              <p className="text-slate-400 dark:text-slate-500 text-xs mt-1">{totalConnections.toLocaleString()} {t('dashboard.activeConns')}</p>
            </div>
            <div className="p-2 bg-emerald-50 dark:bg-emerald-500/10 rounded-lg">
              <div className="w-6 h-6 rounded-full border-2 border-emerald-500 flex items-center justify-center">
                 <div className="w-2 h-2 bg-emerald-500 rounded-full animate-pulse"></div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Charts Row 1 */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white dark:bg-slate-800 p-6 rounded-xl border border-slate-200 dark:border-slate-700 transition-colors">
          <h3 className="text-lg font-semibold text-slate-800 dark:text-white mb-4">{t('dashboard.trafficVolume')}</h3>
          <div className="h-64 w-full">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={data.rps}>
                <defs>
                  <linearGradient id="colorRps" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#3b82f6" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke={chartColors.grid} vertical={false} />
                <XAxis dataKey="time" stroke={chartColors.axis} fontSize={12} tickLine={false} axisLine={false} />
                <YAxis stroke={chartColors.axis} fontSize={12} tickLine={false} axisLine={false} />
                <Tooltip 
                  contentStyle={{ backgroundColor: chartColors.tooltipBg, borderColor: chartColors.tooltipBorder, color: chartColors.tooltipText }}
                  itemStyle={{ color: '#3b82f6' }}
                />
                <Area type="monotone" dataKey="value" stroke="#3b82f6" strokeWidth={2} fillOpacity={1} fill="url(#colorRps)" />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="bg-white dark:bg-slate-800 p-6 rounded-xl border border-slate-200 dark:border-slate-700 transition-colors">
          <h3 className="text-lg font-semibold text-slate-800 dark:text-white mb-4">{t('dashboard.latency')}</h3>
          <div className="h-64 w-full">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={data.latency}>
                <CartesianGrid strokeDasharray="3 3" stroke={chartColors.grid} vertical={false} />
                <XAxis dataKey="time" stroke={chartColors.axis} fontSize={12} tickLine={false} axisLine={false} />
                <YAxis stroke={chartColors.axis} fontSize={12} tickLine={false} axisLine={false} />
                <Tooltip 
                   contentStyle={{ backgroundColor: chartColors.tooltipBg, borderColor: chartColors.tooltipBorder, color: chartColors.tooltipText }}
                   itemStyle={{ color: '#a855f7' }}
                />
                <Line type="monotone" dataKey="value" stroke="#a855f7" strokeWidth={2} dot={false} />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;
