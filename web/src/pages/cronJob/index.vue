<script setup lang="ts">
import { ref} from "vue";
import {apiClient} from "../../axios/common.js"
  let data=ref([])
  async function getData() {
    let res=await apiClient.get("/jobList")
    data.value=res.data
  }
  getData()
</script>

<template>

  <div class="main">
    <nav>
      定时任务管理
    </nav>
    <el-divider></el-divider>
    <div>
      <el-table
          :data="data"
          style="width: 100%">
        <el-table-column
            label="定时任务编号"
            width="120">
          <template #default="scope">
            {{scope.row.ID}}
          </template>
        </el-table-column>
        <el-table-column
            label="定时任务名称"
            width="280">
          <template #default="scope">
            {{scope.row.Name}}
          </template>
        </el-table-column>
        <el-table-column
            label="定时任务时间"
            width="280">
          <template #default="scope">
            {{scope.row.Cron}}
          </template>
        </el-table-column>
        <el-table-column
            label="执行函数">
          <template #default="scope">
            {{scope.row.Func}}
          </template>
        </el-table-column>
      </el-table>
      </div>
  </div>

</template>

<style scoped>
.wt-faq{color:#000}
.wt-faq li{padding:10px 20px;line-height:1.8;list-style: none}
nav{margin-top:20px}
</style>
