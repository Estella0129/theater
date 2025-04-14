<template>
  <div class="people-form">
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="500px"
      @close="resetForm"
    >
      <el-form :model="peopleForm" :rules="rules" ref="peopleFormRef" label-width="100px">
        <el-form-item label="姓名" prop="name">
          <el-input v-model="peopleForm.name"></el-input>
        </el-form-item>
        <el-form-item label="原名" prop="original_name">
          <el-input v-model="peopleForm.original_name"></el-input>
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="peopleForm.role">
            <el-option label="演员" value="actor"></el-option>
            <el-option label="导演" value="director"></el-option>
            <el-option label="编剧" value="writer"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="性别" prop="gender">
          <el-radio-group v-model="peopleForm.gender">
            <el-radio :value="0">未知</el-radio>
            <el-radio :value="1">女</el-radio>
            <el-radio :value="2">男</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="出生日期" prop="birthday">
          <el-date-picker
            v-model="peopleForm.birthday"
            type="date"
            placeholder="选择日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="出生地" prop="place_of_birth">
          <el-input v-model="peopleForm.place_of_birth"></el-input>
        </el-form-item>
        <el-form-item label="简介" prop="biography">
          <el-input
            v-model="peopleForm.biography"
            type="textarea"
            :rows="4"
            placeholder="请输入简介"
          />
        </el-form-item>
        <el-form-item label="头像" prop="profile_path">
          <el-input v-model="peopleForm.profile_path" placeholder="输入图片URL"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { usePeopleStore } from '../../../stores/people'
import { ElMessage } from 'element-plus'

const emit = defineEmits(['refresh'])

const peopleStore = usePeopleStore()

// 对话框数据
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const submitting = ref(false)
const peopleFormRef = ref(null)

const peopleForm = reactive({
  name: '',
  role: 'actor',
  gender: 0,
  original_name: '',
  biography: '',
  birthday: '',
  place_of_birth: '',
  profile_path: ''
})

const rules = {
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ],
  gender: [
    { required: true, message: '请选择性别', trigger: 'change' }
  ]
}

// 打开表单
const open = (title, isEditMode = false, data = null) => {
  dialogTitle.value = title
  isEdit.value = isEditMode
  if (data) {
    Object.assign(peopleForm, data)
  }
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  if (!peopleFormRef.value) return

  await peopleFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (isEdit.value) {
          await peopleStore.updatePeople(peopleForm.id, peopleForm)
          ElMessage.success('更新成功')
        } else {
          await peopleStore.createPeople(peopleForm)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        emit('refresh')
      } catch (error) {
        ElMessage.error(error.message || (isEdit.value ? '更新失败' : '创建失败'))
      } finally {
        submitting.value = false
      }
    }
  })
}

// 重置表单
const resetForm = () => {
  if (peopleFormRef.value) {
    peopleFormRef.value.resetFields()
  }
  Object.assign(peopleForm, {
    name: '',
    role: 'actor',
    gender: 0,
    original_name: '',
    biography: '',
    birthday: '',
    place_of_birth: '',
    profile_path: ''
  })
}

defineExpose({
  open
})
</script>

<style scoped>
.people-form {
  padding: 20px;
}
</style>