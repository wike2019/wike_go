<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<template>
  <div class="page">
    <aside class="nav">
        <h1>WIKE管理系统</h1>
        <nav>
          <p @click="router.push('/api')">接口管理</p>
          <p @click="router.push('/dict')">字典管理</p>
          <p @click="router.push('/role')">角色管理</p>
          <p @click="router.push('/log')">日志管理</p>
          <p @click="router.push('/cronJob')">定时任务管理</p>
          <p @click="router.push('/settings')">系统配置</p>
        </nav>
    </aside>
    <main>
      <router-view></router-view>
    </main>
  </div>
</template>

<style scoped>
.page{display: flex;flex: 1}
.page .nav{width:200px;padding:10px;background:#000;color: #fff;height:100%;min-height:800px;margin-right:10px}
.page h1{font-size:20px;text-align: center}
 nav{margin-top:20px;width:200px}
.nav p{line-height:2;padding-left:10px;font-size:14px;border-bottom:1px solid #ccc;cursor: pointer;;margin:10px 0}
</style>
