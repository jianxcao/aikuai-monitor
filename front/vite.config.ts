import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  build: {
    assetsDir: 'static'
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  css: {
    preprocessorOptions: {
      less: {
        // additionalData: `@import "@/assets/base.less";`
      }
    }
  },
  server: {
    proxy: {
      '^/api/.*': {
        target: 'http://127.0.0.1:7575',
        changeOrigin: true
      }
    }
  }
})
