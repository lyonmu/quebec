
import React, { useEffect, useState, useCallback } from 'react';
import { Plus, ShieldCheck, Lock, Eye, Trash2, Play, Ban } from 'lucide-react';
import { useLanguage } from '../../contexts/LanguageContext';
import { roleService } from '../../services/system/roleService';
import { SystemRoleDetail, YesOrNo } from '../../types';

type RoleStatus = 'ACTIVE' | 'INACTIVE';

interface RoleView {
  id: string;
  name: string;
  users: number;
  status: RoleStatus;
}

const RoleManager: React.FC = () => {
  const { t } = useLanguage();
  const [roles, setRoles] = useState<RoleView[]>([]);
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(10);
  const [total, setTotal] = useState(0);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [editingRole, setEditingRole] = useState<SystemRoleDetail | null>(null);
  const [form, setForm] = useState({ name: '', remark: '' });
  const [actionLoading, setActionLoading] = useState(false);

  const mapYesOrNoToStatus = (status: YesOrNo): RoleStatus =>
    status === 1 ? 'ACTIVE' : 'INACTIVE';

  const mapStatusToYesOrNo = (status: RoleStatus): YesOrNo =>
    status === 'ACTIVE' ? 1 : 2;

  const mapBackendToView = (r: SystemRoleDetail): RoleView => ({
    id: r.id,
    name: r.name,
    users: 0,
    status: mapYesOrNoToStatus(r.status),
  });

  const fetchRoles = useCallback(async () => {
    setLoading(true);
    try {
      const res = await roleService.fetchRolePage({ page, page_size: pageSize });
      if (res.code === 50000 && res.data) {
        setRoles((res.data.items || []).map(mapBackendToView));
        setTotal(res.data.total || 0);
      }
    } catch (e) {
      console.error('Failed to fetch roles', e);
    } finally {
      setLoading(false);
    }
  }, [page, pageSize]);

  useEffect(() => {
    fetchRoles();
  }, [fetchRoles]);
  const handleToggleStatus = async (role: RoleView) => {
    if (actionLoading) return;
    setActionLoading(true);
    try {
      const targetStatus = role.status === 'ACTIVE' ? 2 : 1;
      const res = await roleService.toggleRoleStatus(role.id, { status: targetStatus as YesOrNo });
      if (res.code === 50000) {
        await fetchRoles();
      }
    } catch (e) {
      console.error('Failed to toggle role', e);
    } finally {
      setActionLoading(false);
    }
  };

  const handleDeleteRole = async (roleId: string) => {
    if (!window.confirm(t('roles.confirmDelete'))) return;
    if (actionLoading) return;
    setActionLoading(true);
    try {
      const res = await roleService.deleteRole(roleId);
      if (res.code === 50000) {
        await fetchRoles();
      }
    } catch (e) {
      console.error('Failed to delete role', e);
    } finally {
      setActionLoading(false);
    }
  };

  const openCreateModal = () => {
    setEditingRole(null);
    setForm({ name: '', remark: '' });
    setIsAddModalOpen(true);
  };

  const openEditModal = async (roleId: string) => {
    try {
      const res = await roleService.getRoleDetail(roleId);
      if (res.code === 50000 && res.data) {
        setEditingRole(res.data);
        setForm({ name: res.data.name, remark: res.data.remark || '' });
        setIsAddModalOpen(true);
      }
    } catch (e) {
      console.error('Failed to load role detail', e);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (actionLoading) return;
    setActionLoading(true);
    try {
      if (editingRole) {
        const res = await roleService.updateRole(editingRole.id, { name: form.name, remark: form.remark });
        if (res.code === 50000) {
          setIsAddModalOpen(false);
          await fetchRoles();
        }
      } else {
        const res = await roleService.createRole({ name: form.name, remark: form.remark });
        if (res.code === 50000) {
          setIsAddModalOpen(false);
          await fetchRoles();
        }
      }
    } catch (e) {
      console.error('Failed to submit role', e);
    } finally {
      setActionLoading(false);
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
          onClick={openCreateModal}
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
              {loading ? (
                <tr>
                  <td colSpan={5} className="p-6 text-center text-slate-500 dark:text-slate-400">
                    {t('common.loading')}
                  </td>
                </tr>
              ) : roles.length === 0 ? (
                <tr>
                  <td colSpan={5} className="p-6 text-center text-slate-500 dark:text-slate-400">
                    {t('common.empty') || 'No roles'}
                  </td>
                </tr>
              ) : (
              roles.map(role => (
                <tr key={role.id} className="hover:bg-slate-50 dark:hover:bg-slate-700/30 transition-colors">
                  <td className="p-4">
                    <div className="flex items-center gap-3">
                      <div className="p-2 rounded-lg bg-slate-100 text-slate-600 dark:bg-slate-800 dark:text-slate-400">
                        <ShieldCheck size={18} />
                      </div>
                      {/* 直接展示角色名称，避免依赖固定翻译 key */}
                      <span className="font-medium text-slate-900 dark:text-white">{role.name}</span>
                    </div>
                  </td>
                  <td className="p-4 text-slate-600 dark:text-slate-300">
                    {role.users}
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
                          onClick={() => handleToggleStatus(role)}
                          className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-rose-600 bg-rose-50 hover:bg-rose-100 dark:bg-rose-900/20 dark:hover:bg-rose-900/30 border border-rose-100 dark:border-rose-900/50 rounded transition-colors"
                          disabled={actionLoading}
                          title={t('users.actions.disable')}
                        >
                          <Ban size={12} />
                          {t('users.actions.disable')}
                        </button>
                      ) : (
                        <button 
                          onClick={() => handleToggleStatus(role)}
                          className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-emerald-600 bg-emerald-50 hover:bg-emerald-100 dark:bg-emerald-900/20 dark:hover:bg-emerald-900/30 border border-emerald-100 dark:border-emerald-900/50 rounded transition-colors"
                          disabled={actionLoading}
                          title={t('users.actions.enable')}
                        >
                          <Play size={12} />
                          {t('users.actions.enable')}
                        </button>
                      )}

                      <div className="flex items-center gap-2">
                        <button 
                          onClick={() => openEditModal(role.id)}
                          className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-slate-500 hover:text-blue-600 bg-slate-100 hover:bg-blue-50 dark:bg-slate-700 dark:hover:bg-blue-900/20 dark:text-slate-400 dark:hover:text-blue-400 border border-transparent hover:border-blue-100 dark:hover:border-blue-900/50 rounded transition-all"
                          title={t('common.edit')}
                        >
                          <Eye size={12} />
                        </button>
                        <button 
                          onClick={() => handleDeleteRole(role.id)}
                          className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-slate-500 hover:text-rose-600 bg-slate-100 hover:bg-rose-50 dark:bg-slate-700 dark:hover:bg-rose-900/20 dark:text-slate-400 dark:hover:text-rose-400 border border-transparent hover:border-rose-100 dark:hover:border-rose-900/50 rounded transition-all"
                          title={t('users.actions.delete')}
                        >
                          <Trash2 size={12} />
                        </button>
                      </div>
                    </div>
                  </td>
                </tr>
              )))}
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
            disabled={roles.length === 0 || roles.length < pageSize}
            className="px-3 py-1 border border-slate-200 dark:border-slate-700 rounded hover:bg-white dark:hover:bg-slate-800 disabled:opacity-50"
          >
            {t('common.next')}
          </button>
        </div>
      </div>

      {/* Add/Edit Role Modal */}
      {isAddModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 backdrop-blur-sm p-4">
          <div className="bg-white dark:bg-slate-800 rounded-xl shadow-xl w-full max-w-md border border-slate-200 dark:border-slate-700 overflow-hidden animate-fade-in-up">
            <div className="px-6 py-4 border-b border-slate-200 dark:border-slate-700 flex justify-between items-center bg-slate-50 dark:bg-slate-900/50">
              <h3 className="font-semibold text-slate-900 dark:text-white">
                {editingRole ? t('roles.modal.editTitle') || 'Edit Role' : t('roles.modal.addTitle') || 'Add Role'}
              </h3>
            </div>
            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700 dark:text-slate-300">{t('roles.modal.name') || 'Role Name'}</label>
                <div className="relative">
                  <ShieldCheck size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                  <input 
                    type="text" 
                    className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                    value={form.name}
                    onChange={e => setForm({...form, name: e.target.value})}
                    required
                  />
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700 dark:text-slate-300">{t('roles.modal.remark') || 'Remark'}</label>
                <textarea
                  className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                  rows={3}
                  value={form.remark}
                  onChange={e => setForm({...form, remark: e.target.value})}
                />
              </div>

              <div className="pt-4 flex justify-end gap-3">
                <button 
                  type="button"
                  onClick={() => setIsAddModalOpen(false)}
                  className="px-4 py-2 text-sm font-medium text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 rounded-lg transition-colors"
                >
                  {t('roles.modal.cancel') || t('users.modal.cancel')}
                </button>
                <button 
                  type="submit"
                  className="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors shadow-sm"
                  disabled={actionLoading}
                >
                  {t('roles.modal.save') || t('users.modal.save')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default RoleManager;
