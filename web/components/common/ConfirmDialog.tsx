import React from 'react';
import { AlertTriangle, X } from 'lucide-react';

interface ConfirmDialogProps {
  open: boolean;
  title: string;
  description?: string;
  confirmText: string;
  cancelText: string;
  onConfirm: () => void;
  onCancel: () => void;
  loading?: boolean;
}

const ConfirmDialog: React.FC<ConfirmDialogProps> = ({
  open,
  title,
  description,
  confirmText,
  cancelText,
  onConfirm,
  onCancel,
  loading = false,
}) => {
  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm p-4">
      <div className="bg-slate-900/90 text-slate-50 rounded-2xl shadow-2xl w-full max-w-md border border-slate-700/80 relative overflow-hidden">
        <button
          onClick={onCancel}
          className="absolute top-3 right-3 text-slate-400 hover:text-slate-200 transition-colors"
        >
          <X size={16} />
        </button>
        <div className="px-6 pt-6 pb-4 flex items-start gap-3">
          <div className="mt-1 flex h-9 w-9 items-center justify-center rounded-full bg-amber-500/10 border border-amber-500/40">
            <AlertTriangle className="h-5 w-5 text-amber-400" />
          </div>
          <div className="flex-1">
            <h3 className="text-sm font-semibold text-slate-50">
              {title}
            </h3>
            {description && (
              <p className="mt-2 text-xs text-slate-300 leading-relaxed">
                {description}
              </p>
            )}
          </div>
        </div>
        <div className="px-6 py-4 flex justify-end gap-3 bg-slate-900/80 border-t border-slate-800">
          <button
            type="button"
            onClick={onCancel}
            className="px-4 py-1.5 text-xs font-medium rounded-full border border-slate-600 text-slate-200 hover:bg-slate-800 transition-colors"
            disabled={loading}
          >
            {cancelText}
          </button>
          <button
            type="button"
            onClick={onConfirm}
            className="px-4 py-1.5 text-xs font-medium rounded-full bg-rose-500 hover:bg-rose-600 text-white shadow-sm disabled:opacity-60 disabled:cursor-not-allowed transition-colors"
            disabled={loading}
          >
            {confirmText}
          </button>
        </div>
      </div>
    </div>
  );
};

export default ConfirmDialog;


