<!-- File: vuelang/ui/src/views/DashboardView.vue -->
<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Separator } from '@/components/ui/separator'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

import vuelangLogo from '@/assets/vuelang-logo.svg'

const router = useRouter()
const auth   = useAuthStore()

// Responsive and Collapsible States
const isSidebarCollapsed = ref(false)
const mobileSidebarOpen = ref(false)

// Sidebar dropdown menu state (Accordion Mode)
const isManagementOpen = ref(false)

async function handleLogout() {
  await auth.logout()
  router.push('/login')
}

// Sidebar status toggle function
function toggleSidebar() {
  if (window.innerWidth < 768) {
    mobileSidebarOpen.value = !mobileSidebarOpen.value
  } else {
    isSidebarCollapsed.value = !isSidebarCollapsed.value
  }
}

function getInitials(name?: string) {
  if (!name) return 'VU'
  const parts = name.trim().split(' ')
  if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase()
  return (parts[0][0] + parts[1][0]).toUpperCase()
}
</script>

<template>
  <div class="min-h-screen bg-[#FDFDFC] dark:bg-[#0a0a0a] text-foreground font-sans transition-colors duration-300">
    
    <!-- 1. RESPONSIVE SIDEBAR (Desktop: Collapsible Width, Mobile: Slide-out Drawer) -->
    <aside 
      :class="[
        'fixed top-0 bottom-0 left-0 z-40 bg-[#FAFAFA] dark:bg-[#111] border-r border-border flex flex-col transition-all duration-300',
        isSidebarCollapsed ? 'w-[60px]' : 'w-64',
        mobileSidebarOpen ? 'translate-x-0 w-64' : 'max-md:-translate-x-full md:translate-x-0'
      ]"
    >
      <!-- Sidebar Header (Logo & Brand) -->
      <div class="h-14 flex items-center px-4 gap-2.5 border-b border-border select-none overflow-hidden shrink-0">
        <img :src="vuelangLogo" alt="Vuelang Logo" class="w-5 h-5 object-contain shrink-0" />
        <span 
          v-if="!isSidebarCollapsed" 
          class="font-bold text-sm tracking-tight text-[#1b1b18] dark:text-[#EDEDEC] transition-opacity duration-300"
        >
          Vuelang
        </span>
      </div>

      <!-- Navigation Items -->
      <div class="flex-1 px-3 py-6 space-y-6 overflow-y-auto">
        <div>
          <!-- Section label (Hidden when collapsed) -->
          <div 
            v-if="!isSidebarCollapsed" 
            class="px-2.5 text-[10px] font-bold text-muted-foreground uppercase tracking-wider mb-2 select-none transition-opacity duration-300"
          >
            Platform
          </div>
          
          <div class="space-y-1">
            <!-- Dashboard Main Link -->
            <RouterLink to="/dashboard">
              <Button 
                variant="secondary" 
                :class="[
                  'w-full justify-start gap-2.5 h-9 text-xs shadow-none border-0 bg-white dark:bg-[#1a1a19] text-[#1b1b18] dark:text-[#EDEDEC] font-semibold',
                  isSidebarCollapsed ? 'justify-center px-0' : 'px-3'
                ]"
              >
                <!-- Layout Grid Icon -->
                <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0 text-foreground" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                  <rect x="3" y="3" width="7" height="9" rx="1" />
                  <rect x="14" y="3" width="7" height="5" rx="1" />
                  <rect x="14" y="12" width="7" height="9" rx="1" />
                  <rect x="3" y="16" width="7" height="5" rx="1" />
                </svg>
                <span v-if="!isSidebarCollapsed" class="transition-opacity duration-300">Dashboard</span>
              </Button>
            </RouterLink>

            <!-- Collapsible Sidebar Dropdown Menu (Management) -->
            <div class="space-y-1">
              
              <!-- Case A: Desktop Collapsed (Triggers a floating Dropdown Menu) -->
              <DropdownMenu v-if="isSidebarCollapsed">
                <DropdownMenuTrigger as-child>
                  <Button 
                    variant="ghost" 
                    class="w-full justify-center h-9 px-0 shadow-none border-0 text-muted-foreground hover:text-foreground hover:bg-[#1b1b18]/5 dark:hover:bg-white/5"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                    </svg>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent side="right" align="start" class="w-40 rounded-lg">
                  <DropdownMenuLabel class="text-xs font-semibold text-muted-foreground">Management</DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  <DropdownMenuItem class="cursor-pointer text-xs">Users</DropdownMenuItem>
                  <DropdownMenuItem class="cursor-pointer text-xs">Audit Logs</DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>

              <!-- Case B: Normal Extended Sidebar (Accordion Toggle) -->
              <template v-else>
                <Button 
                  variant="ghost" 
                  @click="isManagementOpen = !isManagementOpen"
                  class="w-full justify-between h-9 text-xs px-3 shadow-none border-0 text-[#1b1b18] dark:text-[#EDEDEC] font-semibold hover:bg-[#1b1b18]/5 dark:hover:bg-white/5"
                >
                  <div class="flex items-center gap-2.5">
                    <!-- Folder/Management Icon -->
                    <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4 shrink-0 text-muted-foreground" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                    </svg>
                    <span class="transition-opacity duration-300">Management</span>
                  </div>
                  <!-- Rotating Chevron Arrow -->
                  <svg 
                    :class="['w-3.5 h-3.5 text-muted-foreground shrink-0 transition-transform duration-200', isManagementOpen ? 'rotate-90' : '']"
                    fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24"
                  >
                    <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
                  </svg>
                </Button>

                <!-- Accordion Inner Sub-Items list -->
                <div 
                  v-if="isManagementOpen" 
                  class="pl-6 space-y-1 relative before:absolute before:left-[19px] before:top-0 before:bottom-2 before:w-[1px] before:bg-border"
                >
                  <RouterLink to="/dashboard">
                    <Button 
                      variant="ghost" 
                      class="w-full h-8 text-[11px] justify-start px-4 text-muted-foreground hover:text-foreground hover:bg-transparent"
                    >
                      Users
                    </Button>
                  </RouterLink>
                  <RouterLink to="/dashboard">
                    <Button 
                      variant="ghost" 
                      class="w-full h-8 text-[11px] justify-start px-4 text-muted-foreground hover:text-foreground hover:bg-transparent"
                    >
                      Audit Logs
                    </Button>
                  </RouterLink>
                </div>
              </template>

            </div>

          </div>
        </div>
      </div>

      <!-- Sidebar Footer (Repository, Docs & User Info) -->
      <div class="p-3 border-t border-border space-y-4 shrink-0">
        <div class="space-y-1">
          <a 
            href="https://github.com/milanz247/vuelang" 
            target="_blank" 
            :class="[
              'flex items-center gap-2.5 py-2 text-xs font-medium text-muted-foreground hover:text-foreground rounded-md transition',
              isSidebarCollapsed ? 'justify-center px-0' : 'px-3'
            ]"
          >
            <!-- Github Icon -->
            <svg class="w-4 h-4 shrink-0" fill="currentColor" viewBox="0 0 24 24">
              <path fill-rule="evenodd" clip-rule="evenodd" d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.554 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"/>
            </svg>
            <span v-if="!isSidebarCollapsed">Repository</span>
          </a>
          <a 
            href="https://github.com/milanz247/vuelang" 
            target="_blank" 
            :class="[
              'flex items-center gap-2.5 py-2 text-xs font-medium text-muted-foreground hover:text-foreground rounded-md transition',
              isSidebarCollapsed ? 'justify-center px-0' : 'px-3'
            ]"
          >
            <!-- Book Icon -->
            <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
            </svg>
            <span v-if="!isSidebarCollapsed">Documentation</span>
          </a>
        </div>

        <!-- User Dropdown Combobox -->
        <DropdownMenu>
          <DropdownMenuTrigger as-child>
            <button 
              :class="[
                'flex items-center gap-3 w-full rounded-lg transition text-left focus:outline-none',
                isSidebarCollapsed ? 'justify-center p-0.5 hover:bg-transparent' : 'p-1.5 hover:bg-muted/80 dark:hover:bg-zinc-800/80'
              ]"
            >
              <Avatar class="w-8 h-8 rounded-md shrink-0">
                <AvatarFallback class="rounded-md text-[11px] font-bold bg-[#1b1b18]/10 dark:bg-white/10 text-foreground select-none">
                  {{ getInitials(auth.user?.name) }}
                </AvatarFallback>
              </Avatar>
              <div v-if="!isSidebarCollapsed" class="flex-1 min-w-0">
                <p class="text-xs font-semibold text-foreground truncate select-none leading-none mb-0.5">
                  {{ auth.user?.name }}
                </p>
                <p class="text-[10px] text-muted-foreground truncate select-none leading-none">
                  {{ auth.user?.email }}
                </p>
              </div>
              <svg v-if="!isSidebarCollapsed" class="w-3.5 h-3.5 text-muted-foreground shrink-0 select-none" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" d="M8 9l4-4 4 4m0 6l-4 4-4-4"/>
              </svg>
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" class="w-56 rounded-lg">
            <DropdownMenuLabel class="text-xs font-normal text-muted-foreground">My Account</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem class="cursor-pointer text-xs">Profile Settings</DropdownMenuItem>
            <DropdownMenuItem class="cursor-pointer text-xs">System Logs</DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem @click="handleLogout" class="cursor-pointer text-xs text-destructive focus:bg-destructive/10 focus:text-destructive">
              Sign out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </aside>

    <!-- Overlay backdrop for mobile collapsible sidebar -->
    <div 
      v-if="mobileSidebarOpen" 
      @click="mobileSidebarOpen = false" 
      class="fixed inset-0 z-30 bg-black/40 md:hidden"
    ></div>

    <!-- 2. MAIN PANELS CONTAINER (Dynamic padding based on collapse state) -->
    <div 
      :class="[
        'flex flex-col min-h-screen transition-all duration-300',
        isSidebarCollapsed ? 'md:pl-[60px]' : 'md:pl-64'
      ]"
    >
      
      <!-- Top header bar containing collapsible triggers -->
      <header class="h-14 flex items-center justify-between px-6 border-b border-border bg-[#FAFAFA] dark:bg-[#111] shrink-0">
        <div class="flex items-center gap-3">
          
          <!-- Desktop/Mobile Sidebar Toggle Button -->
          <Button 
            variant="ghost" 
            size="icon" 
            @click="toggleSidebar" 
            class="shrink-0"
          >
            <!-- Modern Sidebar Column Toggle Icon -->
            <svg class="w-4.5 h-4.5 text-muted-foreground" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
              <rect x="3" y="3" width="18" height="18" rx="2" />
              <path d="M9 3v18" />
            </svg>
          </Button>
          
          <span class="text-sm font-semibold text-foreground select-none">Dashboard</span>
        </div>
      </header>

      <!-- Main Dashboard Content Grid with enlarged statistics -->
      <main class="flex-1 p-6 lg:p-10 space-y-8 max-w-7xl w-full mx-auto">
        
        <!-- Welcome Header -->
        <div class="space-y-1 select-none">
          <h1 class="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
            Welcome, {{ auth.user?.name ?? 'User' }}! 👋
          </h1>
          <p class="text-xs text-muted-foreground">
            Monitor and manage your high-performance Vuelang V2 application runtime.
          </p>
        </div>

        <!-- Metric Stat Cards Grid -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <!-- Card 1 -->
          <Card class="border border-dashed border-border bg-muted/5 shadow-none rounded-xl overflow-hidden relative">
            <CardHeader class="pb-2">
              <CardDescription class="text-xs font-semibold text-muted-foreground select-none">Active Sessions</CardDescription>
              <CardTitle class="text-2xl font-bold tracking-tight text-gray-900 dark:text-white select-none">2,840</CardTitle>
            </CardHeader>
            <CardContent>
              <p class="text-[10px] text-green-600 dark:text-green-400 font-semibold select-none">+12.3% from last week</p>
            </CardContent>
          </Card>

          <!-- Card 2 -->
          <Card class="border border-dashed border-border bg-muted/5 shadow-none rounded-xl overflow-hidden relative">
            <CardHeader class="pb-2">
              <CardDescription class="text-xs font-semibold text-muted-foreground select-none">API Latency</CardDescription>
              <CardTitle class="text-2xl font-bold tracking-tight text-gray-900 dark:text-white select-none">42ms</CardTitle>
            </CardHeader>
            <CardContent>
              <p class="text-[10px] text-muted-foreground select-none">Go runtime optimized</p>
            </CardContent>
          </Card>

          <!-- Card 3 -->
          <Card class="border border-dashed border-border bg-muted/5 shadow-none rounded-xl overflow-hidden relative">
            <CardHeader class="pb-2">
              <CardDescription class="text-xs font-semibold text-muted-foreground select-none">Database Pool</CardDescription>
              <CardTitle class="text-2xl font-bold tracking-tight text-gray-900 dark:text-white select-none">8 / 50</CardTitle>
            </CardHeader>
            <CardContent>
              <p class="text-[10px] text-green-600 dark:text-green-400 font-semibold select-none">Healthy connection pool</p>
            </CardContent>
          </Card>
        </div>

        <!-- Large Content Panel: Recent System Logs & Actions -->
        <Card class="border border-dashed border-border bg-muted/5 shadow-none rounded-xl overflow-hidden">
          <CardHeader class="border-b border-dashed border-border bg-muted/10 p-5">
            <CardTitle class="text-sm font-semibold select-none">Recent System Audit Logs</CardTitle>
            <CardDescription class="text-xs select-none">Real-time state-changing operations recorded by the Go backend.</CardDescription>
          </CardHeader>
          <CardContent class="p-0">
            <div class="divide-y divide-dashed divide-border text-xs">
              <!-- Log Line 1 -->
              <div class="p-4 flex items-center justify-between gap-4">
                <div class="flex items-center gap-3">
                  <span class="w-1.5 h-1.5 rounded-full bg-green-500 shrink-0"></span>
                  <span class="font-medium text-gray-700 dark:text-gray-300">User session started successfully</span>
                </div>
                <span class="text-[10px] text-muted-foreground">Just now</span>
              </div>
              <!-- Log Line 2 -->
              <div class="p-4 flex items-center justify-between gap-4">
                <div class="flex items-center gap-3">
                  <span class="w-1.5 h-1.5 rounded-full bg-[#FF2D20] shrink-0"></span>
                  <span class="font-medium text-gray-700 dark:text-gray-300">Role assignment verified for admin</span>
                </div>
                <span class="text-[10px] text-muted-foreground">2 mins ago</span>
              </div>
              <!-- Log Line 3 -->
              <div class="p-4 flex items-center justify-between gap-4">
                <div class="flex items-center gap-3">
                  <span class="w-1.5 h-1.5 rounded-full bg-blue-500 shrink-0"></span>
                  <span class="font-medium text-gray-700 dark:text-gray-300">Database migration tables scanned</span>
                </div>
                <span class="text-[10px] text-muted-foreground">10 mins ago</span>
              </div>
            </div>
          </CardContent>
        </Card>
        
      </main>
    </div>

  </div>
</template>