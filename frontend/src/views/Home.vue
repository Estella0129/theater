<template>
  <div class="home-container">
    <!-- 新增banner区域 -->
    <div class="banner">
      <div class="banner-content">
        <h1>欢迎来到网上影院</h1>
        <p>这里有海量的电影和人物等你发现。快来探索吧！</p>
        <div class="banner-search">
          <el-input
            v-model="searchQuery"
            placeholder="搜索电影..."
            @keyup.enter="handleSearch"
          >
            <template #append>
              <el-button :icon="Search" @click="handleSearch" />
            </template>
          </el-input>
        </div>
      </div>
    </div>

    <!-- 原有内容保持不变 -->
    <h1>热门电影</h1>
    <div class="movie-grid">
      <div v-for="movie in movies" :key="movie.id" class="movie-card" @click="handleDetail(movie)">
        <el-image :src="getProfileImage(movie.poster_path)" :alt="movie.title" fit="cover" class="movie-poster" :fallback="'https://via.placeholder.com/200x300?text=No+Image'" />
        <h3>{{ movie.title }}</h3>
        <p>{{ movie.release_date ? movie.release_date.split('T')[0] : '' }}</p>
      </div>
    </div>
    <h1>动画电影</h1>
    <div class="movie-grid">
      <div v-for="movie in animationMovies" :key="movie.id" class="movie-card" @click="handleDetail(movie)">
        <el-image :src="getProfileImage(movie.poster_path)" :alt="movie.title" fit="cover" class="movie-poster" :fallback="'https://via.placeholder.com/200x300?text=No+Image'" />
        <h3>{{ movie.title }}</h3>
        <p>{{ movie.release_date ? movie.release_date.split('T')[0] : '' }}</p>
      </div>
    </div>
    <h1>动作电影</h1>
    <div class="movie-grid">
      <div v-for="movie in actionMovies" :key="movie.id" class="movie-card" @click="handleDetail(movie)">
        <el-image :src="getProfileImage(movie.poster_path)" :alt="movie.title" fit="cover" class="movie-poster" :fallback="'https://via.placeholder.com/200x300?text=No+Image'" />
        <h3>{{ movie.title }}</h3>
        <p>{{ movie.release_date ? movie.release_date.split('T')[0] : '' }}</p>
      </div>
    </div>
    <h1>评分最高</h1>
    <div class="movie-grid">
      <div v-for="movie in topRatedMovies" :key="movie.id" class="movie-card" @click="handleDetail(movie)">
        <el-image :src="getProfileImage(movie.poster_path)" :alt="movie.title" fit="cover" class="movie-poster" :fallback="'https://via.placeholder.com/200x300?text=No+Image'" />
        <h3>{{ movie.title }}</h3>
        <p>{{ movie.release_date ? movie.release_date.split('T')[0] : '' }}</p>
      </div>
    </div>
    <h1>趋势榜单</h1>
    <div class="movie-grid">
      <div v-for="movie in trendingMovies" :key="movie.id" class="movie-card" @click="handleDetail(movie)">
        <el-image :src="getProfileImage(movie.poster_path)" :alt="movie.title" fit="cover" class="movie-poster" :fallback="'https://via.placeholder.com/200x300?text=No+Image'" />
        <h3>{{ movie.title }}</h3>
        <p>{{ movie.release_date ? movie.release_date.split('T')[0] : '' }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useMovieStore } from '../stores/movie'
import { useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'

const movieStore = useMovieStore()
const router = useRouter()
const movies = ref([])
const animationMovies = ref([])
const actionMovies = ref([])
const topRatedMovies = ref([])
const trendingMovies = ref([])
const searchQuery = ref('')

const getProfileImage = (path) => {
  return path 
    ? `https://image.tmdb.org/t/p/w200${path}` 
    : 'https://via.placeholder.com/200x300?text=No+Image';
};

onMounted(async () => {
  await movieStore.fetchMovies()
  movies.value = movieStore.movies
  
  // 获取动画类型电影 (genre_id=14)
  const animationData = await movieStore.fetchMovies({genre: 14, sort_by: 'vote_average.desc', page: 1, limit: 20 })
  animationMovies.value = animationData.results.sort((a, b) => b.vote_average - a.vote_average)
  
  // 获取动作类型电影 (genre_id=28)
  const actionData = await movieStore.fetchMovies({genre: 28, sort_by: 'vote_average.desc', page: 1, limit: 20 })
  actionMovies.value = actionData.results.sort((a, b) => b.vote_average - a.vote_average)
  
  // 获取评分最高的20部电影
  const topRatedData = await movieStore.fetchMovies({ sort_by: 'vote_average.desc', page: 1, limit: 20 })
  topRatedMovies.value = topRatedData.results.sort((a, b) => b.vote_average - a.vote_average)
  
  // 获取趋势电影 (按人气排序)
  const trendingData = await movieStore.fetchMovies({ sort_by: 'popularity.desc', page: 1, limit: 20 })
  trendingMovies.value = trendingData.results.sort((a, b) => b.popularity - a.popularity)
})

const handleDetail = (movie) => {
  router.push(`/movies/${movie.id}`)
}

const handleSearch = async () => {
  if (searchQuery.value.trim()) {
    await movieStore.searchMovies(searchQuery.value)
    router.push({ path: '/search', query: { q: searchQuery.value } })
    searchQuery.value = ''
  }
}
</script>

<style scoped>
/* 新增banner样式 */
.banner {
  height: 400px;
  background-image: url('/src/assets/banner.jpg');
  background-size: 100% auto;
  background-position: center;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  text-align: center;
  position: relative;
  margin-bottom: 2rem;
  width: 100%;
}

.banner::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(200, 200, 201, 0.502);
}

.banner-content {
  position: relative;
  z-index: 1;
  max-width: 800px;
  padding: 0 20px;
}

.banner h1 {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.banner p {
  font-size: 1.5rem;
  margin-bottom: 2rem;
}

.banner-search {
  width: 80%;
  max-width: 600px;
  margin: 0 auto;
}

/* 原有样式保持不变 */
.home-container {
  padding: 0 2rem;
}

.movie-grid {
  display: flex;
  overflow-x: auto;
  gap: 20px;
  margin-bottom: 2rem;
  padding: 10px 0;
  scrollbar-width: none;
}

.movie-grid::-webkit-scrollbar {
  display: none;
}

.movie-card {
  flex: 0 0 auto;
  width: 180px;
}
</style>