<script setup lang="ts">
import { ref, computed, nextTick, onBeforeUnmount, onMounted, watch } from 'vue'
import BaseLayout from '../components/BaseLayout.vue'
import { API_BASE_URL, WEBSOCKET_URL } from '../config/api'
import { getWebSocketInstance, type ServerMessage } from '../api/websocket'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/atom-one-dark.css'

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return '<pre class="hljs"><code>' +
               hljs.highlight(str, { language: lang, ignoreIllegals: true }).value +
               '</code></pre>';
      } catch (__) {}
    }
    return '<pre class="hljs"><code>' + md.utils.escapeHtml(str) + '</code></pre>';
  }
})

const renderMarkdown = (content: string) => md.render(content || '')

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
const pendingUserEcho = new Map<string, string>()
const seenStreamEventKeys = new Map<string, number>()

const rememberStreamEvent = (key: string) => {
  const now = Date.now()
  seenStreamEventKeys.set(key, now)

  if (seenStreamEventKeys.size < 400) return
  const cutoff = now - 60_000
  for (const [k, ts] of seenStreamEventKeys) {
    if (ts < cutoff) seenStreamEventKeys.delete(k)
    if (seenStreamEventKeys.size < 300) break
  }
}

const handleWsConnectionChange = (connected: boolean) => {
  isConnected.value = connected
}

const handleWsMessage = (message: ServerMessage) => {
  if (message.type === 'stream_chunk' || message.type === 'stream_end') {
    const payload: any = message.payload
    const content = typeof payload?.content === 'string' ? payload.content : ''
    const isThinking = Boolean(payload?.is_thinking)
    const isFinal = Boolean(payload?.is_final)
    const eventKey = `${message.type}|${message.message_id || ''}|${isThinking ? 1 : 0}|${isFinal ? 1 : 0}|${content}`
    if (seenStreamEventKeys.has(eventKey)) return
    rememberStreamEvent(eventKey)
    upsertAgentMessage(message)
    return
  }
  if (message.type === 'error') {
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
    if (shouldStickToBottom.value) scrollMessagesToBottom(true)
  }
}

const trimTrailingBlankLines = (text: string) => text.replace(/\n{3,}$/g, '\n').replace(/\s+$/g, '')
const trimLeadingBlankLines = (text: string) => text.replace(/^\s*\n+/g, '')

const stripMemoryLines = (text: string) => {
  if (!text) return ''
  return text.split('\n').filter(line => !line.trimStart().startsWith('[MEMORY]')).join('\n')
}

const stripLeadingUserEcho = (messageId: string | undefined, text: string) => {
  if (!messageId) return text
  const userText = pendingUserEcho.get(messageId)
  if (!userText) return text
  const needle = userText.trim()
  if (!needle) { pendingUserEcho.delete(messageId); return text }
  const candidate = trimLeadingBlankLines(text)
  if (!candidate.startsWith(needle)) return text
  const rest = candidate.slice(needle.length)
  if (rest !== '' && !/^\s/.test(rest)) return text
  return trimLeadingBlankLines(rest)
}

const scrollMessagesToBottom = (smooth = false) => {
  nextTick(() => {
    const el = messageListRef.value
    if (!el) return
    el.scrollTo({ top: el.scrollHeight, behavior: smooth ? 'smooth' : 'auto' })
  })
}

const handleMessageListScroll = () => {
  const el = messageListRef.value
  if (!el) return
  const distanceFromBottom = el.scrollHeight - el.scrollTop - el.clientHeight
  shouldStickToBottom.value = distanceFromBottom < 80
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
  if (shouldStickToBottom.value) scrollMessagesToBottom()
}

const hasMeaningfulThinking = (message: any) => {
  if (!message) return false
  const content = typeof message.content === 'string' ? message.content.trim() : ''
  return content && !/^(\.{3}|…+)$/.test(content)
}

