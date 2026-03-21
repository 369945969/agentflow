<script setup lang="ts">
import { ref, onMounted } from 'vue'
import BaseLayout from '../components/BaseLayout.vue'
import { API_BASE_URL } from '../config/api'

// State
const models = ref<any[]>([])
const selectedModelId = ref('')
const thinkingEnabled = ref(true)      // 默认开启
const simplifiedOutput = ref(false)    // 默认关闭 (原“查看中间结果”)

const fetchModels = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/models/`)
    const data = await response.json()
    models.value = data
    // Default select
    const defaultModel = data.find((m: any) => m.is_default)
    if (defaultModel) {
      selectedModelId.value = defaultModel.id
    } else if (data.length > 0) {
      selectedModelId.value = data[0].id
    }
  } catch (error) {
    console.error('Failed to fetch models:', error)
  }
}

onMounted(fetchModels)
</script>

<template>
  <BaseLayout>
    <div class="flex w-full h-full overflow-hidden bg-[#0F1928]">
      
      <!-- Main Chat Area -->
      <main class="flex-1 flex flex-col min-w-0 overflow-hidden">
        <!-- Chat Header -->
        <div class="h-16 border-b border-white/10 flex items-center px-6 bg-white/5 backdrop-blur-md shrink-0">
          <h2 class="text-lg font-semibold text-white/95">单人对话</h2>
        </div>

        <!-- Scrollable Messages -->
        <div class="flex-1 overflow-y-auto p-6 space-y-6 custom-scrollbar">
           <div class="flex flex-col items-center justify-center h-full text-white/30 text-sm">
             <iconify-icon icon="lucide:user" class="text-4xl mb-4"></iconify-icon>
             开始你的单人对话...
           </div>
        </div>

        <!-- Chat Input -->
        <div class="p-6 border-t border-white/10 bg-white/5 backdrop-blur-md shrink-0">
          <div class="flex gap-4">
            <textarea placeholder="输入消息..." rows="1" class="flex-1 bg-white/5 border border-white/10 rounded-xl p-3 text-sm text-white/90 outline-none focus:border-[#3B9BFF]/50 resize-none"></textarea>
            <button class="w-12 h-12 bg-[#3B9BFF] rounded-xl flex items-center justify-center text-white shadow-[0_0_20px_rgba(59,155,255,0.3)]"><iconify-icon icon="lucide:send" class="text-xl"></iconify-icon></button>
          </div>
        </div>
      </main>

      <!-- Right Sidebar: Controls -->
      <aside class="w-[380px] shrink-0 border-l border-white/10 bg-[#1A2536] flex flex-col">
        <div class="p-6 border-b border-white/5">
          <h3 class="text-lg font-semibold text-white/95">会话设置</h3>
        </div>
        
        <div class="p-6 space-y-8 overflow-y-auto custom-scrollbar flex-1">
          <!-- Model Selection -->
          <div class="space-y-3">
            <label class="text-xs text-white/40 font-bold uppercase tracking-widest">选择模型</label>
            <select v-model="selectedModelId" class="w-full bg-white/5 border border-white/10 rounded-xl px-4 py-3 text-sm text-white outline-none">
              <option v-for="m in models" :key="m.id" :value="m.id">{{ m.name }}</option>
            </select>
          </div>

          <!-- Toggles -->
          <div class="space-y-4">
            <div class="flex items-center justify-between p-3 bg-white/5 rounded-xl border border-white/10 cursor-pointer" @click="thinkingEnabled = !thinkingEnabled">
              <span class="text-sm text-white/90">开启 Thinking 模式</span>
              <div class="w-10 h-6 rounded-full p-1 transition-all" :class="thinkingEnabled ? 'bg-[#3B9BFF]' : 'bg-white/10'">
                <div class="w-4 h-4 bg-white rounded-full transition-all" :class="thinkingEnabled ? 'translate-x-4' : 'translate-x-0'"></div>
              </div>
            </div>
            <div class="flex items-center justify-between p-3 bg-white/5 rounded-xl border border-white/10 cursor-pointer" @click="simplifiedOutput = !simplifiedOutput">
              <span class="text-sm text-white/90">简化输出模式</span>
              <div class="w-10 h-6 rounded-full p-1 transition-all" :class="simplifiedOutput ? 'bg-[#3B9BFF]' : 'bg-white/10'">
                <div class="w-4 h-4 bg-white rounded-full transition-all" :class="simplifiedOutput ? 'translate-x-4' : 'translate-x-0'"></div>
              </div>
            </div>
          </div>
        </div>
      </aside>
    </div>
  </BaseLayout>
</template>

<style scoped>
:deep(main) {
  padding: 0 !important;
}
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 155, 255, 0.2); border-radius: 10px; }
textarea { scrollbar-width: none; }
</style>
