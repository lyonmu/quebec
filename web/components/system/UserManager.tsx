
import React, { useState } from 'react';
import { Plus, User as UserIcon, Lock, Shield, Key, Mail, Play, Ban } from 'lucide-react';
import { mockUsers } from '../../services/base/mockData';
import { User } from '../../types';
import { useLanguage } from '../../contexts/LanguageContext';
import forge from 'node-forge';

const UserManager: React.FC = () => {
  const { t } = useLanguage();
  const [users, setUsers] = useState<User[]>(mockUsers);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [isPassModalOpen, setIsPassModalOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);

  // Form states
  const [newUser, setNewUser] = useState({ username: '', password: '', role: 'VIEWER', email: '' });
  const [passForm, setPassForm] = useState({ newPassword: '', confirmPassword: '' });
  const [error, setError] = useState('');

  const handleToggleStatus = (userId: string) => {
    setUsers(users.map(u => {
      if (u.id === userId) {
        return { ...u, status: u.status === 'ACTIVE' ? 'INACTIVE' : 'ACTIVE' };
      }
      return u;
    }));
  };

  const handleOpenPassModal = (user: User) => {
    setSelectedUser(user);
    setPassForm({ newPassword: '', confirmPassword: '' });
    setError('');
    setIsPassModalOpen(true);
  };

  const handleSaveUser = (e: React.FormEvent) => {
    e.preventDefault();

    // Hash username and password
    const mdUser = forge.md.sha256.create();
    mdUser.update(newUser.username);
    const hashedUsername = mdUser.digest().toHex();

    const mdPass = forge.md.sha256.create();
    mdPass.update(newUser.password);
    const hashedPassword = mdPass.digest().toHex();

    console.log('Creating User - Hashed Username:', hashedUsername);
    console.log('Creating User - Hashed Password:', hashedPassword);

    const user: User = {
      id: `u-${Date.now()}`,
      username: newUser.username, // Keeping original for display, but in real app would send hashed
      role: newUser.role as any,
      email: newUser.email,
      status: 'ACTIVE',
      lastLogin: '-'
    };
    setUsers([...users, user]);
    setIsAddModalOpen(false);
    setNewUser({ username: '', password: '', role: 'VIEWER', email: '' });
  };

  const handleUpdatePassword = (e: React.FormEvent) => {
    e.preventDefault();
    if (passForm.newPassword !== passForm.confirmPassword) {
      setError(t('users.modal.passMismatch'));
      return;
    }
    // Simulate API call
    setIsPassModalOpen(false);
    setSelectedUser(null);
  };

  return (
    <div className="p-6 h-full flex flex-col">
      <div className="flex justify-between items-center mb-6">
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
                <th className="p-4 font-medium">{t('users.table.lastLogin')}</th>
                <th className="p-4 font-medium text-right">{t('users.table.actions')}</th>
              </tr>
            </thead>
            <tbody className="text-sm divide-y divide-slate-200 dark:divide-slate-700">
              {users.map(user => (
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
                    <span className={`inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium border ${
                      user.role === 'ADMIN' ? 'bg-purple-50 dark:bg-purple-900/20 text-purple-700 dark:text-purple-300 border-purple-100 dark:border-purple-800' :
                      user.role === 'EDITOR' ? 'bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300 border-blue-100 dark:border-blue-800' :
                      'bg-slate-100 dark:bg-slate-800 text-slate-600 dark:text-slate-400 border-slate-200 dark:border-slate-700'
                    }`}>
                      <Shield size={12} />
                      {t(`users.roles.${user.role.toLowerCase()}`)}
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
                    {user.lastLogin}
                  </td>
                  <td className="p-4 text-right relative">
                    <div className="flex justify-end items-center gap-2">
                      <button 
                        onClick={() => handleOpenPassModal(user)}
                        className="p-1.5 text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded transition-colors"
                        title={t('users.actions.changePass')}
                      >
                        <Key size={16} />
                      </button>
                      
                      {user.status === 'ACTIVE' ? (
                        <button 
                          onClick={() => handleToggleStatus(user.id)}
                          className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-rose-600 bg-rose-50 hover:bg-rose-100 dark:bg-rose-900/20 dark:hover:bg-rose-900/30 border border-rose-100 dark:border-rose-900/50 rounded transition-colors"
                        >
                          <Ban size={12} />
                          {t('users.actions.disable')}
                        </button>
                      ) : (
                        <button 
                          onClick={() => handleToggleStatus(user.id)}
                          className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-emerald-600 bg-emerald-50 hover:bg-emerald-100 dark:bg-emerald-900/20 dark:hover:bg-emerald-900/30 border border-emerald-100 dark:border-emerald-900/50 rounded transition-colors"
                        >
                          <Play size={12} />
                          {t('users.actions.enable')}
                        </button>
                      )}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
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
                    value={newUser.role}
                    onChange={e => setNewUser({...newUser, role: e.target.value})}
                  >
                    <option value="VIEWER">{t('users.roles.viewer')}</option>
                    <option value="EDITOR">{t('users.roles.editor')}</option>
                    <option value="ADMIN">{t('users.roles.admin')}</option>
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
                >
                  {t('users.modal.save')}
                </button>
              </div>
            </form>
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
