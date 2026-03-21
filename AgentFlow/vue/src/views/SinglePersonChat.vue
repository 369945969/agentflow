<script setup lang="ts">
import { ref, computed, nextTick, onBeforeUnmount, onMounted, watch } from 'vue'
import BaseLayout from '../components/BaseLayout.vue'
import { API_BASE_URL, WEBSOCKET_URL } from '../config/api'
import { getWebSocketInstance, type ServerMessage } from '../api/websocket'

// ========== Models ==========
const models = ref<any[]>([])
const selectedModelId = ref('')
const thinkingEnabled = ref(true)
const simplifiedOutput = ref(false)

// ========== Chat ==========
const messages = ref<any[]>([])
const inputText = ref('')
const isConnected = ref(false)

const messageGroups = computed(() => {
  const groups: any[] = []
  const assistantGroups = new Map<string, any>()

  for (const message of messages.value) {
    const baseId = typeof message.streamKey === 'string'
      ? message.streamKey.split(':')[0]
      : message.id

    if (message.type === 'agent' && !message.isError) {
      let group = assistantGroups.get(baseId)
      if (!group) {
        group = {
          id: `assistant:${baseId}`,
          kind: 'assistant',
          thinking: null,
          reply: null
        }
        assistantGroups.set(baseId, group)
        groups.push(group)
      }

      if (message.isThinking) {
        group.thinking = message
      } else {
        group.reply = message
      }
      continue
    }

    groups.push({
      id: message.id,
      kind: 'single',
      message
    })
  }

  return groups
})

// ========== WebSocket ==========
const ws = getWebSocketInstance()
const wsInitialized = ref(false)
const messageListRef = ref<HTMLElement | null>(null)
const thinkingViewports = new Map<string, HTMLElement>()
const thinkingExpanded = ref<Record<string, boolean>>({})
const shouldStickToBottom = ref(true)

const trimTrailingBlankLines = (text: string) => text.replace(/\n{3,}$/g, '\n').replace(/\s+$/g, '')
const trimLeadingBlankLines = (text: string) => text.replace(/^\s*\n+/g, '')

const scrollMessagesToBottom = (smooth = false) => {
  nextTick(() => {
    const el = messageListRef.value
    if (!el) return
    el.scrollTo({
      top: el.scrollHeight,
      behavior: smooth ? 'smooth' : 'auto'
    })
  })
}

const handleMessageListScroll = () => {
  const el = messageListRef.value
  if (!el) return
  const distanceFromBottom = el.scrollHeight - el.scrollTop - el.clientHeight
  shouldStickToBottom.value = distanceFromBottom < 80
}

const isThinkingExpanded = (streamKey: string) => Boolean(thinkingExpanded.value[streamKey])

const toggleThinkingExpanded = (streamKey: string) => {
  thinkingExpanded.value = {
    ...thinkingExpanded.value,
    [streamKey]: !thinkingExpanded.value[streamKey]
  }
  scrollThinkingViewport(streamKey)
}

const setThinkingViewport = (streamKey: string, el: any) => {
  if (el instanceof HTMLElement) {
    thinkingViewports.set(streamKey, el)
    el.scrollTop = el.scrollHeight
    return
  }
  thinkingViewports.delete(streamKey)
}

const scrollThinkingViewport = (streamKey: string) => {
  nextTick(() => {
    const el = thinkingViewports.get(streamKey)
    if (!el) return
    el.scrollTop = el.scrollHeight
  })
  if (shouldStickToBottom.value) {
    scrollMessagesToBottom()
  }
}

