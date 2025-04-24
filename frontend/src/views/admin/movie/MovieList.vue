<template>
  <div class="admin-movie-list">
    <el-card>
      <template #header>
        <div class="header">
          <h2>电影管理</h2>
          <el-button type="primary" @click="openMovieForm()">添加电影</el-button>
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
            <el-button size="small" @click="openMovieForm(scope.row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <MovieForm ref="movieFormRef" @submit-success="onMovieFormSuccess" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useMovieStore } from '../../../stores/movie'
import { ElMessage } from 'element-plus'
import MovieForm from './MovieForm.vue'

const movieStore = useMovieStore()
const movies = ref([])
const movieFormRef = ref(null)

onMounted(async () => {
  await movieStore.fetchAdminMovies() // 修改为调用后台接口
})

const openMovieForm = async (movie = null) => {
  if (movie && movie.id) {
    try {
      const detail = await movieStore.getMovieById(movie.id)
      movieFormRef.value.open(detail)
    } catch (error) {
      ElMessage.error('获取电影详情失败: ' + error.message)
    }
  } else {
    movieFormRef.value.open(null)
  }
}

const onMovieFormSuccess = async () => {
  await movieStore.fetchAdminMovies() // 修改为调用后台接口
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