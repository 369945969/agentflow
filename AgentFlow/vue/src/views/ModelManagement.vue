<script setup lang="ts">
import { ref, onMounted } from 'vue'
import BaseLayout from '../components/BaseLayout.vue'
import { API_BASE_URL } from '../config/api'

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

onMounted(() => {
  fetchModels()
})


const isCreateModalOpen = ref(false)
const newModel = ref({ id: '', name: '', provider: 'OpenAI', base_url: '', api_key: '', model_name: '', description: '', is_default: false })

const openCreateModal = () => {
  isCreateModalOpen.value = true
  isEditing.value = false
  newModel.value = { id: '', name: '', provider: 'OpenAI', base_url: '', api_key: '', model_name: '', description: '', is_default: false }
}

const closeCreateModal = () => {
  isCreateModalOpen.value = false
  newModel.value = { id: '', name: '', provider: 'OpenAI', base_url: '', api_key: '', model_name: '', description: '', is_default: false }
}


const isEditing = ref(false)

const openEditModal = (model: any) => {
  newModel.value = { ...model }
  isCreateModalOpen.value = true
  isEditing.value = true
}

const handleSetDefault = async (model: any) => {
  try {
    const updatedModel = { ...model, is_default: true }
    const response = await fetch(`${API_BASE_URL}/api/models/${model.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(updatedModel)
    })
    if (response.ok) {
      await fetchModels()
    }
  } catch (error) {
    console.error('Failed to set default model:', error)
  }
}

const handleSave = async () => {
  if (isEditing.value) {
    // Update
    try {
      const response = await fetch(`${API_BASE_URL}/api/models/${newModel.value.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newModel.value)
      })
      if (response.ok) {
        await fetchModels()
        closeCreateModal()
      }
    } catch (error) {
      console.error('Failed to update model:', error)
    }
  } else {
    // Create
    await handleCreate()
  }
}

