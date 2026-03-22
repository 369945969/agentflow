<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import BaseLayout from '../components/BaseLayout.vue'
import { API_BASE_URL } from '../config/api'

// State
const models = ref<any[]>([])
const selectedModelId = ref('')
const thinkingEnabled = ref<boolean>(true)
const simplifiedOutput = ref<boolean>(false)
const searchQuery = ref('')
const groups = ref<any[]>([{ id: 'default', name: '所有人 (Default Group)', lastMsg: '欢迎加入群聊', members: [] }])
const activeGroupId = ref('default')
const filteredGroups = ref<any[]>(groups.value)
const expandedGroups = ref<Set<string>>(new Set()) // 存储展开的群聊ID
const showActiveGroupMembers = ref(false) // 控制当前群成员显示
const groupRuleMode = ref<'free' | 'expert' | 'custom'>('free') // 群规模式：自由回答/专家模式/自定义
const customRuleText = ref('') // 自定义群规文本
const CUSTOM_RULE_MAX_LENGTH = 2000
const customRuleLength = computed(() => customRuleText.value.length)
const rightSidebarCollapsed = ref(true) // 控制右侧边栏是否收起
const leftSidebarCollapsed = ref(false) // 控制左侧边栏是否收起，默认展开
const leftSidebarWidth = ref(320) // 左侧边栏宽度，可拖动调整
const isDragging = ref(false) // 是否正在拖动分隔线
const isInitializing = ref(false)

// Computed properties
const activeGroup = computed(() => {
  return groups.value.find(g => g.id === activeGroupId.value) || groups.value[0]
})

// Methods
const toggleGroupExpand = (groupId: string, event: Event) => {
  event.stopPropagation() // 防止触发群聊选择
  if (expandedGroups.value.has(groupId)) {
    expandedGroups.value.delete(groupId)
  } else {
    expandedGroups.value.add(groupId)
  }
}

const selectGroup = (groupId: string) => {
  activeGroupId.value = groupId
}

// Toggle left sidebar
const toggleLeftSidebar = () => {
  leftSidebarCollapsed.value = !leftSidebarCollapsed.value
}

// Drag handling for resizing chat area
const startDrag = (event: MouseEvent) => {
  isDragging.value = true
  document.addEventListener('mousemove', handleDrag)
  document.addEventListener('mouseup', stopDrag)
  // Prevent text selection during drag
  document.body.style.userSelect = 'none'
  document.body.style.cursor = 'col-resize'
  event.preventDefault()
}

const handleDrag = (event: MouseEvent) => {
  if (!isDragging.value) return
  
  // Calculate new width based on mouse position
  const container = document.querySelector('.flex.w-full.h-full.overflow-hidden.relative')
  if (container) {
    const containerRect = container.getBoundingClientRect()
    const mouseX = event.clientX - containerRect.left
    // Set minimum and maximum width constraints
    const minWidth = 200
    const maxWidth = containerRect.width - 400 // Leave space for right sidebar
    const newWidth = Math.max(minWidth, Math.min(mouseX, maxWidth))
    
    // Update left sidebar width
    leftSidebarWidth.value = newWidth
  }
}

const stopDrag = () => {
  isDragging.value = false
  document.removeEventListener('mousemove', handleDrag)
  document.removeEventListener('mouseup', stopDrag)
  // Restore normal cursor and text selection
  document.body.style.userSelect = ''
  document.body.style.cursor = ''
}

// Watch search query to filter groups
watch(searchQuery, (newQuery) => {
  if (!newQuery.trim()) {
    filteredGroups.value = groups.value
    return
  }
  
  const query = newQuery.toLowerCase().trim()
  filteredGroups.value = groups.value.filter(group => 
    group.name.toLowerCase().includes(query)
  )
})

// Create Group Modal State
const isCreateModalOpen = ref(false)
const newGroupName = ref('')
const selectedAgentIds = ref<string[]>([])
const agents = ref<any[]>([]) // 从后端获取的agent列表


const openCreateModal = () => {
  isCreateModalOpen.value = true
}

const closeCreateModal = () => {
  isCreateModalOpen.value = false
  newGroupName.value = ''
  selectedAgentIds.value = []
}

