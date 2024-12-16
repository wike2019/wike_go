import { createRouter, createWebHashHistory } from 'vue-router';

// 基础路由
const routes = [
    {
        path: '/',
        component: () => import('../pages/home/VIndex.vue') // 默认首页
    },
    {
        path: '/:catchAll(.*)', // 捕获所有未匹配的路由
        redirect:"/"
    },
];

const router = createRouter({
    history: createWebHashHistory(),
    routes
});
// 批量加载 `pages` 目录下的所有 .vue 文件
const modules = import.meta.glob('../pages/**/*.vue');
function parseRoutes(routes) {
    return routes.map(route => {
        const parsedRoute = {
            path: route.path,
            name: route.name,
            // 动态加载组件，调整路径格式
            component: modules[`../pages${route.component}.vue`] || (() => Promise.reject(`Component not found: ${route.component}`))
        };

        // 如果存在子路由，递归解析
        if (route.children && route.children.length > 0) {
            parsedRoute.children = parseRoutes(route.children);
        }

        return parsedRoute;
    });
}

export function addDynamicRoutes(routeData) {
    const dynamicRoutes = parseRoutes(routeData);
    dynamicRoutes.forEach(route => {
        router.addRoute(route); // 动态添加到路由
    });
}


export default router;
