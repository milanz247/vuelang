<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 px-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-12 h-12 rounded-xl bg-primary mb-4">
          <span class="text-white font-bold text-xl">V</span>
        </div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Forgot password?</h1>
        <p class="text-gray-500 dark:text-gray-400 mt-1">We'll send a reset link to your email</p>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-8">
        <div v-if="sent" class="text-center py-4">
          <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-green-100 dark:bg-green-900/30 mb-4">
            <svg class="w-8 h-8 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
            </svg>
          </div>
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">Check your email</h3>
          <p class="text-sm text-gray-500 dark:text-gray-400">
            If an account with that email exists, we've sent a password reset link.
          </p>
          <RouterLink to="/login" class="inline-block mt-6 text-primary font-medium hover:underline text-sm">
            Back to sign in
          </RouterLink>
        </div>

        <form v-else @submit.prevent="handleSubmit" class="space-y-5">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Email address</label>
            <input
              v-model="email"
              type="email"
              required
              placeholder="you@example.com"
              class="w-full px-3.5 py-2.5 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition"
            />
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full py-2.5 px-4 rounded-lg bg-primary text-white font-semibold hover:bg-primary/90 disabled:opacity-60 disabled:cursor-not-allowed transition"
          >
            <span v-if="loading">Sending…</span>
            <span v-else>Send reset link</span>
          </button>

          <p class="text-center text-sm text-gray-500 dark:text-gray-400">
            <RouterLink to="/login" class="text-primary font-medium hover:underline">Back to sign in</RouterLink>
          </p>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { authApi } from '@/api/auth'

const email   = ref('')
const loading = ref(false)
const sent    = ref(false)

async function handleSubmit() {
  loading.value = true
  try {
    await authApi.forgotPassword({ email: email.value })
    sent.value = true
  } finally {
    loading.value = false
  }
}
</script>
