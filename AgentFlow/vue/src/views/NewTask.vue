<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import BaseLayout from '../components/BaseLayout.vue'
import { API_BASE_URL } from '../config/api'

const router = useRouter()
const agents = ref<any[]>([])
const taskPrompt = ref('')
const selectedAgentId = ref('')

const fetchAgents = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/agents/`)
    const data = await response.json()
    agents.value = data
    if (data.length > 0) {
      selectedAgentId.value = data[0].id
    }
  } catch (error) {
    console.error('Failed to fetch agents:', error)
  }
}

const handleLaunchTask = () => {
  if (!taskPrompt.value.trim() || !selectedAgentId.value) return
  
  // 逻辑：跳转到单人对话，并自动填充并发送消息
  // 这里我们通过 query 参数传递信息
  router.push({
    path: '/single-person-chat',
    query: {
      agentId: selectedAgentId.value,
      initialPrompt: taskPrompt.value
    }
  })
}

const quickPrompts = [
  '分析市场竞争对手情况',
  '生成本周工作周报草稿',
  '编写一个 Python 爬虫脚本',
  '优化现有的系统架构设计图'
]

onMounted(fetchAgents)
</script>

<template>
  <BaseLayout>
    <div class="w-full h-full flex items-center justify-center bg-[#0F1928] p-6 relative overflow-hidden">
      <!-- Background Decorative Elements -->
      <div class="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-[#3B9BFF]/10 blur-[120px] rounded-full"></div>
      <div class="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] bg-[#50C878]/10 blur-[120px] rounded-full"></div>

      <div class="w-full max-w-3xl flex flex-col items-center z-10 animate-in fade-in zoom-in duration-500">
        <!-- Main Input Container -->
        <div class="w-full bg-white/5 border border-white/10 rounded-[32px] p-3 shadow-2xl backdrop-blur-xl focus-within:border-[#3B9BFF]/50 focus-within:shadow-[0_0_50px_rgba(59,155,255,0.1)] transition-all">
          <textarea
            v-model="taskPrompt"
            placeholder="描述你想让 AI 智能体执行的任务..."
            class="w-full bg-transparent border-none outline-none text-2xl text-white/90 p-6 min-h-[220px] resize-none placeholder:text-white/5 leading-relaxed"
          ></textarea>
          
          <div class="flex items-center justify-between p-4 bg-white/5 rounded-[24px] border border-white/5">
            <div class="flex items-center gap-4">
              <!-- Agent Selector -->
              <div class="relative group">
                <select v-model="selectedAgentId" class="bg-black/40 border border-white/10 rounded-xl px-5 py-2.5 text-sm text-white outline-none hover:border-[#3B9BFF]/50 transition-all cursor-pointer appearance-none pr-10">
                  <option v-for="agent in agents" :key="agent.id" :value="agent.id" class="bg-[#1A2536]">
                    {{ agent.name }}
                  </option>
                </select>
                <iconify-icon icon="lucide:chevron-down" class="absolute right-3 top-1/2 -translate-y-1/2 text-white/30 pointer-events-none group-hover:text-[#3B9BFF] transition-colors"></iconify-icon>
              </div>
              
              <div class="h-5 w-px bg-white/10"></div>
              
              <span class="text-[10px] text-white/20 uppercase tracking-[0.2em] font-bold">Intelligent Dispatching</span>
            </div>

            <button 
              @click="handleLaunchTask"
              :disabled="!taskPrompt.trim() || !selectedAgentId"
              class="h-14 px-10 bg-[#3B9BFF] hover:bg-[#2A7FDB] disabled:opacity-30 disabled:cursor-not-allowed text-white rounded-2xl font-bold shadow-xl transition-all flex items-center gap-3 group"
            >
              <span class="text-base">启动任务</span>
              <iconify-icon icon="lucide:arrow-right" class="text-xl group-hover:translate-x-1 transition-transform"></iconify-icon>
            </button>
          </div>
        </div>

        <!-- Quick Suggestions -->
        <div class="mt-12 flex flex-col items-center gap-4">
          <p class="text-[10px] text-white/20 uppercase tracking-[0.2em] font-bold">快速开始</p>
          <div class="flex flex-wrap justify-center gap-3">
            <button 
              v-for="sug in quickPrompts" 
              :key="sug"
              @click="taskPrompt = sug"
              class="px-5 py-2.5 rounded-full bg-white/5 border border-white/5 text-xs text-white/40 hover:text-white hover:bg-[#3B9BFF]/20 hover:border-[#3B9BFF]/30 transition-all"
            >
              {{ sug }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </BaseLayout>
</template>

<style scoped>
:deep(main) {
  padding: 0 !important;
}

textarea::placeholder {
  font-weight: 300;
}

/* 简单的进入动画 */
.animate-in {
  animation: fadeInZoom 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fadeInZoom {
  from {
    opacity: 0;
    transform: scale(0.95) translateY(10px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}
</style>
