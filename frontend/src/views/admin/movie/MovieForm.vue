<template>
  <el-dialog v-model="dialogVisible" :title="dialogTitle" width="70%">
    <el-form :model="form" label-width="120px" ref="formRef">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本信息" name="basic">
          <el-form-item label="电影名称" prop="title" required>
            <el-input v-model="form.title" />
          </el-form-item>
          <el-form-item label="原名" prop="original_title">
            <el-input v-model="form.original_title" />
          </el-form-item>
          <el-form-item label="语言" prop="original_language">
            <el-input v-model="form.original_language" />
          </el-form-item>
          <el-form-item label="简介" prop="overview">
            <el-input v-model="form.overview" type="textarea" rows="4" />
          </el-form-item>
          <el-form-item label="上映日期" prop="release_date">
            <el-date-picker v-model="form.release_date" type="date" />
          </el-form-item>
          <el-form-item label="时长(分钟)" prop="runtime">
            <el-input-number v-model="form.runtime" :min="0" />
          </el-form-item>
          <el-form-item label="评分" prop="vote_average">
            <el-rate v-model="form.vote_average" :max="10" />
          </el-form-item>
        </el-tab-pane>

        <el-tab-pane label="演职人员" name="credits">
          <div class="credits-container">
            <el-button type="primary" @click="addCredit">添加人员</el-button>
            <el-table :data="form.Credits" style="width: 100%; margin-top: 20px">
              <el-table-column prop="People.name" label="姓名" />
              <el-table-column prop="department" label="部门" />
              <el-table-column prop="job" label="职位" />
              <el-table-column label="操作">
                <template #default="scope">
                  <el-button size="small" @click="editCredit(scope.row)">编辑</el-button>
                  <el-button size="small" type="danger" @click="removeCredit(scope.$index)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>

        <el-tab-pane label="图片" name="images">
          <div class="images-container">
            <el-table :data="form.Images" style="width: 100%; margin-top: 20px">
              <el-table-column label="预览">
                <template #default="{ row }">
                  <el-image :src="'/images'+row.file_path" fit="cover" style="width: 100px; height: 150px" />
                </template>
              </el-table-column>
              <el-table-column label="类型">
  <template #default="{ row, $index }">
    <el-select v-model="row.type" placeholder="选择类型" style="width: 100%">
      <el-option label="海报" value="poster" />
      <el-option label="背景" value="backdrop" />
      <el-option label="剧照" value="still" />
    </el-select>
  </template>
</el-table-column>
              <el-table-column prop="aspect_ratio" label="宽高比" />
              <el-table-column prop="width" label="宽度" />
              <el-table-column prop="height" label="高度" />
              <el-table-column label="操作">
                <template #default="{ row, $index }">
                  <el-button size="small" type="success" @click="handleAddImage($index)">新增</el-button>
                  <el-upload
                    action="#"
                    :show-file-list="false"
                    :auto-upload="false"
                    :on-change="(file) => handleRowImageChange(file, $index)"
                  >
                    <el-button size="small">上传</el-button>
                  </el-upload>
                  <el-button size="small" type="danger" @click="removeImage($index)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>

        <el-tab-pane label="类型" name="genres">
          <div class="genres-container">
            <el-select
              v-model="form.Genres"
              multiple
              filterable
              allow-create
              default-first-option
              placeholder="选择或输入类型"
              :reserve-keyword="false"
              :value-key="'id'"
            >
              <el-option
                v-for="genre in allGenres"
                :key="genre.id"
                :label="genre.name"
                :value="genre"
              />
            </el-select>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-form>

    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="submitForm">确认</el-button>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import axios from 'axios'
import { ElDialog, ElForm, ElFormItem, ElInput, ElButton, ElMessage, ElSelect, ElOption } from 'element-plus'
import { useMovieStore } from '@/stores/movie.js' // Import the store

const dialogVisible = ref(false)
const dialogTitle = ref('')
const activeTab = ref('basic')
const formRef = ref(null)

