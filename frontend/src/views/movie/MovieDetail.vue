<template>
  <div class="movie-detail-container">
    <el-card class="movie-card">
      <template #header>
        <h2>{{ movie.title }}</h2>
      </template>
      
      <div class="movie-info">
        <el-image :src="movie.poster" fit="cover" class="movie-poster"></el-image>
        
        <div class="movie-meta">
          <p><strong>导演:</strong> {{ movie.director }}</p>
          <p><strong>主演:</strong> {{ movie.cast }}</p>
          <p><strong>时长:</strong> {{ movie.duration }}分钟</p>
          <p><strong>评分:</strong> <el-rate v-model="movie.rating" disabled></el-rate></p>
          <p><strong>上映日期:</strong> {{ movie.releaseDate }}</p>
          <p><strong>简介:</strong> {{ movie.description }}</p>
        </div>
      </div>
      
      <div class="movie-actions">
        <el-button type="primary" @click="goBack">返回列表</el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMovieStore } from '../../stores/movie'

const route = useRoute()
const router = useRouter()
const movieStore = useMovieStore()

const movie = ref({})

onMounted(async () => {
  const movieId = route.params.id
  movie.value = await movieStore.getMovieById(movieId)
})

const goBack = () => {
  router.push('/movies')
}
</script>

<style scoped>
.movie-detail-container {
  padding: 20px;
}

.movie-card {
  max-width: 1000px;
  margin: 0 auto;
}

.movie-info {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
}

.movie-poster {
  width: 300px;
  height: 450px;
}

.movie-meta {
  flex: 1;
}

.movie-meta p {
  margin-bottom: 10px;
}

.movie-actions {
  margin-top: 20px;
  text-align: center;
}
</style>