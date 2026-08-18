import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      '/api': { target: 'http://127.0.0.1:8090', changeOrigin: true },
      '/event': { target: 'http://127.0.0.1:8090', changeOrigin: true },
      '/show': { target: 'http://127.0.0.1:8090', changeOrigin: true },
      '/ticket-type': { target: 'http://127.0.0.1:8090', changeOrigin: true },
      '/admin': { target: 'http://127.0.0.1:8090', changeOrigin: true }
    }
  }
})
