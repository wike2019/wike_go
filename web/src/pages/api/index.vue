<script setup lang="ts">
import {reactive, ref} from "vue";
import {apiClient} from "../../axios/common.js"
import {timeFormat} from "../../tools/index.js"
import { marked } from 'marked';
  let dialogVisible=ref(false)
  let form=reactive(
      {
        Group:"",
          Name: '',
          Path: '',
          Method: '',
          page: 1,
          total:0
      }
  )
  function clear() {
    form.Name=""
    form.Path=""
    form.Method=""
    form.page=1
  }
function  handleCurrentChange(currentPage){
  form.page = currentPage;
  getData(form)
}
  let input=ref("")
  let output=ref("")
  let tableData=ref([])
  function show(data){
    input.value=marked(data.Input)
    output.value=marked(data.Output)
    dialogVisible.value=true
    //marked(data.)
    console.log(data)
  }
  async function getData(body) {
    let data=await apiClient.post("/api",body,
        {
          params: {page:body.page,count:10}
        })
     tableData.value=data.data
     console.log(tableData.value[0])
     form.total=data.page.sumCount
  }
  async function init() {
    await getData(form)
  }
  async function search() {
    form.page=1
    await getData(form)
  }
  init()

</script>

<template>
  <el-dialog
      title="文档查询"
      v-model="dialogVisible"
      width="80%">
    <h2 class="title">输入参数</h2>
    <div class="markdown-body" v-html="input">

    </div>
    <h2 class="title">返回参数</h2>
    <div class="markdown-body" v-html="output">

    </div>
    <template #footer>
      <div class="dialog-footer">
        <el-button type="primary" @click="dialogVisible = false">关闭</el-button>
      </div>
    </template>
  </el-dialog>
  <div class="main">
    <nav>
      接口管理
    </nav>
    <header>
      <el-form ref="formRef" :model="form" inline label-width="100px">
        <el-form-item label="接口分组">
          <el-input v-model="form.Group"></el-input>
        </el-form-item>
        <el-form-item label="接口名称">
          <el-input v-model="form.Name"></el-input>
        </el-form-item>
        <el-form-item label="地址">
          <el-input v-model="form.Path"></el-input>
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.Method" placeholder="请选择请求方法" style="width:120px">
            <el-option label="全部" value=""></el-option>
            <el-option label="GET" value="GET"></el-option>
            <el-option label="POST" value="POST"></el-option>
            <el-option label="PUT" value="PUT"></el-option>
            <el-option label="DELETE" value="DELETE"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
            <el-button type="primary" @click="search">搜索</el-button>
            <el-button @click="clear">重置</el-button>
        </el-form-item>
      </el-form>
    </header>
    <el-divider></el-divider>
    <el-table
        :data="tableData"
        style="width: 100%">
      <el-table-column
          label="标题"
          width="250">
        <template #default="scope">
          <span :class="scope.row.Status==2?'del':''">{{scope.row.Group}}-{{scope.row.Name}}</span>
        </template>
      </el-table-column>
      <el-table-column label="请求信息	"  >
        <template #default="scope">
          <el-tag v-if="scope.row.Method=='GET'" type="success">{{scope.row.Method}}</el-tag>
          <el-tag v-if="scope.row.Method=='POST'" type="info">{{scope.row.Method}}</el-tag>
          <el-tag v-if="scope.row.Method=='PUT'" type="warning">{{scope.row.Method}}</el-tag>
          <el-tag v-if="scope.row.Method=='DELETE'" type="danger">{{scope.row.Method}}</el-tag>
          <p class="path"> {{scope.row.Path}}</p>
        </template>
      </el-table-column>
      <el-table-column
          width="180"
          label="修改时间">
        <template #default="scope">
          {{timeFormat(scope.row.UpdatedAt)}}
        </template>
      </el-table-column>
      <el-table-column width="180" label="操作">
        <template #default="scope">
          <el-button @click="show(scope.row)">查看内容</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-divider></el-divider>
    <el-pagination
        background
        @current-change="handleCurrentChange"
        layout="prev, pager, next"
        :total="form.total">
    </el-pagination>
  </div>

</template>

<style scoped>
  nav{line-height: 2;border-bottom:2px solid #000;font-size:16px;margin-bottom:20px}
  .path{margin-left:10px;display: inline}
  .title{margin:10px 0 ;font-weight: bold;font-size: 16px;border-bottom:1px solid #ddd;line-height: 2;color:#646cff}
  .del{    text-decoration: line-through;}
</style>
