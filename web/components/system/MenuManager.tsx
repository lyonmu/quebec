import React, { useEffect, useState, useCallback } from 'react';
import { Plus, Menu as MenuIcon, Trash2, Edit, Eye, Ban, Play, ChevronDown, ChevronRight, X } from 'lucide-react';
import { useLanguage } from '../../contexts/LanguageContext';
import { menuService, roleMenuService } from '../../services/system/menuService';
import { SystemMenu, SystemMenuTreeItem, MenuType, YesOrNo } from '../../types';

type MenuStatus = 'ENABLED' | 'DISABLED';

interface MenuView {
  id: string;
  name: string;
  menu_type: MenuType;
  api_path: string;
  order: number;
  status: MenuStatus;
  children?: MenuView[];
}

const MenuManager: React.FC = () => {
  const { t } = useLanguage();
  const [menus, setMenus] = useState<SystemMenuTreeItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [editingMenu, setEditingMenu] = useState<SystemMenu | null>(null);
  const [form, setForm] = useState({
    name: '',
    menu_type: 1 as MenuType,
    api_path: '',
    api_path_method: '',
    order: 1,
    parent_id: '',
    component: '',
    status: 1 as YesOrNo,
    remark: '',
  });
  const [actionLoading, setActionLoading] = useState(false);
  const [parentMenuOptions, setParentMenuOptions] = useState<Array<{ label: string; value: string }>>([]);

  const mapYesOrNoToStatus = (status: YesOrNo): MenuStatus =>
    status === 1 ? 'ENABLED' : 'DISABLED';

  const mapStatusToYesOrNo = (status: MenuStatus): YesOrNo =>
    status === 'ENABLED' ? 1 : 2;

  const fetchMenus = useCallback(async () => {
    setLoading(true);
    try {
      const res = await menuService.fetchMenuTree();
      if (res.code === 50000 && res.data) {
        setMenus(res.data || []);
      }
    } catch (e) {
      console.error('Failed to fetch menus', e);
    } finally {
      setLoading(false);
    }
  }, []);

  const fetchParentMenus = useCallback(async () => {
    try {
      const res = await menuService.fetchMenuLabels();
      if (res.code === 50000 && res.data) {
        setParentMenuOptions(res.data || []);
      }
    } catch (e) {
      console.error('Failed to fetch parent menus', e);
    }
  }, []);

  useEffect(() => {
    fetchMenus();
    fetchParentMenus();
  }, [fetchMenus, fetchParentMenus]);

  const handleToggleStatus = async (menu: SystemMenuTreeItem) => {
    if (actionLoading) return;
    setActionLoading(true);
    try {
      const targetStatus = menu.status === 1 ? 2 : 1;
      const res = await menuService.toggleMenuStatus(menu.id, { status: targetStatus });
      if (res.code === 50000) {
        await fetchMenus();
      }
    } catch (e) {
      console.error('Failed to toggle menu status', e);
    } finally {
      setActionLoading(false);
    }
  };

  const handleDeleteMenu = async (menuId: string) => {
    if (!window.confirm(t('menus.confirmDelete'))) return;
    if (actionLoading) return;
    setActionLoading(true);
    try {
      const res = await menuService.deleteMenu(menuId);
      if (res.code === 50000) {
        await fetchMenus();
      }
    } catch (e) {
      console.error('Failed to delete menu', e);
    } finally {
      setActionLoading(false);
    }
  };

  const openCreateModal = () => {
    setEditingMenu(null);
    setForm({
      name: '',
      menu_type: 1,
      api_path: '',
      api_path_method: '',
      order: 1,
      parent_id: '',
      component: '',
      status: 1,
      remark: '',
    });
    setIsAddModalOpen(true);
  };

  const openEditModal = async (menuId: string) => {
    try {
      const res = await menuService.getMenuDetail(menuId);
      if (res.code === 50000 && res.data) {
        setEditingMenu(res.data);
        setForm({
          name: res.data.name,
          menu_type: res.data.menu_type,
          api_path: res.data.api_path || '',
          api_path_method: res.data.api_path_method || '',
          order: res.data.order,
          parent_id: res.data.parent_id || '',
          component: res.data.component || '',
          status: res.data.status,
          remark: res.data.remark || '',
        });
        setIsAddModalOpen(true);
      }
    } catch (e) {
      console.error('Failed to load menu detail', e);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (actionLoading) return;
    setActionLoading(true);
    try {
      if (editingMenu) {
        const res = await menuService.updateMenu(editingMenu.id, {
          name: form.name,
          menu_type: form.menu_type,
          api_path: form.api_path || undefined,
          api_path_method: form.api_path_method || undefined,
          order: form.order,
          parent_id: form.parent_id || undefined,
          component: form.component || undefined,
          status: form.status,
          remark: form.remark || undefined,
        });
        if (res.code === 50000) {
          setIsAddModalOpen(false);
          await fetchMenus();
        }
      } else {
        const res = await menuService.createMenu({
          name: form.name,
          menu_type: form.menu_type,
          api_path: form.api_path || undefined,
          api_path_method: form.api_path_method || undefined,
          order: form.order,
          parent_id: form.parent_id || undefined,
          component: form.component || undefined,
          status: form.status,
          remark: form.remark || undefined,
        });
        if (res.code === 50000) {
          setIsAddModalOpen(false);
          await fetchMenus();
        }
      }
    } catch (e) {
      console.error('Failed to submit menu', e);
    } finally {
      setActionLoading(false);
    }
  };

  const getMenuTypeLabel = (type: MenuType): string => {
    switch (type) {
      case 1: return t('menus.type.directory');
      case 2: return t('menus.type.menu');
      case 3: return t('menus.type.button');
      default: return '';
    }
  };

  const renderMenuRow = (menu: SystemMenuTreeItem, level: number = 0) => (
    <React.Fragment key={menu.id}>
      <tr className="hover:bg-slate-50 dark:hover:bg-slate-700/30 transition-colors">
        <td className="p-4">
          <div className="flex items-center gap-3">
            <div
              className="p-2 rounded-lg bg-slate-100 text-slate-600 dark:bg-slate-800 dark:text-slate-400"
              style={{ paddingLeft: `${8 + level * 20}px` }}
            >
              <MenuIcon size={18} />
            </div>
            <span className="font-medium text-slate-900 dark:text-white">{menu.name}</span>
          </div>
        </td>
        <td className="p-4 text-slate-600 dark:text-slate-300">
          {getMenuTypeLabel(menu.menu_type)}
        </td>
        <td className="p-4 text-slate-600 dark:text-slate-300 text-xs font-mono">
          {menu.api_path || '-'}
        </td>
        <td className="p-4 text-slate-600 dark:text-slate-300">
          {menu.order}
        </td>
        <td className="p-4">
          <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${
            menu.status === 1
              ? 'bg-emerald-50 dark:bg-emerald-500/10 text-emerald-600 dark:text-emerald-400'
              : 'bg-slate-100 dark:bg-slate-700 text-slate-500 dark:text-slate-400'
          }`}>
            <span className={`w-1.5 h-1.5 rounded-full mr-1.5 ${
              menu.status === 1 ? 'bg-emerald-500 dark:bg-emerald-400' : 'bg-slate-400'
            }`}></span>
            {menu.status === 1 ? t('menus.status.enabled') : t('menus.status.disabled')}
          </span>
        </td>
        <td className="p-4 text-right">
          <div className="flex justify-end items-center gap-2">
            {menu.status === 1 ? (
              <button
                onClick={() => handleToggleStatus(menu)}
                className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-rose-600 bg-rose-50 hover:bg-rose-100 dark:bg-rose-900/20 dark:hover:bg-rose-900/30 border border-rose-100 dark:border-rose-900/50 rounded transition-colors"
                disabled={actionLoading}
                title={t('users.actions.disable')}
              >
                <Ban size={12} />
              </button>
            ) : (
              <button
                onClick={() => handleToggleStatus(menu)}
                className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-emerald-600 bg-emerald-50 hover:bg-emerald-100 dark:bg-emerald-900/20 dark:hover:bg-emerald-900/30 border border-emerald-100 dark:border-emerald-900/50 rounded transition-colors"
                disabled={actionLoading}
                title={t('users.actions.enable')}
              >
                <Play size={12} />
              </button>
            )}

            <button
              onClick={() => openEditModal(menu.id)}
              className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-slate-500 hover:text-blue-600 bg-slate-100 hover:bg-blue-50 dark:bg-slate-700 dark:hover:bg-blue-900/20 dark:text-slate-400 dark:hover:text-blue-400 border border-transparent hover:border-blue-100 dark:hover:border-blue-900/50 rounded transition-all"
              title={t('common.edit')}
            >
              <Edit size={12} />
            </button>
            <button
              onClick={() => handleDeleteMenu(menu.id)}
              className="flex items-center gap-1 px-2 py-1 text-xs font-medium text-slate-500 hover:text-rose-600 bg-slate-100 hover:bg-rose-50 dark:bg-slate-700 dark:hover:bg-rose-900/20 dark:text-slate-400 dark:hover:text-rose-400 border border-transparent hover:border-rose-100 dark:hover:border-rose-900/50 rounded transition-all"
              title={t('users.actions.delete')}
            >
              <Trash2 size={12} />
            </button>
          </div>
        </td>
      </tr>
      {menu.children && menu.children.map(child => renderMenuRow(child, level + 1))}
    </React.Fragment>
  );

  return (
    <div className="p-6 h-full flex flex-col">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h2 className="text-2xl font-bold text-slate-900 dark:text-white">{t('menus.title')}</h2>
          <p className="text-slate-500 dark:text-slate-400 text-sm mt-1">{t('menus.subtitle')}</p>
        </div>
        <button
          onClick={openCreateModal}
          className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-medium flex items-center gap-2 transition-colors shadow-md shadow-blue-900/10 dark:shadow-blue-900/20"
        >
          <Plus size={18} /> {t('menus.addMenu')}
        </button>
      </div>

      <div className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl overflow-hidden shadow-sm">
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-slate-50 dark:bg-slate-900/50 border-b border-slate-200 dark:border-slate-700 text-slate-500 dark:text-slate-400 text-xs uppercase tracking-wider">
                <th className="p-4 font-medium">{t('menus.table.name')}</th>
                <th className="p-4 font-medium">{t('menus.table.type')}</th>
                <th className="p-4 font-medium">{t('menus.table.api')}</th>
                <th className="p-4 font-medium">{t('menus.table.order')}</th>
                <th className="p-4 font-medium">{t('menus.table.status')}</th>
                <th className="p-4 font-medium text-right">{t('menus.table.actions')}</th>
              </tr>
            </thead>
            <tbody className="text-sm divide-y divide-slate-200 dark:divide-slate-700">
              {loading ? (
                <tr>
                  <td colSpan={6} className="p-6 text-center text-slate-500 dark:text-slate-400">
                    {t('common.loading')}
                  </td>
                </tr>
              ) : menus.length === 0 ? (
                <tr>
                  <td colSpan={6} className="p-6 text-center text-slate-500 dark:text-slate-400">
                    {t('common.empty') || 'No menus'}
                  </td>
                </tr>
              ) : (
                menus.map(menu => renderMenuRow(menu))
              )}
            </tbody>
          </table>
        </div>
      </div>

      {/* Add/Edit Menu Modal */}
      {isAddModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/50 backdrop-blur-sm p-4">
          <div className="bg-white dark:bg-slate-800 rounded-xl shadow-xl w-full max-w-lg border border-slate-200 dark:border-slate-700 overflow-hidden animate-fade-in-up max-h-[90vh] overflow-y-auto">
            <div className="px-6 py-4 border-b border-slate-200 dark:border-slate-700 flex justify-between items-center bg-slate-50 dark:bg-slate-900/50 sticky top-0">
              <h3 className="font-semibold text-slate-900 dark:text-white">
                {editingMenu ? t('menus.modal.editTitle') : t('menus.modal.addTitle')}
              </h3>
              <button
                onClick={() => setIsAddModalOpen(false)}
                className="text-slate-400 hover:text-slate-600 dark:hover:text-slate-300"
              >
                <X size={20} />
              </button>
            </div>
            <form onSubmit={handleSubmit} className="p-6 space-y-4">
              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700 dark:text-slate-300">
                  {t('menus.modal.name')} <span className="text-rose-500">*</span>
                </label>
                <input
                  type="text"
                  className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                  value={form.name}
                  onChange={e => setForm({...form, name: e.target.value})}
                  required
                />
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-700 dark:text-slate-300">
                    {t('menus.modal.type')} <span className="text-rose-500">*</span>
                  </label>
                  <select
                    className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                    value={form.menu_type}
                    onChange={e => setForm({...form, menu_type: Number(e.target.value) as MenuType})}
                  >
                    <option value={1}>{t('menus.modal.type_directory')}</option>
                    <option value={2}>{t('menus.modal.type_menu')}</option>
                    <option value={3}>{t('menus.modal.type_button')}</option>
                  </select>
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-700 dark:text-slate-300">
                    {t('menus.modal.order')}
                  </label>
                  <input
                    type="number"
                    className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                    value={form.order}
                    onChange={e => setForm({...form, order: Number(e.target.value)})}
                    min={1}
                  />
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-700 dark:text-slate-300">
                    {t('menus.modal.apiPath')}
                  </label>
                  <input
                    type="text"
                    className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white font-mono"
                    value={form.api_path}
                    onChange={e => setForm({...form, api_path: e.target.value})}
                    placeholder="/api/example"
                  />
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-700 dark:text-slate-300">
                    {t('menus.modal.apiMethod')}
                  </label>
                  <select
                    className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                    value={form.api_path_method}
                    onChange={e => setForm({...form, api_path_method: e.target.value})}
                  >
                    <option value="">{t('common.select')}</option>
                    <option value="GET">GET</option>
                    <option value="POST">POST</option>
                    <option value="PUT">PUT</option>
                    <option value="DELETE">DELETE</option>
                    <option value="PATCH">PATCH</option>
                  </select>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4">
                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-700 dark:text-slate-300">
                    {t('menus.modal.parent')}
                  </label>
                  <select
                    className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                    value={form.parent_id}
                    onChange={e => setForm({...form, parent_id: e.target.value})}
                  >
                    <option value="">{t('common.select')}</option>
                    {parentMenuOptions.map(opt => (
                      <option key={opt.value} value={opt.value}>{opt.label}</option>
                    ))}
                  </select>
                </div>

                <div className="space-y-2">
                  <label className="text-sm font-medium text-slate-700 dark:text-slate-300">
                    {t('menus.modal.status')}
                  </label>
                  <select
                    className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                    value={form.status}
                    onChange={e => setForm({...form, status: Number(e.target.value) as YesOrNo})}
                  >
                    <option value={1}>{t('menus.status.enabled')}</option>
                    <option value={2}>{t('menus.status.disabled')}</option>
                  </select>
                </div>
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700 dark:text-slate-300">
                  {t('menus.modal.component')}
                </label>
                <input
                  type="text"
                  className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                  value={form.component}
                  onChange={e => setForm({...form, component: e.target.value})}
                  placeholder="/system/users"
                />
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium text-slate-700 dark:text-slate-300">
                  {t('menus.modal.remark')}
                </label>
                <textarea
                  className="w-full bg-slate-50 dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 dark:text-white"
                  rows={2}
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
                  {t('menus.modal.cancel')}
                </button>
                <button
                  type="submit"
                  className="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors shadow-sm"
                  disabled={actionLoading}
                >
                  {t('menus.modal.save')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default MenuManager;
