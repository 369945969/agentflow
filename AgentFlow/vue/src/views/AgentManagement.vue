<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import BaseLayout from '../components/BaseLayout.vue'
import { API_BASE_URL } from '../config/api'

// Layout state: 'grid' or 'list'
const viewLayout = ref<'grid' | 'list'>('grid')

// Search state
const searchQuery = ref('')

// Models state
const models = ref<any[]>([])
const fetchModels = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/models/`)
    const data = await response.json()
    models.value = data
  } catch (error) {
    console.error('Failed to fetch models:', error)
  }
}

// Agents state
const agents = ref<any[]>([])
const fetchAgents = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/agents/`)
    const data = await response.json()
    agents.value = data
  } catch (error) {
    console.error('Failed to fetch agents:', error)
  }
}

const filteredAgents = computed(() => {
  if (!searchQuery.value) return agents.value
  const query = searchQuery.value.toLowerCase()
  return agents.value.filter(agent => 
    agent.name.toLowerCase().includes(query) || 
    agent.description.toLowerCase().includes(query)
  )
})

onMounted(() => {
  fetchModels()
  fetchAgents()
})

// Create/Edit Agent state
const isCreateModalOpen = ref(false)
const isEditing = ref(false)
const newAgent = ref<any>({
  name: '',
  description: '',
  model: '',
  system_prompt: ''
})

const toggleLayout = (type: 'grid' | 'list') => {
  viewLayout.value = type
}

const openCreateModal = () => {
  isEditing.value = false
  newAgent.value = {
    name: '',
    description: '',
    model: '',
    system_prompt: ''
  }
  // Set default model
  const defaultModel = models.value.find((m: any) => m.is_default)
  if (defaultModel) {
    newAgent.value.model = defaultModel.id
  } else if (models.value.length > 0) {
    newAgent.value.model = models.value[0].id
  }
  isCreateModalOpen.value = true
}

const openEditModal = (agent: any) => {
  isEditing.value = true
  newAgent.value = { ...agent }
  isCreateModalOpen.value = true
}

const closeCreateModal = () => {
  isCreateModalOpen.value = false
}

const handleSave = async () => {
  if (isEditing.value) {
    await handleUpdate()
  } else {
    await handleCreate()
  }
}

