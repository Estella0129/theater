import { defineStore } from 'pinia'
import axios from 'axios'

export const useUserStore = defineStore('user', {
  state: () => {
    // 从localStorage初始化用户状态
    const token = localStorage.getItem('token')
    const user = JSON.parse(localStorage.getItem('user'))
    
    // 设置axios默认请求头
    if (token) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`
    }
    
    return {
      user,
      token
    }
  },

  getters: {
    isLoggedIn: (state) => !!state.user,
    currentUser: (state) => state.user
  },

  actions: {
    async register(userData) {
      try {
        const response = await axios.post('/api/v1/frontend/users/register', userData)
        return response.data
      } catch (error) {
        throw new Error(error.response?.data?.error || '注册失败')
      }
    },

    async login(credentials) {
      try {
        const response = await axios.post('/api/v1/frontend/users/login', credentials)
        this.user = response.data.user
        this.token = response.data.token
        // 设置请求头的认证信息
        if (this.token) {
          axios.defaults.headers.common['Authorization'] = `Bearer ${this.token}`
          // 设置cookie，确保保存完整的用户信息
          const userData = {
            id: this.user.id,
            username: this.user.username,
            name: this.user.name,
            email: this.user.email,
            role: this.user.role,
            gender: this.user.gender,
          }
          localStorage.setItem('token', this.token)
          localStorage.setItem('user', JSON.stringify(userData))
        }
        return response.data
      } catch (error) {
        throw new Error(error.response?.data?.error || '登录失败')
      }
    },

    logout() {
      this.user = null
      this.token = null
      delete axios.defaults.headers.common['Authorization']
      // 清除localStorage
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    },

    async updateUser(userId, userData) {
      try {
        const response = await axios.put(`/api/v1/admin/users/${userId}`, userData)
        return response.data
      } catch (error) {
        throw new Error(error.response?.data?.error || '更新用户信息失败')
      }
    },

    async fetchUser(userId) {
      try {
        const response = await axios.get(`/api/v1/frontend/users/${userId}`)
        return response.data
      } catch (error) {
        throw new Error(error.response?.data?.error || '获取用户信息失败')
      }
    },

    async fetchUsers(page = 1, pageSize = 20) {
      try {
        const response = await axios.get('/api/v1/admin/users', {
          params: { page, page_size: pageSize }
        })
        return response.data
      } catch (error) {
        throw new Error(error.response?.data?.error || '获取用户列表失败')
      }
    },

    async deleteUser(userId) {
      try {
        const response = await axios.delete(`/api/v1/admin/users/${userId}`)
        return response.data
      } catch (error) {
        throw new Error(error.response?.data?.error || '删除用户失败')
      }
    }
  }
})