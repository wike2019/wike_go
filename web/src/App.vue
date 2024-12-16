<script setup lang="ts">
import {reactive, ref} from "vue";
import MenuItem from './components/MenuItem.vue'; // 子组件
const menus=reactive({

})
let activeMenu=ref("")
async function fetchMenus() {
  try {
    const savedRoutes = localStorage.getItem('dynamicRoutes');; // 替换为后端接口
    Object.assign(menus, JSON.parse(savedRoutes)); // 将后端返回的数据赋值给菜单
    activeMenu.value=menus[0]
  } catch (error) {
    console.error('获取菜单失败：', error);
  }
}
fetchMenus()
</script>

<template>
  <div>
    <el-menu :default-active="activeMenu" class="menu" mode="vertical">
      <menu-item v-for="menu in menus" :key="menu.path" :menu="menu" />
    </el-menu>
  </div>
  <div><router-view></router-view></div>
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
