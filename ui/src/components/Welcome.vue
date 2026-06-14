<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { Rocket, Zap, Shield, Database, LayoutTemplate, Layers, Github, Terminal, ArrowRight } from 'lucide-vue-next'

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
    icon: Zap,
    title: 'Single binary in production',
    desc: 'The entire app compiles into one static binary. No Node, no Docker required to run.',
  },
  {
    icon: Rocket,
    title: 'Hot-reload in development',
    desc: 'Edit a .vue file → Vite HMR instantly updates. Edit a .go file → Air rebuilds in under a second.',
  },
  {
    icon: Layers,
    title: 'Clean Architecture',
    desc: 'Domain → Usecase → Repository → Delivery. Fully testable in isolation.',
  },
  {
    icon: Shield,
    title: 'Security headers built in',
    desc: 'XSS protection and no-store cache on all API routes — out of the box.',
  },
  {
    icon: LayoutTemplate,
    title: 'shadcn-vue UI library',
    desc: 'Accessible, unstyled components built on Radix. Ready to be customized.',
  },
  {
    icon: Database,
    title: 'MySQL + auto-migration',
    desc: 'Connection pool configured. Schema migrations run on startup. DB is optional.',
  },
]
</script>

<template>
  <div class="min-h-screen bg-background text-foreground font-sans selection:bg-primary/30 relative overflow-hidden">
    
    <!-- Background Decorators -->
    <div class="fixed inset-0 z-0 pointer-events-none">
      <div class="absolute -top-[30%] -left-[10%] w-[70%] h-[70%] rounded-full bg-primary/5 blur-[120px]"></div>
      <div class="absolute top-[20%] -right-[20%] w-[60%] h-[60%] rounded-full bg-blue-500/5 blur-[120px]"></div>
    </div>

    <!-- Navigation -->
    <header class="sticky top-0 z-50 w-full border-b border-border/40 bg-background/60 backdrop-blur-xl supports-[backdrop-filter]:bg-background/60">
      <div class="container mx-auto flex h-16 items-center justify-between px-6">
        <div class="flex items-center gap-3">
          <!-- Custom SVG Logo -->
          <img src="/vuelang-logo.svg" alt="Vuelang Logo" class="w-8 h-8 drop-shadow-md" />
          <span class="text-xl font-bold tracking-tight bg-clip-text text-transparent bg-gradient-to-r from-green-500 to-blue-600">
            Vuelang
          </span>
          <Badge variant="secondary" class="ml-2 text-[10px] uppercase tracking-widest bg-primary/10 hover:bg-primary/20 text-primary">v1.0</Badge>
        </div>
        <nav class="flex items-center gap-6 text-sm font-medium">
          <a href="#how-it-works" class="text-muted-foreground transition-colors hover:text-foreground">Features</a>
          <a href="#demo" class="text-muted-foreground transition-colors hover:text-foreground">Live Demo</a>
          <Button variant="outline" size="sm" class="gap-2 hidden sm:flex" asChild>
            <a href="https://github.com/milanz247/vuelang" target="_blank">
              <Github class="w-4 h-4" />
              GitHub
            </a>
          </Button>
        </nav>
      </div>
    </header>

    <!-- Hero -->
    <section class="relative z-10 pt-24 pb-32 lg:pt-36 lg:pb-40">
      <div class="container mx-auto px-6 text-center">
        <Badge variant="outline" class="mb-6 px-3 py-1 text-sm border-primary/20 bg-primary/5 text-primary backdrop-blur-md">
          <Rocket class="w-3 h-3 mr-2 inline-block" /> 
          The Ultimate Go + Vue Stack
        </Badge>
        
        <h1 class="mx-auto max-w-4xl text-5xl font-extrabold tracking-tight sm:text-6xl lg:text-7xl">
          Build <span class="bg-clip-text text-transparent bg-gradient-to-r from-green-400 to-blue-500">modern web apps</span> faster than ever.
        </h1>
        
        <p class="mx-auto mt-8 max-w-2xl text-lg text-muted-foreground sm:text-xl leading-relaxed">
          A full-stack MVC framework combining the sheer performance of a Go backend with the reactivity of a Vue 3 frontend. Ships as a single binary.
        </p>
        
        <div class="mt-10 flex flex-col sm:flex-row items-center justify-center gap-4">
          <Button size="lg" class="w-full sm:w-auto gap-2 group">
            Get Started
            <ArrowRight class="w-4 h-4 transition-transform group-hover:translate-x-1" />
          </Button>
          <Button size="lg" variant="secondary" class="w-full sm:w-auto gap-2">
            <Terminal class="w-4 h-4" />
            Documentation
          </Button>
        </div>
      </div>
    </section>

    <!-- Live API Demo -->
    <section id="demo" class="relative z-10 py-20 bg-muted/30 border-y border-border/50">
      <div class="container mx-auto px-6">
        <div class="flex flex-col items-center mb-12">
          <h2 class="text-3xl font-bold tracking-tight">Live Backend Interaction</h2>
          <p class="mt-3 text-muted-foreground max-w-xl text-center">
            Watch the Vue frontend seamlessly communicate with the Gin Go API.
          </p>
        </div>

        <div class="mx-auto max-w-2xl relative group">
          <!-- Glow effect behind card -->
          <div class="absolute -inset-1 bg-gradient-to-r from-green-500 to-blue-500 rounded-xl blur opacity-20 group-hover:opacity-40 transition duration-1000 group-hover:duration-200"></div>
          
          <Card class="relative bg-card/80 backdrop-blur-xl border-border/50 shadow-2xl overflow-hidden">
            <div class="h-10 border-b border-border/50 bg-muted/50 flex items-center px-4 gap-2">
              <div class="w-3 h-3 rounded-full bg-red-500/80"></div>
              <div class="w-3 h-3 rounded-full bg-yellow-500/80"></div>
              <div class="w-3 h-3 rounded-full bg-green-500/80"></div>
              <div class="ml-4 text-xs font-mono text-muted-foreground flex items-center">
                <Terminal class="w-3 h-3 mr-2" />
                GET /api/v1/greeting
              </div>
            </div>
            <CardHeader>
              <CardTitle>API Response</CardTitle>
              <CardDescription>Fetched directly from the Go database layer.</CardDescription>
            </CardHeader>
            <CardContent>
              <div class="min-h-[100px] rounded-lg border border-border/50 bg-background/50 p-6 font-mono text-sm relative overflow-hidden">
                <!-- Loader -->
                <div v-if="loading" class="flex items-center text-primary/80 animate-pulse">
                  <Zap class="w-4 h-4 mr-2" /> Fetching data payload...
                </div>
                <!-- Error -->
                <div v-else-if="error" class="text-destructive flex items-center">
                  <Shield class="w-4 h-4 mr-2" /> {{ error }}
                </div>
                <!-- Data -->
                <div v-else class="space-y-3">
                  <div class="flex items-start">
                    <span class="text-green-500 mr-4">message:</span>
                    <span class="text-foreground font-semibold">"{{ greeting.message }}"</span>
                  </div>
                  <div class="flex items-start">
                    <span class="text-blue-500 mr-4">timestamp:</span>
                    <span class="text-muted-foreground">"{{ greeting.timestamp }}"</span>
                  </div>
                </div>
              </div>

              <div class="mt-6 flex justify-end">
                <Button 
                  :disabled="loading" 
                  @click="fetchGreeting"
                  variant="outline"
                  class="gap-2 bg-background/50 backdrop-blur hover:bg-muted"
                >
                  <Rocket :class="['w-4 h-4', loading ? 'animate-spin' : '']" />
                  {{ loading ? 'Processing' : 'Execute Request' }}
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section id="how-it-works" class="relative z-10 py-24">
      <div class="container mx-auto px-6">
        <div class="text-center mb-16">
          <h2 class="text-3xl font-bold tracking-tight sm:text-4xl">Enterprise-Ready Architecture</h2>
          <p class="mt-4 text-muted-foreground">Everything you need to build scalable, secure, and blazing fast applications.</p>
        </div>

        <div class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          <Card 
            v-for="(f, index) in features" 
            :key="index"
            class="bg-card/40 border-border/40 backdrop-blur-sm transition-all duration-300 hover:shadow-xl hover:-translate-y-1 hover:bg-card/60 group"
          >
            <CardHeader>
              <div class="mb-4 inline-flex h-12 w-12 items-center justify-center rounded-lg bg-primary/10 text-primary group-hover:scale-110 transition-transform duration-300">
                <component :is="f.icon" class="h-6 w-6" />
              </div>
              <CardTitle class="text-xl">{{ f.title }}</CardTitle>
              <CardDescription class="mt-2 text-sm leading-relaxed">{{ f.desc }}</CardDescription>
            </CardHeader>
          </Card>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="relative z-10 bg-muted/20 border-t border-border/40 pb-8 pt-16">
      <div class="container mx-auto px-6">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8 mb-12">
          <div class="md:col-span-2">
            <div class="flex items-center gap-2 mb-4">
              <img src="/vuelang-logo.svg" alt="Vuelang" class="w-6 h-6 grayscale opacity-80" />
              <span class="text-lg font-bold tracking-tight text-foreground">Vuelang</span>
            </div>
            <p class="text-sm text-muted-foreground max-w-sm">
              The modern MVC framework bridging the gap between Go's legendary performance and Vue's delightful frontend experience.
            </p>
          </div>
          <div>
            <h3 class="font-semibold mb-4 text-foreground">Resources</h3>
            <ul class="space-y-2 text-sm text-muted-foreground">
              <li><a href="#" class="hover:text-primary transition-colors">Documentation</a></li>
              <li><a href="#" class="hover:text-primary transition-colors">Architecture</a></li>
              <li><a href="#" class="hover:text-primary transition-colors">GitHub Repository</a></li>
            </ul>
          </div>
          <div>
            <h3 class="font-semibold mb-4 text-foreground">Community</h3>
            <ul class="space-y-2 text-sm text-muted-foreground">
              <li><a href="#" class="hover:text-primary transition-colors">Discord Server</a></li>
              <li><a href="#" class="hover:text-primary transition-colors">Twitter</a></li>
            </ul>
          </div>
        </div>
        
        <Separator class="bg-border/50 mb-8" />
        
        <div class="flex flex-col md:flex-row items-center justify-between gap-4 text-xs text-muted-foreground">
          <p>© 2026 Vuelang Framework. Released under the MIT License.</p>
          <div class="flex items-center gap-2">
            <span>Powered by Go + Vue 3 + Tailwind</span>
          </div>
        </div>
      </div>
    </footer>

  </div>
</template>
