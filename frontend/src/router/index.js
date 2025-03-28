import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: () => import('../views/Home.vue')
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('../views/auth/Login.vue')
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('../views/auth/Register.vue')
    },
    {
      path: '/profile',
      name: 'Profile',
      component: () => import('../views/user/Profile.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/movies',
      name: 'MovieList',
      component: () => import('../views/movie/MovieList.vue')
    },
    {
      path: '/movies/:id',
      name: 'MovieDetail',
      component: () => import('../views/movie/MovieDetail.vue')
    },
    {
      path: '/admin',
      name: 'Admin',
      component: () => import('../views/admin/Layout.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: 'users',
          name: 'AdminUsers',
          component: () => import('../views/admin/user/UserList.vue')
        },
        {
          path: 'movies',
          name: 'AdminMovies',
          component: () => import('../views/admin/movie/MovieList.vue')
        }
      ]
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
const userStore = useUserStore()
  
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next('/login')
  } else if (to.meta.requiresAdmin && userStore.currentUser?.role !== 'admin') {
    next('/')
  } else {
    next()
  }
})

export default router