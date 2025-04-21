<template>
  <div class="home-container">
    <h1>电影列表</h1>
    <div class="movie-grid">
      <div v-for="movie in movies" :key="movie.id" class="movie-card" @click="handleDetail(movie)">
        <el-image :src="getProfileImage(movie.poster_path)" :alt="movie.title" fit="cover" class="movie-poster" :fallback="'https://via.placeholder.com/200x300?text=No+Image'" />
        <h3>{{ movie.title }}</h3>
        <p>{{ movie.year }}</p>
      </div>
    </div>
    <h1>动画电影</h1>
    <div class="movie-grid">
      <div v-for="movie in animationMovies" :key="movie.id" class="movie-card" @click="handleDetail(movie)">
        <el-image :src="getProfileImage(movie.poster_path)" :alt="movie.title" fit="cover" class="movie-poster" :fallback="'https://via.placeholder.com/200x300?text=No+Image'" />
        <h3>{{ movie.title }}</h3>
        <p>{{ movie.year }}</p>
      </div>
    </div>
    <h1>动作电影</h1>
    <div class="movie-grid">
      <div v-for="movie in actionMovies" :key="movie.id" class="movie-card" @click="handleDetail(movie)">
        <el-image :src="getProfileImage(movie.poster_path)" :alt="movie.title" fit="cover" class="movie-poster" :fallback="'https://via.placeholder.com/200x300?text=No+Image'" />
        <h3>{{ movie.title }}</h3>
        <p>{{ movie.year }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useMovieStore } from '../stores/movie'
import { useRouter } from 'vue-router'

const movieStore = useMovieStore()
const router = useRouter()
const movies = ref([])
const animationMovies = ref([])
const actionMovies = ref([])

const getProfileImage = (path) => {
  return path 
    ? `https://image.tmdb.org/t/p/w200${path}` 
    : 'https://via.placeholder.com/200x300?text=No+Image';
};

onMounted(async () => {
  await movieStore.fetchMovies()
  movies.value = movieStore.movies
  
  // 获取动画类型电影 (genre_id=16)
  const animationData = await movieStore.fetchMovies({genre: 14})
  animationMovies.value = animationData.results
  //获取动作类型电影 (genre_id=28)
  const actionData = await movieStore.fetchMovies({genre: 28})
  actionMovies.value = actionData.results
})

const handleDetail = (movie) => {
  router.push(`/movies/${movie.id}`)
}
</script>

<style scoped>
.home-container {
  padding: 2rem;
}

.movie-grid {
  display: flex;
  overflow-x: auto;
  gap: 1.5rem;
  padding: 1rem 0;
  scrollbar-width: none;
}

.movie-grid::-webkit-scrollbar {
  display: none;
}

.movie-card {
  cursor: pointer;
  transition: transform 0.2s;
  flex: 0 0 auto;
  width: 150px;
}

.movie-card:hover {
  transform: scale(1.05);
}

.movie-poster {
  width: 150px;
  height: 225px;
  border-radius: 8px;
  object-fit: cover;
}
</style>