
import React from 'react';
import { Lock, RefreshCw, AlertCircle, CheckCircle, Upload } from 'lucide-react';
import { mockCerts } from '../../services/base/mockData';
import { useLanguage } from '../../contexts/LanguageContext';

const CertManager: React.FC = () => {
  const { t } = useLanguage();

  return (
    <div className="p-6 max-w-[1600px] mx-auto space-y-6">
       <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
           <h2 className="text-2xl font-bold text-slate-900 dark:text-white">{t('certs.title')}</h2>
           <p className="text-slate-500 dark:text-slate-400 text-sm mt-1">{t('certs.subtitle')}</p>
        </div>
        <button className="bg-white dark:bg-slate-700 hover:bg-slate-50 dark:hover:bg-slate-600 border border-slate-200 dark:border-slate-600 text-slate-700 dark:text-white px-4 py-2 rounded-lg text-sm font-medium flex items-center gap-2 transition-colors shadow-sm">
          <Upload size={16} /> {t('certs.import')}
        </button>
      </div>

      <div className="grid grid-cols-1 gap-4">
        {mockCerts.map(cert => (
          <div key={cert.id} className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 p-4 rounded-xl flex flex-col md:flex-row items-start md:items-center justify-between shadow-sm">
            <div className="flex items-start gap-4">
              <div className={`mt-1 p-2 rounded-lg ${
                cert.status === 'VALID' ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-500' : 
                cert.status === 'EXPIRING_SOON' ? 'bg-amber-500/10 text-amber-600 dark:text-amber-500' : 
                'bg-rose-500/10 text-rose-600 dark:text-rose-500'
              }`}>
                <Lock size={20} />
              </div>
              <div>
                <h4 className="text-lg font-bold text-slate-900 dark:text-white">{cert.domain}</h4>
                <div className="flex items-center gap-4 text-sm text-slate-500 dark:text-slate-400 mt-1">
                   <span>{t('certs.issuer')}: {cert.issuer}</span>
                   <span className="text-slate-300 dark:text-slate-600">|</span>
                   <span>{t('certs.expires')}: {cert.expiryDate}</span>
                </div>
              </div>
            </div>

            <div className="flex items-center gap-4 mt-4 md:mt-0 w-full md:w-auto justify-between md:justify-end">
              {cert.isAutoRenew && (
                <span className="flex items-center gap-1 text-xs text-blue-600 dark:text-blue-400 bg-blue-50 dark:bg-blue-400/10 border border-blue-100 dark:border-transparent px-2 py-1 rounded">
                  <RefreshCw size={12} /> {t('certs.autoRenew')}
                </span>
              )}
              
              <div className="flex flex-col items-end gap-1">
                 <span className={`flex items-center gap-1 text-sm font-medium ${
                   cert.status === 'VALID' ? 'text-emerald-600 dark:text-emerald-400' : 
                   cert.status === 'EXPIRING_SOON' ? 'text-amber-600 dark:text-amber-400' : 'text-rose-600 dark:text-rose-400'
                 }`}>
                   {cert.status === 'VALID' && <CheckCircle size={14} />}
                   {cert.status === 'EXPIRING_SOON' && <AlertCircle size={14} />}
                   {cert.status === 'EXPIRED' && <AlertCircle size={14} />}
                   {t(`certs.status.${cert.status.toLowerCase().replace('_', '_')}`)}
                 </span>
                 {cert.status !== 'VALID' && (
                   <button className="text-xs text-blue-600 dark:text-blue-400 hover:text-blue-500 dark:hover:text-blue-300 hover:underline">
                     {t('certs.renewNow')}
                   </button>
                 )}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default CertManager;
