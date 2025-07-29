import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  plugins: [
    vue(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  define: {
    // Inject build-time environment variables
    // eslint-disable-next-line no-undef
    'import.meta.env.VITE_COMMIT_REF': JSON.stringify(process.env.VITE_COMMIT_REF || 'dev'),
    // eslint-disable-next-line no-undef
    'import.meta.env.VITE_BUILD_DATE': JSON.stringify(process.env.VITE_BUILD_DATE || new Date().toISOString()),
    // eslint-disable-next-line no-undef
    'import.meta.env.VITE_BUILD_TIME': JSON.stringify(process.env.VITE_BUILD_TIME || new Date().toISOString()),
  },
})
