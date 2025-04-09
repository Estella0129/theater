<template>
  <div class="movie-list-container">
    <el-card class="movie-list-card">
      <template #header>
        <h2>电影列表</h2>
        <div class="genre-tags">
        类型：<el-tag
        v-for="genre in genres"
        :key="genre.id"
        @click="handleGenreClick(genre)"
        type="info"
        class="genre-tag"
        >
        {{ genre.name }}
        </el-tag>
        </div>
      </template>
      

      <el-table :data="movieStore.movies" style="width: 100%">
        <el-table-column prop="title" label="电影名称" width="180" />
        <el-table-column prop="original_title" label="原名" width="180" />
        <el-table-column prop="release_date" label="上映日期">
          <template #default="scope">
            {{ new Date(scope.row.release_date).toLocaleDateString() }}
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button size="small" @click="handleDetail(scope.row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <el-pagination
        background
        layout="prev, pager, next"
        :total="totalPages * 20"
        :current-page="currentPage"
        :page-size="20"
        @current-change="handlePageChange"
        style="margin-top: 20px; justify-content: center;"
      />
      
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
const movies = ref([])
const loading = ref(false)
const currentPage = ref(1)
const totalPages = ref(1)
const queryGenre = ref(1)
const genres = ref([])
const pageSize = ref(20)

const loadMovies = async () => {
  loading.value = true
  try {
    const query = {
      page: currentPage.value,
      pageSize: pageSize.value,
      genre: queryGenre.value
    }
    const data = await movieStore.fetchMovies(query)
    totalPages.value = data.total_pages
    currentPage.value = data.page
    movies.value = data.results
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  currentPage.value = parseInt(route.query.page) || 1
  pageSize.value = parseInt(route.query.pageSize) || 20
  queryGenre.value = route.query.genre || 1
  await loadMovies()
  await loadGenres()
})
const handleGenreClick = (genre) => {
  router.push(`/movies?genre=${genre.id}`)
  queryGenre.value = genre.id
  loadMovies()
}
const handleDetail = (movie) => {
  router.push(`/movies/${movie.id}`)
}

const handlePageChange = (Page) => {
  currentPage.value = Page
  loadMovies()
}

const loadGenres = async () => {
  try {
    genres.value = await movieStore.fetchGenres()
  } catch (error) {
    console.error('Failed to load genres:', error)
  }
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
/* .genre-tags {
  margin: 10px 0;
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
} */

.genre-tag {
  cursor: pointer;
  transition: all 0.3s;
  margin: 5px;
}

.genre-tag:hover {
  transform: translateY(-2px);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.genre-tags {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  margin: 10px 0;
}

.genre-tags::after {
  content: "";
  flex: auto;
}
</style>