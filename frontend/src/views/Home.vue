<template>
  <div class="home-container">
    <h1>电影列表</h1>
    <div class="movie-grid">
      <div v-for="movie in movies" :key="movie.id" class="movie-card">
        <img :src="movie.poster" :alt="movie.title" class="movie-poster" />
        <h3>{{ movie.title }}</h3>
        <p>{{ movie.year }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useMovieStore } from '../stores/movie'

const movieStore = useMovieStore()
const movies = ref([])

onMounted(async () => {
  await movieStore.fetchMovies()
  movies.value = movieStore.movies
})
</script>

<style scoped>
.home-container {
  padding: 2rem;
}

.movie-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 2rem;
}

.movie-card {
  cursor: pointer;
  transition: transform 0.2s;
}

.movie-card:hover {
  transform: scale(1.05);
}

.movie-poster {
  width: 100%;
  height: auto;
  border-radius: 8px;
}
</style>