const fetchModels = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/models/`)
    const data = await response.json()
    models.value = data
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

const fetchAgents = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/agents/`)
    if (response.ok) {
      agents.value = await response.json()
    } else {
      console.warn('Failed to fetch agents')
      agents.value = []
    }
  } catch (error) {
    console.error('Error fetching agents:', error)
    agents.value = []
  }
}

const fetchGroups = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/groups/`)
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    const apiGroups = await response.json()
    
    // 直接使用API返回的群组列表（后端已确保default群组置顶）
    groups.value = apiGroups
    filteredGroups.value = groups.value
    return apiGroups
  } catch (error) {
    console.error('Failed to fetch groups:', error)
    // 如果API调用失败，至少显示默认群组
    groups.value = [{ id: 'default', name: '所有人 (Default Group)', lastMsg: '欢迎加入群聊', members: [] }]
    filteredGroups.value = groups.value
    return null
  }
}

const handleCreateGroup = async () => {
  if (!newGroupName.value) return
  try {
    const response = await fetch(`${API_BASE_URL}/api/groups/`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
       body: JSON.stringify({ name: newGroupName.value, members: selectedAgentIds.value })
    })
    if (response.ok) {
      await fetchGroups()
      closeCreateModal()
    }
  } catch (error) {
    console.error('Failed to create group:', error)
  }
}

const deleteGroup = async (groupId: string, event: Event) => {
  event.stopPropagation() // 防止触发群聊选择
  
  // 二次确认
  if (!confirm(`确定要删除这个群组吗？此操作不可撤销。`)) {
    return
  }
  
  // 防止删除默认群组
  if (groupId === 'default') {
    alert('默认群组不能删除')
    return
  }
  
  try {
    const response = await fetch(`${API_BASE_URL}/api/groups/${groupId}`, {
      method: 'DELETE'
    })
    
    if (response.ok) {
      // 如果删除的是当前选中的群组，切换到默认群组
      if (activeGroupId.value === groupId) {
        activeGroupId.value = 'default'
      }
      await fetchGroups()
    } else {
      const error = await response.json()
      alert(`删除失败: ${error.error || '未知错误'}`)
    }
  } catch (error) {
    console.error('Failed to delete group:', error)
    alert('删除群组时发生错误')
  }
}

// 从当前活跃群组更新设置
const updateSettingsFromActiveGroup = () => {
  const group = activeGroup.value
  if (group) {
    isInitializing.value = true
    groupRuleMode.value = group.group_rule_mode || 'free'
    thinkingEnabled.value = group.thinking_enabled !== undefined ? group.thinking_enabled : true
    simplifiedOutput.value = group.simplified_output || false
    customRuleText.value = group.custom_rule || ''
    
    // 延迟重置标识，确保 watcher 触发后的调用被跳过
    setTimeout(() => {
      isInitializing.value = false
    }, 100)
  }
}

// 保存群组设置到后端
const saveGroupSettings = async () => {
  if (isInitializing.value) return
  
  const groupId = activeGroupId.value
  if (!groupId) return
  
  // 自定义规则长度验证
  if (groupRuleMode.value === 'custom' && customRuleText.value.length > CUSTOM_RULE_MAX_LENGTH) {
    alert(`自定义规则长度不能超过 ${CUSTOM_RULE_MAX_LENGTH} 个字符`)
    return
  }
  
  try {
    const response = await fetch(`${API_BASE_URL}/api/groups/${groupId}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        group_rule_mode: groupRuleMode.value,
        thinking_enabled: thinkingEnabled.value,
        simplified_output: simplifiedOutput.value,
        custom_rule: groupRuleMode.value === 'custom' ? customRuleText.value : ''
      })
    })
    
    if (response.ok) {
      // 更新本地groups数据
      const groupIndex = groups.value.findIndex(g => g.id === groupId)
      if (groupIndex !== -1) {
        groups.value[groupIndex].group_rule_mode = groupRuleMode.value
        groups.value[groupIndex].thinking_enabled = thinkingEnabled.value
        groups.value[groupIndex].simplified_output = simplifiedOutput.value
        groups.value[groupIndex].custom_rule = groupRuleMode.value === 'custom' ? customRuleText.value : ''
      }
      console.log('群组设置已保存')
     } else {
       const error = await response.json()
       console.error('保存失败:', error)
       alert('保存失败: ' + (error.error || '未知错误'))
     }
   } catch (error) {
     console.error('保存群组设置时发生错误:', error)
     alert('保存群组设置时发生错误: ' + (error instanceof Error ? error.message : '未知错误'))
   }
}

