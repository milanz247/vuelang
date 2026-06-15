<!-- File: vuelang/ui/src/views/WelcomeView.vue -->
<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'

// Import the logo using Vite's path alias
import vuelangLogo from '@/assets/vuelang-logo.svg'

const auth = useAuthStore()
</script>

<template>
  <div
    class="flex min-h-screen flex-col items-center bg-[#FDFDFC] p-6 text-[#1b1b18] lg:justify-center lg:p-8 dark:bg-[#0a0a0a]"
  >
    <!-- Top Navigation Header (using Shadcn Button component) -->
    <header class="mb-6 w-full max-w-[335px] text-sm lg:max-w-4xl">
      <nav class="flex items-center justify-end gap-4">
        <!-- If Authenticated -->
        <Button
          v-if="auth.isAuthenticated"
          variant="outline"
          as-child
          class="rounded-sm border border-[#19140035] px-5 py-1.5 text-sm leading-normal text-[#1b1b18] hover:border-[#1915014a] hover:bg-transparent dark:border-[#3E3E3A] dark:text-[#EDEDEC] dark:hover:border-[#62605b]"
        >
          <RouterLink to="/dashboard">Dashboard</RouterLink>
        </Button>

        <!-- If Guest -->
        <template v-else>
          <Button
            variant="ghost"
            as-child
            class="rounded-sm border border-transparent px-5 py-1.5 text-sm leading-normal text-[#1b1b18] hover:border-[#19140035] hover:bg-transparent dark:text-[#EDEDEC] dark:hover:border-[#3E3E3A]"
          >
            <RouterLink to="/login">Log in</RouterLink>
          </Button>
          <Button
            variant="outline"
            as-child
            class="rounded-sm border border-[#19140035] px-5 py-1.5 text-sm leading-normal text-[#1b1b18] hover:border-[#1915014a] hover:bg-transparent dark:border-[#3E3E3A] dark:text-[#EDEDEC] dark:hover:border-[#62605b]"
          >
            <RouterLink to="/register">Register</RouterLink>
          </Button>
        </template>
      </nav>
    </header>

    <div
      class="flex w-full items-center justify-center opacity-100 transition-opacity duration-750 lg:grow starting:opacity-0"
    >
      <!-- Main Content Container (using Shadcn Card) -->
      <Card
        class="flex w-full max-w-[335px] flex-col-reverse overflow-hidden rounded-lg border-0 lg:max-w-4xl lg:flex-row shadow-[inset_0px_0px_0px_1px_rgba(26,26,0,0.16)] dark:shadow-[inset_0px_0px_0px_1px_#fffaed2d] bg-white dark:bg-[#161615]"
      >
        <!-- Left Side: Information & Steps -->
        <div
          class="flex-1 p-6 pb-12 text-[13px] leading-[20px] lg:p-20 text-[#1b1b18] dark:text-[#EDEDEC] flex flex-col justify-center"
        >
          <h1 class="mb-1.5 font-semibold text-[15px] tracking-tight text-gray-900 dark:text-white">
            Welcome to Vuelang
          </h1>
          <p class="mb-6 text-[#706f6c] dark:text-[#A1A09A] leading-relaxed">
            An enterprise-grade full-stack framework bridging Go's blazing performance with the reactive elegance of Vue 3. Compile your entire application into a single executable binary.
          </p>

          <!-- Clean timeline list using Vuelang design patterns -->
          <ul class="mb-4 flex flex-col lg:mb-8">
            <!-- Step 1 -->
            <li
              class="relative flex items-center gap-4 py-2.5 before:absolute before:top-1/2 before:bottom-0 before:left-[0.4rem] before:border-l before:border-[#e3e3e0] dark:before:border-[#3E3E3A]"
            >
              <span class="relative bg-white py-1 dark:bg-[#161615]">
                <span
                  class="flex h-3.5 w-3.5 items-center justify-center rounded-full border border-[#e3e3e0] bg-[#FDFDFC] shadow-[0px_0px_1px_0px_rgba(0,0,0,0.03),0px_1px_2px_0px_rgba(0,0,0,0.06)] dark:border-[#3E3E3A] dark:bg-[#161615]"
                >
                  <span class="h-1.5 w-1.5 rounded-full bg-[#dbdbd7] dark:bg-[#3E3E3A]" />
                </span>
              </span>
              <span>
                Read the comprehensive
                <a
                  href="https://github.com/milanz247/vuelang"
                  target="_blank"
                  class="ml-1 inline-flex items-center space-x-1 font-semibold text-[#0f766e] underline underline-offset-4 dark:text-[#2dd4bf]"
                >
                  <span>Framework Guide</span>
                  <svg
                    width="10"
                    height="11"
                    viewBox="0 0 10 11"
                    fill="none"
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-2.5 w-2.5"
                  >
                    <path
                      d="M7.70833 6.95834V2.79167H3.54167M2.5 8L7.5 3.00001"
                      stroke="currentColor"
                      stroke-linecap="square"
                    />
                  </svg>
                </a>
              </span>
            </li>

            <!-- Step 2 -->
            <li
              class="relative flex items-center gap-4 py-2.5 before:absolute before:top-0 before:bottom-1/2 before:left-[0.4rem] before:border-l before:border-[#e3e3e0] dark:before:border-[#3E3E3A]"
            >
              <span class="relative bg-white py-1 dark:bg-[#161615]">
                <span
                  class="flex h-3.5 w-3.5 items-center justify-center rounded-full border border-[#e3e3e0] bg-[#FDFDFC] shadow-[0px_0px_1px_0px_rgba(0,0,0,0.03),0px_1px_2px_0px_rgba(0,0,0,0.06)] dark:border-[#3E3E3A] dark:bg-[#161615]"
                >
                  <span class="h-1.5 w-1.5 rounded-full bg-[#dbdbd7] dark:bg-[#3E3E3A]" />
                </span>
              </span>
              <span>
                Check out the core code on
                <a
                  href="https://github.com/milanz247/vuelang"
                  target="_blank"
                  class="ml-1 inline-flex items-center space-x-1 font-semibold text-[#0f766e] underline underline-offset-4 dark:text-[#2dd4bf]"
                >
                  <span>GitHub Repository</span>
                  <svg
                    width="10"
                    height="11"
                    viewBox="0 0 10 11"
                    fill="none"
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-2.5 w-2.5"
                  >
                    <path
                      d="M7.70833 6.95834V2.79167H3.54167M2.5 8L7.5 3.00001"
                      stroke="currentColor"
                      stroke-linecap="square"
                    />
                  </svg>
                </a>
              </span>
            </li>
          </ul>

          <!-- Deploy Button (using Shadcn Button) -->
          <ul class="flex gap-3 text-sm leading-normal">
            <li>
              <Button
                variant="default"
                as-child
                class="inline-block rounded-sm border border-black bg-[#1b1b18] px-5 py-1.5 text-sm leading-normal text-white hover:border-black hover:bg-black dark:border-[#eeeeec] dark:bg-[#eeeeec] dark:text-[#1C1C1A] dark:hover:border-white dark:hover:bg-white"
              >
                <RouterLink to="/dashboard">Deploy now</RouterLink>
              </Button>
            </li>
          </ul>
        </div>

        <!-- Right Side: Vuelang Logo Panel (No internal borders, cleanly scaled) -->
        <div
          class="relative -mb-px aspect-[335/364] w-full shrink-0 overflow-hidden bg-[#F8FAFC] lg:mb-0 lg:-ml-px lg:aspect-auto lg:w-[438px] dark:bg-[#111110]"
        >
          <!-- Perfectly Scaled and Centered Vuelang Logo -->
          <div class="absolute inset-0 flex items-center justify-center p-10">
            <img
              :src="vuelangLogo"
              alt="Vuelang Logo"
              class="w-52 h-52 lg:w-64 lg:h-64 object-contain opacity-100 transition-all duration-750 starting:opacity-0 motion-safe:starting:translate-y-6"
            />
          </div>
        </div>
      </Card>
    </div>
    <div class="hidden h-14.5 lg:block"></div>
  </div>
</template>