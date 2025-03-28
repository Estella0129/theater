<template>
  <div class="app">
    <el-container>
      <el-header>
        <nav>
          <div class="nav-left">
            <router-link to="/">首页</router-link>
            <router-link to="/movies">电影</router-link>
          </div>
          <div class="nav-right">
            <template v-if="!isLoggedIn">
              <router-link to="/login" class="login-btn">登录</router-link>
              <router-link to="/register" class="register-btn">注册</router-link>
            </template>
            <template v-else>
              <router-link to="/profile">个人中心</router-link>
              <a @click.prevent="handleLogout" href="#" class="logout-btn">退出</a>
            </template>
          </div>
        </nav>
      </el-header>
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from './stores/user'

const router = useRouter()
const userStore = useUserStore()

const isLoggedIn = computed(() => userStore.isLoggedIn)

const handleLogout = async () => {
  userStore.logout()
  router.push('/login')
}
</script>

<style>
.app {
  font-family: Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

.el-header {
  background-color: #fff;
  border-bottom: 1px solid #eee;
  padding: 1rem;
}

nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.nav-left {
  display: flex;
  gap: 20px;
}

.nav-right {
  display: flex;
  gap: 15px;
}

nav a {
  font-weight: bold;
  color: #2c3e50;
  text-decoration: none;
}

nav a.router-link-exact-active {
  color: #42b983;
}

.login-btn,
.register-btn {
  padding: 8px 16px;
  border-radius: 4px;
}

.login-btn {
  border: 1px solid #409eff;
  color: #409eff;
}

.register-btn {
  background-color: #409eff;
  color: white;
}

.logout-btn {
  color: #f56c6c;
}

.el-main {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}
</style>