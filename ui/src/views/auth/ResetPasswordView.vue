<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 px-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-12 h-12 rounded-xl bg-primary mb-4">
          <span class="text-white font-bold text-xl">V</span>
        </div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Reset your password</h1>
        <p class="text-gray-500 dark:text-gray-400 mt-1">Enter a new strong password</p>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-8">
        <div v-if="success" class="text-center py-4">
          <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-green-100 dark:bg-green-900/30 mb-4">
            <svg class="w-8 h-8 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
            </svg>
          </div>
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">Password updated!</h3>
          <p class="text-sm text-gray-500 dark:text-gray-400">Your password has been reset successfully.</p>
          <RouterLink to="/login" class="inline-block mt-6 text-primary font-medium hover:underline text-sm">
            Sign in with new password
          </RouterLink>
        </div>

        <form v-else @submit.prevent="handleReset" class="space-y-5">
          <div v-if="error" class="bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 rounded-lg p-3 text-sm">
            {{ error }}
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">New password</label>
            <input
              v-model="form.password"
              type="password"
              required
              minlength="8"
              placeholder="Min. 8 characters"
              class="w-full px-3.5 py-2.5 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition"
            />
          </div>

          <button
            type="submit"
            :disabled="loading || !token"
            class="w-full py-2.5 px-4 rounded-lg bg-primary text-white font-semibold hover:bg-primary/90 disabled:opacity-60 disabled:cursor-not-allowed transition"
          >
            <span v-if="loading">Resetting…</span>
            <span v-else>Reset password</span>
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { authApi } from '@/api/auth'

const route   = useRoute()
const token   = ref('')
const loading = ref(false)
const error   = ref('')
const success = ref(false)
const form    = reactive({ password: '' })

onMounted(() => {
  token.value = (route.query.token as string) ?? ''
})

async function handleReset() {
  if (!token.value) {
    error.value = 'Invalid or missing reset token.'
    return
  }
  error.value = ''
  loading.value = true
  try {
    await authApi.resetPassword({ token: token.value, password: form.password })
    success.value = true
  } catch (e: any) {
    error.value = e.response?.data?.message ?? 'Reset failed. The link may have expired.'
  } finally {
    loading.value = false
  }
}
</script>
