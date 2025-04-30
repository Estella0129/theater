<template>
  <div class="people-list">
    <el-card>
      <template #header>
      <div class="header">
        <el-input
          v-model="searchQuery"
          placeholder="搜索姓名"
          style="width: 300px"
          clearable
          @clear="handleSearch"
          @keyup.enter="handleSearch"
        />
        <el-button type="primary" @click="handleAdd">添加人员</el-button>
      </div>
      </template>

      <el-table :data="filteredPeoples" style="width: 100%" border>
        <el-table-column prop="name" label="姓名" width="180" />
        <el-table-column prop="original_name" label="原名" width="180" />
        <el-table-column prop="role" label="角色" width="160">
          <template #default="scope">
            {{ formatRole(scope.row.role) }}
          </template>
        </el-table-column>
        <el-table-column prop="gender" label="性别" width="120">
          <template #default="scope">
            {{ formatGender(scope.row.gender) }}
          </template>
        </el-table-column>
        <el-table-column prop="birthday" label="出生日期" width="120">
          <template #default="scope">
            {{ formatDate(scope.row.birthday) }}
          </template>
        </el-table-column>
        <el-table-column prop="place_of_birth" label="出生地" width="200" />
        <el-table-column label="操作" width="240">
          <template #default="scope">
            <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model="currentPage"
        :page-size="pageSize"
        layout=" prev, pager, next"
        :total="total"
        @current-change="handlePageChange"
      />
    </el-card>

    <PeopleForm ref="peopleFormRef" @refresh="fetchData" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { usePeopleStore } from '../../../stores/people';
import PeopleForm from './PeopleForm.vue';
import { ElMessage } from 'element-plus';

const router = useRouter();

const peopleStore = usePeopleStore();
const peopleFormRef = ref(null);

const currentPage = ref(1);
const pageSize = 20;
const total = ref(0);
const searchQuery = ref('');

onMounted(() => {
  // 从URL参数初始化搜索条件
  const route = router.currentRoute.value;
  if (route.query.search) {
    searchQuery.value = route.query.search;
  }
  if (route.query.page) {
    currentPage.value = parseInt(route.query.page);
  }
  fetchData();
});

const fetchData = async () => {
  try {
    const data = await peopleStore.fetchPeoples(currentPage.value, pageSize, searchQuery.value);
    total.value = data.total;
    // 更新URL参数
    router.push({
      query: {
        ...router.currentRoute.value.query,
        page: currentPage.value,
        search: searchQuery.value
      }
    });
  } catch (error) {
    console.error('获取人员列表失败:', error);
  }
};

const filteredPeoples = computed(() => {
  return peopleStore.Peoples;
});

const handleSearch = () => {
  currentPage.value = 1;
  fetchData();
};

const handleAdd = () => {
  peopleFormRef.value.open('添加人员');
};

const handleEdit = (row) => {
  peopleFormRef.value.open('编辑人员', true, row);
};

const handlePageChange = (page) => {
  currentPage.value = page;
  fetchData();
};

const handleDelete = async (row) => {
  try {
    await peopleStore.deletePeople(row.id);
    await fetchData();
    ElMessage.success('删除成功');
  } catch (error) {
    console.error('删除失败:', error);
  }
};

const formatRole = (role) => {
  const roles = {
    actor: '演员',
    director: '导演',
    writer: '编剧'
  };
  return roles[role] || role;
};

const formatDate = (dateString) => {
  if (!dateString) return '';
  return new Date(dateString).toISOString().split('T')[0];
};

const formatGender = (gender) => {
  const genders = {
    0: '未知',
    1: '女',
    2: '男'
  };
  return genders[gender] || '未知';
};


</script>

<style scoped>
.people-list {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.el-table {
  margin-top: 20px;
}

.el-pagination {
  margin-top: 20px;
  text-align: center;
}
</style>