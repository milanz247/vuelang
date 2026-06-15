<!-- File: vuelang/ui/src/views/auth/ForgotPasswordView.vue -->
<script setup lang="ts">
import { ref } from 'vue'
import { authApi } from '@/api/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'

// Import the logo using Vite's path alias
import vuelangLogo from '@/assets/vuelang-logo.svg'

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
          Forgot password?
        </h1>
        <p class="text-xs text-gray-500 dark:text-gray-400">
          We'll send a reset link to your email
        </p>
      </div>

      <!-- Success State -->
      <div v-if="sent" class="text-center py-4 space-y-4">
        <div class="inline-flex items-center justify-center w-12 h-12 rounded-full bg-gray-50 dark:bg-[#161615] border border-gray-150 dark:border-[#3E3E3A] mb-2">
          <svg class="w-5 h-5 text-gray-600 dark:text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
          </svg>
        </div>
        <div class="space-y-1.5">
          <h3 class="text-sm font-semibold text-gray-900 dark:text-white">Check your email</h3>
          <p class="text-xs text-gray-500 dark:text-gray-400 leading-relaxed">
            If an account with that email exists, we've sent a password reset link.
          </p>
        </div>
        <!-- Back to Sign in Button -->
        <Button variant="outline" as-child class="w-full h-10 border-gray-200 dark:border-[#3E3E3A]">
          <RouterLink to="/login" class="text-xs font-semibold text-gray-700 dark:text-gray-200">
            Back to sign in
          </RouterLink>
        </Button>
      </div>

      <!-- Form State -->
      <form v-else @submit.prevent="handleSubmit" class="space-y-4">
        <!-- Email Address Input -->
        <div class="space-y-1.5">
          <Label for="email" class="text-xs font-semibold text-gray-800 dark:text-gray-200">Email address</Label>
          <Input
            id="email"
            v-model="email"
            type="email"
            required
            placeholder="email@example.com"
            class="h-10 px-3 bg-white dark:bg-[#161615] border-gray-200 dark:border-[#3E3E3A] focus-visible:ring-1 focus-visible:ring-black"
          />
        </div>

        <!-- Submit Button (using Shadcn Button) -->
        <Button
          type="submit"
          class="w-full h-10 bg-black text-white hover:bg-neutral-800 dark:bg-white dark:text-black dark:hover:bg-neutral-200 text-xs font-semibold rounded-lg shadow-sm transition duration-200"
          :disabled="loading"
        >
          <span v-if="loading">Sending…</span>
          <span v-else>Send reset link</span>
        </Button>

        <!-- Footer Link back to Login -->
        <p class="text-center text-xs text-gray-500 dark:text-gray-400 pt-2">
          <RouterLink to="/login" class="text-black dark:text-white font-medium hover:underline">
            Back to sign in
          </RouterLink>
        </p>
      </form>

    </div>
  </div>
</template>