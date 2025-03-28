<template>
  <div class="movie-list-container">
    <el-card class="movie-list-card">
      <template #header>
        <h2>电影列表</h2>
      </template>
      

      <el-table :data="movieStore.movies" style="width: 100%">
        <el-table-column prop="title" label="电影名称" width="180" />
        <el-table-column prop="original_title" label="原名" width="180" />
        <el-table-column prop="release_date" label="上映日期" />
        <el-table-column label="操作">
          <template #default="scope">
            <el-button size="small" @click="handleDetail(scope.row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMovieStore } from '../../stores/movie'

const router = useRouter()
const movieStore = useMovieStore()
const movies = ref([])

onMounted(async () => {
  await movieStore.fetchMovies()
  console.log('Loaded movies:', movieStore.movies)
})

const handleDetail = (movie) => {
  router.push(`/movies/${movie.id}`)
}
</script>

<style scoped>
.movie-list-container {
  padding: 20px;
}

.movie-list-card {
  max-width: 1200px;
  margin: 0 auto;
}
</style>