const handleCreate = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/agents/`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newAgent.value)
    })
    if (response.ok) {
      await fetchAgents()
      closeCreateModal()
    }
  } catch (error) {
    console.error('Failed to create agent:', error)
  }
}

const handleUpdate = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/agents/${newAgent.value.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newAgent.value)
    })
    if (response.ok) {
      await fetchAgents()
      closeCreateModal()
    }
  } catch (error) {
    console.error('Failed to update agent:', error)
  }
}

const handleDelete = async (id: string) => {
  if (!confirm('确定要删除这个Agent吗?')) return
  try {
    const response = await fetch(`${API_BASE_URL}/api/agents/${id}`, {
      method: 'DELETE'
    })
    if (response.ok) {
      await fetchAgents()
    }
  } catch (error) {
    console.error('Failed to delete agent:', error)
  }
}

const getModelName = (modelId: string) => {
  const model = models.value.find(m => m.id === modelId)
  return model ? model.name : modelId
}
</script>

<template>
  <BaseLayout>
    <main class="overflow-x-hidden flex flex-col grow shrink gap-y-8 p-8 relative">
      
      <!-- Header -->
      <div class="flex justify-between items-center">
         <h1 style="color: color-mix( in oklab , #fff 95% , transparent );" class="text-2xl font-semibold">数字员工</h1>
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
        <div class="flex items-center gap-6">
          <div class="flex-1 bg-white/5 border border-white/10 rounded-xl px-4 py-3 flex items-center gap-3 backdrop-blur-md">
            <iconify-icon icon="lucide:search" class="text-white/40"></iconify-icon>
            <input v-model="searchQuery" type="text" placeholder="按名称/描述搜索Agent..." class="bg-transparent border-none outline-none text-sm text-white/70 w-full">
          </div>
        </div>
      </div>

      <!-- Agent Content -->
      <div :class="[viewLayout === 'grid' ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6' : 'flex flex-col gap-4']">
        <div 
          v-for="agent in filteredAgents" :key="agent.id"
          class="group p-6 rounded-2xl border border-white/10 bg-white/5 backdrop-blur-xl hover:shadow-2xl hover:border-[#3B9BFF]/30 transition-all duration-300"
          :class="viewLayout === 'list' ? 'flex items-center gap-6' : ''"
        >
          <div class="flex items-center gap-4" :class="viewLayout === 'list' ? 'w-1/4' : 'mb-4'">
            <div class="w-12 h-12 rounded-xl flex items-center justify-center shrink-0 bg-[#3B9BFF]/20 text-[#3B9BFF]">
              <iconify-icon icon="lucide:user" class="text-2xl"></iconify-icon>
            </div>
            <div>
              <h3 class="text-white/95 font-semibold text-lg">{{ agent.name }}</h3>
              <div class="flex items-center gap-2 mt-1">
                <span class="px-2 py-0.5 rounded text-[10px] bg-white/10 text-white/50">{{ getModelName(agent.model) }}</span>
                <span class="text-[10px] text-[#50C878]">已激活</span>
              </div>
            </div>
          </div>
          
          <p class="text-sm text-white/50 flex-1" :class="viewLayout === 'grid' ? 'mb-6 line-clamp-2' : ''">{{ agent.description }}</p>
          
          <div class="flex items-center gap-3" :class="viewLayout === 'grid' ? 'justify-between' : 'w-48 justify-end'">
            <span v-if="viewLayout === 'grid'" class="text-[10px] text-white/30">Created at {{ new Date(agent.created_at).toLocaleDateString() }}</span>
            <div class="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
              <button @click.stop="openEditModal(agent)" class="w-8 h-8 rounded-lg bg-white/5 border border-white/10 flex items-center justify-center text-white/60 hover:text-white hover:bg-[#3B9BFF]/20" title="编辑"><iconify-icon icon="lucide:edit-3"></iconify-icon></button>
              <button @click.stop="handleDelete(agent.id)" class="w-8 h-8 rounded-lg bg-white/5 border border-white/10 flex items-center justify-center text-white/40 hover:text-red-400 hover:bg-red-500/10 transition-all" title="删除"><iconify-icon icon="lucide:trash-2"></iconify-icon></button>
            </div>
          </div>
        </div>
      </div>

      <!-- Create New Agent Modal (Centered Modal) -->
      <Transition 
        enter-active-class="transition duration-300 ease-out"
        enter-from-class="opacity-0 scale-95"
        enter-to-class="opacity-100 scale-100"
        leave-active-class="transition duration-200 ease-in"
        leave-from-class="opacity-100 scale-100"
        leave-to-class="opacity-0 scale-95"
      >
        <div 
          v-if="isCreateModalOpen"
          class="fixed inset-0 z-[120] flex items-center justify-center p-6"
        >
          <div @click="closeCreateModal" class="absolute inset-0 bg-black/60 backdrop-blur-md"></div>
          <div class="relative w-full max-w-2xl max-h-[90vh] bg-[#1A2536] border border-white/10 rounded-3xl shadow-2xl overflow-hidden flex flex-col">
            <div class="p-6 border-b border-white/10 flex justify-between items-center bg-white/5">
              <h2 class="text-xl font-bold text-white/95">{{ isEditing ? "编辑Agent" : "创建新Agent" }}</h2>
              <button @click="closeCreateModal" class="w-8 h-8 rounded-full flex items-center justify-center hover:bg-white/10 text-white/40"><iconify-icon icon="lucide:x" class="text-xl"></iconify-icon></button>
            </div>

            <div class="flex-1 overflow-y-auto p-8 space-y-6 custom-scrollbar">
              <div class="space-y-4">
                <label class="text-xs text-white/40 block font-medium uppercase tracking-widest">Agent名称</label>
                <input v-model="newAgent.name" type="text" placeholder="例如：智能客服专家" class="w-full bg-black/20 border border-white/10 rounded-xl px-4 py-3 text-white outline-none focus:border-[#3B9BFF]/50 transition-all">
              </div>

              <div class="space-y-4">
                <label class="text-xs text-white/40 block font-medium uppercase tracking-widest">角色描述</label>
                <textarea v-model="newAgent.description" placeholder="简单描述这个Agent的功能..." rows="3" class="w-full bg-black/20 border border-white/10 rounded-xl px-4 py-3 text-white outline-none focus:border-[#3B9BFF]/50 transition-all resize-none"></textarea>
              </div>

              <div class="space-y-4">
                <label class="text-xs text-white/40 block font-medium uppercase tracking-widest">绑定基础模型</label>
                <div class="grid grid-cols-2 gap-3">
                  <div v-for="m in models" :key="m.id" 
                    @click="newAgent.model = m.id"
                    class="p-4 border rounded-2xl cursor-pointer transition-all flex flex-col gap-1 group relative overflow-hidden"
                    :class="newAgent.model === m.id ? 'border-[#3B9BFF] bg-[#3B9BFF]/10' : 'border-white/5 bg-white/5 hover:border-white/20'"
                  >
                    <div class="flex items-center justify-between relative z-10">
                      <span class="text-sm font-semibold" :class="newAgent.model === m.id ? 'text-[#3B9BFF]' : 'text-white/80'">{{ m.name }}</span>
                      <iconify-icon v-if="newAgent.model === m.id" icon="lucide:check-circle-2" class="text-[#3B9BFF]"></iconify-icon>
                    </div>
                    <div class="text-[10px] text-white/40 relative z-10">{{ m.provider }} · {{ m.model_name }}</div>
                    <div v-if="m.is_default" class="absolute top-0 right-0 px-2 py-0.5 bg-[#3B9BFF]/20 text-[#3B9BFF] text-[8px] font-bold uppercase rounded-bl-lg">Default</div>
                  </div>
                </div>
              </div>

              <div class="space-y-4">
                <label class="text-xs text-white/40 block font-medium uppercase tracking-widest">系统提示词 (System Prompt)</label>
                <textarea v-model="newAgent.system_prompt" placeholder="设置Agent的行为指令..." rows="5" class="w-full bg-black/40 border border-white/10 rounded-xl px-4 py-3 text-sm text-white/80 outline-none focus:border-[#3B9BFF]/50 transition-all font-mono"></textarea>
              </div>
            </div>

            <div class="p-6 border-t border-white/10 bg-white/5 flex gap-4">
              <button @click="closeCreateModal" class="flex-1 py-3 rounded-xl border border-white/10 text-white/60 hover:bg-white/5 transition-all text-sm">取消</button>
              <button @click="handleSave" class="flex-1 py-3 rounded-xl bg-[#3B9BFF] text-white font-bold shadow-[0_0_20px_rgba(59,155,255,0.4)] hover:bg-[#2A7FDB] transition-all text-sm">{{ isEditing ? "更新Agent" : "确认创建" }}</button>
            </div>
          </div>
        </div>
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
