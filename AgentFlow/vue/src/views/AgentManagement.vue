<script setup lang="ts">
import { ref } from 'vue'
import BaseLayout from '../components/BaseLayout.vue'

// Layout state: 'grid' or 'list'
const viewLayout = ref<'grid' | 'list'>('grid')

// Create Agent state
const isCreateModalOpen = ref(false)
const newAgent = ref({
  name: '',
  description: '',
  model: 'GPT-4',
  tags: [] as string[]
})

const agents = ref([
  { id: 1, name: '智能客服助手', desc: '专业的客户服务自动化助手，能够处理常见问题咨询。', status: '已激活', model: 'GPT-4', icon: 'lucide:message-circle', color: '#3B9BFF' },
  { id: 2, name: '数据分析专家', desc: '专业的数据分析和可视化助手，擅长处理各种数据集。', status: '已激活', model: 'Claude-3', icon: 'lucide:bar-chart', color: '#87B4FF' },
  { id: 3, name: '代码生成助手', desc: '智能代码生成和优化助手，支持多种编程语言。', status: '草稿', model: 'GPT-4', icon: 'lucide:code', color: '#00E5FF' },
])

const toggleLayout = (type: 'grid' | 'list') => {
  viewLayout.value = type
}

const openCreateModal = () => {
  isCreateModalOpen.value = true
}

const closeCreateModal = () => {
  isCreateModalOpen.value = false
}
</script>

