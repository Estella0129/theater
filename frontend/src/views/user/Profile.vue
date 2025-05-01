<template>
  <div class="profile-container">
    <el-card class="profile-card">
      <template #header>
        <h2>个人中心</h2>
      </template>
      
      <div v-if="!isEditMode">
          <p>个人信息</p>
          <el-form-item label="用户名">
            <el-input v-model="userForm.username" disabled></el-input>
          </el-form-item>
          <el-form-item label="姓名">
            <el-input v-model="userForm.name" disabled></el-input>
          </el-form-item>
          <el-form-item label="邮箱">
            <el-input v-model="userForm.email" disabled></el-input>
          </el-form-item>
          <el-form-item label="角色">
            <el-input v-model="userForm.role" disabled></el-input>
          </el-form-item>
          <el-form-item label="性别">
            <el-input :value="userForm.gender === 'male' ? '男' : '女'" disabled></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="isEditMode = true">编辑信息</el-button>
          </el-form-item>
        </div>
        <div v-else>
          <p>个人信息</p>
          <el-form :model="userForm" :rules="rules" ref="userFormRef" @submit.prevent="handleUpdate">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="userForm.username" disabled></el-input>
            </el-form-item>
            <el-form-item label="姓名" prop="name">
              <el-input v-model="userForm.name"></el-input>
            </el-form-item>
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="userForm.email"></el-input>
            </el-form-item>
            <el-form-item label="角色">
              <el-input v-model="userForm.role" disabled></el-input>
            </el-form-item>
            <el-form-item label="性别" prop="gender">
              <el-select v-model="userForm.gender" placeholder="请选择性别">
                <el-option label="男" value="male"></el-option>
                <el-option label="女" value="female"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" native-type="submit" :loading="loading">保存修改</el-button>
              <el-button @click="isEditMode = false">取消</el-button>
            </el-form-item>
          </el-form>
        </div>
      
      <el-divider />
      <p>修改密码</p>
      <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" @submit.prevent="handlePasswordChange">
        <el-form-item label="当前密码" prop="currentPassword">
          <el-input v-model="passwordForm.currentPassword" type="password"></el-input>
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input v-model="passwordForm.newPassword" type="password"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" native-type="submit" :loading="passwordLoading">修改密码</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '../../stores/user'
import { ElMessage } from 'element-plus'

const isEditMode = ref(false)

const userStore = useUserStore()
const userFormRef = ref(null)
const loading = ref(false)

const userForm = reactive({
  username: '',
  name: '',
  email: '',
  role: '',
  gender: ''
})

const passwordForm = reactive({
  currentPassword: '',
  newPassword: ''
})

const passwordLoading = ref(false)

const rules = {
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

const passwordRules = {
  currentPassword: [
    { required: true, message: '请输入当前密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

onMounted(async () => {
  const user = userStore.currentUser
  if (user) {
    userForm.username = user.username
    userForm.name = user.name
    userForm.email = user.email
    userForm.role = user.role
    userForm.gender = user.gender || ''
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
          username: userForm.username,
          name: userForm.name,
          email: userForm.email,
          gender: userForm.gender
        })
        ElMessage.success('更新成功')
        isEditMode.value = false
      } catch (error) {
        ElMessage.error(error.message || '更新失败')
      } finally {
        loading.value = false
      }
    }
  })
}
const handlePasswordChange = async () => {
  if (!passwordFormRef.value) return

  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      passwordLoading.value = true
      try {
        const user = userStore.currentUser
        await userStore.changePassword({
          current_password: passwordForm.currentPassword,
          new_password: passwordForm.newPassword
        })
        ElMessage.success('密码修改成功')
        passwordForm.currentPassword = ''
        passwordForm.newPassword = ''
      } catch (error) {
        ElMessage.error(error.message || '密码修改失败')
      } finally {
        passwordLoading.value = false
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