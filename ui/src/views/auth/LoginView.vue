<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card'
import { Alert, AlertDescription } from '@/components/ui/alert'

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

<template>
  <div class="min-h-screen flex items-center justify-center bg-background px-4">
    <div class="w-full max-w-md">
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-12 h-12 rounded-xl bg-primary mb-4">
          <span class="text-primary-foreground font-bold text-xl">V</span>
        </div>
        <h1 class="text-2xl font-bold">Welcome back</h1>
        <p class="text-muted-foreground mt-1">Sign in to your Vuelang account</p>
      </div>

      <Card>
        <CardContent class="pt-6">
          <form @submit.prevent="handleLogin" class="space-y-5">
            <Alert v-if="error" variant="destructive">
              <AlertDescription>{{ error }}</AlertDescription>
            </Alert>

            <div class="space-y-2">
              <Label for="email">Email address</Label>
              <Input
                id="email"
                v-model="form.email"
                type="email"
                autocomplete="email"
                required
                placeholder="you@example.com"
              />
            </div>

            <div class="space-y-2">
              <div class="flex justify-between items-center">
                <Label for="password">Password</Label>
                <RouterLink to="/forgot-password" class="text-sm text-primary hover:underline">
                  Forgot password?
                </RouterLink>
              </div>
              <Input
                id="password"
                v-model="form.password"
                type="password"
                autocomplete="current-password"
                required
                placeholder="••••••••"
              />
            </div>

            <Button type="submit" class="w-full" :disabled="loading">
              <span v-if="loading">Signing in…</span>
              <span v-else>Sign in</span>
            </Button>
          </form>

          <p class="mt-6 text-center text-sm text-muted-foreground">
            Don't have an account?
            <RouterLink to="/register" class="text-primary font-medium hover:underline">
              Sign up
            </RouterLink>
          </p>
        </CardContent>
      </Card>
    </div>
  </div>
</template>