<!-- File: vuelang/ui/src/views/auth/ResetPasswordView.vue -->
<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { authApi } from '@/api/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Alert, AlertDescription } from '@/components/ui/alert'

// Import the logo using Vite's path alias
import vuelangLogo from '@/assets/vuelang-logo.svg'

const route   = useRoute()
const token   = ref('')
const loading = ref(false)
const error   = ref('')
const success = ref(false)
const showPassword = ref(false)
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

<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-[#FDFDFC] dark:bg-[#0a0a0a] p-6 font-sans">
    <div class="w-full max-w-[350px] flex flex-col">
      
      <!-- Center Aligned Vuelang Logo -->
      <div class="flex justify-center mb-6">
        <RouterLink to="/">
          <img :src="vuelangLogo" alt="Vuelang Logo" class="w-12 h-12 object-contain" />
        </RouterLink>
      </div>

      <!-- Headers -->
      <div class="text-center mb-6 space-y-1.5">
        <h1 class="text-base font-bold text-gray-900 dark:text-white tracking-tight">
          Reset your password
        </h1>
        <p class="text-xs text-gray-500 dark:text-gray-400">
          Enter a new strong password
        </p>
      </div>

      <!-- Success State -->
      <div v-if="success" class="text-center py-4 space-y-4">
        <div class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-gray-50 dark:bg-[#161615] border border-gray-150 dark:border-[#3E3E3A] mb-2">
          <svg class="w-5 h-5 text-gray-600 dark:text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
          </svg>
        </div>
        <div class="space-y-1.5">
          <h3 class="text-sm font-semibold text-gray-900 dark:text-white">Password updated!</h3>
          <p class="text-xs text-gray-500 dark:text-gray-400 leading-relaxed">
            Your password has been reset successfully.
          </p>
        </div>
        <!-- Sign in with New Password Button -->
        <Button variant="outline" as-child class="w-full h-10 border-gray-200 dark:border-[#3E3E3A]">
          <RouterLink to="/login" class="text-xs font-semibold text-gray-700 dark:text-gray-200">
            Sign in with new password
          </RouterLink>
        </Button>
      </div>

      <!-- Form State -->
      <form v-else @submit.prevent="handleReset" class="space-y-4">
        <!-- Error Alerts -->
        <Alert v-if="error" variant="destructive" class="rounded-lg py-2.5">
          <AlertDescription class="text-xs">{{ error }}</AlertDescription>
        </Alert>

        <!-- New Password Input -->
        <div class="space-y-1.5">
          <Label for="password" class="text-xs font-semibold text-gray-800 dark:text-gray-200">New password</Label>
          <div class="relative">
            <Input
              id="password"
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              required
              minlength="8"
              placeholder="New password (Min. 8 characters)"
              class="h-10 px-3 pr-10 bg-white dark:bg-[#161615] border-gray-200 dark:border-[#3E3E3A] focus-visible:ring-1 focus-visible:ring-black"
            />
            <!-- Show/Hide Password Toggle Button -->
            <button
              type="button"
              @click="showPassword = !showPassword"
              class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 focus:outline-none dark:text-gray-500 dark:hover:text-gray-300"
            >
              <svg v-if="!showPassword" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
              </svg>
              <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 01-1.563-3.029m5.858.908a3 3 0 113.382 3.382m0.74 4.49a9.96 9.96 0 01-5.714-2.134m5.714 2.134l4.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Submit Button (using Shadcn Button) -->
        <Button
          type="submit"
          class="w-full h-10 bg-black text-white hover:bg-neutral-800 dark:bg-white dark:text-black dark:hover:bg-neutral-200 text-xs font-semibold rounded-lg shadow-sm transition duration-200"
          :disabled="loading || !token"
        >
          <span v-if="loading">Resetting…</span>
          <span v-else>Reset password</span>
        </Button>
      </form>

    </div>
  </div>
</template>