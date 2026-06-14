import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],

  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },

  server: {
    port: 5173,
    // During "make dev", Go runs on :8080 and proxies /* to here.
    // But if you open Vite directly on :5173, this proxy
    // forwards your /api calls back to Go on :8080.
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },

  build: {
    outDir: 'dist',
    // Clean the output directory before each build
    emptyOutDir: true,
    // Generate source maps for production debugging (remove if not needed)
    sourcemap: false
  }
})
