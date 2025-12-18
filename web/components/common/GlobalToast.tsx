import React, { useEffect, useState } from 'react';
import { AlertCircle, CheckCircle2, Info, AlertTriangle, X } from 'lucide-react';
import { toastBus, ToastPayload } from '../../services/base/toastBus';

interface ToastState extends ToastPayload {
  id: string;
}

const iconForType = (type: ToastPayload['type']) => {
  switch (type) {
    case 'success':
      return <CheckCircle2 className="w-4 h-4 text-emerald-500" />;
    case 'warning':
      return <AlertTriangle className="w-4 h-4 text-amber-500" />;
    case 'info':
      return <Info className="w-4 h-4 text-blue-500" />;
    case 'error':
    default:
      return <AlertCircle className="w-4 h-4 text-rose-500" />;
  }
};

const GlobalToast: React.FC = () => {
  const [toasts, setToasts] = useState<ToastState[]>([]);

  useEffect(() => {
    const unsubscribe = toastBus.subscribe((toast) => {
      const id = toast.id || `${Date.now()}-${Math.random()}`;
      const next: ToastState = {
        id,
        type: toast.type || 'error',
        message: toast.message,
      };
      setToasts((prev) => [...prev, next]);
      setTimeout(() => {
        setToasts((prev) => prev.filter((t) => t.id !== id));
      }, 4000);
    });
    return unsubscribe;
  }, []);

  if (toasts.length === 0) return null;

  return (
    <div className="fixed top-4 right-4 z-50 space-y-2">
      {toasts.map((toast) => (
        <div
          key={toast.id}
          className="max-w-sm w-full bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-lg shadow-lg px-4 py-3 flex items-start gap-3 animate-fade-in-up"
        >
          <div className="mt-0.5">
            {iconForType(toast.type)}
          </div>
          <div className="flex-1 text-sm text-slate-700 dark:text-slate-200">
            {toast.message}
          </div>
          <button
            onClick={() => setToasts((prev) => prev.filter((t) => t.id !== toast.id))}
            className="ml-2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200"
          >
            <X className="w-3 h-3" />
          </button>
        </div>
      ))}
    </div>
  );
};

export default GlobalToast;


