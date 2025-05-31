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
          children: '',
          parent: '',
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
   'children': [{required: true, message: '子角色名称必须填写', trigger: 'blur'}],
 }

  let tableData=ref([])
  let  parentList=ref([])

  async function getData() {
    let data=await apiClient.get("/roleList")
    parentList.value=data.data.Names

    let res=await apiClient.get("/roleDataList")
    tableData.value=res.data
  }
  async function init() {
    await getData()
  }

  function show(data) {
    router.push({path:"/rule",query:{name :data.children}})

  }
  async function sync(){
    await apiClient.get("/ruleSync")
  }
  async function add(){
    if (!formRef.value) return
    await formRef.value.validate(async (valid) => {
      if (valid) {
        let res=await apiClient.post("/role/create",form)
        if (res.code==200){
          await getData()
        }
      }
    })
  }
  init()

let formRef=ref(null)
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
          await apiClient.delete("/role/delete",
              {
                params: {id:data.ID}
              })
          getData()
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
      角色管理
    </nav>
    <header>
      <el-form ref="formRef" :model="form" :rules="rules" inline label-width="150px">
        <el-form-item label="子角色名称" prop="children">
          <el-input v-model="form.children"></el-input>
        </el-form-item>
        <el-form-item label="继承的角色名称">
          <el-select v-model="form.parent" placeholder="父角色" style="width:120px">
            <el-option label="空角色" value="" ></el-option>
            <el-option :label="parent" :value="parent" :key="parent" v-for="parent in parentList"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
            <el-button  type="success" @click="add">添加角色</el-button>
            <el-button  type="primary" @click="sync">同步权限</el-button>
        </el-form-item>
      </el-form>
    </header>
    <el-divider></el-divider>
    <el-table
        :data="tableData"
        style="width: 100%">
      <el-table-column
          label="角色id"
          width="180">
        <template #default="scope">
          {{scope.row.ID}}
        </template>
      </el-table-column>
      <el-table-column
          label="角色名称"
          >
        <template #default="scope">
          {{scope.row.children}}
        </template>
      </el-table-column>
      <el-table-column
          label="继承角色"
          >
        <template #default="scope">
          {{scope.row.parent?scope.row.parent:"无"}}
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
          <el-button @click="show(scope.row)">配置权限</el-button>
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
  .main{flex: 1;width:100%}
  .path{margin-left:10px;display: inline}
  .title{margin:10px 0 ;font-weight: bold;font-size: 16px;border-bottom:1px solid #ddd;line-height: 2;color:#646cff}
</style>
