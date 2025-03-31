<template>
  <div class="search-results-container">
    <el-card class="search-results-card">
      <template #header>
        <h2>搜索结果</h2>
        <div class="search-query">搜索关键词: "{{ searchQuery }}"</div>
      </template>

      <el-table :data="movieStore.searchResults" style="width: 100%" v-loading="loading">
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
        :total="totalPages * 10"
        :current-page="currentPage"
        @current-change="handlePageChange"
        style="margin-top: 20px; justify-content: center;"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMovieStore } from '../../stores/movie'

const route = useRoute()
const router = useRouter()
const movieStore = useMovieStore()

const searchQuery = ref(route.query.q || '')
const loading = ref(false)
const currentPage = ref(1)
const totalPages = computed(() => movieStore.searchTotalPages)

onMounted(async () => {
  if (searchQuery.value) {
    loading.value = true
    await movieStore.searchMovies(searchQuery.value, currentPage.value)
    loading.value = false
  }
})

const handlePageChange = async (page) => {
  currentPage.value = page
  loading.value = true
  await movieStore.searchMovies(searchQuery.value, page)
  loading.value = false
}

const handleDetail = (movie) => {
  router.push(`/movies/${movie.id}`)
}
</script>

<style scoped>
.search-results-container {
  padding: 20px;
}

.search-results-card {
  max-width: 1200px;
  margin: 0 auto;
}

.search-query {
  margin-top: 10px;
  color: #666;
  font-size: 14px;
}
</style>