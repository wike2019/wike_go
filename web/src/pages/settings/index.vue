<script setup lang="ts">
import {onMounted, onUnmounted, ref} from "vue";
import {apiClient} from "../../axios/common.js"

let data=ref({})
async function getData() {
  let res=await apiClient.get("/systemInfo?log=false")
  data.value=res.data
}
let timer
onMounted(()=>{
  timer=setInterval(getData,1000)
})
onUnmounted(()=>{
  clearInterval(timer)
})
getData()
function avg(data) {
  let total=0;
  for (let i = 0; i < data.length; i++) {
    total+=data[i]
  }
  return (total/data.length).toFixed(2)
}
</script>

<template>

  <div class="main">
    <nav>
      系统管理
    </nav>
    <el-divider></el-divider>
    <div class="system">
        <p>系统:{{data.os.goos}}</p>
        <p>go版本:{{data.os.goVersion}}</p>
        <p>go协程数:{{data.os.numGoroutine}}</p>
        <p>go编译版本:{{data.os.compiler}}</p>
      <p>Cpu核数:{{data.cpu.cores}} 使用率:{{avg(data.cpu.cpus)}}%</p>
      <p>内存:{{(data.ram.usedMb/1024).toFixed(2)}}G/{{(data.ram.totalMb/1024).toFixed(2)}}G 使用率:{{data.ram.usedPercent}}%</p>
      <p>硬盘:{{data.disk.usedGb}}G/{{data.disk.totalGb}}G 使用率:{{data.disk.usedPercent}}%</p>
    </div>
  </div>

</template>

<style scoped>
.system{line-height: 2;font-size:20px;padding:5px 20px}
</style>
