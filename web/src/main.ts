import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css'
import router from './router'
import { addDynamicRoutes } from './router';
import axios from 'axios';
const  app=createApp(App).use(ElementPlus).use(router)

// 模拟从后端获取动态路由
async function fetchRoutes() {
    const response = await axios.get('http://127.0.0.1:9999/core/menu'); // 替换为后端接口
    addDynamicRoutes(response.data)
    localStorage.setItem('dynamicRoutes', JSON.stringify(response.data)); // 保存到 localStorage

}

// 在页面加载时获取动态路由
fetchRoutes().then(() => {
    app.mount('#app');
});
