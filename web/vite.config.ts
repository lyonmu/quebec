import path from 'path';
import { defineConfig, loadEnv } from 'vite';
import react from '@vitejs/plugin-react';

const targetDict = {
  local: "http://127.0.0.1:59024",
};

const proxyConfig = {
  target: targetDict.local,
  changeOrigin: true,
  secure: false,
  ws: true,
};

export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, '.', '');
    return {
      server: {
        port: 3000,
        host: '0.0.0.0',
        proxy: {
          '/core/api': proxyConfig,
        },
      },
      plugins: [react()],
      define: {
        'import.meta.env.VITE_API_BASE_URL': JSON.stringify(env.VITE_API_BASE_URL || '/core/api')
      },
      resolve: {
        alias: {
          '@': path.resolve(__dirname, '.'),
        }
      },
      build: {
        outDir: path.resolve(__dirname, '../cmd/core/internal/web/dist'),
        emptyOutDir: true,
      }
    };
});
