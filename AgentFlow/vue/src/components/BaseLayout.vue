<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

// State for main navigation sidebar - Default collapsed, sync with localStorage
const isSidebarCollapsed = ref(true)
const isAvatarDropdownOpen = ref(false)

// Frequently used items (Top section)
const primaryNavItems = [
  { name: '新建任务', icon: 'lucide:zap', path: '/new-task' },
  { name: '单人对话', icon: 'lucide:user', path: '/single-person-chat' },
  { name: '多人对话', icon: 'lucide:message-square', path: '/multi-person-chat' },
  { name: '编排画布', icon: 'lucide:workflow', path: '/orchestration-list' },
  { name: 'AI智能体', icon: 'lucide:bot', path: '/agent-management' },
]

// History items
const historyNavItems = [
  { name: '历史任务', icon: 'lucide:history', path: '/history-tasks' },
]

// System Management items (Moved to Avatar Dropdown)
const systemNavItems = [
  { name: '系统监控', icon: 'lucide:layout-dashboard', path: '/' },
  { name: '模型管理', icon: 'lucide:cpu', path: '/model-management' },
  { name: '技能管理', icon: 'lucide:zap', path: '/skill-management' },
  { name: '定时任务', icon: 'lucide:calendar', path: '/schedule-tasks' },
  { name: '执行历史', icon: 'lucide:scroll-text', path: '/execution-monitor' },
]

const isActive = (path: string) => route.path === path

const toggleSidebar = () => {
  isSidebarCollapsed.value = !isSidebarCollapsed.value
  localStorage.setItem('sidebar-collapsed', isSidebarCollapsed.value.toString())
}

const toggleAvatarDropdown = (event: Event) => {
  event.stopPropagation()
  isAvatarDropdownOpen.value = !isAvatarDropdownOpen.value
}

const closeDropdowns = () => {
  isAvatarDropdownOpen.value = false
}

onMounted(() => {
  // Restore sidebar state from localStorage
  const savedState = localStorage.getItem('sidebar-collapsed')
  if (savedState !== null) {
    isSidebarCollapsed.value = savedState === 'true'
  }
  
  window.addEventListener('click', closeDropdowns)
})

onUnmounted(() => {
  window.removeEventListener('click', closeDropdowns)
})
</script>

