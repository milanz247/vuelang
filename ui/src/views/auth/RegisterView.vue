<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 px-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-12 h-12 rounded-xl bg-primary mb-4">
          <span class="text-white font-bold text-xl">V</span>
        </div>
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Create your account</h1>
        <p class="text-gray-500 dark:text-gray-400 mt-1">Start building with Vuelang</p>
      </div>

      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-lg p-8">
        <form @submit.prevent="handleRegister" class="space-y-5">
          <div v-if="error" class="bg-red-50 dark:bg-red-900/30 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 rounded-lg p-3 text-sm">
            {{ error }}
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Full name</label>
            <input
              v-model="form.name"
              type="text"
              autocomplete="name"
              required
              minlength="2"
              placeholder="Jane Doe"
              class="w-full px-3.5 py-2.5 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Email address</label>
            <input
              v-model="form.email"
              type="email"
              autocomplete="email"
              required
              placeholder="you@example.com"
              class="w-full px-3.5 py-2.5 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5">Password</label>
            <input
              v-model="form.password"
              type="password"
              autocomplete="new-password"
              required
              minlength="8"
              placeholder="Min. 8 characters"
              class="w-full px-3.5 py-2.5 rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition"
            />
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full py-2.5 px-4 rounded-lg bg-primary text-white font-semibold hover:bg-primary/90 disabled:opacity-60 disabled:cursor-not-allowed transition"
          >
            <span v-if="loading">Creating account…</span>
            <span v-else>Create account</span>
          </button>
        </form>

        <p class="mt-6 text-center text-sm text-gray-500 dark:text-gray-400">
          Already have an account?
          <RouterLink to="/login" class="text-primary font-medium hover:underline">Sign in</RouterLink>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth   = useAuthStore()

const loading = ref(false)
const error   = ref('')
const form    = reactive({ name: '', email: '', password: '' })

async function handleRegister() {
  error.value = ''
  loading.value = true
  try {
    await auth.register(form.name, form.email, form.password)
    router.push('/')
  } catch (e: any) {
    error.value = e.response?.data?.message ?? 'Registration failed. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>