<template>
  <BaseLayout>
    <main class="overflow-x-hidden flex flex-col grow shrink gap-y-8 p-8 relative">
      
      <!-- Header -->
      <div class="flex justify-between items-center">
        <h1 style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-2xl font-semibold">Agents</h1>
        <div class="flex items-center gap-4">
          <!-- Create Button -->
          <button 
            @click="openCreateModal"
            class="hover:bg-[#2A7FDB] flex items-center gap-2 rounded-xl transition-all" 
            style="background-color: rgba(59, 155, 255, 1); box-shadow: 0 0 20px rgba(59, 155, 255, 0.4); padding: 0.6rem 1.5rem;"
          >
            <iconify-icon icon="lucide:plus" class="text-white text-sm"></iconify-icon>
            <span class="text-white text-sm font-medium whitespace-nowrap">创建新Agent</span>
          </button>

          <!-- Hidden Import/Export (as requested) -->
          <!-- <button class="...">批量导入</button> -->
          <!-- <button class="...">批量导出</button> -->

          <!-- Layout Toggle -->
          <div class="flex items-center p-1 bg-white/5 border border-white/10 rounded-xl backdrop-blur-xl">
            <button 
              @click="toggleLayout('grid')"
              :class="['w-10 h-10 flex items-center justify-center rounded-lg transition-all', viewLayout === 'grid' ? 'bg-[#3B9BFF]/20 text-[#3B9BFF]' : 'text-white/40 hover:text-white/60']"
            >
              <iconify-icon icon="lucide:grid-3x3" class="text-lg"></iconify-icon>
            </button>
            <button 
              @click="toggleLayout('list')"
              :class="['w-10 h-10 flex items-center justify-center rounded-lg transition-all', viewLayout === 'list' ? 'bg-[#3B9BFF]/20 text-[#3B9BFF]' : 'text-white/40 hover:text-white/60']"
            >
              <iconify-icon icon="lucide:list" class="text-lg"></iconify-icon>
            </button>
          </div>
        </div>
      </div>

      <!-- Filters & Search -->
      <div style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); border-color: color-mix( in oklab , #3B9BFF 30% , transparent );" class="p-6 border-[1px] border-solid rounded-2xl">
        <div class="flex items-center gap-6 mb-6">
          <div class="flex-1 bg-white/5 border border-white/10 rounded-xl px-4 py-2 flex items-center gap-3 backdrop-blur-md">
            <iconify-icon icon="lucide:search" class="text-white/40"></iconify-icon>
            <input type="text" placeholder="按名称/标签/技能搜索Agent..." class="bg-transparent border-none outline-none text-sm text-white/70 w-full">
          </div>
          <div class="flex items-center gap-3">
            <select class="bg-white/5 border border-white/10 rounded-xl px-4 py-2 text-sm text-white/60 outline-none">
              <option>全部状态</option>
              <option>已激活</option>
              <option>草稿</option>
            </select>
            <select class="bg-white/5 border border-white/10 rounded-xl px-4 py-2 text-sm text-white/60 outline-none">
              <option>全部模型</option>
              <option>GPT-4</option>
              <option>Claude-3</option>
            </select>
          </div>
        </div>
        <div class="flex gap-2">
           <span v-for="tag in ['全部', '对话助手', '数据分析', '代码生成']" :key="tag" 
             class="px-4 py-1.5 rounded-full text-xs border border-white/10 cursor-pointer transition-all hover:bg-white/10"
             :class="tag === '全部' ? 'bg-[#3B9BFF]/20 border-[#3B9BFF]/50 text-[#3B9BFF]' : 'text-white/50 bg-white/5'"
           >{{ tag }}</span>
        </div>
      </div>

      <!-- Agent Content -->
      <div :class="[viewLayout === 'grid' ? 'grid grid-cols-3 gap-6' : 'flex flex-col gap-4']">
        <div 
          v-for="agent in agents" :key="agent.id"
          class="group p-6 rounded-2xl border border-white/10 bg-white/5 backdrop-blur-xl hover:shadow-2xl hover:border-[#3B9BFF]/30 transition-all duration-300"
          :class="viewLayout === 'list' ? 'flex items-center gap-6' : ''"
        >
          <div class="flex items-center gap-4" :class="viewLayout === 'list' ? 'w-1/4' : 'mb-4'">
            <div class="w-12 h-12 rounded-xl flex items-center justify-center shrink-0" :style="{backgroundColor: agent.color+'20', color: agent.color}">
              <iconify-icon :icon="agent.icon" class="text-2xl"></iconify-icon>
            </div>
            <div>
              <h3 class="text-white/95 font-semibold text-lg">{{ agent.name }}</h3>
              <div class="flex items-center gap-2 mt-1">
                <span class="px-2 py-0.5 rounded text-[10px] bg-white/10 text-white/50">{{ agent.model }}</span>
                <span class="text-[10px] text-[#50C878]">{{ agent.status }}</span>
              </div>
            </div>
          </div>
          
          <p class="text-sm text-white/50 flex-1" :class="viewLayout === 'grid' ? 'mb-6 line-clamp-2' : ''">{{ agent.desc }}</p>
          
          <div class="flex items-center gap-3" :class="viewLayout === 'grid' ? 'justify-between' : 'w-48 justify-end'">
            <span v-if="viewLayout === 'grid'" class="text-[10px] text-white/30">Modified 2h ago</span>
            <div class="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
              <button class="w-8 h-8 rounded-lg bg-white/5 border border-white/10 flex items-center justify-center text-white/60 hover:text-white hover:bg-[#3B9BFF]/20"><iconify-icon icon="lucide:edit-3"></iconify-icon></button>
              <button class="w-8 h-8 rounded-lg bg-white/5 border border-white/10 flex items-center justify-center text-white/60 hover:text-white hover:bg-[#3B9BFF]/20"><iconify-icon icon="lucide:play"></iconify-icon></button>
            </div>
          </div>
        </div>
      </div>

      <!-- Create New Agent Modal (Sidebar Style) -->
      <Transition 
        enter-active-class="transition duration-300 ease-out"
        enter-from-class="translate-x-full"
        enter-to-class="translate-x-0"
        leave-active-class="transition duration-200 ease-in"
        leave-from-class="translate-x-0"
        leave-to-class="translate-x-full"
      >
        <aside 
          v-if="isCreateModalOpen"
          class="fixed top-0 right-0 w-[500px] h-full z-[100] bg-[#1A2536]/95 backdrop-blur-3xl border-l border-[#3B9BFF]/30 shadow-[-20px_0_50px_rgba(0,0,0,0.5)] flex flex-col"
        >
          <div class="flex justify-between items-center p-6 border-b border-white/10">
            <h2 class="text-xl font-bold text-white/95">创建新Agent</h2>
            <button @click="closeCreateModal" class="w-10 h-10 flex items-center justify-center text-white/40 hover:text-white hover:bg-white/5 rounded-full"><iconify-icon icon="lucide:x" class="text-xl"></iconify-icon></button>
          </div>

          <div class="flex-1 overflow-y-auto p-8 space-y-8 custom-scrollbar">
            <!-- Icon/Name -->
            <div class="flex items-center gap-6">
              <div class="w-20 h-20 bg-[#3B9BFF]/10 border-2 border-dashed border-[#3B9BFF]/30 rounded-2xl flex flex-col items-center justify-center text-[#3B9BFF] hover:bg-[#3B9BFF]/20 cursor-pointer transition-all">
                <iconify-icon icon="lucide:image-plus" class="text-2xl mb-1"></iconify-icon>
                <span class="text-[10px]">选择图标</span>
              </div>
              <div class="flex-1 space-y-4">
                <div>
                  <label class="text-xs text-white/40 block mb-2 font-medium uppercase tracking-widest">Agent名称</label>
                  <input type="text" placeholder="给你的Agent起个名字..." class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white outline-none focus:border-[#3B9BFF]/50 transition-all">
                </div>
              </div>
            </div>

            <!-- Description -->
            <div class="space-y-4">
              <label class="text-xs text-white/40 block mb-2 font-medium uppercase tracking-widest">角色描述</label>
              <textarea placeholder="描述该Agent的主要职责和专长..." rows="4" class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-white outline-none focus:border-[#3B9BFF]/50 transition-all resize-none"></textarea>
            </div>

            <!-- Model Selection -->
            <div class="space-y-4">
              <label class="text-xs text-white/40 block mb-2 font-medium uppercase tracking-widest">绑定基础模型</label>
              <div class="grid grid-cols-2 gap-4">
                <div v-for="m in ['GPT-4 Turbo', 'Claude-3 Opus', 'Llama-3 (Local)', 'Gemini Pro']" :key="m" 
                  class="p-4 bg-white/5 border rounded-xl cursor-pointer hover:bg-[#3B9BFF]/10 transition-all"
                  :class="m === 'GPT-4 Turbo' ? 'border-[#3B9BFF]/50 bg-[#3B9BFF]/10' : 'border-white/5'"
                >
                  <div class="text-sm text-white/90 font-medium">{{ m }}</div>
                  <div class="text-[10px] text-white/30">Stable & High Performance</div>
                </div>
              </div>
            </div>

            <!-- System Prompt -->
            <div class="space-y-4">
              <div class="flex justify-between items-center">
                <label class="text-xs text-white/40 font-medium uppercase tracking-widest">系统提示词 (System Prompt)</label>
                <button class="text-[10px] text-[#3B9BFF] hover:underline">使用模板</button>
              </div>
              <textarea placeholder="你是一个专业的..." rows="6" class="w-full bg-[#0F1928] border border-white/10 rounded-xl px-4 py-3 text-sm text-white/70 outline-none focus:border-[#3B9BFF]/50 transition-all font-mono"></textarea>
            </div>
          </div>

          <!-- Footer Actions -->
          <div class="p-6 border-t border-white/10 bg-white/5 backdrop-blur-md flex gap-4">
            <button @click="closeCreateModal" class="flex-1 py-3 rounded-xl border border-white/10 text-white/60 hover:bg-white/5 transition-all">取消</button>
            <button class="flex-1 py-3 rounded-xl bg-[#3B9BFF] text-white font-bold shadow-[0_0_20px_rgba(59,155,255,0.4)] hover:bg-[#2A7FDB] transition-all">确认创建</button>
          </div>
        </aside>
      </Transition>

      <!-- Overlay for modal -->
      <div v-if="isCreateModalOpen" @click="closeCreateModal" class="fixed inset-0 bg-black/40 backdrop-blur-sm z-[90]"></div>

    </main>
  </BaseLayout>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(59, 155, 255, 0.2);
  border-radius: 10px;
}
</style>