const handleCreate = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/models/`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newModel.value)
    })
    if (response.ok) {
      await fetchModels()
      closeCreateModal()
    }
  } catch (error) {
    console.error('Failed to create model:', error)
  }
}


const handleDelete = async (id: string) => {
  if (!confirm('确定要删除这个模型吗?')) return
  try {
    const response = await fetch(`${API_BASE_URL}/api/models/${id}`, {
      method: 'DELETE'
    })
    if (response.ok) {
      models.value = models.value.filter(m => m.id !== id)
    }
  } catch (error) {
    console.error('Failed to delete model:', error)
  }
}
</script>

<template>
  <BaseLayout>
    <main class="overflow-x-hidden flex flex-col grow shrink gap-y-8 gap-x-8 p-8">
      <div class="flex justify-between items-center mb-8">
        <h1 class="text-2xl font-semibold text-white/95">模型管理</h1>
        <button @click="openCreateModal" class="bg-[#3B9BFF] hover:bg-[#2A7FDB] text-white px-6 py-2.5 rounded-xl flex items-center gap-2 shadow-[0_0_20px_rgba(59,155,255,0.3)] transition-all">
          <iconify-icon icon="lucide:plus" class="text-sm"></iconify-icon>
          {{ isEditing ? "编辑模型" : "添加新模型" }}
        </button>
      </div>

      <!-- Dynamic Models Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="model in models" :key="model.id" 
          class="p-6 border-[1px] border-solid rounded-2xl bg-[#1A2536]/40 border-white/10 hover:border-[#3B9BFF]/40 hover:bg-[#1A2536]/70 transition-all duration-300 flex flex-col justify-between"
        >
          <div class="flex justify-between items-start mb-4">
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-xl flex items-center justify-center text-xl bg-white/5 relative">
                <iconify-icon icon="lucide:cpu" class="text-[#3B9BFF]"></iconify-icon>
                <div v-if="model.is_default" class="absolute -top-1 -right-1 w-4 h-4 bg-[#3B9BFF] rounded-full flex items-center justify-center border-2 border-[#1A2536]">
                  <iconify-icon icon="lucide:check" class="text-[8px] text-white"></iconify-icon>
                </div>
              </div>
              <div>
                <div class="flex items-center gap-2">
                  <h3 class="text-lg font-semibold text-white/95">{{ model.name }}</h3>
                  <span v-if="model.is_default" class="text-[10px] px-2 py-0.5 rounded-md bg-[#3B9BFF]/20 text-[#3B9BFF] font-medium uppercase tracking-wider">Default</span>
                </div>
                <div class="flex items-center gap-2 mt-1">
                  <span class="text-xs px-2 py-0.5 rounded-full bg-white/5 text-white/60">{{ model.provider }}</span>
                </div>
              </div>
            </div>
            <div class="flex gap-1">
              <button v-if="!model.is_default" @click.stop="handleSetDefault(model)" class="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-yellow-500/10 text-white/40 hover:text-yellow-400 transition-all" title="设为默认">
                <iconify-icon icon="lucide:star" class="text-sm"></iconify-icon>
              </button>
              <button @click.stop="openEditModal(model)" class="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-[#3B9BFF]/10 text-white/40 hover:text-[#3B9BFF] transition-all">
                <iconify-icon icon="lucide:edit-3" class="text-sm"></iconify-icon>
              </button>
              <button @click.stop="handleDelete(model.id)" class="w-8 h-8 flex items-center justify-center rounded-lg hover:bg-red-500/10 text-white/40 hover:text-red-400 transition-all">
                <iconify-icon icon="lucide:trash-2" class="text-sm"></iconify-icon>
              </button>
            </div>
          </div>
          <p class="text-sm text-white/50 mb-6 line-clamp-2">{{ model.description }}</p>
          <div class="text-xs text-white/30 font-mono bg-black/20 p-2 rounded-lg">{{ model.model_name }}</div>
        </div>
      </div>
    
      <!-- Drawer 5: Create Model Modal -->
      <Transition 
        enter-active-class="transition duration-300 ease-out" enter-from-class="opacity-0 scale-95" enter-to-class="opacity-100 scale-100"
        leave-active-class="transition duration-200 ease-in" leave-from-class="opacity-100 scale-100" leave-to-class="opacity-0 scale-95"
      >
        <div v-if="isCreateModalOpen" class="fixed inset-0 z-[120] flex items-center justify-center p-6">
          <div @click="closeCreateModal" class="absolute inset-0 bg-black/60 backdrop-blur-md"></div>
          <div class="relative w-full max-w-lg bg-[#1A2536] border border-white/10 rounded-3xl shadow-2xl overflow-hidden flex flex-col">
            <div class="p-6 border-b border-white/10 flex justify-between items-center bg-white/5">
              <h2 class="text-xl font-bold text-white/95">{{ isEditing ? "编辑模型" : "添加新模型" }}</h2>
              <button @click="closeCreateModal" class="w-8 h-8 rounded-full flex items-center justify-center hover:bg-white/10 text-white/40"><iconify-icon icon="lucide:x" class="text-xl"></iconify-icon></button>
            </div>
            <div class="p-8 space-y-4">
              <input v-model="newModel.name" placeholder="模型显示名称" class="w-full bg-black/20 border border-white/10 rounded-xl px-4 h-[50px] text-sm text-white outline-none">
              <select v-model="newModel.provider" class="w-full bg-black/20 border border-white/10 rounded-xl px-4 py-3 h-[50px] text-sm text-white outline-none">
                <option class="bg-[#1A2536] text-white text-white/50" value="" disabled selected>选择模型提供商</option>
                <option class="bg-[#1A2536] text-white" value="OpenAI">OpenAI</option>
                <option class="bg-[#1A2536] text-white" value="Anthropic">Anthropic</option>
                <option class="bg-[#1A2536] text-white" value="DeepSeek">DeepSeek</option>
                <option class="bg-[#1A2536] text-white" value="Moonshot">Moonshot (Kimi)</option>
                <option class="bg-[#1A2536] text-white" value="Qwen">Qwen (Alibaba)</option>
                <option class="bg-[#1A2536] text-white" value="Ollama">Ollama (Local)</option>
                <option class="bg-[#1A2536] text-white" value="Other">其他</option>
              </select>
              <input v-model="newModel.base_url" placeholder="API Base URL (e.g., https://api.openai.com/v1)" class="w-full bg-black/20 border border-white/10 rounded-xl px-4 h-[50px] text-sm text-white outline-none">
              <input v-model="newModel.model_name" placeholder="模型标识 (e.g. gpt-4)" class="w-full bg-black/20 border border-white/10 rounded-xl px-4 h-[50px] text-sm text-white outline-none">
              <input v-model="newModel.api_key" placeholder="API Key" type="password" class="w-full bg-black/20 border border-white/10 rounded-xl px-4 h-[50px] text-sm text-white outline-none">
              <textarea v-model="newModel.description" placeholder="模型描述..." class="w-full bg-black/20 border border-white/10 rounded-xl px-4 h-[50px] text-sm text-white outline-none h-24"></textarea>
              
              <div class="flex items-center gap-3 p-4 bg-white/5 rounded-xl border border-white/10 cursor-pointer hover:bg-white/10 transition-all" @click="newModel.is_default = !newModel.is_default">
                <div class="w-10 h-6 rounded-full p-1 transition-all duration-300 relative" :class="newModel.is_default ? 'bg-[#3B9BFF]' : 'bg-white/10'">
                  <div class="w-4 h-4 bg-white rounded-full shadow-md transition-all duration-300" :class="newModel.is_default ? 'translate-x-4' : 'translate-x-0'"></div>
                </div>
                <div>
                  <div class="text-sm font-medium text-white/90">设为默认模型</div>
                  <div class="text-xs text-white/40">设置为全局默认使用的语言模型</div>
                </div>
              </div>
            </div>
            <div class="p-6 border-t border-white/10 bg-white/5 flex gap-3">
              <button @click="closeCreateModal" class="flex-1 py-3 rounded-xl border border-white/10 text-sm text-white/60 hover:bg-white/5">取消</button>
              <button @click="handleSave" class="flex-1 py-3 rounded-xl bg-[#3B9BFF] text-white text-sm font-bold shadow-[0_0_20px_rgba(59,155,255,0.4)]">{{ isEditing ? "更新模型" : "保存模型" }}</button>
            </div>
          </div>
        </div>
      </Transition>

    </main>
  </BaseLayout>
</template>
