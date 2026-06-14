<template>
  <div class="min-h-screen bg-background text-foreground">

    <!-- ── Navigation ──────────────────────────────────────────────────── -->
    <header class="border-b border-border bg-card">
      <div class="container mx-auto flex h-14 items-center justify-between px-6">
        <div class="flex items-center gap-2">
          <!-- Wordmark -->
          <span class="text-xl font-bold tracking-tight text-foreground">
            Vue<span class="text-primary/60">lang</span>
          </span>
          <span class="rounded-md border border-border px-1.5 py-0.5 text-[10px] font-semibold uppercase tracking-widest text-muted-foreground">
            v1.0
          </span>
        </div>
        <nav class="flex items-center gap-6 text-sm text-muted-foreground">
          <a href="https://github.com/yourusername/vuelang" target="_blank"
             class="transition-colors hover:text-foreground">GitHub</a>
          <a href="#docs" class="transition-colors hover:text-foreground">Docs</a>
        </nav>
      </div>
    </header>

    <!-- ── Hero ────────────────────────────────────────────────────────── -->
    <section class="border-b border-border bg-card/50">
      <div class="container mx-auto px-6 py-20 text-center">
        <h1 class="mb-4 text-5xl font-bold tracking-tight text-foreground">
          Welcome to <span class="text-primary">Vuelang</span>
        </h1>
        <p class="mx-auto mb-8 max-w-xl text-lg text-muted-foreground">
          A full-stack framework — Go backend + Vue 3 frontend — that ships as a
          <strong class="text-foreground">single binary</strong> with hot-reload
          in development and zero runtime dependencies in production.
        </p>
        <div class="flex items-center justify-center gap-3">
          <a href="https://github.com/yourusername/vuelang"
             class="inline-flex h-10 items-center gap-2 rounded-md bg-primary px-5 text-sm font-medium text-primary-foreground transition hover:bg-primary/90">
            Get Started →
          </a>
          <a href="#how-it-works"
             class="inline-flex h-10 items-center rounded-md border border-border px-5 text-sm font-medium text-foreground transition hover:bg-muted">
            How it works
          </a>
        </div>
      </div>
    </section>

    <!-- ── Live API demo ────────────────────────────────────────────────── -->
    <section class="container mx-auto px-6 py-16">
      <h2 class="mb-2 text-center text-2xl font-semibold tracking-tight">Live backend response</h2>
      <p class="mb-8 text-center text-sm text-muted-foreground">
        This data is fetched from the Go API at <code class="rounded bg-muted px-1 py-0.5 text-xs">/api/v1/greeting</code> on every load.
      </p>

      <div class="mx-auto max-w-md">
        <Card>
          <CardHeader>
            <CardTitle>Greeting from Go</CardTitle>
            <CardDescription>Served by Gin · Clean Architecture · MySQL-ready</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="rounded-md border border-border bg-muted/40 p-4 text-sm">
              <p v-if="loading" class="text-muted-foreground">Fetching from backend…</p>
              <p v-else-if="error" class="text-destructive">{{ error }}</p>
              <template v-else>
                <p class="font-semibold text-foreground">{{ greeting.message }}</p>
                <p class="mt-1 text-xs text-muted-foreground">{{ greeting.timestamp }}</p>
              </template>
            </div>

            <button
              :disabled="loading"
              @click="fetchGreeting"
              class="mt-5 inline-flex h-9 w-full items-center justify-center rounded-md bg-primary text-sm font-medium text-primary-foreground transition hover:bg-primary/90 disabled:opacity-50"
            >
              {{ loading ? 'Loading…' : 'Refresh' }}
            </button>
          </CardContent>
        </Card>
      </div>
    </section>

    <!-- ── Feature cards ─────────────────────────────────────────────────── -->
    <section class="border-t border-border bg-muted/30">
      <div class="container mx-auto px-6 py-16">
        <h2 class="mb-10 text-center text-2xl font-semibold tracking-tight">Built for real products</h2>
        <div class="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
          <Card v-for="f in features" :key="f.title">
            <CardHeader>
              <div class="mb-2 text-2xl">{{ f.icon }}</div>
              <CardTitle class="text-base">{{ f.title }}</CardTitle>
              <CardDescription>{{ f.desc }}</CardDescription>
            </CardHeader>
          </Card>
        </div>
      </div>
    </section>

    <!-- ── Footer ──────────────────────────────────────────────────────── -->
    <footer class="border-t border-border">
      <div class="container mx-auto flex flex-col items-center gap-1 px-6 py-8 text-center text-xs text-muted-foreground sm:flex-row sm:justify-between">
        <span>© 2025 <strong class="text-foreground">Vuelang</strong>. MIT License.</span>
        <span>Go + Vue 3 + shadcn-vue + Tailwind CSS</span>
      </div>
    </footer>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'

interface Greeting {
  message: string
  timestamp: string
}

const greeting = ref<Greeting>({ message: '', timestamp: '' })
const loading  = ref(false)
const error    = ref('')

async function fetchGreeting() {
  loading.value = true
  error.value   = ''
  try {
    const res = await fetch('/api/v1/greeting')
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    greeting.value = await res.json()
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to fetch'
  } finally {
    loading.value = false
  }
}

onMounted(fetchGreeting)

const features = [
  {
    icon: '⚡',
    title: 'Single binary in production',
    desc: 'The entire app — Go API + Vue frontend — compiles into one static binary. No Node, no Docker required to run.',
  },
  {
    icon: '🔥',
    title: 'Hot-reload in development',
    desc: 'Edit a .vue file → Vite HMR updates the browser instantly. Edit a .go file → Air rebuilds in under a second.',
  },
  {
    icon: '🏗️',
    title: 'Clean Architecture',
    desc: 'Domain → Usecase → Repository → Delivery. Each layer is testable in isolation. Add a module by copying the pattern.',
  },
  {
    icon: '🔒',
    title: 'Security headers built in',
    desc: 'X-Content-Type-Options, X-Frame-Options, XSS protection, and no-store cache on all API routes — out of the box.',
  },
  {
    icon: '🎨',
    title: 'shadcn-vue UI library',
    desc: 'Accessible, unstyled components built on Radix. Add new ones with one command: npx shadcn-vue@latest add button.',
  },
  {
    icon: '🗄️',
    title: 'MySQL + auto-migration',
    desc: 'Connection pool configured and ready. Schema migrations run on startup. DB is optional — server starts without it.',
  },
]
</script>