const upsertAgentMessage = (message: ServerMessage) => {
  const payload = message.payload
  if (!payload || typeof payload !== 'object') return

  const isThinking = Boolean(payload.is_thinking)
  const messageId = message.message_id || ''
  const streamKey = `${message.message_id || 'unknown'}:${isThinking ? 'thinking' : 'answer'}`
  const existing = messages.value.find((item: any) => item.type === 'agent' && item.streamKey === streamKey)

  if (message.type === 'stream_chunk') {
    let chunk = stripMemoryLines(payload.content || '')
    if (!isThinking) chunk = stripLeadingUserEcho(messageId, chunk)
    if (!chunk && !existing) return
    if (existing) {
      existing.content = trimLeadingBlankLines(stripMemoryLines(existing.content + chunk))
      if (isThinking) scrollThinkingViewport(streamKey)
      else if (shouldStickToBottom.value) scrollMessagesToBottom()
      return
    }

    const nextMessage = {
      id: streamKey, streamKey, type: 'agent', content: trimLeadingBlankLines(chunk),
      isThinking, isFinal: false, isStreaming: true, timestamp: message.timestamp || Date.now()
    }
    messages.value.push(nextMessage)
    if (isThinking) scrollThinkingViewport(streamKey)
    else if (shouldStickToBottom.value) scrollMessagesToBottom()
    return
  }

  const answerKey = `${message.message_id || 'unknown'}:answer`
  const finalExisting = messages.value.find((item: any) => item.type === 'agent' && item.streamKey === answerKey)
  const finalContent = stripMemoryLines(payload.content || '')
  
  if (finalExisting) {
    finalExisting.content = stripLeadingUserEcho(messageId, trimTrailingBlankLines(trimLeadingBlankLines(stripMemoryLines(finalContent || finalExisting.content || ''))))
    finalExisting.isFinal = true; finalExisting.isStreaming = false
    if (shouldStickToBottom.value) scrollMessagesToBottom(true)
    if (messageId) pendingUserEcho.delete(messageId)
    return
  }

  messages.value.push({
    id: answerKey, streamKey: answerKey, type: 'agent',
    content: stripLeadingUserEcho(messageId, trimTrailingBlankLines(trimLeadingBlankLines(stripMemoryLines(finalContent || '')))),
    isThinking: false, isFinal: true, isStreaming: false, timestamp: message.timestamp || Date.now()
  })
  if (shouldStickToBottom.value) scrollMessagesToBottom(true)
  if (messageId) pendingUserEcho.delete(messageId)
}

const clearChatHistory = () => {
  if (confirm('确定要清空当前对话记录吗？')) {
    messages.value = []
  }
}

