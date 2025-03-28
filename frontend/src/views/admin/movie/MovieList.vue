<template>
  <div class="admin-movie-list">
    <el-card>
      <template #header>
        <div class="header">
          <h2>电影管理</h2>
          <el-button type="primary" @click="handleCreate">添加电影</el-button>
        </div>
      </template>
      
      <el-table :data="movies" style="width: 100%">
        <el-table-column prop="title" label="电影名称" width="180" />
        <el-table-column prop="director" label="导演" width="180" />
        <el-table-column prop="releaseDate" label="上映日期" />
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useMovieStore } from '../stores/movie'

const movieStore = useMovieStore()
const movies = ref([])

onMounted(async () => {
  movies.value = await movieStore.fetchMovies()
})

const handleCreate = () => {
  // TODO: Implement create movie logic
}

const handleEdit = (movie) => {
  // TODO: Implement edit movie logic
}

const handleDelete = (movie) => {
  // TODO: Implement delete movie logic
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