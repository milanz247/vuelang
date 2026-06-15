// File: vuelang/ui/src/router/index.ts
import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    // Public Welcome Page (No Authentication required)
    {
      path: '/',
      name: 'welcome',
      component: () => import('@/views/WelcomeView.vue'),
    },
    
    // Auth routes (Only accessible by guests)
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/auth/LoginView.vue'),
      meta: { guest: true },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/auth/RegisterView.vue'),
      meta: { guest: true },
    },
    {
      path: '/forgot-password',
      name: 'forgot-password',
      component: () => import('@/views/auth/ForgotPasswordView.vue'),
      meta: { guest: true },
    },
    {
      path: '/reset-password',
      name: 'reset-password',
      component: () => import('@/views/auth/ResetPasswordView.vue'),
      meta: { guest: true },
    },

    // Authenticated routes
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/views/DashboardView.vue'),
      meta: { auth: true },
    },

    // Wildcard fallback
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

router.beforeEach(async (to, _from, next) => {
  const auth = useAuthStore()

  // Eagerly load user profile on first navigation if authenticated
  if (auth.isAuthenticated && !auth.user) {
    try { 
      await auth.fetchMe() 
    } catch { 
      auth.clearSession() 
    }
  }

  // If page requires auth and user is not logged in, send to login
  if (to.meta.auth && !auth.isAuthenticated) {
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }

  // If user is logged in and tries to go to welcome/login/register, redirect to dashboard
  if (auth.isAuthenticated && (to.name === 'welcome' || to.meta.guest)) {
    next({ name: 'dashboard' })
    return
  }

  next()
})

export default router