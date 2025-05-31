import { createApp } from 'vue'
import './style.scss'
import App from './App.vue'
import ElementPlus from 'element-plus';
import 'element-plus/dist/index.css'
import router from './router'
import 'github-markdown-css/github-markdown.css';



const  app=createApp(App).use(ElementPlus).use(router)


app.mount('#app');
