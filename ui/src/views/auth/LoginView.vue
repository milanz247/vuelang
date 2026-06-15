<!-- File: vuelang/ui/src/views/auth/LoginView.vue -->
<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Alert, AlertDescription } from '@/components/ui/alert'

// Import the logo using Vite's path alias
import vuelangLogo from '@/assets/vuelang-logo.svg'

const router = useRouter()
const route  = useRoute()
const auth   = useAuthStore()

const loading = ref(false)
const error   = ref('')
const showPassword = ref(false)

const form = reactive({ 
  email: '', 
  password: '',
  remember: false
})

async function handleLogin() {
  error.value = ''
  loading.value = true
  try {
    await auth.login(form.email, form.password)
    const redirect = (route.query.redirect as string) || '/dashboard'
    router.push(redirect)
  } catch (e: any) {
    error.value = e.response?.data?.message ?? 'Login failed. Please try again.'
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
          Log in to your account
        </h1>
        <p class="text-xs text-gray-500 dark:text-gray-400">
          Enter your email and password below to log in
        </p>
      </div>

      <!-- Login Form (Flat Layout matching Breeze) -->
      <form @submit.prevent="handleLogin" class="space-y-4">
        <!-- Error Alerts -->
        <Alert v-if="error" variant="destructive" class="rounded-lg py-2.5">
          <AlertDescription class="text-xs">{{ error }}</AlertDescription>
        </Alert>

        <!-- Email Address Input -->
        <div class="space-y-1.5">
          <Label for="email" class="text-xs font-semibold text-gray-800 dark:text-gray-200">Email address</Label>
          <Input
            id="email"
            v-model="form.email"
            type="email"
            autocomplete="email"
            required
            placeholder="email@example.com"
            class="h-10 px-3 bg-white dark:bg-[#161615] border-gray-200 dark:border-[#3E3E3A] focus-visible:ring-1 focus-visible:ring-black"
          />
        </div>

        <!-- Password Input with Toggle Icon -->
        <div class="space-y-1.5">
          <div class="flex justify-between items-center">
            <Label for="password" class="text-xs font-semibold text-gray-800 dark:text-gray-200">Password</Label>
            <RouterLink to="/forgot-password" class="text-xs text-gray-500 dark:text-gray-400 hover:underline">
              Forgot your password?
            </RouterLink>
          </div>
          <div class="relative">
            <Input
              id="password"
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              autocomplete="current-password"
              required
              placeholder="Password"
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

        <!-- Remember Me Checkbox -->
        <div class="flex items-center gap-3">
          <input
            id="remember_me"
            type="checkbox"
            v-model="form.remember"
            class="h-4 w-4 rounded border-gray-300 dark:border-[#3E3E3A] text-black focus:ring-black dark:bg-[#161615] accent-black cursor-pointer"
          />
          <Label for="remember_me" class="text-xs font-normal text-gray-500 dark:text-gray-400 cursor-pointer select-none">
            Remember me
          </Label>
        </div>

        <!-- Submit Button (using Shadcn Button) -->
        <Button
          type="submit"
          class="w-full h-10 bg-black text-white hover:bg-neutral-800 dark:bg-white dark:text-black dark:hover:bg-neutral-200 text-xs font-semibold rounded-lg shadow-sm transition duration-200"
          :disabled="loading"
        >
          <span v-if="loading">Logging in…</span>
          <span v-else>Log in</span>
        </Button>
      </form>

      <!-- Footer: Link to Registration -->
      <p class="mt-6 text-center text-xs text-gray-500 dark:text-gray-400">
        Don't have an account?
        <RouterLink to="/register" class="text-black dark:text-white font-medium hover:underline">
          Sign up
        </RouterLink>
      </p>

    </div>
  </div>
</template>