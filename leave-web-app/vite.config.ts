import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/api/leave': {
        target: process.env.LEAVE_SERVICE_URL || 'http://localhost:9090',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/leave/, ''),
      },
      '/api/user': {
        target: process.env.USER_SERVICE_URL || 'http://localhost:9091',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api\/user/, ''),
      },
    },
  },
})
