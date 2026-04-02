<script setup lang="ts">
import { ref, computed, onMounted, watch, nextTick, onBeforeUnmount } from 'vue'
import BaseLayout from '../components/BaseLayout.vue'
import { API_BASE_URL, WEBSOCKET_URL } from '../config/api'
import { getWebSocketInstance, type ServerMessage } from '../api/websocket'
import { unified } from 'unified'
import remarkParse from 'remark-parse'
import remarkGfm from 'remark-gfm'
import remarkBreaks from 'remark-breaks'
import remarkRehype from 'remark-rehype'
import rehypeHighlight from 'rehype-highlight'
import rehypeStringify from 'rehype-stringify'
import prettier from 'prettier/standalone'
import prettierMarkdown from 'prettier/plugins/markdown'
import mermaid from 'mermaid'
import 'highlight.js/styles/atom-one-dark.css'

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
const isConnected = ref(false)
const ws = getWebSocketInstance(WEBSOCKET_URL)
const messageListRef = ref<HTMLElement | null>(null)
const chatInput = ref('')
const messagesByGroup = ref<Record<string, any[]>>({})
const pendingPlaceholders = new Map<string, boolean>()
let mermaidInitialized = false
let mermaidTimer: number | null = null

const currentMessages = computed(() => messagesByGroup.value[activeGroupId.value] || [])

const setMessagesForGroup = (groupId: string, nextMessages: any[]) => {
  messagesByGroup.value = { ...messagesByGroup.value, [groupId]: nextMessages }
}

const mdProcessor = unified()
  .use(remarkParse)
  .use(remarkGfm)
  .use(remarkBreaks)
  .use(remarkRehype, { allowDangerousHtml: false })
  .use(rehypeHighlight, { detect: true, ignoreMissing: true })
  .use(rehypeStringify)

