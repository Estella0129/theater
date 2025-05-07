<template>
  <div class="admin-movie-list">
    <MovieForm ref="movieFormRef" />
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
      <el-pagination
      v-model="currentPage"
      :page-size="pageSize"
      layout=" prev, pager, next"
      :total="total"
      @current-change="handlePageChange"
    />
    </el-card>
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
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

defineExpose({
  open(data) {
    dialogTitle.value = data ? '编辑电影' : '添加电影'
    if (data) {
      Object.assign(form, data)
    } else {
      resetForm()
    }
    dialogVisible.value = true
  }
})

onMounted(async () => {
  await fetchMovies()
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

const fetchMovies = async () => {
  try {
    const response = await movieStore.fetchAdminMovies({ 
      page: currentPage.value, 
      pageSize: pageSize.value 
    })
    total.value = response.total
  } catch (error) {
    ElMessage.error('获取电影列表失败: ' + error.message)
  }
}

const handlePageChange = async (page) => {
  currentPage.value = page
  await fetchMovies()
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

.el-pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>