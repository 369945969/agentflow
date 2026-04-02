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

// Skills state
const allSkills = ref<any[]>([])
const fetchAllSkills = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/skills/`)
    const data = await response.json()
    allSkills.value = data
  } catch (error) {
    console.error('Failed to fetch skills:', error)
  }
}

// Digital Employees (Agents) state
const agents = ref<any[]>([])
const fetchAgents = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/agents/`)
    const data = await response.json()
    agents.value = (data || []).map((agent: any) => ({
      ...agent,
      skills: normalizeSkills(agent.skills)
    }))
  } catch (error) {
    console.error('Failed to fetch digital employees:', error)
    agents.value = []
  }
}

const filteredAgents = computed(() => {
  if (!agents.value) return []
  if (!searchQuery.value) return agents.value
  const query = searchQuery.value.toLowerCase()
  return agents.value.filter(agent => 
    (agent.name && agent.name.toLowerCase().includes(query)) || 
    (agent.description && agent.description.toLowerCase().includes(query))
  )
})

onMounted(() => {
  fetchModels()
  fetchAgents()
  fetchAllSkills()
})

// Create/Edit state
const isCreateModalOpen = ref(false)
const isAICreateModalOpen = ref(false)
const isEditing = ref(false)
const isAIGenerating = ref(false)
const aiPrompt = ref('')
const currentEmployee = ref<any>({
  name: '',
  description: '',
  model: '',
  system_prompt: '',
  skills: []
})

const openAICreateModal = () => {
  aiPrompt.value = ''
  isAICreateModalOpen.value = true
}

const handleAICreate = async () => {
  if (!aiPrompt.value.trim()) return
  isAIGenerating.value = true
  
  try {
    const response = await fetch(`${API_BASE_URL}/api/agents/ai-generate`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ prompt: aiPrompt.value })
    })
    
    if (response.ok) {
      const suggestedConfig = await response.json()
      currentEmployee.value = {
        name: suggestedConfig.name || '',
        description: suggestedConfig.description || '',
        model: currentEmployee.value.model,
        system_prompt: suggestedConfig.system_prompt || '',
        skills: suggestedConfig.skills || []
      }
      isAICreateModalOpen.value = false
      isCreateModalOpen.value = true
    } else {
      alert('AI 解析失败，请重试')
    }
  } catch (error) {
    console.error('AI generate error:', error)
    alert('请求失败')
  } finally {
    isAIGenerating.value = false
  }
}

const normalizeSkills = (skills: unknown): string[] => {
  if (Array.isArray(skills)) {
    return [...new Set(skills.filter((skill): skill is string => typeof skill === 'string' && skill.trim().length > 0))]
  }
  if (typeof skills === 'string' && skills.trim().length > 0) {
    return [skills]
  }
  return []
}

const selectedSkills = computed(() => {
  const selected = new Set(normalizeSkills(currentEmployee.value.skills))
  return allSkills.value.filter(skill => selected.has(skill.id))
})

const toggleLayout = (type: 'grid' | 'list') => {
  viewLayout.value = type
}

const openCreateModal = () => {
  isEditing.value = false
  currentEmployee.value = {
    name: '',
    description: '',
    model: '',
    system_prompt: '',
    skills: []
  }
  
  if (models.value && models.value.length > 0) {
    const defaultModel = models.value.find((m: any) => m.is_default)
    if (defaultModel) {
      currentEmployee.value.model = defaultModel.id
    } else {
      currentEmployee.value.model = models.value[0].id
    }
  }
  isCreateModalOpen.value = true
}

const openEditModal = (agent: any) => {
  isEditing.value = true
  currentEmployee.value = {
    ...agent,
    skills: normalizeSkills(agent.skills)
  }
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
      body: JSON.stringify({
        ...currentEmployee.value,
        skills: normalizeSkills(currentEmployee.value.skills)
      })
    })
    if (response.ok) {
      await fetchAgents()
      closeCreateModal()
    }
  } catch (error) {
    console.error('Failed to create digital employee:', error)
  }
}

const handleUpdate = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/agents/${currentEmployee.value.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        ...currentEmployee.value,
        skills: normalizeSkills(currentEmployee.value.skills)
      })
    })
    if (response.ok) {
      await fetchAgents()
      closeCreateModal()
    }
  } catch (error) {
    console.error('Failed to update digital employee:', error)
  }
}

const handleDelete = async (id: string) => {
  if (!confirm('确定要删除这个AI智能体吗?')) return
  try {
    const response = await fetch(`${API_BASE_URL}/api/agents/${id}`, {
      method: 'DELETE'
    })
    if (response.ok) {
      await fetchAgents()
    }
  } catch (error) {
    console.error('Failed to delete digital employee:', error)
  }
}

