<template>
  <div class="user-list">
    <el-card>
    <div class="header">
      <h2>用户管理</h2>
      <el-button type="primary" @click="handleAdd">创建用户</el-button>
    </div>

    <el-table :data="users" v-loading="loading" border>
      <el-table-column prop="id" label="ID" width="80"></el-table-column>
      <el-table-column prop="username" label="用户名"></el-table-column>
      <el-table-column prop="name" label="姓名" width="120"></el-table-column>
      <el-table-column prop="gender" label="性别" width="100">
        <template #default="{ row }">
          {{ row.gender === 'male' ? '男' : '女' }}
        </template>
      </el-table-column>
      <el-table-column prop="email" label="邮箱"></el-table-column>
      <el-table-column prop="role" label="角色"></el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180">
        <template #default="{ row }">
          {{ new Date(row.created_at).toLocaleString() }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.is_frozen ? 'danger' : 'success'">
            {{ row.is_frozen ? '已禁用' : '正常' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="250">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
          <el-button size="small" :type="row.is_frozen ? 'success' : 'warning'" @click="handleToggleFreeze(row)">
            {{ row.is_frozen ? '启用' : '禁用' }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination
        v-model="currentPage"
        :page-size="pageSize"
        layout="prev, pager, next"
        :total="total"
        @current-change="handleCurrentChange"
      />
    </div>
    </el-card>
    <!-- 用户表单对话框 -->
    <el-dialog
      :title="dialogTitle"
      v-model="dialogVisible"
      width="500px"
      @close="resetForm"
    >
      <el-form :model="userForm" :rules="rules" ref="userFormRef" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="userForm.username" :disabled="isEdit"></el-input>
        </el-form-item>
        <el-form-item label="姓名" prop="name">
          <el-input v-model="userForm.name"></el-input>
        </el-form-item>
        <el-form-item label="性别" prop="gender">
          <el-select v-model="userForm.gender" placeholder="请选择性别">
            <el-option label="未知" value=""></el-option>
            <el-option label="男" value="male"></el-option>
            <el-option label="女" value="female"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="userForm.email"></el-input>
        </el-form-item>
        <el-form-item label="密码" prop="password" v-if="!isEdit">
          <el-input v-model="userForm.password" type="password"></el-input>
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="userForm.role">
            <el-option label="普通用户" value="user"></el-option>
            <el-option label="管理员" value="admin"></el-option>
          </el-select>
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
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '../../../stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'

const userStore = useUserStore()

// 表格数据
const users = ref([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 对话框数据
const dialogVisible = ref(false)
const dialogTitle = ref('')
const isEdit = ref(false)
const submitting = ref(false)
const userFormRef = ref(null)

const userForm = reactive({
  username: '',
  name: '',
  gender: '',
  email: '',
  password: '',
  role: 'user'
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6个字符', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

// 加载用户列表
const loadUsers = async () => {
  loading.value = true
  try {
    const response = await userStore.fetchUsers(currentPage.value, pageSize.value)
    users.value = response.results
    total.value = response.total
  } catch (error) {
    ElMessage.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

// 分页处理
const handleSizeChange = (val) => {
  pageSize.value = val
  loadUsers()
}

const handleCurrentChange = (val) => {
  currentPage.value = val
  loadUsers()
}

// 添加用户
const handleAdd = () => {
  isEdit.value = false
  dialogTitle.value = '创建用户'
  dialogVisible.value = true
}

// 编辑用户
const handleEdit = (row) => {
  isEdit.value = true
  dialogTitle.value = '编辑用户'
  Object.assign(userForm, row)
  dialogVisible.value = true
}

// 删除用户
const handleDelete = (row) => {
  ElMessageBox.confirm('确定要删除该用户吗？', '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await userStore.deleteUser(row.id)
      ElMessage.success('删除成功')
      loadUsers()
    } catch (error) {
      ElMessage.error(error.message || '删除失败')
    }
  })
}

const handleToggleFreeze = async (row) => {
  const action = row.is_frozen ? '启用' : '禁用'
  ElMessageBox.confirm(`确定要${action}该用户吗？`, '提示', {
    type: 'warning'
  }).then(async () => {
    try {
      await userStore.toggleFreezeUser(row.id, row.is_frozen ? 0 : 1)
      ElMessage.success(`${action}成功`)
      row.is_frozen = !row.is_frozen
      // 保存状态到localStorage
      localStorage.setItem(`user_${row.id}_frozen`, row.is_frozen ? '1' : '0')
      // 强制更新视图
      users.value = [...users.value]
    } catch (error) {
      ElMessage.error(error.message || `${action}失败`)
    }
  })
}

// 提交表单
const handleSubmit = async () => {
  if (!userFormRef.value) return

  await userFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (isEdit.value) {
          const { is_frozen, ...updateData } = userForm
          await userStore.updateUser(userForm.id, updateData)
          ElMessage.success('更新成功')
          // 保持当前用户的冻结状态
          const currentUser = users.value.find(u => u.id === userForm.id)
          if (currentUser) {
            currentUser.is_frozen = is_frozen // 使用解构出来的原始冻结状态
          }
        } else {
          await userStore.register(userForm)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        // 不再重新加载整个用户列表，避免状态丢失
        users.value = [...users.value]
      } catch (error) {
        if (error.message.includes('UNIQUE constraint failed: users.id')) {
          ElMessage.error('用户ID已存在，请尝试重新创建')
        } else {
          ElMessage.error(error.message || (isEdit.value ? '更新失败' : '创建失败'))
        }
      } finally {
        submitting.value = false
      }
    }
  })
}

// 重置表单
const resetForm = () => {
  if (userFormRef.value) {
    userFormRef.value.resetFields()
  }
  Object.assign(userForm, {
    username: '',
    name: '',
    gender: '',
    email: '',
    password: '',
    role: 'user'
  })
}

onMounted(() => {
  loadUsers().then(() => {
    // 从localStorage恢复冻结状态
    users.value.forEach(user => {
      const frozenState = localStorage.getItem(`user_${user.id}_frozen`)
      if (frozenState !== null) {
        user.is_frozen = frozenState === '1'
      }
    })
  })
})
</script>

<style scoped>
.user-list {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>