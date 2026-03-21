<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

// State for main navigation sidebar
const isSidebarCollapsed = ref(false)

// Frequently used items (Top section)
const primaryNavItems = [
  { name: '系统监控', icon: 'lucide:layout-dashboard', path: '/' },
  { name: '数字员工', icon: 'lucide:bot', path: '/agent-management' },
  { name: '单人对话', icon: 'lucide:user', path: '/single-person-chat' },
  { name: '多人对话', icon: 'lucide:message-square', path: '/multi-person-chat' },
  { name: '编排画布', icon: 'lucide:workflow', path: '/orchestration-editor' },
  { name: '定时任务', icon: 'lucide:calendar', path: '/schedule-tasks' },
]

// Less frequently used / Management items (Bottom section)
const secondaryNavItems = [
  { name: '模型管理', icon: 'lucide:cpu', path: '/model-management' },
  { name: '技能管理', icon: 'lucide:zap', path: '/skill-management' },
  { name: '执行历史', icon: 'lucide:history', path: '/execution-monitor' },
]

const isActive = (path: string) => route.path === path

const toggleSidebar = () => {
  isSidebarCollapsed.value = !isSidebarCollapsed.value
}
</script>

<template>
  <div class="font-[-apple-system,BlinkMacSystemFont,'Segoe UI'] w-full min-h-screen flex flex-col" style="line-height: 1.5; background: rgba(15, 25, 40, 1);">
    <!-- Header -->
    <header style="background: linear-gradient(180deg, rgba(15, 25, 45, 0.85) 0%, rgba(20, 35, 60, 0.85) 100%);" class="w-full shrink-0 z-30">
      <div style="padding: 1rem 2rem;" class="flex justify-between items-center">
        <div class="flex items-center gap-10">
          <div class="flex items-center gap-3">
            <div class="flex justify-center items-center w-6 h-6">
              <iconify-icon style="color: rgba(59, 155, 255, 1);" icon="lucide:workflow" class="text-xl"></iconify-icon>
            </div>
            <span style="color: color-mix(in oklab, #fff 95%, transparent);" class="text-xl font-semibold">AgentFlow</span>
          </div>
          <div style="background-color: color-mix(in oklab, #fff 5%, transparent); backdrop-filter: blur(24px); padding: 0.5rem 1rem; border-color: color-mix(in oklab, #3B9BFF 30%, transparent);" class="flex items-center min-w-[300px] border-[1px] border-solid rounded-xl">
            <div class="flex justify-center items-center w-5 h-5 mr-3">
              <iconify-icon style="color: color-mix(in oklab, #fff 50%, transparent);" icon="lucide:search" class="text-base"></iconify-icon>
            </div>
            <input style="flex-basis: 0%; color: color-mix(in oklab, #fff 70%, transparent);" type="text" placeholder="搜索Agent、编排或任务..." class="bg-transparent grow shrink outline-none">
          </div>
        </div>
        <div class="flex items-center gap-4">
          <button class="hover:border-[#5FB4FF]/50 hover:shadow-[0_0_20px_rgba(59,155,255,0.3)] flex items-center gap-2 border-[1px] border-solid rounded-xl" style="background-color: color-mix(in oklab, #fff 5%, transparent); backdrop-filter: blur(24px); padding: 0.5rem 1rem; border-color: color-mix(in oklab, #3B9BFF 30%, transparent);">
            <div class="flex justify-center items-center w-4 h-4">
              <iconify-icon style="color: rgba(95, 180, 255, 1);" icon="lucide:plus" class="text-sm"></iconify-icon>
            </div>
            <span style="color: color-mix(in oklab, #fff 95%, transparent);" class="text-sm whitespace-nowrap">新建</span>
          </button>
          <div class="flex relative justify-center items-center w-12 h-12">
            <button class="hover:border-[#5FB4FF]/50 flex justify-center items-center w-10 h-10 border-[1px] border-solid rounded-xl" style="background-color: color-mix(in oklab, #fff 5%, transparent); backdrop-filter: blur(24px); border-color: color-mix(in oklab, #3B9BFF 30%, transparent);">
              <div class="flex justify-center items-center w-5 h-5">
                <iconify-icon style="color: color-mix(in oklab, #fff 70%, transparent);" icon="lucide:bell" class="text-base"></iconify-icon>
              </div>
            </button>
            <div style="background-color: rgba(59, 155, 255, 1); box-shadow: 0 0 12px rgba(59, 155, 255, 0.6);" class="flex absolute -top-1 -right-1 justify-center items-center w-6 h-6 rounded-full">
              <span style="color: rgba(255, 255, 255, 1);" class="text-xs font-semibold">3</span>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <img style="border-color: color-mix(in oklab, #3B9BFF 40%, transparent);" alt="User avatar" src="/assets/img/7a0a1950-3737-4f36-8af2-43b0b6314b81.jpeg" class="w-10 h-10 object-cover border-[2px] border-solid rounded-full">
            <button class="bg-transparent flex justify-center items-center w-6 h-6">
              <div class="flex justify-center items-center w-4 h-4">
                <iconify-icon style="color: color-mix(in oklab, #fff 70%, transparent);" icon="lucide:chevron-down" class="text-sm"></iconify-icon>
              </div>
            </button>
          </div>
        </div>
      </div>
    </header>

    <div class="flex w-full flex-1 min-h-0 relative">
      <!-- Main Navigation Sidebar -->
      <aside 
        :class="[
          'shrink-0 transition-all duration-300 z-20 flex flex-col border-r border-white/5 relative',
          isSidebarCollapsed ? 'w-20' : 'w-[260px]'
        ]"
        style="background: linear-gradient(180deg, rgba(20, 30, 50, 0.75) 0%, rgba(25, 40, 65, 0.75) 100%);"
      >
        <!-- Toggle Button -->
        <button 
          @click="toggleSidebar" 
          class="absolute -right-3 top-6 z-30 w-6 h-6 bg-[#1A2536] border border-[#3B9BFF]/30 rounded-full flex items-center justify-center text-[#3B9BFF] hover:bg-[#25344a] shadow-lg transition-transform duration-300"
          :class="isSidebarCollapsed ? 'rotate-180' : ''"
        >
          <iconify-icon icon="lucide:chevron-left" class="text-xs"></iconify-icon>
        </button>

        <div :class="['flex flex-col p-4 pt-10 overflow-y-auto overflow-x-hidden custom-scrollbar h-full', isSidebarCollapsed ? 'items-center' : '']">
          
          <!-- Primary Menu -->
          <div class="flex flex-col gap-1 w-full">
            <router-link v-for="item in primaryNavItems" :key="item.path" :to="item.path" 
              class="flex items-center gap-3 rounded-xl transition-all duration-200 group relative w-full"
              :style="isActive(item.path) ? {
                backgroundColor: 'color-mix(in oklab, #3B9BFF 20%, transparent)',
                padding: '0.75rem 1rem',
                borderColor: 'color-mix(in oklab, #3B9BFF 50%, transparent)',
                borderWidth: '1px',
                borderStyle: 'solid'
              } : {
                padding: '0.75rem 1rem'
              }"
              :class="[!isActive(item.path) ? 'hover:bg-white/5' : '', isSidebarCollapsed ? 'justify-center' : '']"
            >
              <div class="flex justify-center items-center w-5 h-5 shrink-0">
                <iconify-icon :style="{ color: isActive(item.path) ? 'rgba(95, 180, 255, 1)' : 'color-mix(in oklab, #fff 70%, transparent)' }" :icon="item.icon" class="text-base"></iconify-icon>
              </div>
              <span v-if="!isSidebarCollapsed" :style="{ color: isActive(item.path) ? 'color-mix(in oklab, #fff 95%, transparent)' : 'color-mix(in oklab, #fff 70%, transparent)' }" class="text-sm whitespace-nowrap">{{ item.name }}</span>
              <div v-if="isSidebarCollapsed" class="absolute left-full ml-4 px-2 py-1 bg-[#1A2536] text-white text-xs rounded opacity-0 group-hover:opacity-100 pointer-events-none transition-opacity z-50 whitespace-nowrap border border-white/10 shadow-xl">
                {{ item.name }}
              </div>
            </router-link>
          </div>

          <!-- Secondary Menu (Connected directly to primary) -->
          <div class="flex flex-col gap-1 w-full pt-4 border-t border-white/10 mt-4">
            <div v-if="!isSidebarCollapsed" class="px-4 py-2 text-[10px] text-white/20 font-bold uppercase tracking-widest">系统管理</div>
            
            <router-link v-for="item in secondaryNavItems" :key="item.path" :to="item.path" 
              class="flex items-center gap-3 rounded-xl transition-all duration-200 group relative w-full"
              :style="isActive(item.path) ? {
                backgroundColor: 'color-mix(in oklab, #3B9BFF 20%, transparent)',
                padding: '0.75rem 1rem',
                borderColor: 'color-mix(in oklab, #3B9BFF 50%, transparent)',
                borderWidth: '1px',
                borderStyle: 'solid'
              } : {
                padding: '0.75rem 1rem'
              }"
              :class="[!isActive(item.path) ? 'hover:bg-white/5' : '', isSidebarCollapsed ? 'justify-center' : '']"
            >
              <div class="flex justify-center items-center w-5 h-5 shrink-0">
                <iconify-icon :style="{ color: isActive(item.path) ? 'rgba(95, 180, 255, 1)' : 'color-mix(in oklab, #fff 70%, transparent)' }" :icon="item.icon" class="text-base"></iconify-icon>
              </div>
              <span v-if="!isSidebarCollapsed" :style="{ color: isActive(item.path) ? 'color-mix(in oklab, #fff 95%, transparent)' : 'color-mix(in oklab, #fff 70%, transparent)' }" class="text-xs whitespace-nowrap">{{ item.name }}</span>
              <div v-if="isSidebarCollapsed" class="absolute left-full ml-4 px-2 py-1 bg-[#1A2536] text-white text-xs rounded opacity-0 group-hover:opacity-100 pointer-events-none transition-opacity z-50 whitespace-nowrap border border-white/10 shadow-xl">
                {{ item.name }}
              </div>
            </router-link>
          </div>

        </div>
      </aside>

      <!-- Main Content Slot -->
      <main class="flex-1 min-h-0 flex flex-col relative bg-[#0F1928]">
        <slot></slot>
      </main>
    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 155, 255, 0.1); border-radius: 10px; }
</style>
