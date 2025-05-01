<template>
  <div class="People-detail">
    <div class="People-header">
      <div class="People-poster">
        <img :src="getProfileImage(People.profile_path)" alt="People poster" />
      </div>
      <div class="People-info">
        <h1 v-if="People.name">{{ People.name }}</h1>
        <h2 v-if="People.original_name && People.original_name !== People.name">{{ People.original_name }}</h2>
        <div class="People-meta" v-if="People.gender || People.birthday || People.place_of_birth">
          <span v-if="People.gender === 1">女</span>
          <span v-else-if="People.gender === 2">男</span>
          <span v-if="People.birthday">{{ formatDate(People.birthday) }}</span>
          <span v-if="People.place_of_birth">{{ People.place_of_birth }}</span>
        </div>
        <p class="People-bio" v-if="People.biography">{{ People.biography }}</p>
      </div>
    </div>

    <div class="People-content">
      <section class="People-credits">
        <h3>参演作品</h3>
        <div class="credit-list">
          <div v-for="credit in credits" :key="credit.credit_id" class="credit-item">
            <router-link :to="`/movie/${credit.MovieID}`">
              <img :src="getPosterImage(credit.Movie.poster_path)" alt="Movie poster" />
              <div class="credit-info">
                <h4>{{ credit.Movie.title }}</h4>
                <p>{{ translateRole(credit.character || credit.job) }}</p>
                <p>{{ formatDate(credit.Movie.release_date) }}</p>
              </div>
            </router-link>
          </div>
        </div>
      </section>

      <section class="People-images" v-if="images.length> 0">
        <h3>图片</h3>
        <div class="image-gallery">
          <img 
            v-for="image in images" 
            :key="image.file_path" 
            :src="getProfileImage(image.file_path)" 
            alt="People image" 
            @click="openImage(image.file_path)"
          />
        </div>
      </section>
    </div>
    <div class="back-button">
      <el-button type="primary" @click="$router.push('/people')">返回人物列表</el-button>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'PeopleDetail',
  data() {
    return {
      People: {},
      credits: [],
      images: [],
      $http: axios
    }
  },
  methods: {
    getProfileImage(path) {
      return path ? `https://image.tmdb.org/t/p/w500${path}` : '/placeholder-profile.png'
    },
    getPosterImage(path) {
      return path ? `https://image.tmdb.org/t/p/w185${path}` : '/placeholder-poster.png'
    },
    formatDate(dateString) {
      if (!dateString) return '';
      return dateString.split('T')[0];
    },
    openImage(path) {
      // 打开大图预览
    },
    translateRole(role) {
      const translations = {
        'Director': '导演',
        'Writer': '编剧',
        'Producer': '制片人',
        'Actor': '演员',
        'Actress': '演员',
        'Cinematography': '摄影',
        'Editor': '剪辑',
        'Music': '音乐',
        'Sound': '音效',
        'Art': '美术',
        'Costume': '服装',
        'Makeup': '化妆',
        'Visual Effects': '视觉效果'
      };
      return translations[role] || role;
    }
  },
  async created() {
    const PeopleId = this.$route.params.id;
    try {
      // 获取人物详情
      const PeopleResponse = await axios.get(`/api/v1/frontend/peoples/${PeopleId}`);
      this.People = PeopleResponse.data;
      
      this.credits = PeopleResponse.data.Credits
      
    } catch (error) {
      console.error('获取数据失败:', error);
    }
  }
}
</script>

<style scoped>
.People-detail {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
  text-decoration: none
}

.People-header {
  display: flex;
  margin-bottom: 30px;
}

.People-poster {
  width: 300px;
  min-width: 300px;
  margin-right: 30px;
}

.People-poster img {
  width: 100%;
  border-radius: 5px;
}

.People-info h1 {
  font-size: 2.2rem;
  margin-bottom: 10px;
}

.People-info h2 {
  font-size: 1.5rem;
  color: #777;
  margin-bottom: 20px;
}

.People-meta {
  display: flex;
  gap: 15px;
  margin-bottom: 20px;
  color: #666;
}

.People-bio {
  line-height: 1.6;
  margin-bottom: 20px;
}

.People-content {
  margin-top: 40px;
  text-decoration: none;
}

.credit-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 25px;
  padding: 10px 0;
  text-decoration: none;
}

.credit-list a{
  text-decoration: none;
}

.credit-item {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  transition: all 0.3s ease;
}

.credit-item:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0,0,0,0.2);
}

.credit-item img {
  width: 100%;
  height: 300px;
  object-fit: cover;
}

.credit-info {
  padding: 15px;
  color: #333;
  text-decoration: none
}

.credit-info h4 {
  font-size: 1.1rem;
  margin-bottom: 8px;
}

.credit-info p {
  font-size: 0.9rem;
  color: #666;
  margin-bottom: 5px;
}

.image-gallery {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 15px;
}

.image-gallery img {
  width: 100%;
  cursor: pointer;
  transition: transform 0.2s;
}

.image-gallery img:hover {
  transform: scale(1.05);
}

.back-button {
  margin-top: 30px;
  text-align: center;
}
</style>