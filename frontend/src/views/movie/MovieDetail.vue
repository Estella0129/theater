<template>
  <div class="movie-detail-container">
    <el-card class="movie-card">
      <template #header>
        <div class="header-container">
          <h2>{{ movie.title }}</h2>
          <el-button type="info" @click="toggleFavorite" class="favorite-btn">
            <el-icon :color="isFavorite ? '#e6a23c' : '' " :size=" isFavorite?'26px':''"><StarFilled v-if="isFavorite" /><Star v-else /></el-icon>
          </el-button>
        </div>
      </template>

      <div class="movie-info">
        <el-image :src="getProfileImage(movie.poster_path)" fit="cover" class="movie-poster"></el-image>

        <div class="movie-meta">
          <p><strong>导演:</strong> {{ movie.director }}</p>
          <p><strong>主演:</strong> {{ movie.cast }}</p>
          <p><strong>时长:</strong> {{ movie.runtime >= 60 ? Math.floor(movie.runtime / 60) + '小时' + (movie.runtime % 60)
            + '分钟'
            : movie.runtime + '分钟' }}</p>
          <p><strong>评分:</strong> <el-rate v-model="movie.rating" disabled></el-rate></p>
          <p><strong>上映日期:</strong> {{ movie.release_date ? movie.release_date.split('T')[0] : '' }}</p>
          <p><strong>简介:</strong> {{ movie.overview }}</p>
        </div>
      </div>

      <div class="movie-actions">
        <el-button type="primary" @click="goBack">返回列表</el-button>
        <el-button type="info" @click="showStaffDialog = true">查看工作人员</el-button>
        
      </div>

      <div class="cast-list">
        <h3>演员阵容</h3>
        <div class="cast-grid">
          <div class="cast-item" v-for="actor in movie.Credits?.filter(c => c.credit_type === 'cast') || []" :key="actor.credit_id">
            <el-image :src="getProfileImage(actor.People?.profile_path)" fit="cover" class="cast-avatar"></el-image>
            <p class="cast-name">{{ actor.People?.name }}</p>
          </div>
        </div>
      </div>

      <el-dialog v-model="showStaffDialog" title="工作人员" width="50%">
        <div class="staff-grid">
          <div class="staff-item" v-for="person in movie.Credits?.filter(c => c.credit_type === 'crew') || []" :key="person.id">
            <el-image :src="getProfileImage(person.People.profile_path)" fit="cover" class="staff-avatar"></el-image>
            <p class="staff-name">{{ person.People.name }}</p>
            <p class="staff-job">{{ person.job }}</p>
          </div>
        </div>
      </el-dialog>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMovieStore } from '../../stores/movie'
import { Star, StarFilled } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const movieStore = useMovieStore()

const movie = ref({})
const showStaffDialog = ref(false)
const isFavorite = ref(false)

const getProfileImage = (path) => {
      return path 
        ? `https://image.tmdb.org/t/p/w200${path}` 
        : 'https://via.placeholder.com/200x300?text=No+Image';
    };

onMounted(async () => {
  const movieId = route.params.id
  movie.value = await movieStore.getMovieById(movieId)
  // 检查是否已收藏
  checkFavoriteStatus()
})

const checkFavoriteStatus = async () => {
  // 这里需要调用API检查当前用户是否已收藏该电影
  // 暂时模拟已收藏状态
  isFavorite.value = false
}

const toggleFavorite = async () => {
  try {
    if (isFavorite.value) {
      // 调用取消收藏API
    } else {
      // 调用添加收藏API
    }
    isFavorite.value = !isFavorite.value
  } catch (error) {
    console.error('收藏操作失败:', error)
  }
}

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

.header-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.favorite-btn {
  margin-left: 10px;
  font-size: 24px;
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

.cast-list {
  margin-top: 30px;
}

.cast-grid {
  display: flex;
  overflow-x: auto;
  gap: 20px;
  margin-top: 15px;
  padding-bottom: 10px;
}

.cast-item {
  text-align: center;
}

.cast-avatar {
  width: 100px;
  height: 150px;
  border-radius: 4px;
}

.cast-name {
  margin-top: 8px;
  font-size: 14px;
}

.staff-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 20px;
  margin-top: 15px;
}

.staff-item {
  text-align: center;
}

.staff-avatar {
  width: 100px;
  height: 150px;
  border-radius: 4px;
}

.staff-name {
  margin-top: 8px;
  font-size: 14px;
}

.staff-job {
  margin-top: 4px;
  font-size: 12px;
  color: #666;
}
</style>