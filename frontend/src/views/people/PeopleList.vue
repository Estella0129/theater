<template>
  <div class="People-list-container">
    <el-row :gutter="20">
      <el-col 
        v-for="People in displayedPeoples" 
        :key="People.id" 
        :xs="12" :sm="8" :md="6" :lg="6" :xl="6"
      >
        <el-card class="People-card" shadow="hover" @click="handleDetail(People)">
          <div class="People-poster">
            <img :src="getProfileImage(People.profile_path)" alt="People poster" />
          </div>
          <div class="People-info">
            <h3>{{ People.name }}</h3>
            <p v-if="People.original_name !== People.name">{{ People.original_name }}</p>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-pagination
      background
      layout="prev, pager, next"
      :total="totalPeoples"
      :page-size="pageSize"
      :current-page="currentPage"
      @current-change="handlePageChange"
      style="margin-top: 20px; justify-content: center;"
    />
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue';
import { usePeopleStore } from '../../stores/people';
import { useRouter } from 'vue-router';

export default {
  setup() {
    const PeopleStore = usePeopleStore();
    const router = useRouter();
    const currentPage = ref(1);
    const pageSize = 20;

    const totalPeoples = computed(() => PeopleStore.Peoples.length);
    
    const displayedPeoples = computed(() => {
      const start = (currentPage.value - 1) * pageSize;
      const end = start + pageSize;
      return PeopleStore.Peoples.slice(start, end);
    });

    const handleDetail = (People) => {
      router.push(`/people/${People.id}`)
    }

    const handlePageChange = async (page) => {
      currentPage.value = page;
      await PeopleStore.fetchPeoples(page, pageSize);
    };

    const getProfileImage = (path) => {
      return path 
        ? `https://image.tmdb.org/t/p/w200${path}` 
        : 'https://via.placeholder.com/200x300?text=No+Image';
    };

    onMounted(async () => {
      await PeopleStore.fetchPeoples();
    });

    return {
      currentPage,
      pageSize,
      totalPeoples,
      displayedPeoples,
      handlePageChange,
      getProfileImage,
      handleDetail
    };
  }
};
</script>

<style scoped>
.People-list-container {
  padding: 20px;
}

.People-card {
  margin-bottom: 30px;
  cursor: pointer;
  transition: transform 0.3s;
}

.People-card:hover {
  transform: translateY(-5px);
}

.People-poster img {
  width: 100%;
  height: auto;
  border-radius: 4px 4px 0 0;
}

.People-info {
  padding: 10px;
}

.People-info h3 {
  margin: 0;
  font-size: 16px;
  text-align: center;
}

.People-info p {
  margin: 5px 0 0;
  font-size: 14px;
  color: #666;
  text-align: center;
}
</style>