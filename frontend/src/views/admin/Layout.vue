<template>
  <div class="admin-layout">
    <el-container>
      <el-aside width="200px">
        <div class="logo">
          <h2>后台管理</h2>
        </div>
        <el-menu
          :router="true"
          :default-active="$route.path"
          class="el-menu-vertical"
        >
          <el-menu-item index="/admin/users">
            <el-icon><user /></el-icon>
            <span>用户管理</span>
          </el-menu-item>
          <el-menu-item index="/admin/movies">
            <el-icon><film /></el-icon>
            <span>电影管理</span>
          </el-menu-item>
           <el-menu-item index="/admin/people">
            <el-icon><UserFilled /></el-icon>
            <span>人物管理</span>
          </el-menu-item>
          <el-menu-item index="/admin/genre">
            <el-icon><Reading /></el-icon>
            <span>类型管理</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-container>
        <el-header>
          <div class="header-right">
            <el-dropdown>
              <span class="el-dropdown-link">
                {{ userStore.currentUser?.username }}
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleLogout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>
        <el-main>
          <router-view></router-view>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '../../stores/user'
import { User, Film, ArrowDown, UserFilled, Reading} from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()
4

const handleLogout = () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style scoped>
.admin-layout {
  min-height: 100vh;
}

.logo {
  height: 60px;
  display: flex;
  justify-content: center;
  align-items: center;
  color: #fff;
  background-color: #304156;
}

.logo h2 {
  margin: 0;
  font-size: 18px;
}

.el-aside {
  background-color: #304156;
  color: #fff;
}

.el-menu {
  border-right: none;
  background-color: #304156;
}

.el-menu-item {
  color: #bfcbd9;
}

.el-menu-item:hover {
  color: #fff;
  background-color: #263445;
}

.el-menu-item.is-active {
  color: #409eff;
  background-color: #263445;
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 0 20px;
}

.header-right {
  display: flex;
  align-items: center;
}

.el-dropdown-link {
  cursor: pointer;
  display: flex;
  align-items: center;
}

.el-main {
  background-color: #f0f2f5;
  padding: 20px;
}
</style>