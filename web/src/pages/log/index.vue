<script setup lang="ts">
import { ref} from "vue";
import {apiClient} from "../../axios/common.js"
  let log=ref([])
  let data=ref([])
  async function getData() {
    let res=await apiClient.get("/log")
    log.value=res.data
    let page=1
    data.value=log.value.slice((page-1)*10,(page)*10)
  }
  function handleCurrentChange(page) {
    data.value=log.value.slice((page-1)*10,(page)*10)
  }
  getData()
</script>

<template>

  <div class="main">
    <nav>
      日志管理
    </nav>
    <el-divider></el-divider>
    <div>
      <div class="wt-faq">
        <li v-for="(line,index) in data" :key="index" >
          {{line.Text}}
        </li>
      </div>
      <el-pagination
          background
          @current-change="handleCurrentChange"
          layout="prev, pager, next"
          :total="log.length">
      </el-pagination>
    </div>
  </div>

</template>

<style scoped>
.wt-faq{color:#000}
.wt-faq li{padding:10px 20px;line-height:1.8;list-style: none}
nav{margin-top:20px}
</style>
