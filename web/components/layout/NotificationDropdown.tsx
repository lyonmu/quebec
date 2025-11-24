
import React, { useState, useEffect, useRef } from 'react';
import { Bell, Info, AlertTriangle, AlertCircle, CheckCircle, Loader2 } from 'lucide-react';
import { Notification } from '../../types';
import { fetchNotifications } from '../../services/mockData';
import { useLanguage } from '../../contexts/LanguageContext';

const NotificationDropdown: React.FC = () => {
  const { t } = useLanguage();
  const [isOpen, setIsOpen] = useState(false);
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [page, setPage] = useState(1);
  const [isLoading, setIsLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);
  const [hasUnread, setHasUnread] = useState(true);
  
  const dropdownRef = useRef<HTMLDivElement>(null);
  const listRef = useRef<HTMLDivElement>(null);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  // Load initial data or next page
  const loadNotifications = async (pageNum: number) => {
    if (isLoading) return;
    
    setIsLoading(true);
    try {
      const result = await fetchNotifications(pageNum);
      setNotifications(prev => pageNum === 1 ? result.data : [...prev, ...result.data]);
      setHasMore(result.hasMore);
    } catch (error) {
      console.error("Failed to fetch notifications", error);
    } finally {
      setIsLoading(false);
    }
  };

  // Trigger load when opening for the first time
  useEffect(() => {
    if (isOpen && notifications.length === 0) {
      loadNotifications(1);
      setHasUnread(false); // Assume opened clears unread indicator for demo
    }
  }, [isOpen]);

  // Trigger load when page changes
  useEffect(() => {
    if (page > 1) {
      loadNotifications(page);
    }
  }, [page]);

  // Infinite scroll handler
  const handleScroll = () => {
    if (listRef.current) {
      const { scrollTop, scrollHeight, clientHeight } = listRef.current;
      if (scrollHeight - scrollTop <= clientHeight + 20 && !isLoading && hasMore) {
        setPage(prev => prev + 1);
      }
    }
  };

  const getIcon = (type: string) => {
    switch (type) {
      case 'error': return <AlertCircle size={16} className="text-rose-500" />;
      case 'warning': return <AlertTriangle size={16} className="text-amber-500" />;
      case 'success': return <CheckCircle size={16} className="text-emerald-500" />;
      default: return <Info size={16} className="text-blue-500" />;
    }
  };

  const formatTime = (isoString: string) => {
    const date = new Date(isoString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMins / 60);
    const diffDays = Math.floor(diffHours / 24);

    if (diffMins < 60) return `${diffMins}m ago`;
    if (diffHours < 24) return `${diffHours}h ago`;
    return `${diffDays}d ago`;
  };

  return (
    <div className="relative" ref={dropdownRef}>
      <button 
        onClick={() => setIsOpen(!isOpen)}
        className={`relative p-2 transition-colors rounded-lg ${isOpen ? 'bg-slate-100 dark:bg-slate-800 text-slate-900 dark:text-white' : 'text-slate-500 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white'}`}
      >
        <Bell size={20} />
        {hasUnread && (
          <span className="absolute top-2 right-2 w-2 h-2 bg-rose-500 rounded-full border border-white dark:border-slate-900 animate-pulse"></span>
        )}
      </button>

      {isOpen && (
        <div className="absolute right-0 mt-2 w-80 md:w-96 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-700 rounded-xl shadow-xl z-50 overflow-hidden animate-fade-in-up">
          <div className="p-3 border-b border-slate-200 dark:border-slate-800 flex justify-between items-center bg-slate-50/50 dark:bg-slate-950/50">
            <h3 className="font-semibold text-slate-900 dark:text-white text-sm">{t('notifications.title')}</h3>
            <span className="text-xs bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 px-2 py-0.5 rounded-full font-medium">
              {notifications.length}
            </span>
          </div>

          <div 
            ref={listRef}
            onScroll={handleScroll}
            className="max-h-[400px] overflow-y-auto"
          >
            {notifications.length === 0 && !isLoading ? (
              <div className="p-8 text-center text-slate-500 dark:text-slate-400">
                <Bell size={32} className="mx-auto mb-2 opacity-20" />
                <p className="text-sm">{t('notifications.empty')}</p>
              </div>
            ) : (
              <div className="divide-y divide-slate-100 dark:divide-slate-800">
                {notifications.map((notif) => (
                  <div key={notif.id} className={`p-4 hover:bg-slate-50 dark:hover:bg-slate-800/50 transition-colors flex gap-3 ${!notif.read ? 'bg-blue-50/30 dark:bg-blue-900/10' : ''}`}>
                    <div className={`mt-0.5 flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center ${
                      notif.type === 'error' ? 'bg-rose-100 dark:bg-rose-900/20' : 
                      notif.type === 'warning' ? 'bg-amber-100 dark:bg-amber-900/20' : 
                      notif.type === 'success' ? 'bg-emerald-100 dark:bg-emerald-900/20' : 
                      'bg-blue-100 dark:bg-blue-900/20'
                    }`}>
                      {getIcon(notif.type)}
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="flex justify-between items-start mb-1">
                        <h4 className={`text-sm font-medium truncate pr-2 ${!notif.read ? 'text-slate-900 dark:text-white' : 'text-slate-700 dark:text-slate-300'}`}>
                          {notif.title}
                        </h4>
                        <span className="text-[10px] text-slate-400 dark:text-slate-500 flex-shrink-0 whitespace-nowrap">
                          {formatTime(notif.timestamp)}
                        </span>
                      </div>
                      <p className="text-xs text-slate-500 dark:text-slate-400 leading-relaxed line-clamp-2">
                        {notif.message}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            )}
            
            {isLoading && (
              <div className="p-4 flex justify-center items-center text-slate-400">
                <Loader2 size={20} className="animate-spin mr-2" />
                <span className="text-xs">{t('notifications.loading')}</span>
              </div>
            )}

            {!hasMore && notifications.length > 0 && (
              <div className="p-3 text-center text-xs text-slate-400 dark:text-slate-600 bg-slate-50 dark:bg-slate-900/30">
                {t('notifications.noMore')}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default NotificationDropdown;
