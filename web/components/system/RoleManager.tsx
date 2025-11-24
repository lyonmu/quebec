
import React, { useState } from 'react';
import { Plus, ShieldCheck, Lock, Eye, Trash2, Play, Ban } from 'lucide-react';
import { useLanguage } from '../../contexts/LanguageContext';

interface Role {
  id: number;
  name: string;
  users: number;
  permissions: string;
  status: 'ACTIVE' | 'INACTIVE';
}

const RoleManager: React.FC = () => {
  const { t } = useLanguage();
  
  const [roles, setRoles] = useState<Role[]>([
    { id: 1, name: 'ADMIN', users: 2, permissions: 'all', status: 'ACTIVE' },
    { id: 2, name: 'EDITOR', users: 5, permissions: 'read_write', status: 'ACTIVE' },
    { id: 3, name: 'VIEWER', users: 12, permissions: 'read_only', status: 'ACTIVE' },
  ]);

  const handleToggleStatus = (roleId: number) => {
    setRoles(roles.map(role => {
      if (role.id === roleId) {
        return { ...role, status: role.status === 'ACTIVE' ? 'INACTIVE' : 'ACTIVE' };
      }
      return role;
    }));
  };

  const handleDeleteRole = (roleId: number) => {
    if (window.confirm(t('roles.confirmDelete'))) {
       setRoles(roles.filter(role => role.id !== roleId));
    }
  };

  return (
    <div className="p-6 h-full flex flex-col">
      <div className="flex justify-between items-center mb-6">
        <div>
           <h2 className="text-2xl font-bold text-slate-900 dark:text-white">{t('roles.title')}</h2>
           <p className="text-slate-500 dark:text-slate-400 text-sm mt-1">{t('roles.subtitle')}</p>
        </div>
        <button 
          className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-medium flex items-center gap-2 transition-colors shadow-md shadow-blue-900/10 dark:shadow-blue-900/20"
        >
          <Plus size={18} /> {t('roles.addRole')}
        </button>
      </div>

      <div className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl overflow-hidden shadow-sm">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-slate-50 dark:bg-slate-900/50 border-b border-slate-200 dark:border-slate-700 text-slate-500 dark:text-slate-400 text-xs uppercase tracking-wider">
                <th className="p-4 font-medium">{t('roles.table.name')}</th>
                <th className="p-4 font-medium">{t('roles.table.users')}</th>
                <th className="p-4 font-medium">{t('roles.table.permissions')}</th>
                <th className="p-4 font-medium">{t('users.table.status')}</th>
                <th className="p-4 font-medium text-right">{t('roles.table.actions')}</th>
              </tr>
            </thead>
            <tbody className="text-sm divide-y divide-slate-200 dark:divide-slate-700">
              {roles.map(role => (
                <tr key={role.id} className="hover:bg-slate-50 dark:hover:bg-slate-700/30 transition-colors">
                  <td className="p-4">
                    <div className="flex items-center gap-3">
                      <div className={`p-2 rounded-lg ${
                        role.name === 'ADMIN' ? 'bg-purple-100 text-purple-600 dark:bg-purple-900/20 dark:text-purple-300' :
                        role.name === 'EDITOR' ? 'bg-blue-100 text-blue-600 dark:bg-blue-900/20 dark:text-blue-300' :
                        'bg-slate-100 text-slate-600 dark:bg-slate-800 dark:text-slate-400'
                      }`}>
                        <ShieldCheck size={18} />
                      </div>
                      <span className="font-medium text-slate-900 dark:text-white">{t(`users.roles.${role.name.toLowerCase()}`)}</span>
                    </div>
                  </td>
                  <td className="p-4 text-slate-600 dark:text-slate-300">
                    {role.users}
                  </td>
                  <td className="p-4">
                    <span className={`inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium border ${
                       role.permissions === 'all' ? 'bg-rose-50 dark:bg-rose-900/20 text-rose-700 dark:text-rose-300 border-rose-100 dark:border-rose-800' :
                       role.permissions === 'read_write' ? 'bg-amber-50 dark:bg-amber-900/20 text-amber-700 dark:text-amber-300 border-amber-100 dark:border-amber-800' :
                       'bg-emerald-50 dark:bg-emerald-900/20 text-emerald-700 dark:text-emerald-300 border-emerald-100 dark:border-emerald-800'
                    }`}>
                      {role.permissions === 'all' && <Lock size={12} />}
                      {role.permissions === 'read_only' && <Eye size={12} />}
                      {t(`roles.permissions.${role.permissions}`)}
                    </span>
                  </td>
                  <td className="p-4">
                     <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${
                        role.status === 'ACTIVE' ? 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400' : 'bg-slate-100 dark:bg-slate-700 text-slate-500 dark:text-slate-400'
                        }`}>
                        <span className={`w-1.5 h-1.5 rounded-full mr-1.5 ${
                            role.status === 'ACTIVE' ? 'bg-emerald-500 dark:bg-emerald-400' : 'bg-slate-400'
                        }`}></span>
                        {t(`users.status.${role.status.toLowerCase()}`)}
                      </span>
                  </td>
                  <td className="p-4 text-right">
                    <div className="flex justify-end items-center gap-2">
                      {role.status === 'ACTIVE' ? (
                        <button 
                          onClick={() => handleToggleStatus(role.id)}
                          className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-rose-600 bg-rose-50 hover:bg-rose-100 dark:bg-rose-900/20 dark:hover:bg-rose-900/30 border border-rose-100 dark:border-rose-900/50 rounded transition-colors"
                          title={t('users.actions.disable')}
                        >
                          <Ban size={12} />
                          {t('users.actions.disable')}
                        </button>
                      ) : (
                        <button 
                          onClick={() => handleToggleStatus(role.id)}
                          className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-emerald-600 bg-emerald-50 hover:bg-emerald-100 dark:bg-emerald-900/20 dark:hover:bg-emerald-900/30 border border-emerald-100 dark:border-emerald-900/50 rounded transition-colors"
                          title={t('users.actions.enable')}
                        >
                          <Play size={12} />
                          {t('users.actions.enable')}
                        </button>
                      )}

                      <button 
                        onClick={() => handleDeleteRole(role.id)}
                        className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-slate-500 hover:text-rose-600 bg-slate-100 hover:bg-rose-50 dark:bg-slate-700 dark:hover:bg-rose-900/20 dark:text-slate-400 dark:hover:text-rose-400 border border-transparent hover:border-rose-100 dark:hover:border-rose-900/50 rounded transition-all"
                        title={t('users.actions.delete')}
                      >
                        <Trash2 size={12} />
                      </button>
                    </div>
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

export default RoleManager;