const form = reactive({
  title: '',
  original_title: '',
  original_language: '',
  overview: '',
  release_date: '',
  runtime: 0,
  vote_average: 0,
  Credits: [],
  Images: [],
  Genres: [],
  poster_path: '',
  backdrop_path: ''
})

const allGenres = ref([])

const movieStore = useMovieStore() // Instantiate the store
const genres = ref([])

onMounted(async () => {
  try {
    allGenres.value = await movieStore.fetchGenres() // 将获取的类型数据赋值给allGenres
  } catch (error) {
    ElMessage.error('加载类型失败: ' + error.message)
  }
})

const addCredit = () => {
  form.Credits.push({
    credit_id: '',
    credit_type: '',
    department: '',
    job: '',
    People: {
      id: '',
      name: '',
      profile_path: ''
    }
  })
}

const editCredit = (credit) => {
  // 实现编辑逻辑
}

const removeCredit = (index) => {
  form.Credits.splice(index, 1)
}

const imageType = ref('backdrop')

const handleRowImageChange = async (file, index) => {
  const formData = new FormData()
  formData.append('file', file.raw)

  try {
    const resp = await axios.post('/api/v1/admin/upload-image', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })

    // 添加响应拦截处理
    if (resp.status !== 200 || !resp.data?.file_path) {
      throw new Error(resp.data?.message || '上传失败')
    }
    
    form.Images[index] = {
      ...form.Images[index],
      ...resp.data
    }
  } catch (error) {
    ElMessage.error('上传失败: ' + error.response?.data?.error || error.message)
  }
}

const setAsPoster = (index) => {
  form.poster_path = form.Images[index].file_path
}

const removeImage = (index) => {
  form.Images.splice(index, 1)
}

const submitForm = async () => {
  try {
    await formRef.value.validate()
    
    // 准备提交数据
    const submitData = {
      ...form,
      Genres: form.Genres,
      release_date: form.release_date ? new Date(form.release_date).toISOString() : '',
      Credits: form.Credits.map(credit => ({
        ...credit,
        People: credit.People.id ? { id: credit.People.id } : credit.People
      })),
      Images: form.Images.map(image => ({
        ...image,
        file_path: image.file_path.startsWith('blob:') ? '' : image.file_path
      }))
    }

    delete submitData.director;
    
    // 根据是否有ID决定调用创建或更新
    if (form.id) {
      await movieStore.updateMovie(submitData)
      ElMessage.success('电影更新成功')
    } else {
      await movieStore.createMovie(submitData)
      ElMessage.success('电影创建成功')
    }
    
    dialogVisible.value = false
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
  }
}

const open = (movie) => {
  dialogVisible.value = true
  dialogTitle.value = movie ? '编辑电影' : '添加电影'
  if (movie) {
    Object.assign(form, JSON.parse(JSON.stringify(movie)))
    // 类型字段转换为对象数组
    if (Array.isArray(form.Genres) && typeof form.Genres[0] !== 'object') {
      form.Genres = form.Genres.map(id => allGenres.value.find(g => g.id === id) || { id, name: String(id) })
    }
  } else {
    resetForm()
  }
}

const resetForm = () => {
  Object.assign(form, {
    title: '',
    original_title: '',
    original_language: '',
    overview: '',
    release_date: '',
    runtime: 0,
    vote_average: 0,
    Credits: [],
    Images: [],
    Genres: [],
    poster_path: '',
    backdrop_path: ''
  })
}

defineExpose({
  open
})

// 在现有方法后添加新方法
const handleAddImage = (index) => {
  form.Images.splice(index + 1, 0, {
    type: 'poster',
    file_path: '',
    aspect_ratio: 0,
    height: 0,
    width: 0
  })
}
</script>

<style scoped>
.credits-container,
.images-container,
.genres-container {
  padding: 20px;
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.image-item {
  position: relative;
  height: 300px;
}

.image-actions {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  justify-content: center;
  padding: 10px;
  background: rgba(0, 0, 0, 0.5);
}
</style>