// 当群规模式改变时，如果是custom模式则自动保存
watch(groupRuleMode, () => {
  // 任何设置变更都自动保存
  saveGroupSettings()
})

// 监听thinkingEnabled和simplifiedOutput变化
watch([thinkingEnabled, simplifiedOutput], () => {
  saveGroupSettings()
})

// 监听customRuleText变化，但只在custom模式下保存
watch(customRuleText, () => {
  if (groupRuleMode.value === 'custom') {
    saveGroupSettings()
  }
})

// 监听活跃群组变化
watch(activeGroupId, () => {
  updateSettingsFromActiveGroup()
})

onMounted(() => {
  fetchModels()
  fetchAgents()
  fetchGroups().then(() => {
    // 数据加载完成后更新设置
    updateSettingsFromActiveGroup()
  })
})
</script>

<template>
  <BaseLayout>
    <div class="flex w-full h-full overflow-hidden relative">
      
      <!-- Left Sidebar: Conversations -->
      <!-- Collapsed state: small sidebar -->
      <div v-if="leftSidebarCollapsed" class="w-12 shrink-0 border-r border-white/10 bg-[#1A2536]/80 flex flex-col">
        <div class="p-4 border-b border-white/5 flex items-center justify-center">
          <button 
            @click="toggleLeftSidebar"
            class="p-2 rounded-lg hover:bg-white/10 text-white/60 hover:text-white transition-colors"
            title="展开 Conferences"
          >
            <iconify-icon icon="lucide:menu" class="text-xl"></iconify-icon>
          </button>
        </div>
        <div class="p-2 flex-1 flex flex-col items-center gap-3">
          <button 
            @click="openCreateModal"
            class="p-2 rounded-lg hover:bg-white/10 text-[#3B9BFF] hover:text-white transition-colors"
            title="创建群组"
          >
            <iconify-icon icon="lucide:plus" class="text-lg"></iconify-icon>
          </button>
        </div>
      </div>
      
      <!-- Expanded state: full sidebar -->
      <div v-else :style="{ width: leftSidebarWidth + 'px' }" class="shrink-0 border-r border-white/10 bg-[#1A2536]/50 flex flex-col">
        <div class="p-6 border-b border-white/5 flex flex-col gap-4">
          <div class="flex justify-between items-center">
            <h2 class="text-lg font-semibold text-white/95">Conferences</h2>
            <div class="flex items-center gap-2">
              <button @click="openCreateModal" class="text-[#3B9BFF] hover:text-white transition-colors" title="创建群组">
                <iconify-icon icon="lucide:plus" class="text-xl"></iconify-icon>
              </button>
              <button 
                @click="toggleLeftSidebar"
                class="p-2 rounded-full hover:bg-white/10 text-white/40 hover:text-white transition-colors"
                title="收起 Conferences"
              >
                <iconify-icon icon="lucide:chevron-left" class="text-xl"></iconify-icon>
              </button>
            </div>
          </div>
          <div class="bg-black/20 border border-white/10 rounded-xl px-4 py-2.5 flex items-center gap-3">
            <iconify-icon icon="lucide:search" class="text-white/40 text-sm"></iconify-icon>
            <input v-model="searchQuery" type="text" placeholder="搜索群组名称..." class="bg-transparent border-none outline-none text-sm text-white/70 w-full">
          </div>
        </div>
        <div class="p-4 overflow-y-auto flex-1 custom-scrollbar">
           <div v-for="g in filteredGroups" :key="g.id" @click="selectGroup(g.id)" :class="['p-4 rounded-xl cursor-pointer mb-2 transition-all', activeGroupId === g.id ? 'bg-[#3B9BFF]/20 border border-[#3B9BFF]/40' : 'bg-white/5 hover:bg-white/10']">
             <div class="flex items-center justify-between">
              <div class="flex-1 min-w-0">
                <h3 class="text-sm font-medium text-white truncate">{{ g.name }}</h3>
                <p class="text-xs text-white/40 truncate">{{ g.lastMsg }}</p>
              </div>
              <div class="flex items-center gap-1">
                <button v-if="g.id !== 'default'" @click.stop="deleteGroup(g.id, $event)" class="p-1 text-white/40 hover:text-red-400 transition-colors" title="删除群组">
                  <iconify-icon 
                    icon="lucide:trash-2" 
                    class="text-sm"
                  ></iconify-icon>
                </button>
                <button @click.stop="toggleGroupExpand(g.id, $event)" class="p-1 text-white/40 hover:text-white transition-colors" :class="expandedGroups.has(g.id) ? 'text-[#3B9BFF]' : ''">
                  <iconify-icon 
                    :icon="expandedGroups.has(g.id) ? 'lucide:chevron-up' : 'lucide:chevron-down'" 
                    class="text-sm"
                  ></iconify-icon>
                </button>
              </div>
            </div>
            
            <!-- Expanded members list -->
            <div v-if="expandedGroups.has(g.id)" class="mt-3 pt-3 border-t border-white/10">
              <div class="text-xs font-medium text-white/60 mb-2">群成员 ({{ g.members?.length || 0 }})</div>
               <div class="space-y-2">
                 <div v-for="member in g.members" :key="member.id || member" class="flex items-center gap-2">
                   <div class="w-6 h-6 rounded-full bg-white/10 flex items-center justify-center">
                     <span class="text-xs text-white/70">{{ (member.name || member).charAt(0) }}</span>
                   </div>
                   <span class="text-xs text-white/80">{{ member.name || member }}</span>
                 </div>
                  <div v-if="!g.members || g.members.length === 0" class="text-xs text-white/40 italic">
                    暂无成员
                  </div>
                </div>
              </div>
            </div>
          </div>
         </div>
         <!-- Draggable separator between left sidebar and chat area -->
         <div 
           v-if="!leftSidebarCollapsed"
           class="w-1 cursor-col-resize hover:bg-[#3B9BFF]/50 bg-white/10 shrink-0"
           @mousedown="startDrag"
           :class="{ 'bg-[#3B9BFF]': isDragging }"
         ></div>
      
       <!-- Main Chat Area -->
      <main class="flex-1 flex flex-col min-w-0 bg-[#0F1928] overflow-hidden">
        <div class="border-b border-white/10 bg-white/5 backdrop-blur-md shrink-0">
           <!-- Chat Header -->
           <div class="h-16 flex items-center justify-between px-6">
             <h2 class="text-lg font-semibold text-white/95">{{ activeGroup.name }}</h2>
              <button @click="showActiveGroupMembers = !showActiveGroupMembers" class="p-2 text-white/40 hover:text-white transition-colors" :class="showActiveGroupMembers ? 'text-[#3B9BFF]' : ''" :title="showActiveGroupMembers ? '隐藏成员' : '显示成员'">
               <iconify-icon 
                 :icon="showActiveGroupMembers ? 'lucide:users' : 'lucide:users'" 
                 class="text-lg"
               ></iconify-icon>
             </button>
           </div>
          
           <!-- Group Members Bar -->
           <div v-if="showActiveGroupMembers" class="px-6 pb-4">
             <div class="text-xs font-medium text-white/60 mb-2">当前群成员 ({{ activeGroup.members?.length || 0 }})</div>
             <div class="flex flex-wrap gap-2">
               <div v-for="member in activeGroup.members" :key="member.id || member" class="flex items-center gap-1.5 px-3 py-1.5 bg-white/5 rounded-lg border border-white/10">
                 <div class="w-6 h-6 rounded-full bg-white/10 flex items-center justify-center shrink-0">
                   <span class="text-xs text-white/70">{{ (member.name || member).charAt(0) }}</span>
                 </div>
                 <span class="text-sm text-white/80">{{ member.name || member }}</span>
               </div>
               <div v-if="!activeGroup.members || activeGroup.members.length === 0" class="text-sm text-white/40 italic px-3 py-1.5">
                 暂无成员
               </div>
             </div>
           </div>
        </div>
        <div class="flex-1 overflow-y-auto p-6 space-y-6 custom-scrollbar">
           <div class="flex flex-col items-center justify-center h-full text-white/30 text-sm">
             <iconify-icon icon="lucide:message-square" class="text-4xl mb-4"></iconify-icon>
             开始你的多人对话...
           </div>
        </div>
        <div class="p-6 border-t border-white/10 bg-white/5 backdrop-blur-md shrink-0">
          <div class="flex gap-4">
            <textarea placeholder="输入消息..." rows="1" class="flex-1 bg-white/5 border border-white/10 rounded-xl p-3 text-sm text-white/90 outline-none focus:border-[#3B9BFF]/50 resize-none"></textarea>
            <button class="w-12 h-12 bg-[#3B9BFF] rounded-xl flex items-center justify-center text-white shadow-[0_0_20px_rgba(59,155,255,0.3)]"><iconify-icon icon="lucide:send" class="text-xl"></iconify-icon></button>
          </div>
        </div>
      </main>

      <!-- Right Sidebar: Controls -->
      <!-- Collapsed state: small expand button -->
      <div v-if="rightSidebarCollapsed" class="w-12 shrink-0 border-l border-white/10 bg-[#1A2536]/80 flex flex-col items-center py-4">
        <button 
          @click="rightSidebarCollapsed = false"
          class="p-2 rounded-lg hover:bg-white/10 text-white/60 hover:text-white transition-colors"
          title="展开设置"
        >
          <iconify-icon icon="lucide:settings" class="text-xl"></iconify-icon>
        </button>
      </div>
      
      <!-- Expanded state: full sidebar -->
      <aside v-else class="w-[380px] shrink-0 border-l border-white/10 bg-[#1A2536] flex flex-col">
        <div class="p-6 border-b border-white/5 flex justify-between items-center">
          <h3 class="text-lg font-semibold text-white/95">会话设置</h3>
          <button 
            @click="rightSidebarCollapsed = true"
            class="p-2 rounded-full hover:bg-white/10 text-white/40 hover:text-white transition-colors"
            title="收起设置"
          >
            <iconify-icon icon="lucide:chevron-right" class="text-xl"></iconify-icon>
          </button>
        </div>
        
        <div class="p-6 space-y-8 overflow-y-auto custom-scrollbar flex-1">
          <div class="space-y-3">
            <label class="text-xs text-white/40 font-bold uppercase tracking-widest">群规设置</label>
            <div class="flex items-center p-1 bg-white/5 border border-white/10 rounded-xl backdrop-blur-xl">
              <button 
                @click="groupRuleMode = 'free'"
                :class="['flex-1 py-3 flex items-center justify-center rounded-lg transition-all', groupRuleMode === 'free' ? 'bg-[#3B9BFF]/20 text-[#3B9BFF]' : 'text-white/40 hover:text-white/60']"
              >
                <span class="text-sm font-medium">自由回答</span>
              </button>
              <button 
                @click="groupRuleMode = 'expert'"
                :class="['flex-1 py-3 flex items-center justify-center rounded-lg transition-all', groupRuleMode === 'expert' ? 'bg-[#3B9BFF]/20 text-[#3B9BFF]' : 'text-white/40 hover:text-white/60']"
              >
                <span class="text-sm font-medium">专家模式</span>
              </button>
              <button 
                @click="groupRuleMode = 'custom'"
                :class="['flex-1 py-3 flex items-center justify-center rounded-lg transition-all', groupRuleMode === 'custom' ? 'bg-[#3B9BFF]/20 text-[#3B9BFF]' : 'text-white/40 hover:text-white/60']"
              >
                <span class="text-sm font-medium">自定义</span>
              </button>
            </div>
            <div class="text-[10px] text-white/30 pt-1">
              <div v-if="groupRuleMode === 'free'">所有成员可以自由发言，无限制对话</div>
              <div v-if="groupRuleMode === 'expert'">由专家主导，严格控制发言顺序和质量</div>
              <div v-if="groupRuleMode === 'custom'">使用自定义规则控制对话流程和行为</div>
            </div>
            
            <!-- Custom rule input (only shown in custom mode) -->
            <div v-if="groupRuleMode === 'custom'" class="mt-4 space-y-3">
              <label class="text-xs text-white/40 font-medium uppercase tracking-widest">自定义规则</label>
              <textarea 
                v-model="customRuleText" 
                placeholder="请输入自定义规则，例如：每次发言前必须引用相关资料..."
                rows="3" 
                :maxlength="CUSTOM_RULE_MAX_LENGTH"
                class="w-full bg-black/20 border border-white/10 rounded-xl px-4 py-3 text-sm text-white outline-none focus:border-[#3B9BFF]/50 transition-all resize-none"
              ></textarea>
              <div class="flex justify-between items-center mt-1">
                <div class="text-xs text-white/40">最多 {{ CUSTOM_RULE_MAX_LENGTH }} 个字符</div>
                <div class="text-xs" :class="customRuleLength > CUSTOM_RULE_MAX_LENGTH ? 'text-red-400' : 'text-white/40'">{{ customRuleLength }} / {{ CUSTOM_RULE_MAX_LENGTH }}</div>
              </div>
              <button 
                @click="saveGroupSettings"
                class="w-full py-3 rounded-xl bg-[#3B9BFF] text-white font-bold shadow-[0_0_20px_rgba(59,155,255,0.4)] hover:bg-[#2A7FDB] transition-all text-sm"
              >
                保存规则
              </button>
            </div>
          </div>
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

      <!-- Create Group Modal -->
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
              <h2 class="text-xl font-bold text-white/95">创建新群组</h2>
              <button @click="closeCreateModal" class="w-8 h-8 rounded-full flex items-center justify-center hover:bg-white/10 text-white/40"><iconify-icon icon="lucide:x" class="text-xl"></iconify-icon></button>
            </div>

            <div class="flex-1 overflow-y-auto p-8 space-y-6 custom-scrollbar">
              <div class="space-y-4">
                <label class="text-xs text-white/40 block font-medium uppercase tracking-widest">群组名称</label>
                <input v-model="newGroupName" type="text" placeholder="例如：项目讨论组" class="w-full bg-black/20 border border-white/10 rounded-xl px-4 py-3 text-white outline-none focus:border-[#3B9BFF]/50 transition-all">
              </div>

              <div class="space-y-4">
                <label class="text-xs text-white/40 block font-medium uppercase tracking-widest">邀请成员</label>
                <div class="grid grid-cols-2 gap-3">
                  <label v-for="agent in agents" :key="agent.id" 
                    class="p-4 border rounded-2xl cursor-pointer transition-all flex items-center gap-3 group relative overflow-hidden"
                    :class="selectedAgentIds.includes(agent.id) ? 'border-[#3B9BFF] bg-[#3B9BFF]/10' : 'border-white/5 bg-white/5 hover:border-white/20'"
                  >
                    <input type="checkbox" v-model="selectedAgentIds" :value="agent.id" class="accent-[#3B9BFF] hidden">
                    <div class="w-8 h-8 rounded-full bg-white/10 flex items-center justify-center">
                      <span class="text-sm text-white/70">{{ agent.name[0] }}</span>
                    </div>
                    <div class="flex-1">
                      <span class="text-sm font-semibold" :class="selectedAgentIds.includes(agent.id) ? 'text-[#3B9BFF]' : 'text-white/80'">{{ agent.name }}</span>
                      <div class="text-[10px] text-white/40 truncate">{{ agent.description || 'No description' }}</div>
                    </div>
                    <iconify-icon v-if="selectedAgentIds.includes(agent.id)" icon="lucide:check-circle-2" class="text-[#3B9BFF]"></iconify-icon>
                  </label>
                </div>
              </div>
            </div>

            <div class="p-6 border-t border-white/10 bg-white/5 flex gap-4">
              <button @click="closeCreateModal" class="flex-1 py-3 rounded-xl border border-white/10 text-white/60 hover:bg-white/5 transition-all text-sm">取消</button>
              <button @click="handleCreateGroup" class="flex-1 py-3 rounded-xl bg-[#3B9BFF] text-white font-bold shadow-[0_0_20px_rgba(59,155,255,0.4)] hover:bg-[#2A7FDB] transition-all text-sm">确认创建</button>
            </div>
          </div>
        </div>
      </Transition>

      <!-- Overlay for modal -->
      <div v-if="isCreateModalOpen" @click="closeCreateModal" class="fixed inset-0 bg-black/40 backdrop-blur-sm z-[90]"></div>
    </div>
  </BaseLayout>
</template>

<style scoped>
:deep(main) { padding: 0 !important; }
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 155, 255, 0.2); border-radius: 10px; }
textarea { scrollbar-width: none; }
.vertical-text {
  writing-mode: vertical-rl;
  transform: rotate(180deg);
  text-orientation: mixed;
}
</style>
