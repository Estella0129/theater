<template>
  <div class="profile-container">
    <el-card class="profile-card">
      <template #header>
        <h2>个人中心</h2>
      </template>
      <el-form :model="userForm" :rules="rules" ref="userFormRef" @submit.prevent="handleUpdate">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="userForm.username" disabled></el-input>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userForm.email"></el-input>
        </el-form-item>
        <el-form-item label="角色">
          <el-input v-model="userForm.role" disabled></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" native-type="submit" :loading="loading">保存修改</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '../../stores/user'
import { ElMessage } from 'element-plus'

const userStore = useUserStore()
const userFormRef = ref(null)
const loading = ref(false)

const userForm = reactive({
  username: '',
  email: '',
  role: ''
})

const rules = {
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

onMounted(async () => {
  const user = userStore.currentUser
  if (user) {
    userForm.username = user.username
    userForm.email = user.email
    userForm.role = user.role
  }
})

const handleUpdate = async () => {
  if (!userFormRef.value) return

  await userFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const user = userStore.currentUser
        await userStore.updateUser(user.id, {
          email: userForm.email
        })
        ElMessage.success('更新成功')
      } catch (error) {
        ElMessage.error(error.message || '更新失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.profile-container {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding-top: 2rem;
}

.profile-card {
  width: 100%;
  max-width: 500px;
}
</style>