<template>
  <div class="font-[-apple-system,BlinkMacSystemFont,'Segoe UI'] w-full h-screen overflow-hidden flex flex-col" style="line-height: 1.5; background: rgba(15, 25, 40, 1);">
    <!-- Header -->
    <header style="background: linear-gradient(180deg, rgba(15, 25, 45, 0.85) 0%, rgba(20, 35, 60, 0.85) 100%);" class="w-full shrink-0 z-30 border-b border-white/5">
      <div style="padding: 1rem 2rem;" class="flex justify-between items-center">
        <div class="flex items-center gap-10">
          <div class="flex items-center gap-3">
            <div class="flex justify-center items-center w-6 h-6">
              <iconify-icon style="color: rgba(59, 155, 255, 1);" icon="lucide:workflow" class="text-xl"></iconify-icon>
            </div>
            <span style="color: color-mix(in oklab, #fff 95%, transparent);" class="text-xl font-semibold whitespace-nowrap">质量专家Agent Team</span>
          </div>
        </div>
        
        <div class="flex items-center gap-4">
          <!-- User Avatar & Dropdown -->
          <div class="relative">
            <div 
              @click="toggleAvatarDropdown"
              class="flex items-center gap-3 cursor-pointer p-1 rounded-2xl hover:bg-white/5 transition-all group"
            >
              <img style="border-color: color-mix(in oklab, #3B9BFF 40%, transparent);" alt="User avatar" src="/assets/img/7a0a1950-3737-4f36-8af2-43b0b6314b81.jpeg" class="w-10 h-10 object-cover border-[2px] border-solid rounded-full shadow-lg group-hover:scale-105 transition-transform">
              <button class="bg-transparent flex justify-center items-center w-6 h-6">
                <div class="flex justify-center items-center w-4 h-4">
                  <iconify-icon 
                    style="color: color-mix(in oklab, #fff 70%, transparent);" 
                    :icon="isAvatarDropdownOpen ? 'lucide:chevron-up' : 'lucide:chevron-down'" 
                    class="text-sm"
                  ></iconify-icon>
                </div>
              </button>
            </div>

            <!-- Dropdown Menu -->
            <Transition
              enter-active-class="transition duration-200 ease-out"
              enter-from-class="opacity-0 scale-95 -translate-y-2"
              enter-to-class="opacity-100 scale-100 translate-y-0"
              leave-active-class="transition duration-150 ease-in"
              leave-from-class="opacity-100 scale-100 translate-y-0"
              leave-to-class="opacity-0 scale-95 -translate-y-2"
            >
              <div 
                v-if="isAvatarDropdownOpen"
                class="absolute right-0 mt-3 w-64 bg-[#1A2536]/95 backdrop-blur-2xl border border-white/10 rounded-2xl shadow-2xl overflow-hidden z-50 p-2"
                @click.stop
              >
                <div class="px-4 py-3 border-b border-white/5 mb-2">
                  <p class="text-xs font-bold text-white/20 uppercase tracking-widest">系统管理</p>
                </div>
                <div class="flex flex-col gap-1">
                  <router-link 
                    v-for="item in systemNavItems" 
                    :key="item.path" 
                    :to="item.path"
                    @click="isAvatarDropdownOpen = false"
                    class="flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 group"
                    :class="isActive(item.path) ? 'bg-[#3B9BFF]/20 border border-[#3B9BFF]/30' : 'hover:bg-white/5'"
                  >
                    <div class="flex justify-center items-center w-5 h-5">
                      <iconify-icon 
                        :style="{ color: isActive(item.path) ? 'rgba(95, 180, 255, 1)' : 'color-mix(in oklab, #fff 50%, transparent)' }" 
                        :icon="item.icon" 
                        class="text-base group-hover:scale-110 transition-transform"
                      ></iconify-icon>
                    </div>
                    <span 
                      class="text-sm"
                      :style="{ color: isActive(item.path) ? 'color-mix(in oklab, #fff 95%, transparent)' : 'color-mix(in oklab, #fff 70%, transparent)' }"
                    >{{ item.name }}</span>
                  </router-link>
                </div>
                <div class="mt-2 pt-2 border-t border-white/5 px-2 pb-1">
                  <button class="w-full flex items-center gap-3 px-4 py-3 rounded-xl hover:bg-red-500/10 text-red-400 transition-all group">
                    <iconify-icon icon="lucide:log-out" class="text-base"></iconify-icon>
                    <span class="text-sm font-medium">退出登录</span>
                  </button>
                </div>
              </div>
            </Transition>
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
              <span v-if="!isSidebarCollapsed" :style="{ color: isActive(item.path) ? 'color-mix(in oklab, #fff 95%, transparent)' : 'color-mix(in oklab, #fff 70%, transparent)' }" class="text-sm font-medium whitespace-nowrap">{{ item.name }}</span>
              <div v-if="isSidebarCollapsed" class="absolute left-full ml-4 px-2 py-1 bg-[#1A2536] text-white text-xs rounded opacity-0 group-hover:opacity-100 pointer-events-none transition-opacity z-50 whitespace-nowrap border border-white/10 shadow-xl">
                {{ item.name }}
              </div>
            </router-link>
          </div>

          <!-- History Menu (Task Records) -->
          <div class="flex flex-col gap-1 w-full pt-4 border-t border-white/10 mt-4">
            <div v-if="!isSidebarCollapsed" class="px-4 py-2 text-[10px] text-white/20 font-bold uppercase tracking-widest">任务记录</div>
            
            <router-link v-for="item in historyNavItems" :key="item.path" :to="item.path" 
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
              <span v-if="!isSidebarCollapsed" :style="{ color: isActive(item.path) ? 'color-mix(in oklab, #fff 95%, transparent)' : 'color-mix(in oklab, #fff 70%, transparent)' }" class="text-sm font-medium whitespace-nowrap">{{ item.name }}</span>
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
