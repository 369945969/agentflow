<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import BaseLayout from '../components/BaseLayout.vue'
import { API_BASE_URL } from '../config/api'

const historyTasks = ref<any[]>([])
const searchQuery = ref('')

const fetchHistoryTasks = async () => {
  try {
    // In a real scenario, this would call a dedicated history API
    const response = await fetch(`${API_BASE_URL}/api/tasks/history`)
    if (response.ok) {
      historyTasks.value = await response.json()
    } else {
      // Mock data
      historyTasks.value = [
        { id: 'task-1', name: '分析市场竞争', agent: '分析专家', model: 'GPT-4', status: 'completed', time: '2024-03-22 10:30' },
        { id: 'task-2', name: '生成周报草稿', agent: '客服助手', model: 'DeepSeek', status: 'completed', time: '2024-03-21 15:45' }
      ]
    }
  } catch (error) {
    console.error('Failed to fetch history:', error)
    historyTasks.value = [
      { id: 'task-1', name: '分析市场竞争', agent: '分析专家', model: 'GPT-4', status: 'completed', time: '2024-03-22 10:30' },
      { id: 'task-2', name: '生成周报草稿', agent: '客服助手', model: 'DeepSeek', status: 'completed', time: '2024-03-21 15:45' }
    ]
  }
}

const filteredTasks = computed(() => {
  if (!searchQuery.value) return historyTasks.value
  const query = searchQuery.value.toLowerCase()
  return historyTasks.value.filter(t => 
    t.name.toLowerCase().includes(query) || 
    t.agent.toLowerCase().includes(query)
  )
})

onMounted(fetchHistoryTasks)
</script>

<template>
  <BaseLayout>
    <main class="overflow-x-hidden flex flex-col grow shrink gap-y-8 p-8 relative">
      <div>
        <h1 class="text-3xl font-bold text-white mb-2">历史任务</h1>
        <p class="text-white/50 text-sm">查看以往执行的任务记录和结果详情。</p>
      </div>

      <!-- Search Area -->
      <div style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); border-color: color-mix( in oklab , #3B9BFF 30% , transparent );" class="p-6 border-[1px] border-solid rounded-2xl shrink-0">
        <div class="bg-white/5 border border-white/10 rounded-xl px-4 py-3 flex items-center gap-3 backdrop-blur-md w-full">
          <iconify-icon icon="lucide:search" class="text-white/40"></iconify-icon>
          <input v-model="searchQuery" type="text" placeholder="搜索任务名称或执行人..." class="bg-transparent border-none outline-none text-sm text-white/70 w-full">
        </div>
      </div>

      <div class="flex-1 overflow-y-auto custom-scrollbar pr-2">
        <div class="grid grid-cols-1 gap-4">
          <div v-for="task in filteredTasks" :key="task.id" 
            class="p-6 bg-[#1A2536]/40 border border-white/10 rounded-2xl flex items-center justify-between group hover:border-[#3B9BFF]/30 transition-all"
          >
            <div class="flex items-center gap-6 flex-1">
              <div class="w-12 h-12 rounded-xl bg-white/5 flex items-center justify-center text-[#3B9BFF]">
                <iconify-icon icon="lucide:clipboard-check" class="text-2xl"></iconify-icon>
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="text-white/90 font-semibold text-lg truncate">{{ task.name }}</h3>
                <div class="flex items-center gap-3 mt-1">
                  <span class="text-[10px] text-white/30 uppercase tracking-widest font-mono">{{ task.id }}</span>
                  <div class="h-3 w-px bg-white/10"></div>
                  <span class="text-xs text-white/50 flex items-center gap-1">
                    <iconify-icon icon="lucide:user" class="text-[10px]"></iconify-icon>
                    {{ task.agent }}
                  </span>
                </div>
              </div>
            </div>

            <div class="flex items-center gap-10 text-right shrink-0">
              <div>
                <div class="text-[10px] text-white/20 uppercase font-bold tracking-widest mb-1">执行模型</div>
                <div class="text-xs text-white/60 font-mono">{{ task.model }}</div>
              </div>
              <div>
                <div class="text-[10px] text-white/20 uppercase font-bold tracking-widest mb-1">执行时间</div>
                <div class="text-xs text-white/40">{{ task.time }}</div>
              </div>
              <button class="w-10 h-10 rounded-xl bg-[#3B9BFF]/10 text-[#3B9BFF] flex items-center justify-center opacity-0 group-hover:opacity-100 transition-all">
                <iconify-icon icon="lucide:chevron-right" class="text-xl"></iconify-icon>
              </button>
            </div>
          </div>
        </div>
      </div>
    </main>
  </BaseLayout>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 155, 255, 0.2); border-radius: 10px; }
</style>
