export type ToastType = 'success' | 'error' | 'info' | 'warning';

export interface ToastPayload {
  id?: string;
  type?: ToastType;
  message: string;
}

type Listener = (toast: ToastPayload) => void;

const listeners: Listener[] = [];

export const toastBus = {
  subscribe(listener: Listener) {
    listeners.push(listener);
    return () => {
      const idx = listeners.indexOf(listener);
      if (idx >= 0) {
        listeners.splice(idx, 1);
      }
    };
  },
  emit(toast: ToastPayload) {
    listeners.forEach((l) => l(toast));
  },
};

export const showErrorToast = (message: string) => {
  toastBus.emit({ type: 'error', message });
};


