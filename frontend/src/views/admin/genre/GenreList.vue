<template>
  <div class="genre-list-container">
    <div class="genre-list-card">
      <el-card>
        <template #header>
      <div class="header">
        <h2>类型管理</h2>
      <el-button type="primary" @click="showAddDialog = true">添加类型</el-button>
      </div>
      
      </template>
      <el-table :data="genres" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80"></el-table-column>
        <el-table-column prop="name" label="类型名称"></el-table-column>
        <el-table-column label="操作" width="180">
          <template #default="scope">
            <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <el-pagination
        :current-page="currentPage"
        :page-size="pageSize"
        layout="prev, pager, next"
        :total="total"
        @current-change="handlePageChange"
      ></el-pagination>
    </el-card>
    </div>
    
    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="showDialog" :title="dialogTitle">
      <el-form :model="form" label-width="80px">
        <el-form-item label="类型名称">
          <el-input v-model="form.name" autocomplete="off"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useGenreStore } from '../../../stores/genre'

const genreStore = useGenreStore()
const genres = ref([])
const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 对话框相关
const showDialog = ref(false)
const dialogTitle = ref('')
const form = ref({
  id: null,
  name: ''
})
const isEditMode = ref(false)

const loadGenres = async () => {
  loading.value = true
  try {
    const data = await genreStore.fetchGenres()
    genres.value = data
    total.value = data.length
  } finally {
    loading.value = false
  }
}

const handlePageChange = (page) => {
  currentPage.value = 1,
  loadGenres()
}

const handleEdit = (genre) => {
  form.value = { ...genre }
  dialogTitle.value = '编辑类型'
  isEditMode.value = true
  showDialog.value = true
}

const handleDelete = async (genre) => {
  try {
    // 调用删除API
    await genreStore.deleteGenre(genre.id)
    loadGenres()
  } catch (error) {
    console.error('删除失败:', error)
  }
}

const submitForm = async () => {
  try {
    if (isEditMode.value) {
      // 调用更新API
      await genreStore.updateGenre(form.value)
    } else {
      // 调用添加API
      await genreStore.addGenre(form.value)
    }
    showDialog.value = false
    loadGenres()
  } catch (error) {
    console.error('操作失败:', error)
  }
}

onMounted(() => {
  loadGenres()
})
</script>

<style scoped>
.genre-list-container {
  padding: 20px;
}
.header {
  display: flex;
  justify-content: space-between;
  align-items: center; 
}


.el-table {
  margin-top: 20px;
}

.el-button {
  margin-bottom: 20px;
}

.el-pagination {
  margin-top: 20px;
  justify-content: center;
}
</style>