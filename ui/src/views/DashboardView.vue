<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900">
    <!-- Navbar -->
    <nav class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16 items-center">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-primary flex items-center justify-center">
              <span class="text-white font-bold text-sm">V</span>
            </div>
            <span class="font-semibold text-gray-900 dark:text-white">Vuelang</span>
          </div>

          <div class="flex items-center gap-4">
            <span class="text-sm text-gray-500 dark:text-gray-400 hidden sm:block">
              {{ auth.user?.email }}
            </span>
            <div class="flex gap-1.5">
              <span
                v-for="role in auth.user?.roles"
                :key="role"
                class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-primary/10 text-primary"
              >
                {{ role }}
              </span>
            </div>
            <button
              @click="handleLogout"
              class="text-sm text-gray-500 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white transition"
            >
              Sign out
            </button>
          </div>
        </div>
      </div>
    </nav>

    <!-- Content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
          Welcome, {{ auth.user?.name ?? 'User' }}! 👋
        </h1>
        <p class="text-gray-500 dark:text-gray-400 mt-1">
          Your Vuelang V2 dashboard is ready.
        </p>
      </div>

      <!-- Stats grid -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-10">
        <div
          v-for="stat in stats"
          :key="stat.label"
          class="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6"
        >
          <p class="text-sm font-medium text-gray-500 dark:text-gray-400">{{ stat.label }}</p>
          <p class="text-2xl font-bold text-gray-900 dark:text-white mt-1">{{ stat.value }}</p>
          <p class="text-xs text-green-600 dark:text-green-400 mt-1">{{ stat.sub }}</p>
        </div>
      </div>

      <!-- Framework features -->
      <div class="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Vuelang V2 Features</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <div
            v-for="feature in features"
            :key="feature.title"
            class="flex items-start gap-3 p-4 rounded-lg bg-gray-50 dark:bg-gray-700/50"
          >
            <div class="flex-shrink-0 w-8 h-8 rounded-lg bg-primary/10 flex items-center justify-center text-primary text-sm">
              {{ feature.icon }}
            </div>
            <div>
              <p class="text-sm font-medium text-gray-900 dark:text-white">{{ feature.title }}</p>
              <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">{{ feature.desc }}</p>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth   = useAuthStore()

const stats = [
  { label: 'Framework',  value: 'Vuelang V2', sub: 'Enterprise ready'  },
  { label: 'Backend',    value: 'Go + Gin',   sub: 'High performance'  },
  { label: 'Frontend',   value: 'Vue 3 + Vite', sub: 'Instant HMR'    },
  { label: 'Auth',       value: 'JWT + RBAC', sub: 'Production secure' },
]

const features = [
  { icon: '🔐', title: 'JWT Authentication',     desc: 'Access + refresh tokens, automatic rotation'   },
  { icon: '🛡️', title: 'RBAC',                  desc: 'Role-based access control with middleware'      },
  { icon: '⚡', title: 'Rate Limiting',           desc: 'IP-based token bucket per endpoint group'     },
  { icon: '📦', title: 'Service Layer',           desc: 'Business logic separated from controllers'    },
  { icon: '🗄️', title: 'Repository Pattern',     desc: 'Clean data access layer, easy to test'        },
  { icon: '🔒', title: 'Security Headers',        desc: 'CSP, HSTS, X-Frame-Options and more'          },
  { icon: '🐳', title: 'Docker Ready',            desc: 'Multi-stage build, single binary deployment'  },
  { icon: '📊', title: 'Audit Logs',              desc: 'Track every state-changing operation'         },
  { icon: '🛠️', title: 'CLI Scaffolding',        desc: 'vuelang make:model, make:controller and more' },
]

async function handleLogout() {
  await auth.logout()
  router.push('/login')
}
</script>
