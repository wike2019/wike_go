<template>
  {{menu}}
  <el-sub-menu v-if="menu.children && menu.children.length > 0" :index="menu.path">
    <template #title>
      <span @click="rediect(menu.name)">{{ menu.name }}</span>
    </template>
    <menu-item
        v-for="child in menu.children"
        :key="child.path"
        :menu="child"
    />
  </el-sub-menu>
  <el-menu-item v-else :index="menu.path" @click="rediect(menu.name)">
    {{ menu.name }}
  </el-menu-item>
</template>

<script setup>
import router from "../router";

function rediect(name){
  router.push({name})
}
defineProps({
  menu: {
    type: Object,
    required: true,
  },
});

</script>