const updateAgentSettings = async () => {
  if (!selectedUserId.value || !selectedUser.value) return
  try {
    await fetch(`${API_BASE_URL}/api/agents/${selectedUserId.value}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        ...selectedUser.value,
        thinking_enabled: thinkingEnabled.value,
        simplified_output: simplifiedOutput.value
      })
    })
  } catch (error) { console.error(error) }
}

const setupWebSocket = () => {
  if (wsInitialized.value) return
  wsInitialized.value = true
  ws.on('connection-change', handleWsConnectionChange)
  ws.on('message', handleWsMessage)
  if (!ws.getConnectionStatus()) ws.connect()
}

const sendMessage = () => {
  const text = inputText.value.trim()
  if (!text || !selectedUserId.value) return
  messages.value.push({ id: Date.now().toString(), type: 'user', content: text, timestamp: Date.now() })
  const messageId = ws.sendText(text, selectedUserId.value)
  if (messageId) pendingUserEcho.set(messageId, text)
  inputText.value = ''
  shouldStickToBottom.value = true
  scrollMessagesToBottom(true)
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
  } catch (error) { console.error(error) }
}

const users = ref<any[]>([])
const searchQuery = ref('')
const selectedUserId = ref('')
const selectedUser = computed(() => users.value.find(u => u.id === selectedUserId.value))
const activeUserName = computed(() => selectedUser.value?.name || '--')

const fetchUsers = async () => {
  try {
    const response = await fetch(`${API_BASE_URL}/api/agents/`)
    if (response.ok) {
      users.value = await response.json()
      if (!selectedUserId.value && users.value.length > 0) selectedUserId.value = users.value[0].id
    }
  } catch (error) { console.error(error) }
}

const filteredUsers = computed(() => {
  const query = searchQuery.value.toLowerCase().trim()
  return users.value.filter(u => u.name.toLowerCase().includes(query))
})

const leftSidebarCollapsed = ref(false)
const leftSidebarWidth = ref(320)
const isDragging = ref(false)
const rightSidebarCollapsed = ref(true)

const toggleLeftSidebar = () => { leftSidebarCollapsed.value = !leftSidebarCollapsed.value }

const startDrag = (event: MouseEvent) => {
  isDragging.value = true
  document.addEventListener('mousemove', handleDrag)
  document.addEventListener('mouseup', stopDrag)
  document.body.style.userSelect = 'none'
  document.body.style.cursor = 'col-resize'
}

const handleDrag = (event: MouseEvent) => {
  if (!isDragging.value) return
  const container = document.querySelector('.flex.w-full.h-full')
  if (container) {
    const rect = container.getBoundingClientRect()
    leftSidebarWidth.value = Math.max(200, Math.min(event.clientX - rect.left, rect.width - 400))
  }
}

const stopDrag = () => {
  isDragging.value = false
  document.removeEventListener('mousemove', handleDrag)
  document.removeEventListener('mouseup', stopDrag)
  document.body.style.userSelect = ''
  document.body.style.cursor = ''
}

watch(selectedUserId, () => {
  const user = selectedUser.value
  if (user) {
    thinkingEnabled.value = user.thinking_enabled ?? true
    simplifiedOutput.value = user.simplified_output ?? false
    ws.setUser(user.id); ws.clearSession(); messages.value = []
  }
  setupWebSocket()
})

watch([thinkingEnabled, simplifiedOutput], updateAgentSettings)

onMounted(() => { fetchModels(); fetchUsers(); setupWebSocket() })
</script>

<template>
  <BaseLayout>
    <div class="flex w-full flex-1 min-h-0 overflow-hidden bg-[#0F1928]">
      <aside v-if="!leftSidebarCollapsed" :style="{ width: leftSidebarWidth + 'px' }" class="shrink-0 border-r border-white/10 bg-[#1A2536]/50 flex flex-col">
        <div class="p-6 border-b border-white/5 flex flex-col gap-4">
          <div class="flex justify-between items-center">
            <h2 class="text-lg font-semibold text-white/95">AI智能体</h2>
            <button @click="toggleLeftSidebar" class="p-2 rounded-full hover:bg-white/10 text-white/40 hover:text-white transition-colors">
              <iconify-icon icon="lucide:chevron-left" class="text-xl"></iconify-icon>
            </button>
          </div>
          <div class="bg-black/20 border border-white/10 rounded-xl px-4 py-2.5 flex items-center gap-3">
            <iconify-icon icon="lucide:search" class="text-white/40 text-sm"></iconify-icon>
            <input v-model="searchQuery" type="text" placeholder="搜索AI智能体..." class="bg-transparent border-none outline-none text-sm text-white/70 w-full">
          </div>
        </div>
        <div class="p-4 overflow-y-auto flex-1 custom-scrollbar">
          <div v-for="user in filteredUsers" :key="user.id" @click="selectedUserId = user.id" :class="['p-4 rounded-xl cursor-pointer mb-2 transition-all', selectedUserId === user.id ? 'bg-[#3B9BFF]/20 border border-[#3B9BFF]/40' : 'bg-white/5 hover:bg-white/10']">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 rounded-full bg-white/10 flex items-center justify-center font-bold text-[#3B9BFF]">{{ user.name[0] }}</div>
              <div class="flex-1 min-w-0">
                <h3 class="text-sm font-medium text-white truncate">{{ user.name }}</h3>
                <p class="text-xs text-white/40 truncate">{{ user.description }}</p>
              </div>
            </div>
          </div>
        </div>
        <div class="w-1 cursor-col-resize hover:bg-[#3B9BFF]/50 bg-white/10 shrink-0" @mousedown="startDrag" :class="{ 'bg-[#3B9BFF]': isDragging }"></div>
      </aside>

      <main class="flex-1 flex flex-col min-w-0 min-h-0 overflow-hidden">
        <div class="h-16 border-b border-white/10 bg-white/5 flex items-center justify-between px-6 backdrop-blur-md">
          <h2 class="text-lg font-semibold text-white/95">{{ activeUserName }}</h2>
          <button @click="clearChatHistory" class="text-xs text-white/40 hover:text-white flex items-center gap-1"><iconify-icon icon="lucide:trash-2"></iconify-icon>清空记录</button>
        </div>
        <div ref="messageListRef" class="message-scroll flex-1 p-6 space-y-6 overflow-y-auto custom-scrollbar" @scroll="handleMessageListScroll">
          <div v-for="group in messageGroups" :key="group.id" class="flex" :class="group.kind === 'single' && group.message.type === 'user' ? 'justify-end' : 'justify-start'">
            <div v-if="group.kind === 'assistant'" class="w-full max-w-[760px] space-y-3">
              <div class="relative rounded-2xl p-4 border bg-white/10 text-white/90 border-white/10">
                <div class="mb-3 flex items-center gap-2">
                  <div class="w-8 h-8 rounded-lg bg-[#3B9BFF]/20 flex items-center justify-center text-[#3B9BFF] font-bold text-xs">{{ selectedUser?.name[0] }}</div>
                  <span class="text-xs font-semibold text-white/80">{{ selectedUser?.name }}</span>
                </div>
                <div v-if="group.thinking && hasMeaningfulThinking(group.thinking)" class="mb-3 p-3 bg-amber-500/10 border border-amber-500/20 rounded-lg text-xs text-amber-100/70">
                  <div class="font-bold mb-1 flex items-center gap-2"><iconify-icon icon="lucide:brain-circuit"></iconify-icon>推理过程</div>
                  <div class="whitespace-pre-wrap" v-html="renderMarkdown(group.thinking.content)"></div>
                </div>
                <div class="whitespace-pre-wrap" v-html="renderMarkdown(group.reply?.content || '')"></div>
                <div class="text-[10px] text-white/30 mt-2 text-right">{{ new Date(group.reply?.timestamp || Date.now()).toLocaleTimeString() }}</div>
              </div>
            </div>
            <div v-else-if="group.message.type === 'user'" class="flex flex-row-reverse items-start gap-4 w-full">
              <div class="flex flex-col items-center gap-1 shrink-0">
                <div class="w-10 h-10 rounded-full bg-[#50C878]/20 flex items-center justify-center text-[#50C878]"><iconify-icon icon="lucide:user" class="text-xl"></iconify-icon></div>
                <span class="text-[10px] text-white/40">Human</span>
              </div>
              <div class="max-w-[760px] rounded-2xl p-4 border bg-[#3B9BFF] text-white border-[#3B9BFF]">
                <div class="whitespace-pre-wrap">{{ group.message.content }}</div>
                <div class="text-[10px] text-white/40 mt-2 text-right">{{ new Date(group.message.timestamp).toLocaleTimeString() }}</div>
              </div>
            </div>
          </div>
        </div>
        <div class="p-6 border-t border-white/10 bg-white/5 backdrop-blur-md">
          <div class="flex gap-4">
            <textarea v-model="inputText" @keydown="handleTextareaKeydown" placeholder="输入消息..." rows="1" class="flex-1 bg-white/5 border border-white/10 rounded-xl p-3 text-sm text-white outline-none focus:border-[#3B9BFF]/50 resize-none"></textarea>
            <button @click="sendMessage" class="w-12 h-12 bg-[#3B9BFF] rounded-xl flex items-center justify-center text-white shadow-lg transition-all active:scale-95"><iconify-icon icon="lucide:send" class="text-xl"></iconify-icon></button>
          </div>
        </div>
      </main>

      <aside class="w-[320px] shrink-0 border-l border-white/10 bg-[#1A2536] flex flex-col p-6 space-y-8">
        <h3 class="text-lg font-semibold text-white/95">会话设置</h3>
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
      </aside>
    </div>
  </BaseLayout>
</template>

<style scoped>
:deep(main) { padding: 0 !important; }
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(59, 155, 255, 0.2); border-radius: 10px; }
.message-scroll { scrollbar-gutter: stable; }
</style>