const getModelName = (modelId: string) => {
  const model = models.value.find(m => m.id === modelId)
  return model ? model.name : (modelId || '未绑定模型')
}

const getSkillName = (skillId: string) => {
  const skill = allSkills.value.find(s => s.id === skillId)
  return skill ? skill.name : skillId
}

const toggleSkill = (skillId: string) => {
  currentEmployee.value.skills = normalizeSkills(currentEmployee.value.skills)
  const index = currentEmployee.value.skills.indexOf(skillId)
  if (index === -1) {
    currentEmployee.value.skills.push(skillId)
  } else {
    currentEmployee.value.skills.splice(index, 1)
  }
}
</script>

<template>
  <BaseLayout>
    <main class="overflow-x-hidden flex flex-col grow shrink gap-y-8 p-8 relative">
      <!-- Header Area -->
      <div class="flex justify-between items-center shrink-0">
        <div>
          <h1 class="text-3xl font-bold text-white mb-2">AI智能体管理</h1>
          <p class="text-white/50 text-sm">创建和管理您的企业AI智能体，为其分配模型和专业技能。</p>
        </div>
        <div class="flex gap-4">
          <div class="flex bg-white/5 rounded-xl p-1 border border-white/10 h-fit">
            <button @click="toggleLayout('grid')" :class="['p-2 rounded-lg transition-all', viewLayout === 'grid' ? 'bg-[#3B9BFF] text-white shadow-lg' : 'text-white/40 hover:text-white']">
              <iconify-icon icon="lucide:layout-grid"></iconify-icon>
            </button>
            <button @click="toggleLayout('list')" :class="['p-2 rounded-lg transition-all', viewLayout === 'list' ? 'bg-[#3B9BFF] text-white shadow-lg' : 'text-white/40 hover:text-white']">
              <iconify-icon icon="lucide:list"></iconify-icon>
            </button>
          </div>
          <button @click="openAICreateModal" class="px-6 py-3 bg-white/5 hover:bg-white/10 text-white rounded-xl font-bold border border-white/10 transition-all flex items-center gap-2 h-fit group">
            <iconify-icon icon="lucide:sparkles" class="text-[#3B9BFF] group-hover:animate-pulse"></iconify-icon>
            AI 创建智能体
          </button>
          <button @click="openCreateModal" class="px-6 py-3 bg-[#3B9BFF] hover:bg-[#2A7FDB] text-white rounded-xl font-bold shadow-[0_0_20px_rgba(59,155,255,0.4)] transition-all flex items-center gap-2 h-fit">
            <iconify-icon icon="lucide:plus"></iconify-icon>
            创建智能体
          </button>
        </div>
      </div>

      <!-- Search Bar -->
      <div style="background-color: color-mix( in oklab , #fff 5% , transparent ); backdrop-filter: blur(24px); border-color: color-mix( in oklab , #3B9BFF 30% , transparent );" class="p-6 border-[1px] border-solid rounded-2xl shrink-0">
        <div class="bg-white/5 border border-white/10 rounded-xl px-4 py-3 flex items-center gap-3 backdrop-blur-md w-full">
          <iconify-icon icon="lucide:search" class="text-white/40"></iconify-icon>
          <input v-model="searchQuery" type="text" placeholder="按名称/描述搜索AI智能体..." class="bg-transparent border-none outline-none text-sm text-white/70 w-full">
        </div>
      </div>

      <!-- List/Grid Content -->
      <div class="flex-1 min-h-0 overflow-y-auto custom-scrollbar pr-2">
        <div v-if="filteredAgents.length > 0" :class="[viewLayout === 'grid' ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6' : 'flex flex-col gap-4']">
          <div 
            v-for="agent in filteredAgents" :key="agent.id"
            class="group p-6 rounded-2xl border border-white/10 bg-[#1A2536]/40 backdrop-blur-xl hover:shadow-2xl hover:border-[#3B9BFF]/30 transition-all duration-300"
            :class="viewLayout === 'list' ? 'flex items-center gap-6' : 'flex flex-col'"
          >
            <div class="flex items-center gap-4" :class="viewLayout === 'list' ? 'w-1/4 shrink-0' : 'mb-4'">
              <div class="w-12 h-12 rounded-xl flex items-center justify-center shrink-0 bg-[#3B9BFF]/20 text-[#3B9BFF]">
                <iconify-icon icon="lucide:user" class="text-2xl"></iconify-icon>
              </div>
              <div class="min-w-0">
                <h3 class="text-white/95 font-semibold text-lg truncate">{{ agent.name }}</h3>
                <div class="flex items-center gap-2 mt-1">
                  <span class="px-2 py-0.5 rounded text-[10px] bg-white/10 text-white/50 truncate max-w-[100px]">{{ getModelName(agent.model) }}</span>
                  <span class="text-[10px] text-[#50C878] shrink-0">已激活</span>
                </div>
              </div>
            </div>
            
            <div class="flex-1 min-w-0">
              <p class="text-sm text-white/50 mb-4" :class="viewLayout === 'grid' ? 'line-clamp-2 h-10' : ''">{{ agent.description || '暂无职责描述' }}</p>
              <!-- Display Skills -->
              <div v-if="agent.skills && agent.skills.length > 0" class="flex flex-wrap gap-2 mb-4">
                <span v-for="skillId in agent.skills" :key="skillId" class="px-2 py-0.5 bg-[#3B9BFF]/10 text-[#3B9BFF] text-[10px] rounded-md border border-[#3B9BFF]/20 whitespace-nowrap">
                  {{ getSkillName(skillId) }}
                </span>
              </div>
            </div>
            
            <div class="flex items-center gap-3 shrink-0" :class="viewLayout === 'grid' ? 'justify-between pt-4 border-t border-white/5' : 'w-48 justify-end'">
              <span class="text-[10px] text-white/30">ID: {{ agent.id.substring(0, 8) }}...</span>
              <div class="flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                <button @click.stop="openEditModal(agent)" class="w-8 h-8 rounded-lg bg-white/5 border border-white/10 flex items-center justify-center text-white/60 hover:text-white hover:bg-[#3B9BFF]/20" title="编辑"><iconify-icon icon="lucide:edit-3"></iconify-icon></button>
                <button @click.stop="handleDelete(agent.id)" class="w-8 h-8 rounded-lg bg-white/5 border border-white/10 flex items-center justify-center text-white/40 hover:text-red-400 hover:bg-red-500/10 transition-all" title="删除"><iconify-icon icon="lucide:trash-2"></iconify-icon></button>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else class="py-32 flex flex-col items-center justify-center text-white/20 border-2 border-dashed border-white/5 rounded-3xl">
          <div class="w-20 h-20 bg-white/5 rounded-full flex items-center justify-center mb-6">
            <iconify-icon icon="lucide:user-x" class="text-5xl"></iconify-icon>
          </div>
          <p class="text-xl font-medium mb-2">暂无AI智能体</p>
          <p class="text-sm mb-8 text-white/10">点击上方按钮开始创建您的第一位AI智能体</p>
          <button @click="openCreateModal" class="px-6 py-2 border border-white/10 rounded-xl hover:bg-white/5 transition-all text-white/60">立即创建</button>
        </div>
      </div>

      <!-- Create/Edit Modal -->
      <Transition 
        enter-active-class="transition duration-300 ease-out"
        enter-from-class="opacity-0 scale-95"
        enter-to-class="opacity-100 scale-100"
        leave-active-class="transition duration-200 ease-in"
        leave-from-class="opacity-100 scale-100"
        leave-to-class="opacity-0 scale-95"
      >
        <div v-if="isCreateModalOpen" class="fixed inset-0 z-[120] flex items-center justify-center p-6">
          <div @click="closeCreateModal" class="absolute inset-0 bg-black/60 backdrop-blur-md"></div>
          <div class="relative w-full max-w-4xl max-h-[90vh] bg-[#1A2536] border border-white/10 rounded-3xl shadow-2xl overflow-hidden flex flex-col">
            <div class="p-6 border-b border-white/10 flex justify-between items-center bg-white/5">
              <h2 class="text-xl font-bold text-white/95">{{ isEditing ? "编辑AI智能体" : "创建AI智能体" }}</h2>
              <button @click="closeCreateModal" class="w-8 h-8 rounded-full flex items-center justify-center hover:bg-white/10 text-white/40"><iconify-icon icon="lucide:x" class="text-xl"></iconify-icon></button>
            </div>

            <div class="flex-1 overflow-y-auto p-8 space-y-8 custom-scrollbar">
              <div class="grid grid-cols-1 lg:grid-cols-2 gap-10">
                <!-- Left: Basic Info -->
                <div class="space-y-6">
                  <div class="space-y-4">
                    <label class="text-xs text-white/40 block font-bold uppercase tracking-widest">员工名称</label>
                    <input v-model="currentEmployee.name" type="text" placeholder="例如：智能客服专家" class="w-full bg-black/20 border border-white/10 rounded-xl px-4 py-3 text-white outline-none focus:border-[#3B9BFF]/50 transition-all">
                  </div>

                  <div class="space-y-4">
                    <label class="text-xs text-white/40 block font-bold uppercase tracking-widest">职责描述</label>
                    <textarea v-model="currentEmployee.description" placeholder="描述该AI智能体的主要职责..." rows="3" class="w-full bg-black/20 border border-white/10 rounded-xl px-4 py-3 text-white outline-none focus:border-[#3B9BFF]/50 transition-all resize-none"></textarea>
                  </div>

                  <div class="space-y-4">
                    <label class="text-xs text-white/40 block font-bold uppercase tracking-widest">基础模型</label>
                    <div class="grid grid-cols-1 gap-3 max-h-[220px] overflow-y-auto custom-scrollbar pr-2">
                      <div v-for="m in models" :key="m.id" 
                        @click="currentEmployee.model = m.id"
                        class="p-4 border rounded-2xl cursor-pointer transition-all flex flex-col gap-1 group relative overflow-hidden"
                        :class="currentEmployee.model === m.id ? 'border-[#3B9BFF] bg-[#3B9BFF]/10' : 'border-white/5 bg-white/5 hover:border-white/20'"
                      >
                        <div class="flex items-center justify-between relative z-10">
                          <span class="text-sm font-semibold" :class="currentEmployee.model === m.id ? 'text-[#3B9BFF]' : 'text-white/80'">{{ m.name }}</span>
                          <iconify-icon v-if="currentEmployee.model === m.id" icon="lucide:check-circle-2" class="text-[#3B9BFF]"></iconify-icon>
                        </div>
                        <div class="text-[10px] text-white/40 relative z-10">{{ m.provider }} · {{ m.model_name }}</div>
                        <div v-if="m.is_default" class="absolute top-0 right-0 px-2 py-0.5 bg-[#3B9BFF]/20 text-[#3B9BFF] text-[8px] font-bold uppercase rounded-bl-lg">Default</div>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Right: Skills -->
                <div class="space-y-6">
                  <div class="flex items-center justify-between gap-4">
                    <label class="text-xs text-white/40 block font-bold uppercase tracking-widest">绑定专业技能（可多选）</label>
                    <span class="text-xs text-[#3B9BFF]">{{ selectedSkills.length }} 项已选择</span>
                  </div>
                  <div v-if="selectedSkills.length > 0" class="flex flex-wrap gap-2">
                    <span v-for="skill in selectedSkills" :key="skill.id" class="px-3 py-1 rounded-full bg-[#3B9BFF]/10 text-[#3B9BFF] text-xs border border-[#3B9BFF]/20">
                      {{ skill.name }}
                    </span>
                  </div>
                  <div class="grid grid-cols-1 gap-2 max-h-[450px] overflow-y-auto custom-scrollbar pr-2">
                    <div v-for="skill in allSkills" :key="skill.id" 
                      @click="toggleSkill(skill.id)"
                      class="p-3 border rounded-xl cursor-pointer transition-all flex items-center gap-3"
                      :class="normalizeSkills(currentEmployee.skills).includes(skill.id) ? 'border-[#3B9BFF] bg-[#3B9BFF]/10' : 'border-white/5 bg-white/5 hover:border-white/10'"
                    >
                      <div class="w-5 h-5 rounded border flex items-center justify-center shrink-0"
                        :class="normalizeSkills(currentEmployee.skills).includes(skill.id) ? 'border-[#3B9BFF] bg-[#3B9BFF] text-white' : 'border-white/15 bg-black/20 text-transparent'">
                        <iconify-icon icon="lucide:check" class="text-xs"></iconify-icon>
                      </div>
                      <div class="w-8 h-8 rounded-lg flex items-center justify-center shrink-0" :style="{ backgroundColor: skill.color + '20', color: skill.color }">
                        <iconify-icon :icon="skill.icon || 'lucide:zap'"></iconify-icon>
                      </div>
                      <div class="flex-1 min-w-0 text-left">
                        <div class="text-sm font-medium text-white/90 truncate">{{ skill.name }}</div>
                        <div class="text-[10px] text-white/30 truncate uppercase">{{ skill.type }}</div>
                      </div>
                    </div>
                    <div v-if="allSkills.length === 0" class="text-center py-12 bg-white/5 rounded-2xl border border-white/5 border-dashed">
                      <p class="text-xs text-white/20">暂无可用技能</p>
                      <router-link to="/skill-management" class="text-[10px] text-[#3B9BFF] hover:underline mt-2 inline-block">前往添加</router-link>
                    </div>
                  </div>
                </div>
              </div>

              <!-- System Prompt -->
              <div class="space-y-4 pt-4 border-t border-white/5">
                <label class="text-xs text-white/40 block font-bold uppercase tracking-widest">系统提示词 (System Prompt)</label>
                <textarea v-model="currentEmployee.system_prompt" placeholder="设置该AI智能体的行为指令，例如：你是一个资深的 Java 开发工程师..." rows="5" class="w-full bg-black/40 border border-white/10 rounded-xl px-4 py-3 text-sm text-white/80 outline-none focus:border-[#3B9BFF]/50 transition-all font-mono"></textarea>
              </div>
            </div>

            <div class="p-6 border-t border-white/10 bg-white/5 flex gap-4">
              <button @click="closeCreateModal" class="flex-1 py-3 rounded-xl border border-white/10 text-white/60 hover:bg-white/5 transition-all text-sm">取消</button>
              <button @click="handleSave" class="flex-1 py-3 rounded-xl bg-[#3B9BFF] text-white font-bold shadow-[0_0_20px_rgba(59,155,255,0.4)] hover:bg-[#2A7FDB] transition-all text-sm">{{ isEditing ? "保存变更" : "确认创建" }}</button>
            </div>
          </div>
        </div>
      </Transition>

      <!-- AI Create Agent Modal (Manus Style) -->
      <Transition
        enter-active-class="transition duration-300 ease-out"
        enter-from-class="opacity-0 scale-105"
        enter-to-class="opacity-100 scale-100"
        leave-active-class="transition duration-200 ease-in"
        leave-from-class="opacity-100 scale-100"
        leave-to-class="opacity-0 scale-105"
      >
        <div v-if="isAICreateModalOpen" class="fixed inset-0 z-[200] flex items-center justify-center p-6 bg-[#0F1928]/90 backdrop-blur-2xl">
          <div class="absolute top-8 right-8">
            <button @click="isAICreateModalOpen = false" class="w-12 h-12 rounded-full flex items-center justify-center bg-white/5 hover:bg-white/10 text-white/40 hover:text-white transition-all">
              <iconify-icon icon="lucide:x" class="text-2xl"></iconify-icon>
            </button>
          </div>

          <div class="w-full max-w-2xl flex flex-col items-center">
            <div class="flex items-center gap-4 mb-8">
              <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-[#3B9BFF] to-[#50C878] flex items-center justify-center text-white shadow-lg">
                <iconify-icon icon="lucide:sparkles" class="text-2xl"></iconify-icon>
              </div>
              <h2 class="text-3xl font-bold text-white tracking-tight">AI 智能创建</h2>
            </div>

            <p class="text-white/40 text-center mb-10 max-w-md">只需描述你想要什么样的 AI 智能体，系统将自动为你生成最佳配置、职责和提示词。</p>

            <div class="w-full bg-white/5 border border-white/10 rounded-[32px] p-2 shadow-2xl focus-within:border-[#3B9BFF]/50 transition-all">
              <textarea
                v-model="aiPrompt"
                placeholder="例如：我想要一个专门分析财务报表的专家，能够识别风险并提供优化建议..."
                class="w-full bg-transparent border-none outline-none text-lg text-white/90 p-6 min-h-[180px] resize-none placeholder:text-white/10"
              ></textarea>
              
              <div class="flex items-center justify-between p-4 bg-white/5 rounded-[24px]">
                <div class="flex items-center gap-2 pl-2">
                  <iconify-icon icon="lucide:info" class="text-white/20 text-xs"></iconify-icon>
                  <span class="text-[10px] text-white/20 uppercase tracking-[0.1em] font-bold">Powered by UNIIOC</span>
                </div>

                <button 
                  @click="handleAICreate"
                  :disabled="!aiPrompt.trim() || isAIGenerating"
                  class="h-12 px-8 bg-[#3B9BFF] hover:bg-[#2A7FDB] disabled:opacity-30 disabled:cursor-not-allowed text-white rounded-2xl font-bold shadow-lg transition-all flex items-center gap-2"
                >
                  <template v-if="isAIGenerating">
                    <iconify-icon icon="lucide:loader-2" class="animate-spin text-lg"></iconify-icon>
                    <span class="text-sm">正在分析中...</span>
                  </template>
                  <template v-else>
                    <span class="text-sm">立即生成</span>
                    <iconify-icon icon="lucide:arrow-right" class="text-lg"></iconify-icon>
                  </template>
                </button>
              </div>
            </div>
          </div>
        </div>
      </Transition>
    </main>
  </BaseLayout>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 155, 255, 0.2); border-radius: 10px; }
</style>
