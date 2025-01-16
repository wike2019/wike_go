import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import eslintPlugin from "vite-plugin-eslint";
import { fileURLToPath, URL } from 'node:url';

export default defineConfig({
  plugins: [vue(), eslintPlugin()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)), // 配置 @ 为 src 目录
    },
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    proxy: {
      // 代理配置
      '/core': {
        target: 'http://127.0.0.1:9999',  // 后端服务地址
        changeOrigin: true,  // 需要虚拟主机站点
      },
    }
  },
});




