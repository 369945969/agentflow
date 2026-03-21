<script setup lang="ts">
import { ref } from 'vue'
import BaseLayout from '../components/BaseLayout.vue'

const isLeftSidebarOpen = ref(true)
const isRightSidebarOpen = ref(true)

const toggleLeftSidebar = () => {
  isLeftSidebarOpen.value = !isLeftSidebarOpen.value
}

const toggleRightSidebar = () => {
  isRightSidebarOpen.value = !isRightSidebarOpen.value
}
</script>

<template>
  <BaseLayout>
    <div class="flex w-full h-full overflow-hidden relative">
      
      <!-- Toggle Left Button -->
      <button 
        v-if="!isLeftSidebarOpen"
        @click="toggleLeftSidebar"
        class="absolute left-4 top-4 z-20 w-8 h-8 bg-[#1A2536] border border-[#3B9BFF]/30 rounded-lg flex items-center justify-center text-[#3B9BFF] shadow-lg"
      >
        <iconify-icon icon="lucide:chevron-right"></iconify-icon>
      </button>

      <!-- Left Sidebar: Conversations -->
      <aside 
        :class="[
          'shrink-0 border-r border-white/10 transition-all duration-300 overflow-hidden flex flex-col',
          isLeftSidebarOpen ? 'w-[320px]' : 'w-0 border-none'
        ]"
        style="background-color: color-mix( in oklab , #fff 3% , transparent );"
      >
        <div class="w-[320px] flex flex-col h-full">
          <div class="flex justify-between items-center p-6 border-b border-white/5">
            <h2 class="text-lg font-semibold text-white/95">Conversations</h2>
            <button @click="toggleLeftSidebar" class="text-white/50 hover:text-white">
              <iconify-icon icon="lucide:chevron-left"></iconify-icon>
            </button>
          </div>
          
          <div class="p-4 space-y-4 overflow-y-auto custom-scrollbar flex-1">
             <div style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(12px); padding: 0.5rem 0.75rem; border-color: color-mix( in oklab , #3B9BFF 30% , transparent );" class="flex items-center border-[1px] border-solid rounded-xl">
                <iconify-icon icon="lucide:search" class="text-sm text-white/50 mr-2"></iconify-icon>
                <input type="text" placeholder="搜索会话..." class="text-sm bg-transparent grow outline-none text-white/70">
              </div>

              <div class="space-y-3">
                <div v-for="i in 3" :key="i" :class="['p-3 border rounded-xl cursor-pointer transition-colors', i === 1 ? 'bg-[#3B9BFF]/10 border-[#3B9BFF]/40' : 'bg-white/5 border-white/10 hover:bg-white/10']">
                  <div class="flex justify-between items-start mb-2">
                    <h3 class="text-sm font-medium text-white/95">客户服务优化讨论</h3>
                    <span class="text-[10px] text-white/40">2m ago</span>
                  </div>
                  <p class="text-xs text-white/50 line-clamp-2">我们需要分析当前的客户反馈数据...</p>
                </div>
              </div>
          </div>
        </div>
      </aside>

      <!-- Main Chat Area -->
      <main class="flex-1 flex flex-col min-w-0 bg-[#0F1928] overflow-hidden">
        <!-- Chat Header -->
        <div class="h-16 border-b border-white/10 flex items-center justify-between px-6 bg-white/5 backdrop-blur-md shrink-0">
          <div class="flex items-center gap-4">
            <input value="客户服务优化讨论" class="bg-transparent border-none outline-none text-lg font-semibold text-white/95">
          </div>
          <div class="flex items-center gap-3">
            <button class="flex items-center gap-2 px-3 py-1.5 bg-white/5 border border-white/10 rounded-lg text-xs text-white/70 hover:bg-white/10">
              <iconify-icon icon="lucide:trash-2"></iconify-icon>
              清空对话
            </button>
          </div>
        </div>

        <!-- Scrollable Messages -->
        <div class="flex-1 overflow-y-auto p-6 space-y-6 custom-scrollbar">
           <!-- Participant List (Floating style) -->
           <div class="flex flex-wrap gap-2 mb-6">
              <div v-for="a in ['客服专员', '数据分析师']" :key="a" class="flex items-center gap-2 px-3 py-1 bg-[#3B9BFF]/10 border border-[#3B9BFF]/30 rounded-full">
                <div class="w-2 h-2 bg-[#3B9BFF] rounded-full"></div>
                <span class="text-xs text-white/80">{{ a }}</span>
              </div>
           </div>

           <!-- Messages -->
           <div class="flex gap-4 max-w-[80%]">
              <div class="w-10 h-10 rounded-full bg-[#3B9BFF]/20 flex items-center justify-center text-[#3B9BFF] shrink-0"><iconify-icon icon="lucide:user"></iconify-icon></div>
              <div class="space-y-1">
                <div class="bg-white/10 p-4 rounded-2xl rounded-tl-none text-sm text-white/90 border border-white/5">
                  我们需要分析最近一个月的客户投诉数据，特别是关于响应时间和解决效率的反馈。
                </div>
                <span class="text-[10px] text-white/40 ml-1">14:32</span>
              </div>
           </div>

           <div class="flex gap-4 max-w-[80%] ml-auto flex-row-reverse">
              <div class="w-10 h-10 rounded-full bg-[#50C878]/20 flex items-center justify-center text-[#50C878] shrink-0"><iconify-icon icon="lucide:bot"></iconify-icon></div>
              <div class="space-y-1 text-right">
                <div class="bg-[#3B9BFF]/20 p-4 rounded-2xl rounded-tr-none text-sm text-white/90 border border-[#3B9BFF]/30">
                  我已经完成了最近30天的投诉数据分析。平均响应时间为2.3小时，较上月增长了12%。
                </div>
                <span class="text-[10px] text-white/40 mr-1">14:35</span>
              </div>
           </div>
        </div>

        <!-- Chat Input -->
        <div class="p-6 border-t border-white/10 bg-white/5 backdrop-blur-md shrink-0">
          <div class="flex items-center gap-3 mb-3">
             <button class="p-2 hover:bg-white/5 rounded-lg text-white/50"><iconify-icon icon="lucide:paperclip"></iconify-icon></button>
             <button class="px-3 py-1.5 bg-[#3B9BFF]/20 border border-[#3B9BFF]/50 rounded-lg text-xs text-[#3B9BFF] flex items-center gap-2">
               <iconify-icon icon="lucide:play"></iconify-icon>
               触发编排
             </button>
          </div>
          <div class="flex gap-4">
            <textarea placeholder="输入消息..." rows="1" class="flex-1 bg-white/5 border border-white/10 rounded-xl p-3 text-sm text-white/90 outline-none focus:border-[#3B9BFF]/50 resize-none"></textarea>
            <button class="w-12 h-12 bg-[#3B9BFF] rounded-xl flex items-center justify-center text-white shadow-[0_0_20px_rgba(59,155,255,0.3)]"><iconify-icon icon="lucide:send" class="text-xl"></iconify-icon></button>
          </div>
        </div>
      </main>

      <!-- Right Sidebar: Context -->
      <aside 
        :class="[
          'shrink-0 border-l border-white/10 transition-all duration-300 overflow-hidden flex flex-col',
          isRightSidebarOpen ? 'w-[380px]' : 'w-0 border-none'
        ]"
        style="background-color: color-mix( in oklab , #fff 3% , transparent );"
      >
        <div class="w-[380px] flex flex-col h-full">
          <div class="flex justify-between items-center p-6 border-b border-white/5">
            <button @click="toggleRightSidebar" class="text-white/50 hover:text-white">
              <iconify-icon icon="lucide:chevron-right"></iconify-icon>
            </button>
            <h3 class="text-lg font-semibold text-white/95">Context & Memory</h3>
          </div>
          
          <div class="p-6 overflow-y-auto custom-scrollbar flex-1 space-y-6">
             <div class="flex p-1 bg-white/5 rounded-xl border border-white/10">
                <button class="flex-1 py-1.5 text-xs bg-[#3B9BFF]/20 text-white rounded-lg border border-[#3B9BFF]/50">参与者</button>
                <button class="flex-1 py-1.5 text-xs text-white/50 hover:text-white">对话记忆</button>
             </div>

             <div class="space-y-4">
                <div v-for="a in ['客服专员', '数据分析师']" :key="a" class="p-4 bg-white/5 border border-white/10 rounded-xl">
                  <div class="flex gap-3 items-center">
                    <div class="w-10 h-10 rounded-lg bg-[#3B9BFF]/20 flex items-center justify-center text-[#3B9BFF]"><iconify-icon icon="lucide:bot"></iconify-icon></div>
                    <div class="flex-1">
                      <div class="text-sm font-medium text-white/95">{{ a }}</div>
                      <div class="text-[10px] text-white/50">GPT-4 Turbo</div>
                    </div>
                    <div class="w-2 h-2 bg-[#50C878] rounded-full"></div>
                  </div>
                </div>
             </div>
          </div>
        </div>
      </aside>

      <!-- Toggle Right Button -->
      <button 
        v-if="!isRightSidebarOpen"
        @click="toggleRightSidebar"
        class="absolute right-4 top-4 z-20 w-8 h-8 bg-[#1A2536] border border-[#3B9BFF]/30 rounded-lg flex items-center justify-center text-[#3B9BFF] shadow-lg"
      >
        <iconify-icon icon="lucide:chevron-left"></iconify-icon>
      </button>

    </div>
  </BaseLayout>
</template>

<style scoped>
:deep(main) {
  padding: 0 !important;
}

.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(59, 155, 255, 0.2);
  border-radius: 10px;
}

textarea {
  scrollbar-width: none;
}
</style>
