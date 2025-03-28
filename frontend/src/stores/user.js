import { defineStore } from 'pinia'
import axios from 'axios'

export const useUserStore = defineStore('user', {
  state: () => ({
    user: null,
    token: null
  }),

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