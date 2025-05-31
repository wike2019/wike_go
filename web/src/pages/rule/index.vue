<script setup lang="ts">
import { ref} from "vue";
import {apiClient} from "../../axios/common.js"
import {timeFormat} from "../../tools/index.js"
import {useRoute} from "vue-router";
import {ElMessage} from "element-plus";
const route = useRoute();
  let tableData=ref([])
  let value=ref([])
  let list=ref([])
  let parentList=ref([])
const props = { multiple: true , value: 'value',
  label: 'label',
  children: 'children',
  emitPath: false,
  checkStrictly: true,}
  console.log(route)
  async function getData() {
    let res=await apiClient.get("/getApiFront")
    tableData.value=res.data

    let info=await apiClient.get("/ruleInfo",{params:{ role:route.query.name}})
    value.value=info.data.apiIDs
    console.log(info)
    list.value=info.data.apiList
    parentList.value=info.data.apiParentList
  }
getData()
async function submit() {
    let data={
      role:route.query.name,
      apiId:value.value
    }
  await apiClient.post("/ruleCreate",data)
  ElMessage.success('提交成功，请耐心等待审核结果');
  await getData()
}
const handleChange = (value,selectedData) => {
  console.log(value,selectedData)
}
</script>

<template>

  <div class="main">
    <nav>
      权限管理
    </nav>
    <el-divider></el-divider>
    <el-cascader-panel
        style="width: fit-content"
        :show-all-levels="false"
        v-model="value"
        :options="tableData"
        :props="props"
        @change="handleChange"
    >
      <template #default="{data,node}">
        <span>{{ data.label }}</span>
        <span v-if="node.isLeaf" style="margin-left: 8px; color: gray; font-size: 12px">
        ({{ data.Method }} {{ data.Path }})
      </span>
      </template>
    </el-cascader-panel>
    <el-button @click="submit" style="margin-top:20px">确定</el-button>
    <el-divider></el-divider>
    <h2>自身接口</h2>
    <el-table
        :data="list"
        style="width: 100%">
      <el-table-column
          label="id编号"
          width="120">
        <template #default="scope">
          {{scope.row.ID}}
        </template>
      </el-table-column>
      <el-table-column
          label="角色名称">
        <template #default="scope">
          {{scope.row.role}}
        </template>
      </el-table-column>
      <el-table-column
          label="请求类型">
        <template #default="scope">
          {{scope.row.method}}
        </template>
      </el-table-column>
      <el-table-column
          label="请求路径">
        <template #default="scope">
          {{scope.row.path}}
        </template>
      </el-table-column>
      <el-table-column
          label="修改时间">
        <template #default="scope">
          {{timeFormat(scope.row.UpdatedAt)}}
        </template>
      </el-table-column>
    </el-table>
    <el-divider></el-divider>
    <h2>继承接口</h2>
    <el-table
        :data="parentList"
        style="width: 100%">
      <el-table-column
          label="id编号"
          width="120">
        <template #default="scope">
          {{scope.row.ID}}
        </template>
      </el-table-column>
      <el-table-column
          label="角色名称">
        <template #default="scope">
          {{scope.row.role}}
        </template>
      </el-table-column>
      <el-table-column
          label="请求类型">
        <template #default="scope">
          {{scope.row.method}}
        </template>
      </el-table-column>
      <el-table-column
          label="请求路径">
        <template #default="scope">
          {{scope.row.path}}
        </template>
      </el-table-column>
      <el-table-column
          label="修改时间">
        <template #default="scope">
          {{timeFormat(scope.row.UpdatedAt)}}
        </template>
      </el-table-column>
    </el-table>
  </div>

</template>

<style scoped>
  nav{line-height: 2;border-bottom:2px solid #000;font-size:16px;margin-bottom:20px}

  .path{margin-left:10px;display: inline}
  .title{margin:10px 0 ;font-weight: bold;font-size: 16px;border-bottom:1px solid #ddd;line-height: 2;color:#646cff}
  h2{font-size:16px;padding:10px 5px;border-bottom: #1a1a1a solid 1px}
</style>
