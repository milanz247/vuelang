<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 px-4">
    <div class="w-full max-w-md">
      <!-- Logo -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-12 h-12 rounded-xl bg-primary mb-4">
          <span class="text-white font-bold text-xl">V</span>
        </div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Welcome back</h1>
        <p class="text-gray-500 dark:text-gray-400 mt-1">Sign in to your Vuelang account</p>
      </div>

      <!-- Card -->
      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-8">
        <form @submit.prevent="handleLogin" class="space-y-5">
          <!-- Error alert -->
          <div v-if="error" class="bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 rounded-lg p-3 text-sm">
            {{ error }}
          </div>

          <!-- Email -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">
              Email address
            </label>
            <input
              v-model="form.email"
              type="email"
              autocomplete="email"
              required
              placeholder="you@example.com"
              class="w-full px-3.5 py-2.5 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition"
            />
          </div>

          <!-- Password -->
          <div>
            <div class="flex justify-between items-center mb-1.5">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                Password
              </label>
              <RouterLink to="/forgot-password" class="text-sm text-primary hover:underline">
                Forgot password?
              </RouterLink>
            </div>
            <input
              v-model="form.password"
              type="password"
              autocomplete="current-password"
              required
              placeholder="••••••••"
              class="w-full px-3.5 py-2.5 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition"
            />
          </div>

          <!-- Submit -->
          <button
            type="submit"
            :disabled="loading"
            class="w-full py-2.5 px-4 rounded-lg bg-primary text-white font-semibold hover:bg-primary/90 disabled:opacity-60 disabled:cursor-not-allowed transition"
          >
            <span v-if="loading">Signing in…</span>
            <span v-else>Sign in</span>
          </button>
        </form>

        <p class="mt-6 text-center text-sm text-gray-500 dark:text-gray-400">
          Don't have an account?
          <RouterLink to="/register" class="text-primary font-medium hover:underline">
            Sign up
          </RouterLink>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route  = useRoute()
const auth   = useAuthStore()

const loading = ref(false)
const error   = ref('')
const form    = reactive({ email: '', password: '' })

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(form.email, form.password)
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (e: any) {
    error.value = e.response?.data?.message ?? 'Login failed. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>
