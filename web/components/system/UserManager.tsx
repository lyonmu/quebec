
import React, { useEffect, useState, useCallback } from 'react';
import { Plus, User as UserIcon, Lock, Shield, Key, Mail, Play, Ban, Eye } from 'lucide-react';
import { User } from '../../types';
import { useLanguage } from '../../contexts/LanguageContext';
import { userService } from '../../services/system/userService';
import { roleService } from '../../services/system/roleService';
import { Options, YesOrNo, SystemUserDetail } from '../../types';
import { OPERATION_TYPE_MAP, DEFAULT_OPERATION_TYPE } from '../../services/base/translations';
import forge from 'node-forge';

const UserManager: React.FC = () => {
  const { t } = useLanguage();
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState<boolean>(false);
  const [page, setPage] = useState<number>(1);
  const [pageSize] = useState<number>(10);
  const [total, setTotal] = useState<number>(0);
  const [roleOptions, setRoleOptions] = useState<Options[]>([]);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [isPassModalOpen, setIsPassModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);

  // Form states
  const [newUser, setNewUser] = useState({
    username: '',
    nickname: '',
    password: '',
    roleId: '',
    email: '',
  });
  const [passForm, setPassForm] = useState({ prePassword: '', newPassword: '', confirmPassword: '' });
  const [error, setError] = useState('');
  const [actionLoading, setActionLoading] = useState(false);
  const [editUserId, setEditUserId] = useState<string | null>(null);
  const [editForm, setEditForm] = useState({
    username: '',
    nickname: '',
    email: '',
    roleId: '',
  });
  const [editLoading, setEditLoading] = useState(false);
  const [detailOpen, setDetailOpen] = useState(false);
  const [detailUser, setDetailUser] = useState<SystemUserDetail | null>(null);

  const formatTimestamp = (sec?: number) => {
    if (!sec) return '-';
    return new Date(sec * 1000).toLocaleString();
  };

  const getOperationType = (operationType?: number) => {
    if (typeof operationType !== 'number') return '-';
    const key = OPERATION_TYPE_MAP[operationType] || DEFAULT_OPERATION_TYPE;
    return t(key);
  };

  const mapStatusToYesOrNo = (status: 'ACTIVE' | 'INACTIVE'): YesOrNo =>
    status === 'ACTIVE' ? 1 : 2;

  const mapYesOrNoToStatus = (status: YesOrNo): 'ACTIVE' | 'INACTIVE' =>
    status === 1 ? 'ACTIVE' : 'INACTIVE';

  const mapBackendUserToFrontend = (u: SystemUserDetail): User => ({
    id: u.id,
    username: u.username,
    nickname: u.nickname,
    // 直接使用后端返回的角色名称做展示，避免依赖固定的枚举翻译
    role: (u.role_name as any) || ('-' as any),
    status: mapYesOrNoToStatus(u.status),
    lastLogin: u.last_operation_time ? new Date(u.last_operation_time * 1000).toLocaleString() : '-',
    email: u.email,
    lastPasswordChange: u.last_password_change ? u.last_password_change * 1000 : undefined,
  });

  const fetchRoles = useCallback(async () => {
    try {
      const res = await roleService.fetchRoleLabels();
      if (res.code === 50000 && res.data) {
        setRoleOptions(res.data);
      }
    } catch (e) {
      console.error('Failed to fetch role labels', e);
    }
  }, []);

  const fetchUsers = useCallback(async () => {
    setLoading(true);
    try {
      const res = await userService.fetchUserPage({ page, page_size: pageSize });
      if (res.code === 50000 && res.data) {
        setUsers((res.data.items || []).map(mapBackendUserToFrontend));
        setTotal(res.data.total || 0);
      }
    } catch (e) {
      console.error('Failed to fetch users', e);
    } finally {
      setLoading(false);
    }
  }, [page, pageSize]);

  useEffect(() => {
    fetchRoles();
  }, [fetchRoles]);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  const handleToggleStatus = async (userId: string, currentStatus: 'ACTIVE' | 'INACTIVE') => {
    if (actionLoading) return;
    setActionLoading(true);
    try {
      const payload = { status: mapStatusToYesOrNo(currentStatus) === 1 ? (2 as YesOrNo) : (1 as YesOrNo) };
      const res = await userService.toggleUserStatus(userId, payload);
      if (res.code === 50000) {
        await fetchUsers();
      }
    } catch (e) {
      console.error('Failed to toggle user status', e);
    } finally {
      setActionLoading(false);
    }
  };

  const handleOpenPassModal = (user: User) => {
    setSelectedUser(user);
    setPassForm({ prePassword: '', newPassword: '', confirmPassword: '' });
    setError('');
    setIsPassModalOpen(true);
  };

  const handleOpenDetailModal = async (userId: string) => {
    setError('');
    try {
      const res = await userService.getUserDetail(userId);
      if (res.code === 50000 && res.data) {
        setDetailUser(res.data);
        setDetailOpen(true);
      }
    } catch (e) {
      console.error('Failed to load user detail', e);
      setError(t('common.error'));
    }
  };

  const handleOpenEditModal = async (userId: string) => {
    setError('');
    setEditLoading(true);
    setIsEditModalOpen(true);
    try {
      const res = await userService.getUserDetail(userId);
      if (res.code === 50000 && res.data) {
        const u = res.data;
        setEditUserId(u.id);
        setEditForm({
          username: u.username,
          nickname: u.nickname,
          email: u.email,
          roleId: u.role_id,
        });
      }
    } catch (e) {
      console.error('Failed to load user detail', e);
      setError(t('common.error'));
    } finally {
      setEditLoading(false);
    }
  };

  const handleSaveUser = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!newUser.roleId) {
      setError(t('users.modal.roleRequired') || 'Role is required');
      return;
    }

    setActionLoading(true);
    setError('');
    try {
      const payload = {
        username: newUser.username,
        nickname: newUser.nickname || newUser.username,
        email: newUser.email,
        password: newUser.password,
        role_id: newUser.roleId,
      };
      const res = await userService.createUser(payload);
      if (res.code === 50000) {
        setIsAddModalOpen(false);
        setNewUser({ username: '', nickname: '', password: '', roleId: '', email: '' });
        await fetchUsers();
      }
    } catch (e: any) {
      console.error('Failed to create user', e);
      setError(e?.message || t('common.error') || 'Error');
    } finally {
      setActionLoading(false);
    }
  };

  const handleUpdatePassword = async (e: React.FormEvent) => {
    e.preventDefault();
    if (passForm.newPassword !== passForm.confirmPassword) {
      setError(t('users.modal.passMismatch'));
      return;
    }

    if (!selectedUser) return;

    setActionLoading(true);
    try {
      // Hash passwords with SHA256
      const hashPrePassword = forge.md.sha256.create();
      hashPrePassword.update(passForm.prePassword);
      const hashedPrePassword = hashPrePassword.digest().toHex();

      const hashNewPassword = forge.md.sha256.create();
      hashNewPassword.update(passForm.newPassword);
      const hashedNewPassword = hashNewPassword.digest().toHex();

      const hashConfirmPassword = forge.md.sha256.create();
      hashConfirmPassword.update(passForm.confirmPassword);
      const hashedConfirmPassword = hashConfirmPassword.digest().toHex();

      const payload = {
        pre_password: hashedPrePassword,
        new_password: hashedNewPassword,
        confirm_password: hashedConfirmPassword,
      };
      const res = await userService.updateUserPassword(selectedUser.id, payload);
      if (res.code === 50000) {
        setIsPassModalOpen(false);
        setSelectedUser(null);
      }
    } catch (e: any) {
      console.error('Failed to update password', e);
      setError(e?.message || t('common.error') || 'Error');
    } finally {
      setActionLoading(false);
    }
  };

  const handleUpdateUser = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editUserId) return;
    if (!editForm.roleId) {
      setError(t('users.modal.roleRequired') || 'Role is required');
      return;
    }
    setActionLoading(true);
    setError('');
    try {
      const payload = {
        username: editForm.username,
        nickname: editForm.nickname,
        email: editForm.email,
        role_id: editForm.roleId,
      };
      const res = await userService.updateUser(editUserId, payload);
      if (res.code === 50000) {
        setIsEditModalOpen(false);
        setEditUserId(null);
        await fetchUsers();
      }
    } catch (e: any) {
      console.error('Failed to update user', e);
      setError(e?.message || t('common.error') || 'Error');
    } finally {
      setActionLoading(false);
    }
  };

  return (
    <div className="p-6 max-w-[1600px] mx-auto space-y-6">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
        <div>
           <h2 className="text-2xl font-bold text-slate-900 dark:text-white">{t('users.title')}</h2>
           <p className="text-slate-500 dark:text-slate-400 text-sm mt-1">{t('users.subtitle')}</p>
        </div>
        <button 
          onClick={() => setIsAddModalOpen(true)}
          className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-medium flex items-center gap-2 transition-colors shadow-md shadow-blue-900/10 dark:shadow-blue-900/20"
        >
          <Plus size={18} /> {t('users.addUser')}
        </button>
      </div>

      <div className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl overflow-hidden shadow-sm flex-1">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-slate-50 dark:bg-slate-900/50 border-b border-slate-200 dark:border-slate-700 text-slate-500 dark:text-slate-400 text-xs uppercase tracking-wider">
                <th className="p-4 font-medium">{t('users.table.username')}</th>
                <th className="p-4 font-medium">{t('users.table.role')}</th>
                <th className="p-4 font-medium">{t('users.table.status')}</th>
                <th className="p-4 font-medium">{t('users.table.lastPasswordChange')}</th>
                <th className="p-4 font-medium text-right">{t('users.table.actions')}</th>
              </tr>
            </thead>
            <tbody className="text-sm divide-y divide-slate-200 dark:divide-slate-700">
              {loading ? (
                <tr>
                  <td colSpan={5} className="p-6 text-center text-slate-500 dark:text-slate-400">
                    {t('common.loading')}
                  </td>
                </tr>
              ) : users.length === 0 ? (
                <tr>
                  <td colSpan={5} className="p-6 text-center text-slate-500 dark:text-slate-400">
                    {t('common.empty') || 'No users'}
                  </td>
                </tr>
              ) : (
              users.map(user => (
                <tr key={user.id} className="hover:bg-slate-50 dark:hover:bg-slate-700/30 transition-colors">
                  <td className="p-4">
                    <div className="flex items-center gap-3">
                      <div className="w-8 h-8 rounded-full bg-slate-100 dark:bg-slate-700 flex items-center justify-center text-slate-500 dark:text-slate-400">
                        <UserIcon size={16} />
                      </div>
                      <div>
                        <div className="font-medium text-slate-900 dark:text-white">{user.username}</div>
                        <div className="text-xs text-slate-500">{user.email}</div>
                      </div>
                    </div>
                  </td>
                  <td className="p-4">
                    <span className="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium border bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 border-slate-200 dark:border-slate-700">
                      <Shield size={12} />
                      {user.role || '-'}
                    </span>
                  </td>
                  <td className="p-4">
                     <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${
                        user.status === 'ACTIVE' ? 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400' : 'bg-slate-100 dark:bg-slate-700 text-slate-500 dark:text-slate-400'
                        }`}>
                        <span className={`w-1.5 h-1.5 rounded-full mr-1.5 ${
                            user.status === 'ACTIVE' ? 'bg-emerald-500 dark:bg-emerald-400' : 'bg-slate-400'
                        }`}></span>
                        {t(`users.status.${user.status.toLowerCase()}`)}
                      </span>
                  </td>
                  <td className="p-4 text-slate-500 dark:text-slate-400 font-mono text-xs">
                    {user.lastPasswordChange ? new Date(user.lastPasswordChange).toLocaleString() : '-'}
                  </td>
                  <td className="p-4 text-right relative">
                    <div className="flex justify-end items-center gap-2">
                      <button
                        onClick={() => handleOpenDetailModal(user.id)}
                        className="p-1.5 text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded transition-colors"
                        title={t('common.view')}
                      >
                        <Eye size={16} />
                      </button>
                      <button
                        onClick={() => handleOpenEditModal(user.id)}
                        className="p-1.5 text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded transition-colors"
                        title={t('common.edit')}
                      >
                        <UserIcon size={16} />
                      </button>
                      <button 
                        onClick={() => handleOpenPassModal(user)}
                        className="p-1.5 text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded transition-colors"
                        title={t('users.actions.changePass')}
                      >
                        <Key size={16} />
                      </button>
                      
                      {user.status === 'ACTIVE' ? (
                        <button 
                          onClick={() => handleToggleStatus(user.id, user.status)}
                          className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-rose-600 bg-rose-50 hover:bg-rose-100 dark:bg-rose-900/20 dark:hover:bg-rose-900/30 border border-rose-100 dark:border-rose-900/50 rounded transition-colors"
                          disabled={actionLoading}
                        >
                          <Ban size={12} />
                          {t('users.actions.disable')}
                        </button>
                      ) : (
                        <button 
                          onClick={() => handleToggleStatus(user.id, user.status)}
                          className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-emerald-600 bg-emerald-50 hover:bg-emerald-100 dark:bg-emerald-900/20 dark:hover:bg-emerald-900/30 border border-emerald-100 dark:border-emerald-900/50 rounded transition-colors"
                          disabled={actionLoading}
                        >
                          <Play size={12} />
                          {t('users.actions.enable')}
                        </button>
                      )}
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
            disabled={users.length === 0 || users.length < pageSize}
            className="px-3 py-1 border border-slate-200 dark:border-slate-700 rounded hover:bg-white dark:hover:bg-slate-800 disabled:opacity-50"
          >
            {t('common.next')}
          </button>
        </div>
      </div>

      {/* Add User Modal */}
      {isAddModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 backdrop-blur-sm p-4">
          <div className="bg-white dark:bg-slate-800 rounded-xl shadow-xl w-full max-w-md border border-slate-200 dark:border-slate-700 overflow-hidden animate-fade-in-up">
            <div className="px-6 py-4 border-b border-slate-200 dark:border-slate-700 flex justify-between items-center bg-slate-50 dark:bg-slate-900/50">
              <h3 className="font-semibold text-slate-900 dark:text-white">{t('users.modal.addTitle')}</h3>
            </div>
            <form onSubmit={handleSaveUser} className="p-6 space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700 dark:text-slate-300">{t('users.modal.username')}</label>
                <div className="relative">
                  <UserIcon size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                  <input 
                    type="text" 
                    className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                    value={newUser.username}
                    onChange={e => setNewUser({...newUser, username: e.target.value})}
                    required
                  />
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700 dark:text-slate-300">{t('users.modal.nickname') || 'Nickname'}</label>
                <div className="relative">
                  <UserIcon size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                  <input 
                    type="text" 
                    className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                    value={newUser.nickname}
                    onChange={e => setNewUser({...newUser, nickname: e.target.value})}
                  />
                </div>
              </div>
              
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700 dark:text-slate-300">{t('users.modal.email')}</label>
                <div className="relative">
                  <Mail size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                  <input 
                    type="email" 
                    className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                    value={newUser.email}
                    onChange={e => setNewUser({...newUser, email: e.target.value})}
                  />
                </div>
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700 dark:text-slate-300">{t('users.modal.role')}</label>
                <div className="relative">
                  <Shield size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                  <select 
                    className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white appearance-none"
                    value={newUser.roleId}
                    onChange={e => setNewUser({...newUser, roleId: e.target.value})}
                  >
                    <option value="">{t('common.select') || 'Select'}</option>
                    {roleOptions.map(option => (
                      option.children ? (
                        <optgroup key={option.value} label={option.label}>
                          {option.children.map(child => (
                            <option key={child.value} value={child.value}>{child.label}</option>
                          ))}
                        </optgroup>
                      ) : (
                        <option key={option.value} value={option.value}>{option.label}</option>
                      )
                    ))}
                  </select>
                </div>
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700 dark:text-slate-300">{t('users.modal.password')}</label>
                <div className="relative">
                  <Lock size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                  <input 
                    type="password" 
                    className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                    value={newUser.password}
                    onChange={e => setNewUser({...newUser, password: e.target.value})}
                    required
                  />
                </div>
              </div>

              <div className="pt-4 flex justify-end gap-3">
                <button 
                  type="button"
                  onClick={() => setIsAddModalOpen(false)}
                  className="px-4 py-2 text-sm font-medium text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 rounded-lg transition-colors"
                >
                  {t('users.modal.cancel')}
                </button>
                <button 
                  type="submit"
                  className="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors shadow-sm"
                  disabled={actionLoading}
                >
                  {t('users.modal.save')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Edit User Modal */}
      {isEditModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 backdrop-blur-sm p-4">
          <div className="bg-white dark:bg-slate-800 rounded-xl shadow-xl w-full max-w-md border border-slate-200 dark:border-slate-700 overflow-hidden animate-fade-in-up">
            <div className="px-6 py-4 border-b border-slate-200 dark:border-slate-700 flex justify-between items-center bg-slate-50 dark:bg-slate-900/50">
              <h3 className="font-semibold text-slate-900 dark:text-white">
                {t('users.modal.editUserTitle') || t('common.edit')}
              </h3>
            </div>
            <form onSubmit={handleUpdateUser} className="p-6 space-y-4">
              {editLoading ? (
                <div className="text-center text-sm text-slate-500 dark:text-slate-400">
                  {t('common.loading')}
                </div>
              ) : (
                <>
                  <div className="space-y-2">
                    <label className="text-sm font-medium text-slate-700 dark:text-slate-300">{t('users.modal.username')}</label>
                    <div className="relative">
                      <UserIcon size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                      <input 
                        type="text" 
                        className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                        value={editForm.username}
                        onChange={e => setEditForm({...editForm, username: e.target.value})}
                        required
                      />
                    </div>
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium text-slate-700 dark:text-slate-300">{t('users.modal.nickname') || 'Nickname'}</label>
                    <div className="relative">
                      <UserIcon size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                      <input 
                        type="text" 
                        className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                        value={editForm.nickname}
                        onChange={e => setEditForm({...editForm, nickname: e.target.value})}
                      />
                    </div>
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium text-slate-700 dark:text-slate-300">{t('users.modal.email')}</label>
                    <div className="relative">
                      <Mail size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                      <input 
                        type="email" 
                        className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                        value={editForm.email}
                        onChange={e => setEditForm({...editForm, email: e.target.value})}
                      />
                    </div>
                  </div>
                  <div className="space-y-2">
                    <label className="text-sm font-medium text-slate-700 dark:text-slate-300">{t('users.modal.role')}</label>
                    <div className="relative">
                      <Shield size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                      <select 
                        className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white appearance-none"
                        value={editForm.roleId}
                        onChange={e => setEditForm({...editForm, roleId: e.target.value})}
                      >
                        <option value="">{t('common.select') || 'Select'}</option>
                        {roleOptions.map(option => (
                          option.children ? (
                            <optgroup key={option.value} label={option.label}>
                              {option.children.map(child => (
                                <option key={child.value} value={child.value}>{child.label}</option>
                              ))}
                            </optgroup>
                          ) : (
                            <option key={option.value} value={option.value}>{option.label}</option>
                          )
                        ))}
                      </select>
                    </div>
                  </div>
                </>
              )}

              <div className="pt-4 flex justify-end gap-3">
                <button 
                  type="button"
                  onClick={() => setIsEditModalOpen(false)}
                  className="px-4 py-2 text-sm font-medium text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 rounded-lg transition-colors"
                >
                  {t('users.modal.cancel')}
                </button>
                <button 
                  type="submit"
                  className="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors shadow-sm"
                  disabled={actionLoading || editLoading}
                >
                  {t('users.modal.update')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Detail Modal */}
      {detailOpen && detailUser && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 backdrop-blur-sm p-4">
          <div className="bg-white dark:bg-slate-800 rounded-xl shadow-xl w-full max-w-md border border-slate-200 dark:border-slate-700 overflow-hidden animate-fade-in-up">
            <div className="px-6 py-4 border-b border-slate-200 dark:border-slate-700 flex justify-between items-center bg-slate-50 dark:bg-slate-900/50">
              <h3 className="font-semibold text-slate-900 dark:text-white">
                {t('users.modal.detailTitle') || t('common.view')}
              </h3>
            </div>
            <div className="p-6 grid grid-cols-2 gap-4 text-xs text-slate-700 dark:text-slate-200">
              <div>
                <div className="text-slate-400 mb-1">ID</div>
                <div className="font-mono break-all">{detailUser.id}</div>
              </div>
              <div>
                <div className="text-slate-400 mb-1">{t('users.table.status')}</div>
                <div className={`inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium border ${
                  mapYesOrNoToStatus(detailUser.status) === 'ACTIVE'
                    ? 'bg-emerald-50 dark:bg-emerald-900/20 text-emerald-700 dark:text-emerald-300 border-emerald-200 dark:border-emerald-800'
                    : 'bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-300 border-slate-200 dark:border-slate-700'
                }`}>
                  <span className={`w-1.5 h-1.5 rounded-full ${
                    mapYesOrNoToStatus(detailUser.status) === 'ACTIVE'
                      ? 'bg-emerald-500 dark:bg-emerald-400'
                      : 'bg-slate-400'
                  }`}></span>
                  {t(`users.status.${mapYesOrNoToStatus(detailUser.status).toLowerCase()}`)}
                </div>
              </div>
              <div>
                <div className="text-slate-400 mb-1">{t('users.modal.username')}</div>
                <div>{detailUser.username}</div>
              </div>
              <div>
                <div className="text-slate-400 mb-1">{t('users.modal.nickname')}</div>
                <div>{detailUser.nickname}</div>
              </div>
              <div>
                <div className="text-slate-400 mb-1">{t('users.modal.email')}</div>
                <div>{detailUser.email}</div>
              </div>
              <div>
                <div className="text-slate-400 mb-1">{t('users.table.role')}</div>
                <div>{detailUser.role_name}</div>
              </div>
              <div>
                <div className="text-slate-400 mb-1">{t('users.table.lastPasswordChange')}</div>
                <div>{formatTimestamp(detailUser.last_password_change)}</div>
              </div>
              <div>
                <div className="text-slate-400 mb-1">{t('online_users.table.lastTime')}</div>
                <div>{formatTimestamp(detailUser.last_operation_time)}</div>
              </div>
              <div>
                <div className="text-slate-400 mb-1">{t('online_users.table.type')}</div>
                <div>{getOperationType(detailUser.operation_type as unknown as number)}</div>
              </div>
              <div className="col-span-2">
                <div className="text-slate-400 mb-1">{t('logs.table.details')}</div>
                <div className="mt-1 whitespace-pre-wrap break-words">
                  {detailUser.remark || '-'}
                </div>
              </div>
            </div>
            <div className="px-6 py-3 border-t border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-900/70 flex justify-end">
              <button
                onClick={() => setDetailOpen(false)}
                className="px-4 py-1.5 text-xs font-medium text-slate-700 dark:text-slate-200 rounded-lg border border-slate-300 dark:border-slate-600 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors"
              >
                {t('users.modal.cancel')}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Password Modal */}
      {isPassModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 backdrop-blur-sm p-4">
           <div className="bg-white dark:bg-slate-800 rounded-xl shadow-xl w-full max-w-md border border-slate-200 dark:border-slate-700 overflow-hidden animate-fade-in-up">
            <div className="px-6 py-4 border-b border-slate-200 dark:border-slate-700 flex justify-between items-center bg-slate-50 dark:bg-slate-900/50">
              <h3 className="font-semibold text-slate-900 dark:text-white">
                {t('users.modal.editPassTitle')} - <span className="text-blue-600">{selectedUser?.username}</span>
              </h3>
            </div>
            <form onSubmit={handleUpdatePassword} className="p-6 space-y-4">
              {error && (
                <div className="p-3 bg-rose-50 dark:bg-rose-900/20 text-rose-600 dark:text-rose-400 text-sm rounded-lg flex items-center gap-2">
                  <Shield size={16} />
                  {error}
                </div>
              )}
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700 dark:text-slate-300">{t('users.modal.prePassword') || 'Current Password'}</label>
                <div className="relative">
                  <Lock size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                  <input 
                    type="password" 
                    className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                    value={passForm.prePassword}
                    onChange={e => setPassForm({...passForm, prePassword: e.target.value})}
                    required
                  />
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700 dark:text-slate-300">{t('users.modal.newPassword')}</label>
                <div className="relative">
                  <Lock size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                  <input 
                    type="password" 
                    className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                    value={passForm.newPassword}
                    onChange={e => setPassForm({...passForm, newPassword: e.target.value})}
                    required
                  />
                </div>
              </div>
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700 dark:text-slate-300">{t('users.modal.confirmPassword')}</label>
                <div className="relative">
                  <Lock size={18} className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400" />
                  <input 
                    type="password" 
                    className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                    value={passForm.confirmPassword}
                    onChange={e => setPassForm({...passForm, confirmPassword: e.target.value})}
                    required
                  />
                </div>
              </div>

              <div className="pt-4 flex justify-end gap-3">
                <button 
                  type="button"
                  onClick={() => setIsPassModalOpen(false)}
                  className="px-4 py-2 text-sm font-medium text-slate-700 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-700 rounded-lg transition-colors"
                >
                  {t('users.modal.cancel')}
                </button>
                <button 
                  type="submit"
                  className="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors shadow-sm"
                >
                  {t('users.modal.update')}
                </button>
              </div>
            </form>
           </div>
        </div>
      )}
    </div>
  );
};

export default UserManager;