const sanitizeMarkdown = (content: string) => {
  let text = String(content || '')
  if (!text) return ''
  text = text.replace(/\r\n/g, '\n')
  const lines = text.split('\n')
  const out: string[] = []
  let inFence = false

  const splitInlineMarkers = (line: string) => {
    let s = line
    s = s.replace(/([^\n])(```)/g, '$1\n\n```')
    s = s.replace(/([^\n])(#{1,6})(?=[^\s#])/g, '$1\n\n$2')
    s = s.replace(/([^\n])-(?=[\u4E00-\u9FFFA-Za-z0-9])/g, '$1\n-')
    return s
  }

  for (const line of lines) {
    const isFenceLine = line.trimStart().startsWith('```')
    if (isFenceLine) {
      inFence = !inFence
      out.push(line)
      continue
    }
    if (inFence) {
      out.push(line)
      continue
    }

    const expanded = splitInlineMarkers(line)
    if (expanded.includes('\n')) {
      expanded.split('\n').forEach(l => out.push(l))
      continue
    }
    out.push(line)
  }

  return out
    .map(l => l.replace(/^(\s*[-*+])(\S)/, '$1 $2').replace(/^(\s*\d+)\.(\S)/, '$1. $2').replace(/^(#{1,6})(\S)/, '$1 $2'))
    .join('\n')
}

const formatMarkdown = (content: string) => {
  const input = sanitizeMarkdown(String(content || ''))
  if (!input) return ''
  try {
    const out = prettier.format(input, {
      parser: 'markdown',
      plugins: [prettierMarkdown],
      proseWrap: 'preserve',
      printWidth: 120
    })
    return typeof out === 'string' ? out : input
  } catch {
    return input
  }
}

const renderMarkdown = (content: string) => {
  const input = formatMarkdown(content || '')
  try {
    return String(mdProcessor.processSync(input))
  } catch {
    const esc: Record<string, string> = {
      '&': '&amp;',
      '<': '&lt;',
      '>': '&gt;',
      '"': '&quot;',
      "'": '&#39;'
    }
    return `<pre class="hljs"><code>${input.replace(/[&<>"']/g, (s: string) => esc[s] || s)}</code></pre>`
  }
}

const scrollToBottom = (instant = false) => {
  nextTick(() => {
    const el = messageListRef.value
    if (!el) return
    el.scrollTo({ top: el.scrollHeight, behavior: instant ? 'auto' : 'smooth' })
  })
}

const scheduleMermaidRender = () => {
  if (typeof window === 'undefined') return
  if (mermaidTimer) window.clearTimeout(mermaidTimer)
  mermaidTimer = window.setTimeout(() => {
    nextTick(() => {
      const root = messageListRef.value
      if (!root) return
      const codeNodes = root.querySelectorAll('pre code')
      codeNodes.forEach(code => {
        const className = (code.getAttribute('class') || '').toLowerCase()
        const isMermaidLang =
          className.includes('language-mermaid') ||
          className.includes('lang-mermaid') ||
          className.includes('language-mermaidgraph') ||
          className.includes('lang-mermaidgraph')
        if (!isMermaidLang) return
        const pre = code.parentElement
        if (!pre || pre.tagName !== 'PRE') return
        const anyPre = pre as HTMLElement
        if (anyPre.dataset.mermaidDone === '1') return
        anyPre.dataset.mermaidDone = '1'
        const container = document.createElement('div')
        container.className = 'mermaid'
        let graphText = (code.textContent || '').trimStart()
        const trimmed = graphText.trimStart()
        if (trimmed.toLowerCase().startsWith('mermaidgraph')) {
          graphText = trimmed.slice('mermaidgraph'.length).trimStart()
        } else {
          graphText = trimmed
        }
        if (/^(TD|LR|RL|BT)\b/.test(graphText)) {
          graphText = `graph ${graphText}`
        }
        container.textContent = graphText
        pre.replaceWith(container)
      })
      if (!mermaidInitialized) return
      const nodes = root.querySelectorAll('.mermaid')
      if (nodes.length === 0) return
      try {
        mermaid.run({ nodes: Array.from(nodes) as any })
      } catch {
      }
    })
  }, 60)
}

const typeOut = (messageObj: any, field: 'thinking' | 'reply', finalText: string) => {
  const text = String(finalText || '')
  const durationMs = Math.min(1600, Math.max(350, text.length * 2))
  const start = performance.now()
  messageObj[field] = messageObj[field] || {}
  messageObj[field].isAnimating = true
  messageObj[field].display = ''
  messageObj[field].content = ''
  messageObj[field].final = text
  return new Promise<void>((resolve) => {
    const step = (now: number) => {
      const t = Math.min(1, (now - start) / durationMs)
      const count = Math.min(text.length, Math.max(0, Math.floor(text.length * t)))
      messageObj[field].display = text.slice(0, count)
      messageObj.timestamp = Date.now()
      scrollToBottom(true)
      if (t >= 1) {
        messageObj[field].isAnimating = false
        messageObj[field].content = text
        messageObj[field].display = ''
        resolve()
        scheduleMermaidRender()
        return
      }
      requestAnimationFrame(step)
    }
    requestAnimationFrame(step)
  })
}

const getSenderName = (senderId: string) => {
  if (!senderId) return 'AI'
  const group = activeGroup.value
  const member = (group?.members || []).find((m: any) => (m.id || m) === senderId)
  if (!member) return senderId
  return member.name || member.id || member
}

const ensurePlaceholder = (groupId: string, messageId: string) => {
  if (pendingPlaceholders.has(`${groupId}:${messageId}`)) return
  pendingPlaceholders.set(`${groupId}:${messageId}`, true)
  const next = [...(messagesByGroup.value[groupId] || [])]
  next.push({
    id: `placeholder:${messageId}`,
    type: 'agent',
    senderId: '',
    isPlaceholder: true,
    content: thinkingEnabled.value ? '正在推理中...' : '正在输入中...',
    timestamp: Date.now()
  })
  setMessagesForGroup(groupId, next)
  scrollToBottom(true)
}

const removePlaceholder = (groupId: string, messageId: string) => {
  const key = `${groupId}:${messageId}`
  if (!pendingPlaceholders.has(key)) return
  pendingPlaceholders.delete(key)
  const next = (messagesByGroup.value[groupId] || []).filter((m: any) => m.id !== `placeholder:${messageId}`)
  setMessagesForGroup(groupId, next)
}

const handleWsConnectionChange = (connected: boolean) => {
  isConnected.value = connected
}

const handleWsMessage = (message: ServerMessage) => {
  const groupId = message.session_id || ''
  if (!groupId) return
  if (!messagesByGroup.value[groupId]) {
    setMessagesForGroup(groupId, [])
  }

  if (message.type === 'stream_chunk') {
    ensurePlaceholder(groupId, message.message_id || '')
    return
  }
  if (message.type === 'stream_end') {
    const payload: any = message.payload || {}
    const mid = message.message_id || ''
    removePlaceholder(groupId, mid)

    const senderId = typeof payload.sender_id === 'string' ? payload.sender_id : ''
    const thinkingText = typeof payload.thinking_content === 'string' ? payload.thinking_content : ''
    const replyText = typeof payload.content === 'string' ? payload.content : ''

    const msgObj: any = {
      id: `agent:${mid}:${senderId}:${Date.now()}`,
      type: 'agent',
      senderId,
      thinking: { content: '', isAnimating: false },
      reply: { content: '', isAnimating: false },
      timestamp: message.timestamp || Date.now()
    }

    const next = [...(messagesByGroup.value[groupId] || []), msgObj]
    setMessagesForGroup(groupId, next)
    scrollToBottom(true)

    ;(async () => {
      if (thinkingText) {
        await typeOut(msgObj, 'thinking', thinkingText)
      }
      await typeOut(msgObj, 'reply', replyText)
    })()
    return
  }
}

const setupWebSocket = () => {
  ws.on('connection-change', handleWsConnectionChange)
  ws.on('message', handleWsMessage)
  ws.connect()
}

onBeforeUnmount(() => {
  ws.off('connection-change', handleWsConnectionChange)
  ws.off('message', handleWsMessage)
})

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
  ws.setSession(activeGroupId.value)
  scrollToBottom(true)
  scheduleMermaidRender()
})

const sendChatMessage = () => {
  const text = chatInput.value.trim()
  const groupId = activeGroupId.value
  if (!text || !groupId) return

  ws.setSession(groupId)
  ws.setUser('group_user')
  const messageId = ws.sendGroupText(text, 'group_user', groupId, groupId, {
    override_config: {
      stream_mode: 'final_only',
      enable_thinking: thinkingEnabled.value,
      return_strategy: { type: 'final_only', final_result_marker: '', chunk_size: 0, debounce_ms: 0 },
      timeout_ms: 300000,
      max_retries: 0
    }
  })

  const next = [...(messagesByGroup.value[groupId] || [])]
  next.push({
    id: `user:${messageId}`,
    type: 'user',
    content: text,
    timestamp: Date.now()
  })
  setMessagesForGroup(groupId, next)
  ensurePlaceholder(groupId, messageId)
  chatInput.value = ''
  scrollToBottom(true)
}

const handleChatKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Enter' && !e.shiftKey) {
    e.preventDefault()
    sendChatMessage()
  }
}

onMounted(() => {
  if (!mermaidInitialized) {
    mermaid.initialize({
      startOnLoad: false,
      theme: 'dark',
      securityLevel: 'strict'
    })
    mermaidInitialized = true
  }
  setupWebSocket()
  fetchModels()
  fetchAgents()
  fetchGroups().then(() => {
    // 数据加载完成后更新设置
    updateSettingsFromActiveGroup()
    ws.setSession(activeGroupId.value)
    scrollToBottom(true)
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
        <div ref="messageListRef" class="flex-1 overflow-y-auto p-6 custom-scrollbar text-[13px]">
          <div v-if="currentMessages.length === 0" class="flex flex-col items-center justify-center h-full text-white/30 text-[13px]">
            <iconify-icon icon="lucide:message-square" class="text-4xl mb-4"></iconify-icon>
            开始你的多人对话...
          </div>
          <div v-else class="space-y-6">
            <div v-for="m in currentMessages" :key="m.id" class="flex" :class="m.type === 'user' ? 'justify-end' : 'justify-start'">
              <div v-if="m.type === 'user'" class="max-w-[760px] rounded-2xl p-4 border bg-[#3B9BFF] text-white border-[#3B9BFF] text-[13px]">
                <div class="whitespace-pre-wrap break-words">{{ m.content }}</div>
                <div class="text-xs mt-2 opacity-60 text-right">{{ new Date(m.timestamp).toLocaleTimeString() }}</div>
              </div>
              <div v-else class="max-w-[760px] rounded-2xl p-4 border bg-white/10 text-white/90 border-white/10">
                <div class="mb-3 flex items-center gap-2">
                  <div class="w-8 h-8 rounded-lg bg-[#3B9BFF]/20 flex items-center justify-center text-[#3B9BFF] font-bold text-xs">
                    {{ getSenderName(m.senderId).charAt(0) }}
                  </div>
                  <span class="text-xs font-semibold text-white/80">{{ getSenderName(m.senderId) }}</span>
                </div>

                <div v-if="m.thinking?.isAnimating || m.thinking?.content" class="mb-3 p-3 bg-amber-500/10 border border-amber-500/20 rounded-lg text-amber-100/70">
                  <div class="font-bold mb-1 flex items-center gap-2 text-xs"><iconify-icon icon="lucide:brain-circuit"></iconify-icon>推理过程</div>
                  <pre v-if="m.thinking?.isAnimating" class="markdown-typing">{{ m.thinking.display }}</pre>
                  <div v-else class="markdown-body break-words text-[13px]" v-html="renderMarkdown(m.thinking.content)"></div>
                </div>

                <div v-if="m.isPlaceholder" class="mt-2 inline-flex items-center gap-2 text-xs opacity-70">
                  <span class="typing-dot"></span>
                  <span class="typing-dot"></span>
                  <span class="typing-dot"></span>
                  <span class="ml-1">{{ m.content }}</span>
                </div>
                <pre v-else-if="m.reply?.isAnimating" class="markdown-typing">{{ m.reply.display }}</pre>
                <div v-else class="markdown-body break-words text-[13px]" v-html="renderMarkdown(m.reply?.content || '')"></div>
                <div class="text-xs mt-2 opacity-60">{{ new Date(m.timestamp).toLocaleTimeString() }}</div>
              </div>
            </div>
          </div>
        </div>
        <div class="p-6 border-t border-white/10 bg-white/5 backdrop-blur-md shrink-0">
          <div class="flex gap-4">
            <textarea v-model="chatInput" placeholder="输入消息..." rows="1" class="flex-1 bg-white/5 border border-white/10 rounded-xl p-3 text-[13px] text-white/90 outline-none focus:border-[#3B9BFF]/50 resize-none" @keydown="handleChatKeydown"></textarea>
            <button class="w-12 h-12 bg-[#3B9BFF] rounded-xl flex items-center justify-center text-white shadow-[0_0_20px_rgba(59,155,255,0.3)]" @click="sendChatMessage"><iconify-icon icon="lucide:send" class="text-xl"></iconify-icon></button>
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

:deep(.markdown-body) {
  line-height: 1.8;
  word-break: break-word;
  font-size: inherit;
}

:deep(.markdown-body ul),
:deep(.markdown-body ol) {
  margin: 0.5em 0 0.5em 1.1em;
  padding-left: 1.1em;
  list-style-position: outside;
}

:deep(.markdown-body ul) { list-style-type: disc; }
:deep(.markdown-body ol) { list-style-type: decimal; }

:deep(.markdown-body li) {
  margin: 0.25em 0;
}

:deep(.markdown-body pre) {
  margin: 0.8em 0;
  padding: 0.9em 1em;
  border-radius: 0.9rem;
  background: rgba(0, 0, 0, 0.35);
  overflow: auto;
}

:deep(.markdown-body pre code) {
  padding: 0;
  background: transparent;
  font-size: 0.9em;
}

:deep(.markdown-body .mermaid) {
  margin: 0.8em 0;
  padding: 0.9em 1em;
  border-radius: 0.9rem;
  background: rgba(0, 0, 0, 0.25);
  overflow: auto;
}

:deep(.markdown-body .mermaid svg) {
  max-width: 100%;
  height: auto;
}

.markdown-typing {
  margin: 0.25em 0;
  white-space: pre-wrap;
  word-break: break-word;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
  font-size: 0.92em;
  line-height: 1.7;
  color: rgba(255, 255, 255, 0.9);
}

.typing-dot {
  width: 6px;
  height: 6px;
  border-radius: 9999px;
  background: rgba(255, 255, 255, 0.6);
  display: inline-block;
  animation: typingBounce 1s infinite ease-in-out;
}

.typing-dot:nth-child(1) { animation-delay: 0ms; }
.typing-dot:nth-child(2) { animation-delay: 120ms; }
.typing-dot:nth-child(3) { animation-delay: 240ms; }

@keyframes typingBounce {
  0%, 80%, 100% { transform: translateY(0); opacity: 0.4; }
  40% { transform: translateY(-3px); opacity: 1; }
}
</style>
