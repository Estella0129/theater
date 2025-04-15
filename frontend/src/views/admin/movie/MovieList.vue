<template>
  <div class="admin-movie-list">
    <el-card>
      <template #header>
        <div class="header">
          <h2>电影管理</h2>
          <el-button type="primary" @click="handleCreate">添加电影</el-button>
        </div>
      </template>
      
      <el-table :data="movieStore.movies" style="width: 100%">
        <el-table-column prop="title" label="电影名称" width="180" />
        <el-table-column prop="director" label="导演" width="180" />
        <el-table-column prop="release_date" label="上映日期">
          <template #default="scope">
            {{ formatDate(scope.row.release_date) }}
          </template>
        </el-table-column>
        <el-table-column prop="rating" label="评分">
          <template #default="scope">
            <el-rate v-model="scope.row.rating" disabled></el-rate>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="50%">
      <el-form :model="form" label-width="120px">
        <el-form-item label="电影名称" prop="title">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="导演" prop="director">
          <el-input v-model="form.director" />
        </el-form-item>
        <el-form-item label="上映日期" prop="release_date">
          <el-date-picker v-model="form.release_date" type="date" />
        </el-form-item>
        <el-form-item label="评分" prop="rating">
          <el-rate v-model="form.rating" />
        </el-form-item>
        <el-form-item label="简介" prop="description">
          <el-input v-model="form.description" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useMovieStore } from '../../../stores/movie'
import { ElMessage } from 'element-plus'

const movieStore = useMovieStore()
const movies = ref([])

onMounted(async () => {
  await movieStore.fetchMovies()
})

const dialogVisible = ref(false)
const dialogTitle = ref('')
const form = ref({
  title: '',
  director: '',
  releaseDate: '',
  rating: 0,
  description: ''
})

const handleCreate = () => {
  dialogTitle.value = '添加电影'
  form.value = {
    title: '',
    director: '',
    releaseDate: '',
    rating: 0,
    description: ''
  }
  dialogVisible.value = true
}

const handleEdit = (movie) => {
  dialogTitle.value = '编辑电影'
  form.value = {
    title: movie.title,
    director: movie.director,
    releaseDate: movie.releaseDate,
    rating: movie.rating,
    description: movie.description || ''
  }
  dialogVisible.value = true
}

const handleDelete = async (movie) => {
  try {
    await movieStore.deleteMovie(movie.id)
    await movieStore.fetchMovies()
    ElMessage.success('删除成功')
  } catch (error) {
    ElMessage.error('删除失败: ' + error.message)
  }
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toISOString().split('T')[0]
}

const submitForm = async () => {
  try {
    if (dialogTitle.value === '添加电影') {
      await movieStore.createMovie(form.value)
    } else {
      await movieStore.updateMovie(form.value)
    }
    dialogVisible.value = false
    await movieStore.fetchMovies()
    ElMessage.success('操作成功')
  } catch (error) {
    ElMessage.error('操作失败: ' + error.message)
  }
}
</script>

<style scoped>
.admin-movie-list {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>