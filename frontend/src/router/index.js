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
      path: '/search',
      name: 'SearchResults',
      component: () => import('../views/movie/SearchResults.vue')
    },
    {
      path: '/People',
      name: 'People',
      component: () => import('../views/people/PeopleList.vue')
    },
    {
      path: '/People/:id',
      name: 'PeopleDetail',
      component: () => import('../views/People/PeopleDetail.vue')
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
        },
        {
          path: 'people',
          name: 'AdminPeople',
          component: () => import('../views/admin/people/PeopleList.vue')
        }
        ,
        {
          path: 'genre',
          name: 'AdminGenre',
          component: () => import('../views/admin/genre/GenreList.vue')
        }
      ]
    }
  ]
})

// 路由守卫
router.beforeEach((to, from, next) => {
const userStore = useUserStore()
  
  // 未登录用户只能访问首页、登录和注册页面
  const allowedRoutes = ['Home', 'Login', 'Register']
  if (!userStore.isLoggedIn && !allowedRoutes.includes(to.name)) {
    next('/')
    return
  }
  
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next('/login')
  } else if (to.meta.requiresAdmin && userStore.currentUser?.role !== 'admin') {
    next('/')
  } else {
    next()
  }
})

export default router