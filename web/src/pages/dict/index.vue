<script setup lang="ts">
import {reactive, ref} from "vue";
import {apiClient} from "../../axios/common.js"
import {timeFormat} from "../../tools/index.js"
import {ElMessageBox} from "element-plus";
import {useRouter} from "vue-router";
  let router= useRouter()
  let dialogVisible=ref(false)
  let form=reactive(
      {
          status:0,
          name: '',
          type: '',
          page: 1,
          total:0
      }
  )

  let data=reactive({
    ID:0,
    status:2,
    name: '',
    type: '',
    desc:""
  })
 let rules= {
   'name': [{required: true, message: '字典名称必须填写', trigger: 'blur'}],
   'type': [{required: true, message: '字典类型必须填写', trigger: 'blur'}],
 }
  function clear() {
    form.name=""
    form.status=0
    form.type=""
    form.page=1
  }
function  handleCurrentChange(currentPage){
  form.page = currentPage;
  getData(form)
}
  let tableData=ref([])
  function show(data){
    console.log(data)
    router.push({path:"/dictDetail",query:{id:data.ID}})
  }
  function add() {
    data.ID=0
    dialogVisible.value=true
  }
  function edit(row) {
    data.ID=row.ID
    data.status=row.status
    data.name=row.name
    data.type=row.type
    data.desc=row.desc
    dialogVisible.value=true
  }
  async function getData(body) {
    let data=await apiClient.post("/dictionaryList",body,
        {
          params: {page:body.page,count:10}
        })
     tableData.value=data.data
     console.log(tableData.value[0])
     form.total=data.sumCount
  }
  async function init() {
    await getData(form)
  }
  async function search() {
    form.page=1
    await getData(form)
  }
  init()
let ruleFormRef=ref(null)
async function submit() {
    if (!ruleFormRef.value) return
    await ruleFormRef.value.validate(async (valid) => {
      if (valid) {
        if (data.ID===0){
          delete  data.ID
          let res=await apiClient.post("/dictionary/create",data)
          if (res.code==200){
            dialogVisible.value=false
            clear()
            search()
          }
        }else {
          console.log(data)
          let res=await apiClient.post("/dictionary/update?id="+data.ID,data)
          if (res.code==200){
            dialogVisible.value=false
            clear()
            search()
          }
        }

      }
    })

}
   function del(data) {
    ElMessageBox.confirm(
        '你确定要删除吗',
        'Warning',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        }
    )
        .then(async () => {
          await apiClient.get("/dictionaryDelete",
              {
                params: {ID:data.ID}
              })
          search()
        })
        .catch(() => {
        })

  }
</script>

<template>
  <el-dialog
      title="字典类型管理"
      v-model="dialogVisible"
      width="80%">
    <el-form ref="ruleFormRef" :model="data"  :rules="rules" label-width="100px">
      <el-form-item label="字典名称" prop="name" >
        <el-input v-model="data.name"></el-input>
      </el-form-item>
      <el-form-item label="字典类型" prop="type">
        <el-input v-model="data.type"></el-input>
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="data.status" placeholder="状态" style="width:120px">
          <el-option label="正常" :value="2"></el-option>
          <el-option label="停用" :value="1"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="备注">
        <el-input v-model="data.desc" type="textarea"></el-input>
      </el-form-item>
    </el-form>
    <template #footer>
      <div class="dialog-footer">
        <el-button type="primary" @click="submit">确定</el-button>
        <el-button type="info" @click="dialogVisible = false">取消</el-button>
      </div>
    </template>
  </el-dialog>
  <div class="main">
    <nav>
      字典管理
    </nav>
    <header>
      <el-form ref="formRef" :model="form" inline label-width="100px">
        <el-form-item label="字典名称">
          <el-input v-model="form.name"></el-input>
        </el-form-item>
        <el-form-item label="字典类型">
          <el-input v-model="form.type"></el-input>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status" placeholder="状态" style="width:120px">
            <el-option label="全部" :value="0"></el-option>
            <el-option label="正常" :value="2"></el-option>
            <el-option label="停用" :value="1"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
            <el-button type="primary" @click="search">搜索</el-button>
            <el-button @click="clear">重置</el-button>
            <el-button  type="success" @click="add">添加字典</el-button>
        </el-form-item>
      </el-form>
    </header>
    <el-divider></el-divider>
    <el-table
        :data="tableData"
        style="width: 100%">
      <el-table-column
          label="字典编号"
          width="120">
        <template #default="scope">
          {{scope.row.ID}}
        </template>
      </el-table-column>
      <el-table-column
          label="字典名称"
          width="280">
        <template #default="scope">
          {{scope.row.name}}
        </template>
      </el-table-column>
      <el-table-column
          label="字典类型"
          width="250">
        <template #default="scope">
          {{scope.row.type}}
        </template>
      </el-table-column>
      <el-table-column
          label="字典状态"
          width="250">
        <template #default="scope">
          {{scope.row.status==2?"正常":"停用"}}
        </template>
      </el-table-column>
      <el-table-column
          label="字典备注"
          width="250">
        <template #default="scope">
          {{scope.row.desc}}
        </template>
      </el-table-column>
      <el-table-column
          width="180"
          label="修改时间">
        <template #default="scope">
          {{timeFormat(scope.row.UpdatedAt)}}
        </template>
      </el-table-column>
      <el-table-column width="280" label="操作">
        <template #default="scope">
          <el-button @click="show(scope.row)">添加</el-button>
          <el-button @click="edit(scope.row)" type="warning">编辑</el-button>
          <el-button type="danger" @click="del(scope.row)">删除</el-button>
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
</style>
