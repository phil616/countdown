import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/Login.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/oauth/callback',
      name: 'oauth-callback',
      component: () => import('@/views/OAuthCallback.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      component: () => import('@/views/Layout.vue'),
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          redirect: '/dashboard',
        },
        {
          path: 'dashboard',
          name: 'dashboard',
          component: () => import('@/views/Dashboard.vue'),
        },
        {
          path: 'units',
          name: 'units',
          component: () => import('@/views/Units/UnitList.vue'),
        },
        {
          path: 'units/:id',
          name: 'unit-detail',
          component: () => import('@/views/Units/UnitDetail.vue'),
        },
        {
          path: 'projects',
          name: 'projects',
          component: () => import('@/views/Projects/ProjectList.vue'),
        },
        {
          path: 'projects/:id',
          name: 'project-detail',
          component: () => import('@/views/Projects/ProjectDetail.vue'),
        },
        {
          path: 'todos',
          name: 'todos',
          component: () => import('@/views/Todos/TodoList.vue'),
        },
        {
          path: 'notifications',
          name: 'notifications',
          component: () => import('@/views/Notifications.vue'),
        },
        {
          path: 'settings',
          name: 'settings',
          component: () => import('@/views/Settings/SettingsPage.vue'),
        },
      ],
    },
  ],
})

router.beforeEach((to) => {
  const auth = useAuthStore()
  if (to.meta.requiresAuth !== false && !auth.isLoggedIn) {
    return '/login'
  }
  if (to.path === '/login' && auth.isLoggedIn) {
    return '/dashboard'
  }
})

export default router
