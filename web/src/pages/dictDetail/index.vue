<script setup lang="ts">
import {reactive, ref} from "vue";
import {apiClient} from "../../axios/common.js"
import {timeFormat} from "../../tools/index.js"
import {ElMessageBox} from "element-plus";
import {useRoute} from "vue-router";
  let dialogVisible=ref(false)
  const route = useRoute();
  let data=reactive({
    ID:0,
    sysDictionaryID:0,
    status:2,
    extend: '',
    value: '',
    label: '',
    sort: 1,
    desc:""
  })
 let rules= {
   'value': [{required: true, message: '字典值必须填写', trigger: 'blur'}],
   'label': [{required: true, message: '字典名称必须填写', trigger: 'blur'}],
   'sysDictionaryID': [{required: true, message: '字典关联id必须填写', trigger: 'blur'}],
 }
  let tableData=ref([])
  function add() {
    data.ID=0
    data.status=2
    data.sysDictionaryID=route.query.id*1
    data.value=""
    data.label =""
    data.desc=""
    data.extend=""
    dialogVisible.value=true
  }
  function edit(row) {
    data.ID=row.ID
    data.status=row.status
    data.sysDictionaryID=row.sysDictionaryID
    data.value=row.value
    data.label =row.label
    data.desc=row.desc
    data.extend=row.extend
    dialogVisible.value=true
  }
  async function getData() {
    let data=await apiClient.get("/dictionaryItem",
        {
          params: {id:route.query.id}
        })
    tableData.value=data.data.sysDictionaryDetails
    console.log(tableData.value[0])
  }
  async function init() {
    await getData()
  }
  async function search() {
    await getData()
  }
  init()
let ruleFormRef=ref(null)
async function submit() {
    if (!ruleFormRef.value) return
    await ruleFormRef.value.validate(async (valid) => {
      if (valid) {
        if (data.ID===0){
          let res=await apiClient.post("/dictionaryDetail/create",data)
          if (res.code==200){
            dialogVisible.value=false
            search()
          }
        }else {
          console.log(data)
          let res=await apiClient.post("/dictionaryDetail/update?id="+data.ID,data)
          if (res.code==200){
            dialogVisible.value=false
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
          await apiClient.delete("/dictionaryDetail/delete",
              {
                params: {id:data.ID}
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
    <el-form ref="ruleFormRef" :model="data"  :rules="rules" label-width="150px">
      <el-form-item label="字典id" prop="ID" v-show="data.ID!=0">
        <el-input v-model="data.ID" :disabled="data.ID!=0" ></el-input>
      </el-form-item>
      <el-form-item label="字典关联id" prop="sysDictionaryID" >
        <el-input v-model="data.sysDictionaryID" disabled></el-input>
      </el-form-item>
      <el-form-item label="字典名称" prop="label" >
        <el-input v-model="data.label" :disabled="data.ID!=0"></el-input>
      </el-form-item>
      <el-form-item label="字典值" prop="value">
        <el-input v-model="data.value" :disabled="data.ID!=0"></el-input>
      </el-form-item>
      <el-form-item label="字典扩展属性" prop="extend">
        <el-input v-model="data.extend" ></el-input>
      </el-form-item>
      <el-form-item label="字典排序" prop="sort">
        <el-input-number v-model="data.sort" ></el-input-number>
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

            <el-button  type="success" @click="add">添加字典</el-button>
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
          label="字典名称">
        <template #default="scope">
          {{scope.row.label}}
        </template>
      </el-table-column>
      <el-table-column
          label="字典值">
        <template #default="scope">
          {{scope.row.value}}
        </template>
      </el-table-column>
      <el-table-column
          label="字典排序">
        <template #default="scope">
          {{scope.row.sort}}
        </template>
      </el-table-column>
      <el-table-column
          label="字典扩展值">
        <template #default="scope">
          {{scope.row.extend}}
        </template>
      </el-table-column>
      <el-table-column
          label="字典状态">
        <template #default="scope">
          {{scope.row.status==2?"正常":"停用"}}
        </template>
      </el-table-column>
      <el-table-column
          label="字典备注"
          >
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
          <el-button @click="edit(scope.row)" type="warning">编辑</el-button>
          <el-button type="danger" @click="del(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-divider></el-divider>
  </div>

</template>

<style scoped>
  nav{line-height: 2;border-bottom:2px solid #000;font-size:16px;margin-bottom:20px}
  .path{margin-left:10px;display: inline}
  .title{margin:10px 0 ;font-weight: bold;font-size: 16px;border-bottom:1px solid #ddd;line-height: 2;color:#646cff}
</style>