const upsertAgentMessage = (message: ServerMessage) => {
  const payload = message.payload
  if (!payload || typeof payload !== 'object') return

  const isThinking = Boolean(payload.is_thinking)
  const streamKey = `${message.message_id || 'unknown'}:${isThinking ? 'thinking' : 'answer'}`
  const existing = messages.value.find((item: any) => item.type === 'agent' && item.streamKey === streamKey)

  if (message.type === 'stream_chunk') {
    if (existing) {
      existing.timestamp = message.timestamp || Date.now()
      existing.isFinal = false
      existing.isStreaming = true
      if (isThinking) {
        existing.content += payload.content || ''
        existing.content = trimLeadingBlankLines(existing.content)
        scrollThinkingViewport(streamKey)
      } else {
        existing.content += payload.content || ''
        existing.content = trimLeadingBlankLines(existing.content)
        if (shouldStickToBottom.value) {
          scrollMessagesToBottom()
        }
      }
      return
    }

    const nextMessage = {
      id: streamKey,
      streamKey,
      type: 'agent',
      content: isThinking ? '' : trimLeadingBlankLines(payload.content || ''),
      isThinking,
      isFinal: false,
      isStreaming: true,
      timestamp: message.timestamp || Date.now()
    }
    messages.value.push(nextMessage)
    if (isThinking) {
      nextMessage.content = trimLeadingBlankLines(payload.content || '')
      scrollThinkingViewport(streamKey)
    } else if (shouldStickToBottom.value) {
      scrollMessagesToBottom()
    }
    return
  }

  const answerKey = `${message.message_id || 'unknown'}:answer`
  const thinkingKey = `${message.message_id || 'unknown'}:thinking`
  const thinkingMessage = messages.value.find((item: any) => item.type === 'agent' && item.streamKey === thinkingKey)
  if (thinkingMessage) {
    thinkingMessage.isStreaming = false
    scrollThinkingViewport(thinkingKey)
  }
  const finalExisting = messages.value.find((item: any) => item.type === 'agent' && item.streamKey === answerKey)
  if (finalExisting) {
    finalExisting.content = trimTrailingBlankLines(trimLeadingBlankLines(payload.content || finalExisting.content || ''))
    finalExisting.timestamp = message.timestamp || Date.now()
    finalExisting.isThinking = false
    finalExisting.isFinal = true
    finalExisting.isStreaming = false
    if (shouldStickToBottom.value) {
      scrollMessagesToBottom(true)
    }
    return
  }

  messages.value.push({
    id: answerKey,
    streamKey: answerKey,
    type: 'agent',
    content: trimTrailingBlankLines(trimLeadingBlankLines(payload.content || '')),
    isThinking: false,
    isFinal: true,
    isStreaming: false,
    timestamp: message.timestamp || Date.now()
  })
  if (shouldStickToBottom.value) {
    scrollMessagesToBottom(true)
  }
}

// Update agent settings in database
const updateAgentSettings = async () => {
  const userId = selectedUserId.value
  if (!userId) return
  
  const user = selectedUser.value
  if (!user) return
  
  try {
    const response = await fetch(`${API_BASE_URL}/api/agents/${userId}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        ...user,
        thinking_enabled: thinkingEnabled.value,
        simplified_output: simplifiedOutput.value
      })
    })
    
    if (response.ok) {
      // Update local user data
      const userIndex = users.value.findIndex(u => u.id === userId)
      if (userIndex !== -1) {
        users.value[userIndex].thinking_enabled = thinkingEnabled.value
        users.value[userIndex].simplified_output = simplifiedOutput.value
      }
      console.log('数字员工设置已更新')
    } else {
      const error = await response.json()
      console.error('更新失败:', error)
    }
  } catch (error) {
    console.error('更新数字员工设置时发生错误:', error)
  }
}

// ========== WebSocket Functions ==========
const setupWebSocket = () => {
  if (wsInitialized.value) return
  wsInitialized.value = true
  console.log('[SinglePersonChat] Initializing WebSocket listeners', { url: WEBSOCKET_URL })
  
  // Set up event listeners if not already set up
  ws.on('connection-change', (connected: boolean) => {
    isConnected.value = connected
    console.log('[SinglePersonChat] WebSocket connection changed:', connected)
  })
  
  ws.on('message', (message: ServerMessage) => {
    console.log('[SinglePersonChat] Received message:', message)
    if (message.type === 'stream_chunk' || message.type === 'stream_end') {
      upsertAgentMessage(message)
    } else if (message.type === 'error') {
      console.error('[SinglePersonChat] WebSocket error payload:', message.payload)
      const payload = message.payload
      const errorText = payload?.message || '服务暂时不可用，请稍后重试'
      messages.value.push({
        id: `error:${message.message_id || Date.now()}`,
        type: 'agent',
        content: `错误：${errorText}`,
        isThinking: false,
        isFinal: true,
        isError: true,
        timestamp: message.timestamp || Date.now()
      })
      if (shouldStickToBottom.value) {
        scrollMessagesToBottom(true)
      }
    }
  })
  
  // Connect if not already connected
  if (!ws.getConnectionStatus()) {
    ws.connect()
  }
}

const sendMessage = () => {
  if (!inputText.value.trim() || !selectedUserId.value) return
  console.log('[SinglePersonChat] sendMessage called', {
    selectedUserId: selectedUserId.value,
    isConnected: ws.getConnectionStatus(),
    sessionId: ws.getSessionId(),
    textLength: inputText.value.trim().length
  })
  
  const userMessage = {
    id: Date.now().toString(),
    type: 'user',
    content: inputText.value.trim(),
    timestamp: Date.now()
  }
  
  messages.value.push(userMessage)
  shouldStickToBottom.value = true
  scrollMessagesToBottom(true)
  
  // Send via WebSocket
  const messageId = ws.sendText(inputText.value.trim(), selectedUserId.value)
  
  console.log('[SinglePersonChat] Message dispatched with ID:', messageId)
  inputText.value = ''
}

const handleTextareaKeydown = (event: KeyboardEvent) => {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    sendMessage()
  }
}

const fetchModels = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/models/`)
    const data = await response.json()
    models.value = data
    const defaultModel = data.find((m: any) => m.is_default)
    if (defaultModel) selectedModelId.value = defaultModel.id
    else if (data.length > 0) selectedModelId.value = data[0].id
  } catch (error) {
    console.error('Failed to fetch models:', error)
  }
}

// ========== Users ==========
const users = ref<any[]>([])
const searchQuery = ref('')


// Fetch users sorted by ID descending (newest first)
const fetchUsers = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/agents/`)
    if (response.ok) {
      users.value = await response.json()
      
      // Auto-select first user as default conversation user
      if (!selectedUserId.value && users.value.length > 0) {
        selectedUserId.value = users.value[0].id
      }
    } else {
      console.warn('Failed to fetch users')
      users.value = []
    }
  } catch (error) {
    console.error('Error fetching users:', error)
    users.value = []
  }
}

// Filtered users based on search query
const filteredUsers = computed(() => {
  if (!searchQuery.value.trim()) return users.value

  const query = searchQuery.value.toLowerCase().trim()
  return users.value.filter(user =>
    user.name.toLowerCase().includes(query)
  )
})

// Selected user for this conversation
const selectedUserId = ref('')
const selectedUser = computed(() => {
  if (!selectedUserId.value) return null
  return users.value.find(u => u.id === selectedUserId.value)
})
const activeUserName = computed(() => selectedUser.value?.name || '--')

// Users visible in collapsed sidebar mode (max 5, always include selected user)
const visibleUsersInCollapsedMode = computed(() => {
  const filtered = filteredUsers.value
  if (filtered.length <= 5) return filtered
  
  // Always include selected user if exists in filtered results
  const selectedUserObj = filtered.find(u => u.id === selectedUserId.value)
  const otherUsers = filtered.filter(u => u.id !== selectedUserId.value)
  
  // If selected user exists, include it plus 4 others
  if (selectedUserObj) {
    return [selectedUserObj, ...otherUsers.slice(0, 4)]
  } else {
    return otherUsers.slice(0, 5)
  }
})

// Sidebar state
const leftSidebarCollapsed = ref(false) // 控制左侧边栏是否收起，默认展开
const leftSidebarWidth = ref(320) // 左侧边栏宽度，可拖动调整
const isDragging = ref(false) // 是否正在拖动分隔线
const rightSidebarCollapsed = ref(true) // 控制右侧边栏是否收起，默认收起

// Toggle left sidebar
const toggleLeftSidebar = () => {
  leftSidebarCollapsed.value = !leftSidebarCollapsed.value
}

// Select user from collapsed sidebar
const selectUserFromCollapsed = (userId: string) => {
  selectedUserId.value = userId
}

// Drag handling for resizing sidebar
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
  const container = document.querySelector('.flex.w-full.h-full.overflow-hidden')
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

// Watch for selected user changes to update settings
watch(selectedUserId, () => {
  const user = selectedUser.value
  if (user) {
    thinkingEnabled.value = user.thinking_enabled ?? true
    simplifiedOutput.value = user.simplified_output ?? false
    
    // Update WebSocket user
    ws.setUser(user.id)
    ws.clearSession()
    console.log('[SinglePersonChat] Selected user changed', {
      userId: user.id,
      userName: user.name
    })
    
    // Clear previous messages when switching users
    messages.value = []
    thinkingExpanded.value = {}
    shouldStickToBottom.value = true
  }
  
  // Ensure WebSocket is connected and set up
  setupWebSocket()
})

// Watch for settings changes to save to database
watch([thinkingEnabled, simplifiedOutput], () => {
  updateAgentSettings()
})

onMounted(() => {
  console.log('[SinglePersonChat] mounted', { websocketUrl: WEBSOCKET_URL })
  fetchModels()
  fetchUsers()
  setupWebSocket()
})

onBeforeUnmount(() => {
  thinkingViewports.clear()
})
</script>

<template>
  <BaseLayout>
    <div class="flex w-full h-full overflow-hidden bg-[#0F1928]">
      <!-- Left Sidebar: Agent List -->
       <!-- Collapsed state: small sidebar with user avatars -->
       <div v-if="leftSidebarCollapsed" class="w-12 shrink-0 border-r border-white/10 bg-[#1A2536]/80 flex flex-col">
         <div class="p-4 border-b border-white/5 flex items-center justify-center">
           <button 
             @click="toggleLeftSidebar"
             class="p-2 rounded-lg hover:bg-white/10 text-white/60 hover:text-white transition-colors"
              title="展开 数字员工"
           >
             <iconify-icon icon="lucide:menu" class="text-xl"></iconify-icon>
           </button>
         </div>
         <div class="flex-1 overflow-y-auto py-2 custom-scrollbar">
           <div class="flex flex-col items-center gap-2">
              <button 
                v-for="user in visibleUsersInCollapsedMode" 
                :key="user.id"
                @click="selectUserFromCollapsed(user.id)"
                :title="user.name"
                :class="['w-8 h-8 rounded-full flex items-center justify-center text-xs border transition-all hover:scale-110', selectedUserId === user.id ? 'bg-[#3B9BFF] text-white border-[#3B9BFF]' : 'bg-white/10 text-white/70 border-transparent hover:bg-white/20']"
              >
                {{ user.name.charAt(0) }}
              </button>
           </div>
            <div v-if="filteredUsers.length > 5" class="text-center mt-2">
              <div class="text-[10px] text-white/30">+{{ filteredUsers.length - 5 }}</div>
           </div>
         </div>
       </div>
      
      <!-- Expanded state: full sidebar -->
      <div v-else :style="{ width: leftSidebarWidth + 'px' }" class="shrink-0 border-r border-white/10 bg-[#1A2536]/50 flex flex-col">
        <div class="p-6 border-b border-white/5 flex flex-col gap-4">
          <div class="flex justify-between items-center">
             <h2 class="text-lg font-semibold text-white/95">数字员工</h2>
            <button 
              @click="toggleLeftSidebar"
              class="p-2 rounded-full hover:bg-white/10 text-white/40 hover:text-white transition-colors"
               title="收起 数字员工"
            >
              <iconify-icon icon="lucide:chevron-left" class="text-xl"></iconify-icon>
            </button>
          </div>
          <div class="bg-black/20 border border-white/10 rounded-xl px-4 py-2.5 flex items-center gap-3">
            <iconify-icon icon="lucide:search" class="text-white/40 text-sm"></iconify-icon>
             <input v-model="searchQuery" type="text" placeholder="搜索数字员工名称..." class="bg-transparent border-none outline-none text-sm text-white/70 w-full">
          </div>
        </div>
        <div class="p-4 overflow-y-auto flex-1 custom-scrollbar">
          <div v-for="user in filteredUsers" :key="user.id" 
               @click="selectedUserId = user.id"
               :class="['p-4 rounded-xl cursor-pointer mb-2 transition-all', selectedUserId === user.id ? 'bg-[#3B9BFF]/20 border border-[#3B9BFF]/40' : 'bg-white/5 hover:bg-white/10']">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-full bg-white/10 flex items-center justify-center">
                <span class="text-sm text-white/70">{{ user.name[0] }}</span>
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="text-sm font-medium text-white truncate">{{ user.name }}</h3>
                <p class="text-xs text-white/40 truncate">{{ user.description || 'No description' }}</p>
              </div>
            </div>
          </div>
          <div v-if="filteredUsers.length === 0" class="text-center py-8">
            <iconify-icon icon="lucide:users" class="w-8 h-8 mx-auto mb-2 text-white/20" />
            <div class="text-white/30 text-sm">{{ searchQuery ? '未找到匹配Agent' : '暂无Agent' }}</div>
          </div>
        </div>
        <!-- Draggable separator between left sidebar and chat area -->
        <div 
          class="w-1 cursor-col-resize hover:bg-[#3B9BFF]/50 bg-white/10 shrink-0"
          @mousedown="startDrag"
          :class="{ 'bg-[#3B9BFF]': isDragging }"
        ></div>
      </div>

      <!-- Main Chat Area -->
      <main class="flex-1 flex flex-col min-w-0 overflow-hidden">
        <div class="border-b border-white/10 bg-white/5 backdrop-blur-md shrink-0">
          <div class="h-16 flex items-center px-6">
            <div class="flex items-center gap-3">
              <div class="w-2 h-2 rounded-full" :class="isConnected ? 'bg-green-500' : 'bg-red-500'" :title="isConnected ? '已连接到服务器' : '未连接'"></div>
              <h2 class="text-lg font-semibold text-white/95">{{ activeUserName }}</h2>
            </div>
          </div>
        </div>
        <div
          ref="messageListRef"
          class="flex-1 overflow-y-auto p-6 space-y-6 custom-scrollbar"
          @scroll="handleMessageListScroll"
        >
          <div v-if="messageGroups.length === 0" class="flex flex-col items-center justify-center h-full text-white/30 text-sm">
            <iconify-icon icon="lucide:user" class="text-4xl mb-4"></iconify-icon>
            开始你的单人对话...
          </div>
          <div v-else class="space-y-6">
            <div
              v-for="group in messageGroups"
              :key="group.id"
              class="flex"
              :class="group.kind === 'single' && group.message.type === 'user' ? 'justify-end' : 'justify-start'"
            >
              <div v-if="group.kind === 'assistant'" class="w-full max-w-[760px] space-y-3">
                <div v-if="group.thinking" class="w-full text-amber-50">
                  <div class="rounded-2xl border border-amber-200/12 bg-amber-500/10 shadow-[0_12px_32px_rgba(251,191,36,0.08)]">
                    <button
                      type="button"
                      class="flex w-full items-center justify-between gap-3 px-3 py-2 text-left"
                      @click="toggleThinkingExpanded(group.thinking.streamKey)"
                    >
                      <div class="min-w-0 flex items-center gap-2">
                        <span class="inline-flex h-6 w-6 items-center justify-center rounded-full bg-amber-300/14 text-amber-200">
                          <iconify-icon icon="lucide:brain-circuit" class="text-sm"></iconify-icon>
                        </span>
                        <div class="min-w-0">
                          <div class="flex items-center gap-2">
                            <div class="text-[11px] font-semibold uppercase tracking-[0.22em] text-amber-200/90">Thinking</div>
                            <span v-if="group.thinking.isStreaming" class="h-2 w-2 rounded-full bg-amber-200/80 animate-pulse"></span>
                          </div>
                          <div class="text-xs text-amber-50/60 break-words">
                            {{ group.thinking.isStreaming ? '推理中（实时）' : '推理完成' }}
                          </div>
                        </div>
                      </div>
                      <div class="ml-3 flex shrink-0 items-center gap-2 text-amber-100/70">
                        <span class="text-[11px]">{{ isThinkingExpanded(group.thinking.streamKey) ? '收起' : '展开' }}</span>
                        <iconify-icon
                          icon="lucide:chevron-down"
                          class="text-base transition-transform duration-200"
                          :class="isThinkingExpanded(group.thinking.streamKey) ? 'rotate-180' : ''"
                        ></iconify-icon>
                      </div>
                    </button>
                    <div
                      :ref="el => setThinkingViewport(group.thinking.streamKey, el)"
                      class="thinking-viewport w-full whitespace-pre-wrap break-words px-3 pb-3 text-sm leading-7 text-amber-50/88"
                      :class="isThinkingExpanded(group.thinking.streamKey) ? 'thinking-viewport-expanded' : 'thinking-viewport-collapsed'"
                    >{{ group.thinking.content || (group.thinking.isStreaming ? '思考中…' : '') }}</div>
                  </div>
                </div>

                <div
                  v-if="group.reply"
                  class="rounded-2xl p-4 border bg-white/10 text-white/90 border-white/10"
                >
                  <div class="mb-2 inline-flex items-center gap-2 rounded-full px-2.5 py-1 text-[11px] uppercase tracking-[0.18em] bg-emerald-300/15 text-emerald-200">
                    Reply
                  </div>
                  <div class="whitespace-pre-wrap">{{ group.reply.content }}</div>
                  <div
                    v-if="group.reply.isStreaming"
                    class="mt-3 inline-flex items-center gap-2 text-xs opacity-70"
                  >
                    <span class="h-2 w-2 rounded-full bg-current animate-pulse"></span>
                    <span>生成中...</span>
                  </div>
                  <div class="text-xs mt-2 opacity-60">{{ new Date(group.reply.timestamp).toLocaleTimeString() }}</div>
                </div>
              </div>

              <div
                v-else-if="group.message.type === 'user'"
                class="max-w-[760px] rounded-2xl p-4 border bg-[#3B9BFF] text-white border-[#3B9BFF]"
              >
                <div class="whitespace-pre-wrap">{{ group.message.content }}</div>
                <div class="text-xs mt-2 opacity-60">{{ new Date(group.message.timestamp).toLocaleTimeString() }}</div>
              </div>

              <div
                v-else
                class="max-w-[760px] rounded-2xl p-4 border bg-red-500/10 text-red-50 border-red-400/30"
              >
                <div class="mb-2 inline-flex items-center gap-2 rounded-full px-2.5 py-1 text-[11px] uppercase tracking-[0.18em] bg-red-400/15 text-red-200">
                  Error
                </div>
                <div class="whitespace-pre-wrap">{{ group.message.content }}</div>
                <div class="text-xs mt-2 opacity-60">{{ new Date(group.message.timestamp).toLocaleTimeString() }}</div>
              </div>
            </div>
          </div>
        </div>
        <div class="p-6 border-t border-white/10 bg-white/5 backdrop-blur-md shrink-0">
          <div class="flex gap-4">
            <textarea v-model="inputText" @keydown="handleTextareaKeydown" placeholder="输入消息..." rows="1" class="chat-input flex-1 bg-white/5 border border-white/10 rounded-xl p-3 text-sm text-white/90 outline-none focus:border-[#3B9BFF]/50 resize-none overflow-y-auto max-h-32"></textarea>
            <button @click="sendMessage" class="w-12 h-12 bg-[#3B9BFF] rounded-xl flex items-center justify-center text-white shadow-[0_0_20px_rgba(59,155,255,0.3)]"><iconify-icon icon="lucide:send" class="text-xl"></iconify-icon></button>
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
.chat-input {
  scrollbar-width: thin;
  scrollbar-color: rgba(59, 155, 255, 0.35) rgba(255, 255, 255, 0.06);
}
.chat-input::-webkit-scrollbar {
  width: 6px;
}
.chat-input::-webkit-scrollbar-thumb {
  background: rgba(59, 155, 255, 0.35);
  border-radius: 9999px;
}
.chat-input::-webkit-scrollbar-thumb:hover {
  background: rgba(95, 180, 255, 0.5);
}
.chat-input::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.06);
  border-radius: 9999px;
}
.thinking-viewport {
  overflow-y: auto;
  padding-right: 0.25rem;
  transition: max-height 180ms ease;
}
.thinking-viewport-collapsed {
  max-height: 6.5em;
  overflow-y: hidden;
  mask-image: linear-gradient(to bottom, black 62%, transparent 100%);
}
.thinking-viewport-expanded {
  max-height: 18em;
  mask-image: none;
}
.thinking-viewport::-webkit-scrollbar {
  width: 4px;
}
.thinking-viewport::-webkit-scrollbar-thumb {
  background: rgba(251, 191, 36, 0.28);
  border-radius: 9999px;
}
</style>
